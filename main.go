package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/risor-io/risor"
)

type CrawlerFunc func(crawler *Crawler, query string) error

type Crawler struct {
	Response string
	Status   int
}

type CrawlerRegistry struct {
	crawlers map[string]CrawlerFunc
	mutex    sync.RWMutex
}

func NewCrawlerRegistry() *CrawlerRegistry {
	return &CrawlerRegistry{
		crawlers: make(map[string]CrawlerFunc),
	}
}

func (cr *CrawlerRegistry) Register(name string, callback CrawlerFunc) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	cr.crawlers[name] = callback
}

func (cr *CrawlerRegistry) Call(name string, query string) (*Crawler, error) {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()
	if callback, ok := cr.crawlers[name]; ok {
		crawler := &Crawler{}
		if err := callback(crawler, query); err != nil {
			return nil, err
		}
		return crawler, nil
	}
	return nil, fmt.Errorf("crawler not found: %s", name)
}

var defaultScript = `
print("script running...")
registry.register("google", func(crawler, query) {
	response := fetch("https://www.google.com/search?q=" + query)
	crawler.response = response.text()
	crawler.status = response.status_code
	printf("crawl complete for \"%s\" (status: %d)\n", query, crawler.status)
});
print("crawling...")
result := registry.call("google", "animals")
print("status:", result.status, "response len:", len(result.response))
`

func main() {
	var script string
	flag.StringVar(&script, "script", defaultScript, "path to the script file")
	flag.Parse()

	app := NewCrawlerRegistry()

	ctx := context.Background()

	_, err := risor.Eval(ctx, script, risor.WithGlobals(map[string]interface{}{
		"registry": NewCrawlerRegistryObject(app),
	}))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

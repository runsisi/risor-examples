package main

import (
	"fmt"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

const CrawlerType object.Type = "crawler.crawler"

type CrawlerObject struct {
	value *Crawler
}

func (c *CrawlerObject) Type() object.Type {
	return CrawlerType
}

func (c *CrawlerObject) Inspect() string {
	return "crawler.crawler()"
}

func (c *CrawlerObject) Interface() interface{} {
	return c.value
}

func (c *CrawlerObject) IsTruthy() bool {
	return true
}

func (c *CrawlerObject) Cost() int {
	return 0
}

func (c *CrawlerObject) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("type error: unable to marshal crawler.crawler")
}

func (c *CrawlerObject) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.Errorf("eval error: unsupported operation for %s: %v", CrawlerType, opType)
}

func (c *CrawlerObject) Equals(other object.Object) object.Object {
	return object.NewBool(c == other)
}

func (c *CrawlerObject) SetAttr(name string, value object.Object) error {
	switch name {
	case "response":
		value, err := object.AsString(value)
		if err != nil {
			return err.Value()
		}
		c.value.Response = value
		return nil
	case "status":
		value, err := object.AsInt(value)
		if err != nil {
			return err.Value()
		}
		c.value.Status = int(value)
		return nil
	default:
		return fmt.Errorf("attribute error: %s object has no attribute %q", CrawlerType, name)
	}
}

func (c *CrawlerObject) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "response":
		return object.NewString(c.value.Response), true
	case "status":
		return object.NewInt(int64(c.value.Status)), true
	}
	return nil, false
}

func NewCrawlerObject(c *Crawler) *CrawlerObject {
	return &CrawlerObject{value: c}
}

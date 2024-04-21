package main

import (
	"context"
	"fmt"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

const CrawlerRegistryType object.Type = "crawler.registry"

type CrawlerRegistryObject struct {
	value *CrawlerRegistry
}

func (r *CrawlerRegistryObject) Type() object.Type {
	return CrawlerRegistryType
}

func (r *CrawlerRegistryObject) Inspect() string {
	return "crawler.registry()"
}

func (r *CrawlerRegistryObject) Interface() interface{} {
	return r.value
}

func (r *CrawlerRegistryObject) IsTruthy() bool {
	return true
}

func (r *CrawlerRegistryObject) Cost() int {
	return 0
}

func (r *CrawlerRegistryObject) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("type error: unable to marshal crawler.registry")
}

func (r *CrawlerRegistryObject) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.Errorf("eval error: unsupported operation for %s: %v", CrawlerRegistryType, opType)
}

func (r *CrawlerRegistryObject) Equals(other object.Object) object.Object {
	return object.NewBool(r == other)
}

func (r *CrawlerRegistryObject) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("attribute error: %s object has no attribute %q", CrawlerRegistryType, name)
}

func (r *CrawlerRegistryObject) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "register":
		return object.NewBuiltin("crawler.registry.register", func(ctx context.Context, args ...object.Object) object.Object {
			if len(args) != 2 {
				return object.Errorf("wrong number of arguments. got=%d, want=2", len(args))
			}
			name, errObj := object.AsString(args[0])
			if errObj != nil {
				return errObj
			}
			fn, ok := args[1].(*object.Function)
			if !ok {
				return object.Errorf("argument error: expected function, got %s", args[1].Type())
			}
			callFunc, ok := object.GetCallFunc(ctx)
			if !ok {
				return object.Errorf("unable to get call function")
			}
			r.value.Register(name, func(crawler *Crawler, query string) error {
				c := NewCrawlerObject(crawler)
				_, err := callFunc(ctx, fn, []object.Object{c, object.NewString(query)})
				return err
			})
			return object.Nil
		}), true
	case "call":
		return object.NewBuiltin("crawler.registry.call", func(ctx context.Context, args ...object.Object) object.Object {
			if len(args) != 2 {
				return object.Errorf("wrong number of arguments. got=%d, want=2", len(args))
			}
			name, errObj := object.AsString(args[0])
			if errObj != nil {
				return errObj
			}
			query, errObj := object.AsString(args[1])
			if errObj != nil {
				return errObj
			}
			crawler, err := r.value.Call(name, query)
			if err != nil {
				return object.Errorf("crawler error: %v", err)
			}
			return NewCrawlerObject(crawler)
		}), true
	}
	return nil, false
}

func NewCrawlerRegistryObject(r *CrawlerRegistry) *CrawlerRegistryObject {
	return &CrawlerRegistryObject{value: r}
}

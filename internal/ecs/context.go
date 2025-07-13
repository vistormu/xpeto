package ecs

import (
	"fmt"
	"reflect"
)

type Context struct {
	resources map[reflect.Type]any
}

func NewContext() *Context {
	return &Context{
		resources: make(map[reflect.Type]any),
	}
}

func AddResource[T any](ctx *Context, resource T) {
	key := reflect.TypeOf(resource)
	ctx.resources[key] = resource
}

func GetResource[T any](ctx *Context) (T, bool) {
	res, ok := ctx.resources[reflect.TypeOf((*T)(nil)).Elem()]
	if !ok {
		var zero T
		return zero, false
	}
	return res.(T), true
}

func MustResource[T any](ctx *Context) T {
	res, ok := GetResource[T](ctx)
	if !ok {
		panic(fmt.Sprintf("resource of type %T not found in context", new(T)))
	}
	return res
}

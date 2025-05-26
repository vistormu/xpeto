package ecs

import (
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
	ctx.resources[reflect.TypeOf((*T)(nil)).Elem()] = resource
}

func GetResource[T any](ctx *Context) (T, bool) {
	res, ok := ctx.resources[reflect.TypeOf((*T)(nil)).Elem()]
	if !ok {
		var zero T
		return zero, false
	}
	return res.(T), true
}

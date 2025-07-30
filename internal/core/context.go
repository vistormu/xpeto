package core

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
	ctx.resources[reflect.TypeFor[T]()] = resource
}

func AddResourceByType(ctx *Context, resource any, type_ reflect.Type) {
	if type_ == nil {
		panic("type cannot be nil")
	}
	ctx.resources[type_] = resource
}

func GetResource[T any](ctx *Context) (T, bool) {
	v, ok := ctx.resources[reflect.TypeFor[T]()]
	if !ok {
		var zero T
		return zero, false
	}

	out, ok := v.(T)
	if !ok {
		var zero T
		return zero, false
	}

	return out, true
}

func MustResource[T any](ctx *Context) T {
	res, ok := GetResource[T](ctx)
	if !ok {
		panic(fmt.Sprintf("resource of type %T not found in context", new(T)))
	}
	return res
}

package xp

import (
	"github.com/vistormu/xpeto/core/ecs"
)

// #############
// CORE FEATURES
// #############

// ===
// ecs
// ===

// the world is the central
type World = ecs.World

// ------
// entity
// ------

// an entity is...
// packed [index:gen]
type Entity = ecs.Entity

// this function checks
var AddEntity = ecs.AddEntity

// tjisnmoljkn
var RemoveEntity = ecs.RemoveEntity

// flhnfkljf
var HasEntity = ecs.HasEntity

// ---------
// component
// ---------

func AddComponent[T any](w *World, e Entity, c T) bool {
	return ecs.AddComponent(w, e, c)
}

func GetComponent[T any](w *World, e Entity) (*T, bool) {
	return ecs.GetComponent[T](w, e)
}

func RemoveComponent[T any](w *World, e Entity) bool {
	return ecs.RemoveComponent[T](w, e)
}

// ------
// system
// ------

// a system operates on the components of entitues
type System = ecs.System

// ---------
// resources
// ---------

// add a resource to the world
// it is recommended to add the resource by value and not by reference
// the resource will be stored internally as value
// as it uses go generics, there can only be one value per type (ecs global singletons)
func AddResource[T any](w *World, r T) {
	ecs.AddResource(w, r)
}

var AddResourceByType = ecs.AddResourceByType

// the type `T` of the function cannot be a reference to a type
// if it is, it will return false
// the result will always be a reference to the resource so that
// the user can mutate its value
func GetResource[T any](w *World) (*T, bool) {
	return ecs.GetResource[T](w)
}

// the type `T` of the function cannot be a reference to a type
// if it is, it will return false
func RemoveResorce[T any](w *World) bool {
	return ecs.RemoveResource[T](w)
}

// -----
// query
// -----

type Filter = ecs.Filter

func With[T any]() Filter {
	return ecs.With[T]()
}

func Without[T any]() Filter {
	return ecs.Without[T]()
}

var Or = ecs.Or

func Query1[A any](w *World, filters ...Filter) *ecs.Query1[A] {
	return ecs.NewQuery1[A](w, filters...)
}

func Query2[A, B any](w *World, filters ...Filter) *ecs.Query2[A, B] {
	return ecs.NewQuery2[A, B](w, filters...)
}

func Query3[A, B, C any](w *World, filters ...Filter) *ecs.Query3[A, B, C] {
	return ecs.NewQuery3[A, B, C](w, filters...)
}

func Query4[A, B, C, D any](w *World, filters ...Filter) *ecs.Query4[A, B, C, D] {
	return ecs.NewQuery4[A, B, C, D](w, filters...)
}

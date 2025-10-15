package xp

import (
	"github.com/vistormu/xpeto/app"
	"github.com/vistormu/xpeto/pkg"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"

	"github.com/vistormu/xpeto/core/pkg/event"
	"github.com/vistormu/xpeto/core/pkg/state"
	"github.com/vistormu/xpeto/core/pkg/time"
	"github.com/vistormu/xpeto/core/pkg/transform"
	"github.com/vistormu/xpeto/core/pkg/window"

	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/font"
	"github.com/vistormu/xpeto/pkg/image"
	"github.com/vistormu/xpeto/pkg/input"
	"github.com/vistormu/xpeto/pkg/sprite"
	"github.com/vistormu/xpeto/pkg/text"
	"github.com/vistormu/xpeto/pkg/vector"
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

// description
var RemoveEntity = ecs.RemoveEntity

// description
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

// a system operates on the components of entities
type System = ecs.System

// `GetSystemId` retrieves the unique identifier of the current running system
//
// the id is provided by the scheduler at system registration time (`xp.AddSystem`)
var GetSystemId = ecs.GetSystemId

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
func RemoveResource[T any](w *World) bool {
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

// ========
// schedule
// ========

// ---------
// condition
// ---------

// a condition function is defined as
//
// func(*World) bool
//
// it takes the ecs world as a parameter and returns a boolean indicating if the system
// should or should not run
type ConditionFn = schedule.ConditionFn

// description
var Once = schedule.Once

// description
var OnceWhen = schedule.OnceWhen

// ---------
// scheduler
// ---------

// the `Scheduler` is the main struct to run
//
// methods:
// - `RunIf(condition ConditionFn)`: add a condition for the execution of a system
type Scheduler = schedule.Scheduler

// adds a system to an schedule
var AddSystem = schedule.AddSystem

// -----
// stage
// -----

// description
var PreStartup = schedule.PreStartup

// description
var Startup = schedule.Startup

// description
var PostStartup = schedule.PostStartup

// description
var First = schedule.First

// description
var PreUpdate = schedule.PreUpdate

// description
var FixedFirst = schedule.FixedFirst

// description
var FixedPreUpdate = schedule.FixedPreUpdate

// description
var FixedUpdate = schedule.FixedUpdate

// description
var FixedPostUpdate = schedule.FixedPostUpdate

// description
var FixedLast = schedule.FixedLast

// description
var Update = schedule.Update

// description
var PostUpdate = schedule.PostUpdate

// description
var Last = schedule.Last

// description
var PreDraw = schedule.PreDraw

// description
var Draw = schedule.Draw

// description
var PostDraw = schedule.PostDraw

// #############
// CORE PACKAGES
// #############

// =====
// event
// =====

// description
func AddEvent[T any](w *World, ev T) {
	event.AddEvent(w, ev)
}

// description
func GetEvents[T any](w *World) ([]T, bool) {
	return event.GetEvents[T](w)
}

// =====
// state
// =====

// ---------
// condition
// ---------

// description
func InState[T comparable](s T) ConditionFn {
	return state.InState(s)
}

// -----
// stage
// -----

// description
func OnExit[T comparable](s T) schedule.StageFn {
	return state.OnExit(s)
}

// description
func OnTransition[T comparable](from, to T) schedule.StageFn {
	return state.OnTransition(from, to)
}

// description
func OnEnter[T comparable](s T) schedule.StageFn {
	return state.OnEnter(s)
}

// -----
// state
// -----

// description
func AddStateMachine[T comparable](sch *Scheduler, initial T) {
	state.AddStateMachine(sch, initial)
}

// description
func GetState[T comparable](w *World) (T, bool) {
	return state.GetState[T](w)
}

// description
func SetNextState[T comparable](w *World, s T) bool {
	return state.SetNextState(w, s)
}

// ====
// time
// ====

// -----
// clock
// -----

// description
type ClockSettings = time.ClockSettings

// description
type RealClock = time.RealClock

// description
type VirtualClock = time.VirtualClock

// description
type FixedClock = time.FixedClock

// ---------
// condition
// ---------

// =========
// transform
// =========

// description
type Transform = transform.Transform

// ======
// window
// ======

// description
type Layout = window.Layout

// description
type Screen = window.Screen

// description
type WindowSettings = window.WindowSettings

// ###
// APP
// ###

// ===
// app
// ===

// description
var NewApp = app.NewApp

// =======
// runners
// =======

// description
const Ebiten = app.Ebiten

// description
const Headless = app.Headless

// ################
// DEFAULT PACKAGES
// ################

var DefaultPkgs = pkg.DefaultPkgs

// =====
// asset
// =====

// ---------
// condition
// ---------

// description
func IsAssetLoaded[B any]() ConditionFn {
	return asset.IsAssetLoaded[B]()
}

// -----
// event
// -----

// ------
// handle
// ------

// description
type Handle = asset.Handle

// ------
// loader
// ------

// ------
// server
// ------

// description
var SetFileSystem = asset.SetFileSystem

// description
var AddAssetLoader = asset.AddAssetLoader

// description
func AddAssets[T, B any](w *World) {
	asset.AddAssets[T, B](w)
}

// description
func GetAsset[T any](w *World, handle Handle) (T, bool) {
	return asset.GetAsset[T](w, handle)
}

// =====
// audio
// =====

// ====
// font
// ====

// description
type Font = font.Font

// =====
// image
// =====

// description
type Image = image.Image

// =====
// input
// =====

// ------
// events
// ------

// -------
// gamepad
// -------

// --------
// keyboard
// --------

// description
type Key = input.Key

const (
	KeyA = input.KeyA
	KeyD = input.KeyD
	KeyS = input.KeyS
	KeyW = input.KeyW

	KeyEnter = input.KeyEnter

	KeyArrowDown = input.KeyArrowDown
	KeyArrowUp   = input.KeyArrowUp
)

// description
type Keyboard = input.Keyboard

// -----
// mouse
// -----

// description
type Mouse = input.Mouse

// ======
// render
// ======

// ======
// sprite
// ======

// description
type Sprite = sprite.Sprite

// ====
// text
// ====

const (
	AlignStart  = text.AlignStart
	AlignCenter = text.AlignCenter
	AlignEnd    = text.AlignEnd
)

type Text = text.Text

// ======
// vector
// ======

// description
type Circle = vector.Circle

// description
type Rect = vector.Rect

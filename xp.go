package xp

import (
	"github.com/vistormu/go-dsa/constraints"
	"github.com/vistormu/xpeto/app"

	"github.com/vistormu/xpeto/core"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"

	"github.com/vistormu/xpeto/core/event"
	"github.com/vistormu/xpeto/core/log"
	"github.com/vistormu/xpeto/core/state"
	"github.com/vistormu/xpeto/core/time"
	"github.com/vistormu/xpeto/core/window"

	"github.com/vistormu/xpeto/pkg"
	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/geometry"
	"github.com/vistormu/xpeto/pkg/input"
	"github.com/vistormu/xpeto/pkg/render"
	"github.com/vistormu/xpeto/pkg/shape"
	"github.com/vistormu/xpeto/pkg/sprite"
	"github.com/vistormu/xpeto/pkg/text"
	"github.com/vistormu/xpeto/pkg/transform"
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

// description
func AddComponent[T any](w *World, e Entity, c T) bool {
	return ecs.AddComponent(w, e, c)
}

// description
func GetComponent[T any](w *World, e Entity) (*T, bool) {
	return ecs.GetComponent[T](w, e)
}

// description
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

// descriptopm
var Exit = schedule.Exit

// #############
// CORE PACKAGES
// #############

type Pkg = core.Pkg

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

// ===
// log
// ===
type LogLevel = log.Level

const (
	Debug   = log.Debug
	Info    = log.Info
	Warning = log.Warning
	Error   = log.Error
	Fatal   = log.Fatal
)

var SetLogLevel = log.SetLogLevel

var LogDebug = log.LogDebug

var LogInfo = log.LogInfo

var LogWarning = log.LogWarning

var LogError = log.LogError

var LogFatal = log.LogFatal

var F = log.F

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

// description
var SetTPS = time.SetTPS

// description
var SetFixedDelta = time.SetFixedDelta

// description
var PauseClock = time.PauseClock

// ---------
// condition
// ---------

// ======
// window
// ======

// -------
// scaling
// -------

// description
type ScalingMode = window.ScalingMode

const (
	ScalingFree    ScalingMode = window.ScalingFree
	ScalingInteger ScalingMode = window.ScalingInteger
	ScalingHiDPI   ScalingMode = window.ScalingHiDPI
)

// description
var SetScalingMode = window.SetScalingMode

// description
var SetPixelSnap = window.SetPixelSnap

// description
// var GetDesiredVirtualSize = window.GetDesiredVirtualSize

// --------
// viewport
// --------

// description
type Viewport = window.Viewport

// description
var ComputeViewport = window.ComputeViewport

// ------
// window
// ------

// description
type ResizingMode = window.ResizingMode

const (
	ResizingDisabled   ResizingMode = window.ResizingModeDisabled
	ResizingEnabled    ResizingMode = window.ResizingModeEnabled
	ResizingFullscreen ResizingMode = window.ResizingModeOnlyFullscreenEnabled
)

// description
type WindowAction = window.WindowAction

const (
	ActionNone     WindowAction = window.ActionNone
	ActionMaximize WindowAction = window.ActionMaximize
	ActionMinimize WindowAction = window.ActionMinimize
	ActionRestore  WindowAction = window.ActionRestore
)

// description
type RealWindow = window.RealWindow

// description
type VirtualWindow = window.VirtualWindow

// description
var SetRealWindowSize = window.SetRealWindowSize

// description
func GetRealWindowSize[T constraints.Number](w *World) (width, height T) {
	return window.GetRealWindowSize[T](w)
}

// description
var SetFullScreen = window.SetFullScreen

// description
var SetAntiAliasing = window.SetAntiAliasing

// description
var SetVSync = window.SetVSync

// description
var SetRunnableOnUnfocused = window.SetRunnableOnUnfocused

// description
var SetResizingMode = window.SetResizingMode

// description
var SetWindowSizeLimits = window.SetWindowSizeLimits

// description
var MaximizeWindow = window.MaximizeWindow

// description
var MinimizeWindow = window.MinimizeWindow

// description
var RestoreWindow = window.RestoreWindow

// description
var SetVirtualWindowSize = window.SetVirtualWindowSize

// description
func GetVirtualWindowSize[T constraints.Number](w *World) (width, height T) {
	return window.GetVirtualWindowSize[T](w)
}

// ###
// APP
// ###

// description
var NewApp = app.NewApp

// description
type ExitApp = app.ExitApp

// ################
// DEFAULT PACKAGES
// ################

var DefaultPkgs = pkg.DefaultPkgs

// =====
// asset
// =====

// -----
// asset
// -----

// description
type Asset = asset.Asset

// ------
// loader
// ------

type LoaderFn[T any] = asset.LoaderFn[T]

// ------
// server
// ------

// description
var AddStaticFS = asset.AddStaticFS

// description
func AddLoaderFn[T any](w *World, fn asset.LoaderFn[T], extensions ...string) {
	asset.AddLoaderFn(w, fn, extensions...)
}

// description
func AddAsset[T any](w *World) {
	asset.AddAsset[T](w)
}

// description
func GetAsset[T any](w *World, a Asset) (*T, bool) {
	return asset.GetAsset[T](w, a)
}

// description
func RemoveAsset[T any](w *World, a Asset) bool {
	return asset.RemoveAsset[T](w, a)
}

// =====
// audio
// =====

// ====
// text
// ====

type Align = text.Align

const (
	AlignStart  Align = text.AlignStart
	AlignCenter Align = text.AlignCenter
	AlignEnd    Align = text.AlignEnd
)

type Text = text.Text

type DefaultFonts = text.DefaultFonts

// ========
// geometry
// ========

type Vector[T constraints.Number] = geometry.Vector[T]

// ======
// sprite
// ======

// description
type Sprite = sprite.Sprite

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
	KeyA Key = input.KeyA
	KeyD Key = input.KeyD
	KeyS Key = input.KeyS
	KeyW Key = input.KeyW

	KeyEnter Key = input.KeyEnter

	KeyArrowDown Key = input.KeyArrowDown
	KeyArrowUp   Key = input.KeyArrowUp
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

// ------
// anchor
// ------

// description
type Anchor = render.Anchor

const (
	AnchorCenter      Anchor = render.AnchorCenter
	AnchorTopLeft     Anchor = render.AnchorTopLeft
	AnchorTop         Anchor = render.AnchorTop
	AnchorTopRight    Anchor = render.AnchorTopRight
	AnchorLeft        Anchor = render.AnchorLeft
	AnchorRight       Anchor = render.AnchorRight
	AnchorBottomLeft  Anchor = render.AnchorBottomLeft
	AnchorBottom      Anchor = render.AnchorBottom
	AnchorBottomRight Anchor = render.AnchorBottomRight
)

// =====
// shape
// =====

// description
type Shape = shape.Shape

// description
type EllipeShape = shape.Ellipse

var NewCircleShape = shape.NewCircle

// description
type RectShape = shape.Rect

var NewRectShape = shape.NewRect

// description
type PathShape = shape.Path

var NewPathShape = shape.NewPath

// desctiption
type SegmentShape = shape.Segment

var NewSegmentShape = shape.NewSegment

// =========
// transform
// =========

type Transform = transform.Transform

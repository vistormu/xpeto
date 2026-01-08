package xp

import (
	"github.com/vistormu/go-dsa/constraints"
	"github.com/vistormu/go-dsa/geometry"
	"github.com/vistormu/xpeto/app"

	"github.com/vistormu/xpeto/core"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"

	"github.com/vistormu/xpeto/core/clock"
	"github.com/vistormu/xpeto/core/event"
	"github.com/vistormu/xpeto/core/log"
	"github.com/vistormu/xpeto/core/window"

	"github.com/vistormu/xpeto/pkg"
	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/input"
	"github.com/vistormu/xpeto/pkg/render"
	"github.com/vistormu/xpeto/pkg/shape"
	"github.com/vistormu/xpeto/pkg/sprite"
	"github.com/vistormu/xpeto/pkg/text"
	"github.com/vistormu/xpeto/pkg/transform"
)

// #############
// core features
// #############

// ===
// ecs
// ===

// stores all entities, components, and resources
//
// acts as the shared context passed into systems
type World = ecs.World

// ------
// entity
// ------

// identifies an entity inside a world
//
// packs index and generation into one value
type Entity = ecs.Entity

// spawns a new entity and returns its id
//
// complexity: O(1)
var AddEntity = ecs.AddEntity

// despawns an entity and removes all its components
//
// returns false if the entity is not alive
//
// complexity: O(C), where C is the number of registered component types
var RemoveEntity = ecs.RemoveEntity

// checks if an entity is alive in the world
//
// complexity: O(1)
var HasEntity = ecs.HasEntity

// ---------
// component
// ---------

// attaches a component value of type T to an entity
//
// returns false if the entity is not alive
//
// returns false if T is a pointer type
//
// complexity: amortized O(1)
func AddComponent[T any](w *World, e Entity, c T) bool {
	return ecs.AddComponent(w, e, c)
}

// removes the component of type T from an entity
//
// returns false if the entity is not alive or the component is missing
//
// complexity: O(1)
func RemoveComponent[T any](w *World, e Entity) bool {
	return ecs.RemoveComponent[T](w, e)
}

// retrieves a pointer to the component value of type T
//
// returns false if the entity is not alive or the component is missing
//
// pointer can become stale after structural changes that grow the store for T
//
// complexity: O(1)
func GetComponent[T any](w *World, e Entity) (*T, bool) {
	return ecs.GetComponent[T](w, e)
}

// checks if the entity has a component of type T
//
// returns false if the entity is not alive
//
// complexity: O(1)
func HasComponent[T any](w *World, e Entity) bool {
	return ecs.HasComponent[T](w, e)
}

// ------
// system
// ------

// defines a unit of work executed by the scheduler
type System = ecs.System

// --------
// resource
// --------

// stores a resource value of type T in the world
//
// returns false if T is a pointer type
//
// overwrites any previous resource of the same type
//
// complexity: O(1)
func AddResource[T any](w *World, r T) bool {
	return ecs.AddResource(w, r)
}

// ensures a resource of type T exists, initializing it once when missing
//
// returns a pointer to the stored resource
//
// returns nil only if T is a pointer type (resource storage rejects pointer types)
//
// complexity: O(1)
func EnsureResource[T any](w *World, init func() T) *T {
	return ecs.EnsureResource(w, init)
}

// removes a stored resource of type T
//
// returns false if T is a pointer type or the resource is missing
//
// complexity: O(1)
func RemoveResource[T any](w *World) bool {
	return ecs.RemoveResource[T](w)
}

// retrieves a pointer to a stored resource of type T
//
// returns false if T is a pointer type or the resource is missing
//
// complexity: O(1)
func GetResource[T any](w *World) (*T, bool) {
	return ecs.GetResource[T](w)
}

// checks if a resource of type T exists
//
// returns false if T is a pointer type
//
// complexity: O(1)
func HasResource[T any](w *World) bool {
	return ecs.HasResource[T](w)
}

// -----
// query
// -----

// selects entities during query iteration
//
// note: filters must be non nil
type Filter = ecs.Filter

// keeps entities that have a component of type T
//
// complexity: O(1)
func With[T any]() Filter {
	return ecs.With[T]()
}

// keeps entities that do not have a component of type T
//
// complexity: O(1)
func Without[T any]() Filter {
	return ecs.Without[T]()
}

// combines filters with logical or
//
// complexity: O(k), where k is the number of filters passed to or
var Or = ecs.Or

// iterates entities that have component A
//
// complexity: O(m · k), where m is the number of entities that have A and k is the number of filters
func Query1[A any](w *World, filters ...Filter) *ecs.Query1[A] {
	return ecs.NewQuery1[A](w, filters...)
}

// iterates entities that have components A and B
//
// complexity: O(m · k), where m is the number of entities in the smallest store among A and B and k is the number of filters
func Query2[A, B any](w *World, filters ...Filter) *ecs.Query2[A, B] {
	return ecs.NewQuery2[A, B](w, filters...)
}

// iterates entities that have components A, B, and C
//
// complexity: O(m · k), where m is the number of entities in the smallest store among A, B, and C and k is the number of filters
func Query3[A, B, C any](w *World, filters ...Filter) *ecs.Query3[A, B, C] {
	return ecs.NewQuery3[A, B, C](w, filters...)
}

// iterates entities that have components A, B, C, and D
//
// complexity: O(m · k), where m is the number of entities in the smallest store among A, B, C, and D and k is the number of filters
func Query4[A, B, C, D any](w *World, filters ...Filter) *ecs.Query4[A, B, C, D] {
	return ecs.NewQuery4[A, B, C, D](w, filters...)
}

// ========
// schedule
// ========

// ---------
// condition
// ---------

// a condition function returns true when a system should run
type ConditionFn = schedule.ConditionFn

// returns a condition that evaluates to true only once
var Once = schedule.Once

// returns a condition that evaluates to true once, when the provided predicate becomes true
var OnceWhen = schedule.OnceWhen

// returns true when the current state of T equals s
//
// it returns false if the state machine is missing
func IsInState[T comparable](s T) ConditionFn {
	return schedule.IsInState(s)
}

// ---------
// scheduler
// ---------
// compiles and runs systems in stages
type Scheduler = schedule.Scheduler

// adds a system to a stage in a scheduler
var AddSystem = schedule.AddSystem

// provides fluent options for AddSystem
var SystemOpt = schedule.SystemOpt

// -----
// stage
// -----

var PreStartup = schedule.PreStartup
var Startup = schedule.Startup
var PostStartup = schedule.PostStartup
var First = schedule.First
var PreUpdate = schedule.PreUpdate

func OnExit[T comparable](s T) schedule.Stage              { return schedule.OnExit(s) }
func OnTransition[T comparable](from, to T) schedule.Stage { return schedule.OnTransition(from, to) }
func OnEnter[T comparable](s T) schedule.Stage             { return schedule.OnEnter(s) }

var FixedFirst = schedule.FixedFirst
var FixedPreUpdate = schedule.FixedPreUpdate
var FixedUpdate = schedule.FixedUpdate
var FixedPostUpdate = schedule.FixedPostUpdate
var FixedLast = schedule.FixedLast
var Update = schedule.Update
var PostUpdate = schedule.PostUpdate
var Last = schedule.Last
var PreDraw = schedule.PreDraw
var Draw = schedule.Draw
var PostDraw = schedule.PostDraw
var Exit = schedule.Exit

// -----
// state
// -----

// registers the state machine for T and wires transition systems
func AddStateMachine[T comparable](sch *Scheduler, initial T) {
	schedule.AddStateMachine(sch, initial)
}

// returns the current state
func GetState[T comparable](w *World) (T, bool) {
	return schedule.GetState[T](w)
}

// requests a transition for the next update
//
// returns false if the state machine is missing
func SetNextState[T comparable](w *World, s T) bool {
	return schedule.SetNextState(w, s)
}

// #############
// CORE PACKAGES
// #############

type Pkg = core.Pkg

// =====
// clock
// =====

// --------
// settings
// --------

// selects how the clock produces deltas and fixed steps
type ClockMode = clock.ClockMode

const (
	// `uses `ClockSettings.FixedDelta` to produce fixed steps
	ClockModeFixed ClockMode = clock.ModeFixed

	// indicates that a backend drives the tick rate and should keep `FixedClock.Delta` updated
	ClockModeSyncWithFPS ClockMode = clock.ModeSyncWithFPS
)

// stores clock intent and safety limits
type ClockSettings = clock.ClockSettings

// sets the clock mode
var SetClockMode = clock.SetMode

// configures the clock to run in fixed mode with tps steps per second
//
// note: tps must be > 0
var SetTPS = clock.SetTPS

// configures the clock to run in fixed mode with the provided step delta
//
// d <= 0 falls back to defaultfixeddelta
var SetFixedDelta = clock.SetFixedDelta

// sets the virtual time scale
var SetClockScale = clock.SetScale

// pauses or resumes virtual time
var PauseClock = clock.PauseClock

// sets the maximum number of fixed steps that can be produced per tick
var SetMaxSteps = clock.SetMaxSteps

// clamps the real delta before scaling
var SetMaxDelta = clock.SetMaxDelta

// clamps the scaled virtual delta, 0 disables the clamp
var SetMaxVirtualDelta = clock.SetMaxVirtualDelta

// -----
// clock
// -----

// stores real time measured by the clock tick
type RealClock = clock.RealClock

// stores virtual time derived from clamped real time and scale
type VirtualClock = clock.VirtualClock

// stores fixed step state derived from virtual time
type FixedClock = clock.FixedClock

// ---------
// condition
// ---------

// returns true every n virtual frames
//
// uses `VirtualClock.Frame`
var EveryNFrames = clock.EveryNFrames

// returns true every n fixed frames
//
// uses `FixedClock.Frame`
var EveryNFixedFrames = clock.EveryNFixedFrames

// returns a condition that becomes true once, when virtual elapsed reaches d
var OnceAfterElapsed = clock.OnceAfterElapsed

// returns true while virtual elapsed is at least d
var AfterElapsed = clock.AfterElapsed

// returns a condition that becomes true once, when real elapsed reaches d
var OnceAfterRealElapsed = clock.OnceAfterRealElapsed

// returns true while real elapsed is at least d
var AfterRealElapsed = clock.AfterRealElapsed

// returns true at a fixed cadence based on virtual elapsed
//
// the cadence is anchored to the clock start, not to the first evaluation
var EveryDuration = clock.EveryDuration

// returns true when the last tick produced at least one fixed step
var EveryFixedSteps = clock.EveryFixedSteps

// =====
// event
// =====

// sends an event value of type T into the event bus
//
// events are buffered by type and consumed per system id
//
// emission is safe to call from other goroutines, but consumption should happen inside scheduled systems
func AddEvent[T any](w *World, ev T) {
	event.AddEvent(w, ev)
}

// retrieves unread events of type T for the currently running system
//
// it returns ok=false when there are no unread events for this system
//
// events are isolated per system id, so different systems can read the same event stream independently
//
// note: it requires ecs.RunningSystem to exist in the world
func GetEvents[T any](w *World) ([]T, bool) {
	return event.GetEvents[T](w)
}

// ===
// log
// ===

// the severity level used by the logger
type LogLevel = log.Level

const (
	// intended for development diagnostics
	Debug = log.Debug

	// intended for high level engine or game state messages
	Info = log.Info

	// intended for unexpected but recoverable situations
	Warning = log.Warning

	// intended for failures that affect behaviour
	Error = log.Error

	// intended for unrecoverable failures. it does not exit the process, it only records the message
	Fatal = log.Fatal
)

// sets the minimum level that will be recorded
var SetLogLevel = log.SetLogLevel

// description
var FlushLoggerManually = log.FlushLoggerManually

// description
var FlushLoggerEveryFrame = log.FlushLoggerEveryFrame

// description
var FlushLoggerEveryNFrames = log.FlushLoggerEveryNFrames

// description
var FlushLoggerEveryNRecords = log.FlushLoggerEveryNRecords

// description
var SilenceLogLevels = log.SilenceLevels

// description
var UnsilenceLogLevels = log.UnsilenceLevels

// records a Debug message
var LogDebug = log.LogDebug

// records an Info message
var LogInfo = log.LogInfo

// records a Warning message
var LogWarning = log.LogWarning

// records an Error message
var LogError = log.LogError

// records a Fatal message
var LogFatal = log.LogFatal

// creates a structured field in the form "key: value"
var F = log.F

// appends a log sink that will receive flushed records
var AddLogSink = log.AddSink

// removes all log sinks
var ClearLogSinks = log.ClearSinks

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

// composes the engine, loads core packages, loads backend, loads user packages, and runs the backend loop
//
// order:
// - creates world and scheduler
// - loads core packages
// - constructs backend (backend can register systems)
// - loads user packages
// - runs startup
// - runs backend loop
// - runs exit on return
type App = app.App

// builds a backend instance for a given world and scheduler
type BackendFactory = app.BackendFactory

// backend main loop
type Backend = app.Backend

// app options builder
var AppOpt = app.AppOpt

// creates a new app using the provided backend factory
var NewApp = app.NewApp

// request to exit the app loop
//
// backends may observe this event and stop the loop
type ExitAppEvent = app.ExitAppEvent

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

// describes the lifecycle of an asset handle
type AssetState = asset.AssetState

const (
	AssetRequested AssetState = asset.AssetRequested
	AssetLoading   AssetState = asset.AssetLoading
	AssetLoaded    AssetState = asset.AssetLoaded
	AssetFailed    AssetState = asset.AssetFailed
)

// identifies an asset instance
type Asset = asset.Asset

// ------
// loader
// ------

// decodes raw bytes into an asset value
type LoaderFn[T any] = asset.LoaderFn[T]

// ----------
// conditions
// ----------

// returns true when asset a is loaded
var WhenAssetLoaded = asset.WhenAssetLoaded

// returns true when asset a is failed
var WhenAssetFailed = asset.WhenAssetFailed

// returns true when asset a matches st
var WhenAssetState = asset.WhenAssetState

// returns true when all assets are loaded, false when the list is empty
var WhenAllAssetsLoaded = asset.WhenAllAssetsLoaded

// returns true when any asset is failed, false when the list is empty
var WhenAnyAssetFailed = asset.WhenAnyAssetFailed

// returns true when all Asset fields in bundle resource T are loaded
//
// requirements:
// - T must be a struct resource stored by value in the world
// - all fields of type xp.asset must be non zero and loaded
func WhenBundleLoaded[T any]() ConditionFn {
	return asset.WhenBundleLoaded[T]()
}

// returns true when any Asset field in bundle resource T is failed
//
// note: asset.whenbundlefailed currently has a bug and always returns false.
// fix it in pkg/asset before relying on this api.
func WhenBundleFailed[T any]() ConditionFn {
	return asset.WhenBundleFailed[T]()
}

// ------
// server
// ------

// mounts an fs under a base name used in paths like "base/path/to/file.ext"
var AddStaticFS = asset.AddStaticFS

// registers a loader for one or more extensions (".png", "png", ...)
func AddLoaderFn[T any](w *World, fn asset.LoaderFn[T], extensions ...string) {
	asset.AddLoaderFn(w, fn, extensions...)
}

// requests loading all assets declared as `asset.Asset` fields in bundle T.
//
// each field must include a struct tag: `path:"base/rel.ext"`
//
// the created bundle is stored as a resource of type T.
func AddAsset[T any](w *World) {
	asset.AddAsset[T](w)
}

// removes the asset from the server and invalidates the handle
func RemoveAsset[T any](w *World, a Asset) bool {
	return asset.RemoveAsset[T](w, a)
}

// retrieves the loaded asset value of type T by handle
func GetAsset[T any](w *World, a Asset) (*T, bool) {
	return asset.GetAsset[T](w, a)
}

// reads the current state for an asset handle
var GetAssetState = asset.GetAssetState

// convenience: returns true when asset is loaded
var IsAssetLoaded = asset.IsAssetLoaded

// returns the original request path for an asset handle
var GetAssetPath = asset.GetAssetPath

// returns the error for a failed asset handle
var GetAssetError = asset.GetAssetError

// =====
// audio
// =====

// ====
// text
// ====

// ----
// align
// ----

type Align = text.Align

const (
	AlignStart  Align = text.AlignStart
	AlignCenter Align = text.AlignCenter
	AlignEnd    Align = text.AlignEnd
)

// ----
// wrap
// ----

type WrapMode = text.WrapMode

const (
	WrapNone WrapMode = text.WrapNone
	WrapWord WrapMode = text.WrapWord
	WrapRune WrapMode = text.WrapRune
)

// ----
// text
// ----

// 2d text component rendered by the active backend
//
// notes:
//
// - if Font == 0, the backend will try to use the default font bundle.
//
// - MaxWidth and Wrap are backend dependent in v0.1.0 (ebiten font currently treats text as unwrapped).
type Text = text.Text

// creates a Text component with sane defaults
//
// defaults:
//
// - size: 18
//
// - color: white
//
// - align: start
//
// - wrap: none
//
// - maxwidth: 0
//
// - orderkey: layer 0, order 0, tie 0
var NewText = text.NewText

// fluent options for NewText
var TextOpt = text.TextOpt

// default font bundle, loaded through the asset system
type DefaultFont = text.DefaultFont

// retrieves the default font asset handle when loaded
var GetDefaultFont = text.GetDefaultFont

// ========
// geometry
// ========

type Vector[T constraints.Number] = geometry.Vector[T]

// ======
// sprite
// ======

// 2d sprite component rendered by the active backend
//
// notes:
//
// - Image is an asset handle (use the asset plugin to load it).
//
// - OrderKey controls render ordering inside the sprite render phase.
//
// - it uses a default center anchor, and it reads render.anchor when present.
//
// - pixel snapping follows window.Scaling.SnapPixels.
type Sprite = sprite.Sprite

// creates a Sprite with sane defaults
//
// defaults:
// - orderkey: layer 0, order 0, tie 0
var NewSprite = sprite.NewSprite

// =====
// input
// =====

// -----
// focus
// -----

// emitted by backends when the window focus changes
type FocusChangedEvent = input.FocusChangedEvent

// --------
// keyboard
// --------

type Key = input.Key

const (
	// letters
	KeyA Key = input.KeyA
	KeyB Key = input.KeyB
	KeyC Key = input.KeyC
	KeyD Key = input.KeyD
	KeyE Key = input.KeyE
	KeyF Key = input.KeyF
	KeyG Key = input.KeyG
	KeyH Key = input.KeyH
	KeyI Key = input.KeyI
	KeyJ Key = input.KeyJ
	KeyK Key = input.KeyK
	KeyL Key = input.KeyL
	KeyM Key = input.KeyM
	KeyN Key = input.KeyN
	KeyO Key = input.KeyO
	KeyP Key = input.KeyP
	KeyQ Key = input.KeyQ
	KeyR Key = input.KeyR
	KeyS Key = input.KeyS
	KeyT Key = input.KeyT
	KeyU Key = input.KeyU
	KeyV Key = input.KeyV
	KeyW Key = input.KeyW
	KeyX Key = input.KeyX
	KeyY Key = input.KeyY
	KeyZ Key = input.KeyZ

	// digits
	Key0 Key = input.Key0
	Key1 Key = input.Key1
	Key2 Key = input.Key2
	Key3 Key = input.Key3
	Key4 Key = input.Key4
	Key5 Key = input.Key5
	Key6 Key = input.Key6
	Key7 Key = input.Key7
	Key8 Key = input.Key8
	Key9 Key = input.Key9

	// numpad
	KeyNumpad0        Key = input.KeyNumpad0
	KeyNumpad1        Key = input.KeyNumpad1
	KeyNumpad2        Key = input.KeyNumpad2
	KeyNumpad3        Key = input.KeyNumpad3
	KeyNumpad4        Key = input.KeyNumpad4
	KeyNumpad5        Key = input.KeyNumpad5
	KeyNumpad6        Key = input.KeyNumpad6
	KeyNumpad7        Key = input.KeyNumpad7
	KeyNumpad8        Key = input.KeyNumpad8
	KeyNumpad9        Key = input.KeyNumpad9
	KeyNumpadAdd      Key = input.KeyNumpadAdd
	KeyNumpadDecimal  Key = input.KeyNumpadDecimal
	KeyNumpadDivide   Key = input.KeyNumpadDivide
	KeyNumpadEnter    Key = input.KeyNumpadEnter
	KeyNumpadEqual    Key = input.KeyNumpadEqual
	KeyNumpadMultiply Key = input.KeyNumpadMultiply
	KeyNumpadSubtract Key = input.KeyNumpadSubtract

	// punctuation and symbols
	KeyBracketLeft   Key = input.KeyBracketLeft
	KeyBracketRight  Key = input.KeyBracketRight
	KeyComma         Key = input.KeyComma
	KeyBackspace     Key = input.KeyBackspace
	KeyBackslash     Key = input.KeyBackslash
	KeyEqual         Key = input.KeyEqual
	KeyBackquote     Key = input.KeyBackquote
	KeyIntlBackslash Key = input.KeyIntlBackslash
	KeyMinus         Key = input.KeyMinus
	KeyPeriod        Key = input.KeyPeriod
	KeyQuote         Key = input.KeyQuote
	KeySemicolon     Key = input.KeySemicolon
	KeySlash         Key = input.KeySlash
	KeySpace         Key = input.KeySpace
	KeyTab           Key = input.KeyTab

	// arrows
	KeyArrowDown  Key = input.KeyArrowDown
	KeyArrowLeft  Key = input.KeyArrowLeft
	KeyArrowRight Key = input.KeyArrowRight
	KeyArrowUp    Key = input.KeyArrowUp

	// modifiers (left and right)
	KeyAltLeft      Key = input.KeyAltLeft
	KeyAltRight     Key = input.KeyAltRight
	KeyControlLeft  Key = input.KeyControlLeft
	KeyControlRight Key = input.KeyControlRight
	KeyMetaLeft     Key = input.KeyMetaLeft
	KeyMetaRight    Key = input.KeyMetaRight
	KeyShiftLeft    Key = input.KeyShiftLeft
	KeyShiftRight   Key = input.KeyShiftRight

	// modifiers (generic)
	KeyAlt     Key = input.KeyAlt
	KeyControl Key = input.KeyControl
	KeyShift   Key = input.KeyShift
	KeyMeta    Key = input.KeyMeta

	// locks and system
	KeyCapsLock    Key = input.KeyCapsLock
	KeyContextMenu Key = input.KeyContextMenu
	KeyDelete      Key = input.KeyDelete
	KeyEnd         Key = input.KeyEnd
	KeyEnter       Key = input.KeyEnter
	KeyEscape      Key = input.KeyEscape
	KeyHome        Key = input.KeyHome
	KeyInsert      Key = input.KeyInsert
	KeyNumLock     Key = input.KeyNumLock
	KeyPageDown    Key = input.KeyPageDown
	KeyPageUp      Key = input.KeyPageUp
	KeyPause       Key = input.KeyPause
	KeyPrintScreen Key = input.KeyPrintScreen
	KeyScrollLock  Key = input.KeyScrollLock

	// function keys
	KeyF1  Key = input.KeyF1
	KeyF2  Key = input.KeyF2
	KeyF3  Key = input.KeyF3
	KeyF4  Key = input.KeyF4
	KeyF5  Key = input.KeyF5
	KeyF6  Key = input.KeyF6
	KeyF7  Key = input.KeyF7
	KeyF8  Key = input.KeyF8
	KeyF9  Key = input.KeyF9
	KeyF10 Key = input.KeyF10
	KeyF11 Key = input.KeyF11
	KeyF12 Key = input.KeyF12
	KeyF13 Key = input.KeyF13
	KeyF14 Key = input.KeyF14
	KeyF15 Key = input.KeyF15
	KeyF16 Key = input.KeyF16
	KeyF17 Key = input.KeyF17
	KeyF18 Key = input.KeyF18
	KeyF19 Key = input.KeyF19
	KeyF20 Key = input.KeyF20
	KeyF21 Key = input.KeyF21
	KeyF22 Key = input.KeyF22
	KeyF23 Key = input.KeyF23
	KeyF24 Key = input.KeyF24
)

type KeyEvent = input.KeyEvent

type Keyboard = input.Keyboard

// -----
// mouse
// -----

type MouseButton = input.MouseButton

const (
	MouseButton0      MouseButton = input.MouseButton0
	MouseButton1      MouseButton = input.MouseButton0
	MouseButton2      MouseButton = input.MouseButton0
	MouseButton3      MouseButton = input.MouseButton0
	MouseButton4      MouseButton = input.MouseButton0
	MouseButtonLeft   MouseButton = input.MouseButtonLeft
	MouseButtonMiddle MouseButton = input.MouseButtonMiddle
	MouseButtonRight  MouseButton = input.MouseButtonRight
	MouseButtonMax    MouseButton = input.MouseButtonMax
)

type MouseButtonEvent = input.MouseButtonEvent

type MouseMoveEvent = input.MouseMoveEvent

type MouseWheelEvent = input.MouseWheelEvent

type Mouse = input.Mouse

// -------
// gamepad
// -------

type GamepadId = input.GamepadId

type GamepadInfo = input.GamepadInfo

type GamepadButton = input.GamepadButton

const (
	GamepadButtonSouth     GamepadButton = input.GamepadButtonSouth
	GamepadButtonEast      GamepadButton = input.GamepadButtonEast
	GamepadButtonWest      GamepadButton = input.GamepadButtonWest
	GamepadButtonNorth     GamepadButton = input.GamepadButtonNorth
	GamepadButtonL1        GamepadButton = input.GamepadButtonL1
	GamepadButtonR1        GamepadButton = input.GamepadButtonR1
	GamepadButtonL2        GamepadButton = input.GamepadButtonL2
	GamepadButtonR2        GamepadButton = input.GamepadButtonR2
	GamepadButtonSelect    GamepadButton = input.GamepadButtonSelect
	GamepadButtonStart     GamepadButton = input.GamepadButtonStart
	GamepadButtonLStick    GamepadButton = input.GamepadButtonLStick
	GamepadButtonRStick    GamepadButton = input.GamepadButtonRStick
	GamepadButtonDpadUp    GamepadButton = input.GamepadButtonDpadUp
	GamepadButtonDpadDown  GamepadButton = input.GamepadButtonDpadDown
	GamepadButtonDpadLeft  GamepadButton = input.GamepadButtonDpadLeft
	GamepadButtonDpadRight GamepadButton = input.GamepadButtonDpadRight
)

type GamepadAxis = input.GamepadAxis

type GamepadEventKind = input.GamepadEventKind

const (
	GamepadConnected    GamepadEventKind = input.GamepadConnected
	GamepadDisconnected GamepadEventKind = input.GamepadDisconnected
)

type GamepadConnectionEvent = input.GamepadConnectionEvent

type GamepadButtonEvent = input.GamepadButtonEvent

type GamepadAxisEvent = input.GamepadAxisEvent

type Gamepad = input.Gamepad

type Gamepads = input.Gamepads

// ----------
// text input
// ----------

type TextInputEvent = input.TextInputEvent

// ======
// render
// ======

// ----------
// renderstage
// ----------

// defines the order of render phases inside the draw system
type RenderStage = render.RenderStage

const (
	Opaque      RenderStage = render.Opaque
	Transparent RenderStage = render.Transparent
	Ui          RenderStage = render.Ui
	PostFx      RenderStage = render.PostFx
)

// ------
// anchor
// ------

// anchor component used by sprite, text, and shape backends
type Anchor = render.Anchor

type AnchorType = render.AnchorType

const (
	AnchorCenter      AnchorType = render.AnchorCenter
	AnchorTopLeft     AnchorType = render.AnchorTopLeft
	AnchorTop         AnchorType = render.AnchorTop
	AnchorTopRight    AnchorType = render.AnchorTopRight
	AnchorLeft        AnchorType = render.AnchorLeft
	AnchorRight       AnchorType = render.AnchorRight
	AnchorBottomLeft  AnchorType = render.AnchorBottomLeft
	AnchorBottom      AnchorType = render.AnchorBottom
	AnchorBottomRight AnchorType = render.AnchorBottomRight
)

// -------
// orderkey
// -------

// sortable key used to order render items inside a stage
type OrderKey = render.OrderKey

// packs layer, order, tie into a 64 bit order key
var NewOrderKey = render.NewOrderKey

// -----------------
// render registration
// -----------------

// registers the extraction function for type T
//
// extraction should be cheap: gather and copy the minimum data needed to draw.
func AddExtractionFn[C, T any](w *World, fn func(*World) []T) {
	render.AddExtractionFn[C](w, fn)
}

// registers the sort function for type T
//
// note: the renderer sorts ascending. if a stage needs reverse order, encode it in the returned key.
func AddSortFn[C, T any](w *World, fn func(T) uint64) {
	render.AddSortFn[C](w, fn)
}

// registers a render function for type T in a stage
//
// requirements:
//
// - an extraction fn for T must be registered first
//
// - a sort fn for T must be registered first
func AddRenderFn[C, T any](w *World, stage RenderStage, fn func(canvas *C, v T)) {
	render.AddRenderFn(w, stage, fn)
}

// =====
// shape
// =====

// 2d shape component rendered by the active backend
//
// notes:
//
// - shapes are positioned by transform.x and transform.y.
//
// - default anchor is center, and render.anchor overrides it.
//
// - orderkey controls render ordering inside the shape render phase.
type Shape = shape.Shape

// selects the geometry kind stored inside Shape
type ShapeKind = shape.ShapeKind

const (
	ShapeNone    ShapeKind = shape.None
	ShapeArrow   ShapeKind = shape.Arrow
	ShapeCapsule ShapeKind = shape.Capsule
	ShapeEllipse ShapeKind = shape.Ellipse
	ShapeLine    ShapeKind = shape.Line
	ShapePath    ShapeKind = shape.Path
	ShapePolygon ShapeKind = shape.Polygon
	ShapeRay     ShapeKind = shape.Ray
	ShapeRect    ShapeKind = shape.Rect
	ShapeSegment ShapeKind = shape.Segment
)

// fill types for shapes
type FillType = shape.FillType

const (
	FillSolid          FillType = shape.FillSolid
	FillLinearGradient FillType = shape.FillLinearGradient
	FillRadialGradient FillType = shape.FillRadialGradient
	FillImage          FillType = shape.FillImage
	FillVideo          FillType = shape.FillVideo
)

// shape fill definition
type Fill = shape.Fill

// shape stroke definition
type Stroke = shape.Stroke

// creates a Shape with defaults and applies options
var NewShape = shape.NewShape

// fluent options for NewShape and for building geometry
//
// geometry:
// - arrow, capsule, ellipse, line, path, polygon, ray, rect, segment
//
// styling:
// - fillsolid, stroke, order
var ShapeOpt = shape.ShapeOpt

// =========
// transform
// =========

type Transform = transform.Transform

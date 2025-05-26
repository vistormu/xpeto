package xp

import (
	"github.com/vistormu/xpeto/audio"
	"github.com/vistormu/xpeto/ecs"
	"github.com/vistormu/xpeto/engine"
	"github.com/vistormu/xpeto/event"
	"github.com/vistormu/xpeto/font"
	"github.com/vistormu/xpeto/image"
	"github.com/vistormu/xpeto/input"
	"github.com/vistormu/xpeto/scene"
	"github.com/vistormu/xpeto/state"

	"github.com/vistormu/xpeto/internal/geometry"

	"github.com/vistormu/xpeto/engine/component"
	audiosys "github.com/vistormu/xpeto/engine/system/audio"
	colsys "github.com/vistormu/xpeto/engine/system/collision"
	inputsys "github.com/vistormu/xpeto/engine/system/input"
	scenesys "github.com/vistormu/xpeto/engine/system/scene"
)

// =====
// audio
// =====
type Audio = *audio.Audio
type AudioHandle = audio.Handle
type AudioManager = *audio.Manager

func RegisterAudio(ctx Context, path string) audio.Handle {
	am, ok := ecs.GetResource[AudioManager](ctx)
	if !ok {
		return audio.Handle{}
	}

	return am.Register(path)

}

func LoadAudios(ctx Context, audios ...AudioHandle) {
	am, ok := ecs.GetResource[AudioManager](ctx)
	if !ok {
		return
	}

	for _, audio := range audios {
		am.Load(audio)
	}
}

func UnloadAudios(ctx Context, audios ...AudioHandle) {
	am, ok := ecs.GetResource[AudioManager](ctx)
	if !ok {
		return
	}

	for _, audio := range audios {
		am.Unload(audio)
	}
}

func GetAudio(ctx Context, handle AudioHandle) (Audio, bool) {
	am, ok := ecs.GetResource[AudioManager](ctx)
	if !ok {
		return nil, false
	}

	audio, ok := am.Audio(handle)
	return audio, ok
}

type AudioPlay = audiosys.AudioPlay
type AudioStop = audiosys.AudioStop
type AudioPause = audiosys.AudioPause
type AudioResume = audiosys.AudioResume

// ===
// ecs
// ===
type Entity = ecs.Entity
type System = ecs.System
type Context = *ecs.Context
type Filter = ecs.Filter
type SystemFilter = ecs.SystemFilter
type EntityManager = *ecs.EntityManager
type SystemManager = *ecs.SystemManager

// resources
func AddResource[T any](ctx Context, resource T) {
	ecs.AddResource(ctx, resource)
}

func GetResource[T any](ctx Context) (T, bool) {
	res, ok := ecs.GetResource[T](ctx)
	return res, ok
}

// queries
func Query(ctx Context, filter Filter) []Entity {
	em, ok := ecs.GetResource[ecs.EntityManager](ctx)
	if !ok {
		return nil
	}

	return em.Query(filter)
}

func Has[T any]() ecs.Filter {
	return ecs.Has[T]()
}

var Not = ecs.Not
var And = ecs.And
var Or = ecs.Or

// components
func GetComponent[T any](ctx Context, id Entity) (T, bool) {
	em, ok := ecs.GetResource[EntityManager](ctx)
	if !ok {
		var zero T
		return zero, false
	}

	return ecs.GetComponent[T](em, id)
}

func AddComponent[T any](ctx Context, entity Entity, component T) {
	em, ok := ecs.GetResource[EntityManager](ctx)
	if !ok {
		return
	}

	ecs.AddComponent(em, entity, component)
}

func RemoveComponent[T any](ctx Context, entity Entity) {
	em, ok := ecs.GetResource[EntityManager](ctx)
	if !ok {
		return
	}

	ecs.RemoveComponent[T](em, entity)
}

// entities
func CreateEntity(ctx Context) Entity {
	em, ok := ecs.GetResource[EntityManager](ctx)
	if !ok {
		return ecs.Entity{}
	}

	return em.Create()
}

func DestroyEntity(ctx Context, entity Entity) {
	em, ok := ecs.GetResource[EntityManager](ctx)
	if !ok {
		return
	}

	em.Destroy(entity)
}

func DestroyAllEntities(ctx Context) {
	em, ok := ecs.GetResource[EntityManager](ctx)
	if !ok {
		return
	}

	em.DestroyAll()
}

// ======
// engine
// ======
type Engine = engine.Engine

var NewEngine = engine.NewEngine

// =====
// event
// =====
type EventManager = *event.Manager

func Subscribe[T any](ctx Context, callback func(event T)) {
	em, ok := ecs.GetResource[EventManager](ctx)
	if !ok {
		return
	}

	event.Subscribe(em, callback)
}

func Unsubscribe[T any](ctx Context, callback func(event T)) {
	em, ok := ecs.GetResource[EventManager](ctx)
	if !ok {
		return
	}

	event.Unsubscribe(em, callback)
}

func Publish[T any](ctx Context, event T) {
	em, ok := ecs.GetResource[EventManager](ctx)
	if !ok {
		return
	}

	em.Publish(event)
}

// ====
// font
// ====
type Font = font.Font
type FontHandle = font.Handle

// =====
// image
// =====
type Image = *image.Image
type ImageHandle = image.Handle
type ImageManager = *image.Manager

func RegisterImage(ctx Context, path string) image.Handle {
	im, ok := ecs.GetResource[ImageManager](ctx)
	if !ok {
		return image.Handle{}
	}

	return im.Register(path)
}

func LoadImages(ctx Context, images ...ImageHandle) {
	im, ok := ecs.GetResource[ImageManager](ctx)
	if !ok {
		return
	}

	for _, image := range images {
		im.Load(image)
	}
}

func UnloadImages(ctx Context, images ...ImageHandle) {
	im, ok := ecs.GetResource[ImageManager](ctx)
	if !ok {
		return
	}

	for _, image := range images {
		im.Unload(image)
	}
}

func GetImage(ctx Context, handle ImageHandle) (Image, bool) {
	im, ok := ecs.GetResource[ImageManager](ctx)
	if !ok {
		return nil, false
	}

	return im.Image(handle), true
}

// =====
// input
// =====
type Key = input.Key

var (
	KeyA     = input.KeyA
	KeyB     = input.KeyB
	KeyC     = input.KeyC
	KeyD     = input.KeyD
	KeyE     = input.KeyE
	KeyF     = input.KeyF
	KeyG     = input.KeyG
	KeyH     = input.KeyH
	KeyI     = input.KeyI
	KeyJ     = input.KeyJ
	KeyK     = input.KeyK
	KeyL     = input.KeyL
	KeyM     = input.KeyM
	KeyN     = input.KeyN
	KeyO     = input.KeyO
	KeyP     = input.KeyP
	KeyQ     = input.KeyQ
	KeyR     = input.KeyR
	KeyS     = input.KeyS
	KeyT     = input.KeyT
	KeyU     = input.KeyU
	KeyV     = input.KeyV
	KeyW     = input.KeyW
	KeyX     = input.KeyX
	KeyY     = input.KeyY
	KeyZ     = input.KeyZ
	KeySpace = input.KeySpace
)

type KeyPress = inputsys.KeyPress
type KeyMaintain = inputsys.KeyMaintain
type KeyRelease = inputsys.KeyRelease

// =====
// scene
// =====
type Scene = scene.Scene
type SceneManager = *scene.Manager
type ScenePush = scenesys.ScenePush
type ScenePop = scenesys.ScenePop

// =====
// state
// =====
type State = state.State
type StateManager = *state.Manager

func AddState(ctx Context, state State) {
	sm, ok := ecs.GetResource[StateManager](ctx)
	if !ok {
		return
	}

	sm.Add(state)
}

func RemoveState(ctx Context, state State) {
	sm, ok := ecs.GetResource[StateManager](ctx)
	if !ok {
		return
	}

	sm.Remove(state)
}

// ==========
// structures
// ==========
type Vector = geometry.Vector[float32]
type Size = geometry.Size[float32]
type Rect = geometry.Rect[float32]
type Point = geometry.Point[float32]

func NewVector(x, y float32) Vector {
	return geometry.NewVector[float32](x, y)
}

func NewSize(width, height float32) Size {
	return geometry.NewSize[float32](width, height)
}

func NewRect(x, y, width, height float32) Rect {
	return geometry.NewRect[float32](x, y, width, height)
}

func NewPoint(x, y float32) Point {
	return geometry.NewPoint[float32](x, y)
}

// =====
// other
// =====
type Animation = *component.Animation
type Collidable = *component.Collidable
type Dynamic = *component.Dynamic
type Playable = *component.Playable
type Renderable = *component.Renderable
type Transform = *component.Transform

func NewTransform(posX, posY, scaleX, scaleY, rotation float32) Transform {
	return &component.Transform{
		Position: geometry.NewVector[float32](posX, posY),
		Scale:    geometry.NewVector[float32](scaleX, scaleY),
		Rotation: rotation,
	}
}

func NewRenderable(image ImageHandle, layer int) Renderable {
	return &component.Renderable{
		Sprite: image,
		Layer:  layer,
	}
}

// filters
func Always(ctx Context) bool {
	return true
}

func WhenActive(scene Scene) SystemFilter {
	return func(context *ecs.Context) bool {
		sm, ok := ecs.GetResource[SceneManager](context)
		if !ok {
			return false
		}

		return sm.IsActive(scene)
	}
}

// collisions
type CollisionEnter = colsys.CollisionEnter
type CollisionExit = colsys.CollisionExit
type CollisionStay = colsys.CollisionStay

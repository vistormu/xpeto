package xp

import (
	"fmt"
	"reflect"

	"github.com/vistormu/xpeto/internal/ecs"
	"github.com/vistormu/xpeto/internal/engine"
	"github.com/vistormu/xpeto/internal/event"
	"github.com/vistormu/xpeto/internal/geometry"
	"github.com/vistormu/xpeto/internal/state"
	"github.com/vistormu/xpeto/pkg/audio"
	"github.com/vistormu/xpeto/pkg/core"
	"github.com/vistormu/xpeto/pkg/render"
)

var NewGame = engine.NewGame

// ===
// ecs
// ===
type (
	Context   = ecs.Context
	Entity    = ecs.Entity
	Component = ecs.Component
	Archetype = ecs.Archetype
	Filter    = ecs.Filter
)

func AddResource[T any](ctx *Context, resource T) {
	ecs.AddResource(ctx, resource)
}

func GetResource[T any](ctx *Context) (T, bool) {
	return ecs.GetResource[T](ctx)
}

func MustResource[T any](ctx *Context) T {
	return ecs.MustResource[T](ctx)
}

func CreateEntity(ctx *Context, archetype Archetype) Entity {
	return MustResource[*ecs.Manager](ctx).Create(archetype)
}

func DestroyEntity(ctx *Context, entity Entity) {
	MustResource[*ecs.Manager](ctx).Destroy(entity)
}

func DestroyAllEntities(ctx *Context) {
	MustResource[*ecs.Manager](ctx).DestroyAll()
}

func Query(ctx *Context, filter Filter) []Entity {
	return MustResource[*ecs.Manager](ctx).Query(filter)
}

func GetComponent[T any](ctx *Context, entity Entity) (T, bool) {
	return ecs.GetComponent[T](MustResource[*ecs.Manager](ctx), entity)
}

func AddComponent[T any](ctx *Context, entity Entity, component T) {
	ecs.AddComponent(MustResource[*ecs.Manager](ctx), entity, component)
}

func RemoveComponent[T any](ctx *Context, entity Entity) {
	ecs.RemoveComponent[T](MustResource[*ecs.Manager](ctx), entity)
}

func Has[T any]() ecs.Filter {
	return ecs.Has[T]()
}

var (
	Not = ecs.Not
	And = ecs.And
	Or  = ecs.Or
)

func RegisterAssets[T any](ctx *Context, assets *T) {
	structValue := reflect.ValueOf(assets).Elem()
	structType := structValue.Type()

	imageType := reflect.TypeOf((*render.Image)(nil)).Elem()
	audioType := reflect.TypeOf((*audio.Audio)(nil)).Elem()

	for i := 0; i < structValue.NumField(); i++ {
		fieldValue := structValue.Field(i)
		fieldStruct := structType.Field(i)

		fieldType := fieldValue.Type()
		fieldTypeName := fieldType.String()

		switch fieldType {
		case imageType:
			renderManager, ok := GetResource[*render.Manager](ctx)
			if !ok {
				panic("To use images, you must enable the render plugin.")
			}

			assetPath := fieldStruct.Tag.Get("path")
			if assetPath == "" {
				panic(fmt.Sprintf("missing path tag on field %s", fieldStruct.Name))
			}

			handle := renderManager.Register(assetPath)
			if fieldValue.CanSet() {
				fieldValue.Set(reflect.ValueOf(handle))
			} else {
				panic(fmt.Sprintf("cannot set field %s", fieldStruct.Name))
			}

		case audioType:
			audioManager, ok := GetResource[*audio.Manager](ctx)
			if !ok {
				panic("To use audio, you must enable the audio plugin.")
			}

			assetPath := fieldStruct.Tag.Get("path")
			if assetPath == "" {
				panic(fmt.Sprintf("missing path tag on field %s", fieldStruct.Name))
			}

			handle := audioManager.Register(assetPath)
			if fieldValue.CanSet() {
				fieldValue.Set(reflect.ValueOf(handle))
			} else {
				panic(fmt.Sprintf("cannot set field %s", fieldStruct.Name))
			}

		default:
			panic(fmt.Sprintf("unsupported asset type %s on field %s", fieldTypeName, fieldStruct.Name))
		}
	}

	// Insert the fully populated struct as a resource into the context
	AddResource(ctx, *assets)
}

// ======
// events
// ======
type Event = event.Event

func Subscribe[T any](ctx *Context, callback func(data T)) Event {
	return event.Subscribe(MustResource[*event.Manager](ctx), callback)
}

func Unsubscribe(ctx *Context, ev Event) {
	MustResource[*event.Manager](ctx).Unsubscribe(ev)
}

func Publish[T any](ctx *Context, data T) {
	MustResource[*event.Manager](ctx).Publish(data)
}

// =====
// state
// =====
type (
	State           = state.State
	StateActivate   = state.StateActivate
	StateDeactivate = state.StateDeactivate
)

const (
	OnEnter     = state.OnEnter
	OnExit      = state.OnExit
	Update      = state.Update
	FixedUpdate = state.FixedUpdate
)

var InitialState = state.InitialState

// ===============
// data structures
// ===============
type (
	Vector = geometry.Vector[float32]
	Size   = geometry.Size[float32]
	Rect   = geometry.Rect[float32]
	Point  = geometry.Point[float32]
)

// other
type (
	Transform = core.Transform
)

// =====
// audio
// =====
type (
	Audio       = audio.Audio
	AudioPause  = audio.AudioPause
	AudioPlay   = audio.AudioPlay
	AudioResume = audio.AudioResume
	AudioStop   = audio.AudioStop
	AudioPlugin = audio.Plugin
)

func LoadAudios(ctx *Context, audios ...Audio) {
	am, ok := GetResource[*audio.Manager](ctx)
	if !ok {
		fmt.Printf("to use audio, you must use the audio plugin")
		return
	}

	am.Load(audios...)
}

func UnloadAudios(ctx *Context, audios ...Audio) {
	am, ok := GetResource[*audio.Manager](ctx)
	if !ok {
		fmt.Printf("to use audio, you must use the audio plugin")
		return
	}

	am.Unload(audios...)
}

// ======
// render
// ======
type (
	Image      = render.Image
	Renderable = render.Renderable
)

func LoadImages(ctx *Context, images ...Image) {
	rm, ok := GetResource[*render.Manager](ctx)
	if !ok {
		fmt.Printf("to use render, you must use the render plugin")
		return
	}

	rm.Load(images...)
}

func UnloadImages(ctx *Context, images ...Image) {
	rm, ok := GetResource[*render.Manager](ctx)
	if !ok {
		fmt.Printf("to use render, you must use the render plugin")
		return
	}

	rm.Unload(images...)
}

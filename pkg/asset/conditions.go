package asset

import (
	"reflect"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func WhenAssetLoaded(a Asset) schedule.ConditionFn {
	return func(w *ecs.World) bool {
		return IsAssetLoaded(w, a)
	}
}

func WhenAssetFailed(a Asset) schedule.ConditionFn {
	return func(w *ecs.World) bool {
		st, ok := GetAssetState(w, a)
		return ok && st == AssetFailed
	}
}

func WhenAssetState(a Asset, st AssetState) schedule.ConditionFn {
	return func(w *ecs.World) bool {
		got, ok := GetAssetState(w, a)
		return ok && got == st
	}
}

func WhenAllAssetsLoaded(assets ...Asset) schedule.ConditionFn {
	cp := append([]Asset(nil), assets...)
	return func(w *ecs.World) bool {
		if len(cp) == 0 {
			return false
		}
		for _, a := range cp {
			if !IsAssetLoaded(w, a) {
				return false
			}
		}
		return true
	}
}

func WhenAnyAssetFailed(assets ...Asset) schedule.ConditionFn {
	cp := append([]Asset(nil), assets...)
	return func(w *ecs.World) bool {
		if len(cp) == 0 {
			return false
		}
		for _, a := range cp {
			st, ok := GetAssetState(w, a)
			if ok && st == AssetFailed {
				return true
			}
		}
		return false
	}
}

func WhenBundleLoaded[T any]() schedule.ConditionFn {
	t := reflect.TypeFor[T]()
	t = baseType(t)

	return func(w *ecs.World) bool {
		b, ok := ecs.GetResource[T](w)
		if !ok {
			return false
		}

		if t == nil || t.Kind() != reflect.Struct {
			return false
		}

		v := reflect.ValueOf(b)
		if v.Kind() == reflect.Pointer {
			v = v.Elem()
		}
		// if T is stored by value, v is a struct
		if v.Kind() != reflect.Struct {
			return false
		}

		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)

			// only care about Asset fields
			if f.Type() != reflect.TypeFor[Asset]() {
				continue
			}

			a, ok := f.Interface().(Asset)
			if !ok || a == Asset(0) {
				return false
			}
			if !IsAssetLoaded(w, a) {
				return false
			}
		}

		return true
	}
}

func WhenBundleFailed[T any]() schedule.ConditionFn {
	t := reflect.TypeFor[T]()
	t = baseType(t)

	return func(w *ecs.World) bool {
		b, ok := ecs.GetResource[T](w)
		if !ok {
			return false
		}

		if t == nil || t.Kind() != reflect.Struct {
			return false
		}

		v := reflect.ValueOf(b)
		if v.Kind() == reflect.Pointer {
			v = v.Elem()
		}
		if v.Kind() != reflect.Struct {
			return false
		}

		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			if f.Type() != reflect.TypeFor[Asset]() {
				continue
			}

			a, ok := f.Interface().(Asset)
			if !ok || a == Asset(0) {
				continue
			}

			st, ok := GetAssetState(w, a)
			if ok && st == AssetFailed {
				return true
			}
		}

		return false
	}
}

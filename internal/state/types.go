package state

import (
	"github.com/vistormu/xpeto/internal/ecs"
)

type Hook = int

const (
	OnEnter Hook = iota
	OnExit
	OnTransition
	Update
	FixedUpdate
)

type OnEnterFn = func(*ecs.Context)
type OnExitFn = func(*ecs.Context)
type OnTransitionFn = func(*ecs.Context)
type UpdateFn = func(*ecs.Context, float32)
type FixedUpdateFn = func(*ecs.Context, float32)

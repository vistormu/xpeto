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
	Draw
)

type ContextFn = func(*ecs.Context)
type UpdateFn = func(*ecs.Context, float32)

type StateFn interface {
	ContextFn | UpdateFn
}

type StateManager interface {
	Register(hook Hook, state any, fn any)
}

type StateSystem interface {
	OnEnter(ctx *ecs.Context)
	Update(ctx *ecs.Context)
}

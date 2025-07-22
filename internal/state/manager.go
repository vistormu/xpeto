package state

import (
	"fmt"
)

type Manager[S comparable] struct {
	active          S
	previous        S
	onEnterFns      map[S][]ContextFn
	onExitFns       map[S][]ContextFn
	onTransitionFns map[[2]S][]ContextFn
	updateFns       map[S][]UpdateFn
	fixedUpdateFns  map[S][]UpdateFn
}

func NewManager[S comparable](initial S) *Manager[S] {
	return &Manager[S]{
		active:          initial,
		previous:        initial,
		onEnterFns:      make(map[S][]ContextFn),
		onExitFns:       make(map[S][]ContextFn),
		onTransitionFns: make(map[[2]S][]ContextFn),
		updateFns:       make(map[S][]UpdateFn),
		fixedUpdateFns:  make(map[S][]UpdateFn),
	}
}

func (m *Manager[S]) Register(hook Hook, state S, fn any) {
	switch hook {
	case OnEnter:
		function, ok := fn.(ContextFn)
		if !ok {
			fmt.Println("OnEnter")
			return
		}
		m.onEnterFns[state] = append(m.onEnterFns[state], function)

	case OnExit:
		function, ok := fn.(ContextFn)
		if !ok {
			fmt.Println("OnExit")
			return
		}
		m.onExitFns[state] = append(m.onExitFns[state], function)

	case Update:
		function, ok := fn.(UpdateFn)
		if !ok {
			fmt.Println("Update")
			return
		}
		m.updateFns[state] = append(m.updateFns[state], function)

	case FixedUpdate:
		function, ok := fn.(UpdateFn)
		if !ok {
			fmt.Println("FixedUpdate")
			return
		}
		m.fixedUpdateFns[state] = append(m.fixedUpdateFns[state], function)
	}
}

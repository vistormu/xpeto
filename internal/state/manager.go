package state

import (
	"fmt"
)

type Manager[S comparable] struct {
	active          S
	previous        S
	onEnterFns      map[S][]OnEnterFn
	onExitFns       map[S][]OnExitFn
	onTransitionFns map[[2]S][]OnTransitionFn
	updateFns       map[S][]UpdateFn
	fixedUpdateFns  map[S][]FixedUpdateFn
}

func NewManager[S comparable](initial S) *Manager[S] {
	return &Manager[S]{
		active:          initial,
		previous:        initial,
		onEnterFns:      make(map[S][]OnEnterFn),
		onExitFns:       make(map[S][]OnExitFn),
		onTransitionFns: make(map[[2]S][]OnTransitionFn),
		updateFns:       make(map[S][]UpdateFn),
		fixedUpdateFns:  make(map[S][]FixedUpdateFn),
	}
}

func (m *Manager[S]) Register(hook Hook, state S, fn any) {
	switch hook {
	case OnEnter:
		function, ok := fn.(OnEnterFn)
		if !ok {
			fmt.Println("OnEnter")
			return
		}
		m.onEnterFns[state] = append(m.onEnterFns[state], function)

	case OnExit:
		function, ok := fn.(OnExitFn)
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
		function, ok := fn.(FixedUpdateFn)
		if !ok {
			fmt.Println("FixedUpdate")
			return
		}
		m.fixedUpdateFns[state] = append(m.fixedUpdateFns[state], function)
	}
}

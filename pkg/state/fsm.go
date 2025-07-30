package state

import (
	"fmt"
)

type Fsm[S comparable] struct {
	active          S
	previous        S
	onEnterFns      map[S][]ContextFn
	onExitFns       map[S][]ContextFn
	onTransitionFns map[[2]S][]ContextFn
	updateFns       map[S][]UpdateFn
	fixedUpdateFns  map[S][]UpdateFn
}

func NewFsm[S comparable](initial S) *Fsm[S] {
	return &Fsm[S]{
		active:          initial,
		previous:        initial,
		onEnterFns:      make(map[S][]ContextFn),
		onExitFns:       make(map[S][]ContextFn),
		onTransitionFns: make(map[[2]S][]ContextFn),
		updateFns:       make(map[S][]UpdateFn),
		fixedUpdateFns:  make(map[S][]UpdateFn),
	}
}

func (fsm *Fsm[S]) Register(hook Hook, state S, fn any) {
	switch hook {
	case OnEnter:
		function, ok := fn.(ContextFn)
		if !ok {
			fmt.Println("OnEnter")
			return
		}
		fsm.onEnterFns[state] = append(fsm.onEnterFns[state], function)

	case OnExit:
		function, ok := fn.(ContextFn)
		if !ok {
			fmt.Println("OnExit")
			return
		}
		fsm.onExitFns[state] = append(fsm.onExitFns[state], function)

	case Update:
		function, ok := fn.(UpdateFn)
		if !ok {
			fmt.Println("Update")
			return
		}
		fsm.updateFns[state] = append(fsm.updateFns[state], function)

	case FixedUpdate:
		function, ok := fn.(UpdateFn)
		if !ok {
			fmt.Println("FixedUpdate")
			return
		}
		fsm.fixedUpdateFns[state] = append(fsm.fixedUpdateFns[state], function)
	}
}

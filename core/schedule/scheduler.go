package schedule

import (
	"fmt"
	"slices"

	"github.com/vistormu/xpeto/core/ecs"
)

// ========
// schedule
// ========
type Scheduler struct {
	logger     *logger
	labelIndex *labelIndex

	builder  *builder
	store    *storage
	compiler *compiler
	runner   *runner
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		logger:     newLogger(),
		labelIndex: newLabelIndex(),
		builder:    newBuilder(),
		store:      newStorage(),
		compiler:   newCompiler(),
		runner:     newRunner(),
	}
}

func (sch *Scheduler) flushPending() {
	if sch == nil || sch.builder == nil || sch.store == nil {
		return
	}
	if sch.builder.node == nil {
		return
	}

	n := sch.builder.build()
	if n == nil {
		return
	}
	sch.store.add(n)
}

func (sch *Scheduler) compile(w *ecs.World) {
	// systems
	sch.flushPending()
	sch.compiler.compileDirty(sch.store, sch.labelIndex, sch.logger)

	// resources
	if !ecs.HasResource[RunningSystem](w) {
		ecs.AddResource(w, RunningSystem{})
	}

	if !ecs.HasResource[transitionEvent](w) {
		ecs.AddResource(w, newTransitionEvent())
	}
}

// =======
// options
// =======
type option = func(*Scheduler)

type systemOpt struct{}

var SystemOpt systemOpt

func AddSystem(sch *Scheduler, stage Stage, system ecs.System, opts ...option) {
	sch.flushPending()

	id := sch.builder.nextId
	st := stage(sch.store, id)

	sch.builder.add(sch.logger, st, system)

	for _, opt := range opts {
		if opt != nil {
			opt(sch)
		}
	}
}

func (systemOpt) Label(label string) option {
	return func(sch *Scheduler) {
		sch.builder.label(sch.logger, sch.labelIndex, label)
	}
}

func (systemOpt) RunIf(conditions ...ConditionFn) option {
	return func(sch *Scheduler) {
		sch.builder.runIf(sch.logger, conditions...)
	}
}

func (systemOpt) RunBefore(labels ...string) option {
	return func(sch *Scheduler) {
		sch.builder.before(sch.logger, labels...)
	}
}

func (systemOpt) RunAfter(labels ...string) option {
	return func(sch *Scheduler) {
		sch.builder.after(sch.logger, labels...)
	}
}

// ===========
// backend API
// ===========
func RunStartup(w *ecs.World, sch *Scheduler) {
	sch.compile(w)
	sch.runner.runStages(w, sch.store, preStartup, startup, postStartup)
}

func RunUpdate(w *ecs.World, sch *Scheduler) {
	sch.compile(w)

	tr, _ := ecs.GetResource[transitionEvent](w)

	// first pass
	sch.runner.runStages(w, sch.store, first, preUpdate, stateTransition)

	// state transitions
	sch.runner.runIds(w, sch.store, tr.onExit)
	sch.runner.runIds(w, sch.store, tr.onTransition)
	sch.runner.runIds(w, sch.store, tr.onEnter)
	tr.clear()

	// fixed update
	steps := max(sch.runner.fixedStepsFn(w), 0)
	for range steps {
		sch.runner.runStages(w, sch.store, fixedFirst, fixedPreUpdate, fixedUpdate, fixedPostUpdate, fixedLast)
	}

	// last pass
	sch.runner.runStages(w, sch.store, update, postUpdate, last)
}

func RunDraw(w *ecs.World, sch *Scheduler) {
	sch.compile(w)
	sch.runner.runStages(w, sch.store, preDraw, draw, postDraw)
}

func RunExit(w *ecs.World, sch *Scheduler) {
	sch.compile(w)
	sch.runner.runStages(w, sch.store, exit)
}

func SetFixedStepsFn(sch *Scheduler, fn func(*ecs.World) int) {
	if fn != nil {
		sch.runner.fixedStepsFn = fn
	}
}

// ========
// debuging
// ========
func Diagnostics(sch *Scheduler) []Diagnostic {
	return sch.logger.diagnostics.ToSlice()
}

func Plan(sch *Scheduler) string {
	stages := make([]stage, 0, len(sch.store.stages))
	for st := range sch.store.stages {
		stages = append(stages, st)
	}
	slices.Sort(stages)

	var out string
	for _, st := range stages {
		ids, ok := sch.store.plan[st]
		if !ok || len(ids) == 0 {
			ids = sch.store.stages[st]
		}

		if len(ids) == 0 {
			continue
		}

		out += fmt.Sprintf("== %s ==\n", st.String())

		for i, id := range ids {
			n, ok := sch.store.get(id)
			if !ok || n == nil {
				continue
			}

			label := n.label
			if label == "" {
				label = "<unlabeled>"
			}

			out += fmt.Sprintf("  %02d: id=%d label=%s\n", i, n.id, label)
		}
	}

	return out
}

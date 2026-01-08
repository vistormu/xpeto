package schedule

import (
	"github.com/vistormu/xpeto/core/ecs"
)

// ===========
// system info
// ===========
type RunningSystem struct {
	Id    uint64
	Label string
}

// ======
// runner
// ======
type runner struct {
	fixedStepsFn func(*ecs.World) int
}

func newRunner() *runner {
	return &runner{
		fixedStepsFn: func(*ecs.World) int { return 0 },
	}
}

func (r *runner) runStages(w *ecs.World, store *storage, stages ...stage) {
	for _, stage := range stages {
		ids, ok := store.plan[stage]
		if !ok || len(ids) == 0 {
			ids = store.stages[stage]
		}
		r.runIds(w, store, ids)
	}
}

func (r *runner) runIds(w *ecs.World, store *storage, ids []uint64) {
	rs, _ := ecs.GetResource[RunningSystem](w)

	for _, id := range ids {
		n, ok := store.get(id)
		if !ok || n == nil || n.system == nil {
			continue
		}

		if !conditionsPass(w, n.conditions) {
			continue
		}

		rs.Id = n.id
		rs.Label = n.label
		n.system(w)
	}
}

func conditionsPass(w *ecs.World, conds []ConditionFn) bool {
	for _, c := range conds {
		if c == nil {
			continue
		}
		if !c(w) {
			return false
		}
	}
	return true
}

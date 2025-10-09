package schedule

import "github.com/vistormu/xpeto/core/ecs"

type schedule struct {
	id         uint64
	stage      stage
	system     ecs.System
	before     []uint64
	after      []uint64
	conditions []ConditionFn
}

func newSchedule() *schedule {
	return &schedule{
		before:     make([]uint64, 0),
		after:      make([]uint64, 0),
		conditions: make([]ConditionFn, 0),
	}
}

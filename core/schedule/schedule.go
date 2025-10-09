package schedule

import "github.com/vistormu/xpeto/core/ecs"

type Schedule struct {
	Id         uint64
	stage      Stage
	System     ecs.System
	before     []uint64
	after      []uint64
	Conditions []ConditionFn
}

func newSchedule() *Schedule {
	return &Schedule{
		before:     make([]uint64, 0),
		after:      make([]uint64, 0),
		Conditions: make([]ConditionFn, 0),
	}
}

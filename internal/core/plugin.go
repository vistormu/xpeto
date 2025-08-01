package core

// ========
// schedule
// ========
type Schedule struct {
	Name      string
	Stage     Stage
	System    func(*Context)
	Before    []string
	After     []string
	Condition func(*Context) bool
}

func (sb *Schedule) WithSystem(name string, stage Stage, system func(*Context)) *Schedule {
	sb.Name = name
	sb.Stage = stage
	sb.System = system

	return sb
}

func (sb *Schedule) RunIf(condition func(*Context) bool) *Schedule {
	sb.Condition = condition
	return sb
}

func (sb *Schedule) RunBefore(systems ...string) *Schedule {
	sb.Before = append(sb.Before, systems...)
	return sb
}

func (sb *Schedule) RunAfter(systems ...string) *Schedule {
	sb.After = append(sb.After, systems...)
	return sb
}

// =========
// scheduler
// =========
type ScheduleBuilder struct {
	Schedules []*Schedule
}

func (sb *ScheduleBuilder) NewSchedule() *Schedule {
	s := &Schedule{
		Name:   "",
		Stage:  PreUpdate,
		System: nil,
	}

	sb.Schedules = append(sb.Schedules, s)
	return s
}

// ======
// plugin
// ======
type Plugin func(*Context, *ScheduleBuilder)

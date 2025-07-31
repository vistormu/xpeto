package core

type ScheduleBuilder struct {
	Name      string
	Stage     Stage
	System    func(*Context)
	Before    []string
	After     []string
	Condition func(*Context) bool
}

func (sb *ScheduleBuilder) WithSystem(name string, stage Stage, system func(*Context)) *ScheduleBuilder {
	sb.Name = name
	sb.Stage = stage
	sb.System = system

	return sb
}

func (sb *ScheduleBuilder) RunIf(condition func(*Context) bool) *ScheduleBuilder {
	sb.Condition = condition
	return sb
}

func (sb *ScheduleBuilder) RunBefore(systems ...string) *ScheduleBuilder {
	sb.Before = append(sb.Before, systems...)
	return sb
}

func (sb *ScheduleBuilder) RunAfter(systems ...string) *ScheduleBuilder {
	sb.After = append(sb.After, systems...)
	return sb
}

type Plugin interface {
	Build(*Context, *ScheduleBuilder)
}

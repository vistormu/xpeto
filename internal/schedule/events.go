package schedule

type EventStateTransition[T comparable] struct {
	Exited  *T
	Entered *T
}

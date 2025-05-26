package ecs

type Filter interface {
	Match(em *EntityManager, e Entity) bool
}

type has[T any] struct{}

func (h has[T]) Match(em *EntityManager, e Entity) bool {
	_, ok := GetComponent[T](em, e)
	return ok
}
func Has[T any]() Filter { return has[T]{} }

type not struct{ f Filter }

func (n not) Match(em *EntityManager, e Entity) bool { return !n.f.Match(em, e) }
func Not(f Filter) Filter                            { return not{f} }

type and struct{ fs []Filter }

func (a and) Match(em *EntityManager, e Entity) bool {
	for _, f := range a.fs {
		if !f.Match(em, e) {
			return false
		}
	}
	return true
}
func And(fs ...Filter) Filter { return and{fs} }

type or struct{ fs []Filter }

func (o or) Match(em *EntityManager, e Entity) bool {
	for _, f := range o.fs {
		if f.Match(em, e) {
			return true
		}
	}
	return false
}
func Or(fs ...Filter) Filter { return or{fs} }

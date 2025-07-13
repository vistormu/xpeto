package ecs

type Filter interface {
	Match(m *Manager, e Entity) bool
}

// has
type has[T any] struct{}

func (h has[T]) Match(m *Manager, e Entity) bool {
	_, ok := GetComponent[T](m, e)
	return ok
}
func Has[T any]() Filter { return has[T]{} }

// not
type not struct{ f Filter }

func (n not) Match(m *Manager, e Entity) bool { return !n.f.Match(m, e) }
func Not(f Filter) Filter                     { return not{f} }

// and
type and struct{ fs []Filter }

func (a and) Match(m *Manager, e Entity) bool {
	for _, f := range a.fs {
		if !f.Match(m, e) {
			return false
		}
	}
	return true
}
func And(fs ...Filter) Filter { return and{fs} }

// or
type or struct{ fs []Filter }

func (o or) Match(m *Manager, e Entity) bool {
	for _, f := range o.fs {
		if f.Match(m, e) {
			return true
		}
	}
	return false
}
func Or(fs ...Filter) Filter { return or{fs} }

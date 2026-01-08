package schedule

import (
	"github.com/vistormu/xpeto/core/ecs"
)

// ====
// node
// ====
type node struct {
	id         uint64
	stage      stage
	label      string
	system     ecs.System
	conditions []ConditionFn
	before     []string
	after      []string
	insertion  uint32
}

func newNode() *node {
	return &node{
		conditions: make([]ConditionFn, 0),
		before:     make([]string, 0),
		after:      make([]string, 0),
	}
}

// =======
// builder
// =======
type builder struct {
	nextId uint64
	node   *node
}

func newBuilder() *builder {
	return &builder{
		nextId: 1,
	}
}

func (b *builder) add(l *logger, stage stage, system ecs.System) {
	if system == nil {
		l.add("attempted to add a nil system", 0, "", stage)
		return
	}

	n := newNode()

	n.id = b.nextId
	b.nextId++

	n.stage = stage
	n.system = system

	b.node = n
}

func (b *builder) runIf(l *logger, conditions ...ConditionFn) {
	if b.node == nil {
		l.add("executed RunIf before adding a system", 0, "", empty)
		return
	}

	b.node.conditions = append(b.node.conditions, conditions...)
}

func (b *builder) label(l *logger, li *labelIndex, label string) {
	if b.node == nil {
		l.add("executed Label before adding a system", 0, "", empty)
		return
	}

	if label == "" {
		return
	}

	if !li.add(label, b.node.id, b.node.stage) {
		l.add("duplicate label detected", b.node.id, label, b.node.stage)
		return
	}

	b.node.label = label
}

func (b *builder) before(l *logger, labels ...string) {
	if b.node == nil {
		l.add("executed Before before adding a system", 0, "", empty)
		return
	}

	b.node.before = append(b.node.before, labels...)
}

func (b *builder) after(l *logger, labels ...string) {
	if b.node == nil {
		l.add("executed After before adding a system", 0, "", empty)
		return
	}

	b.node.after = append(b.node.after, labels...)
}

func (b *builder) build() *node {
	n := b.node
	b.node = nil

	return n
}

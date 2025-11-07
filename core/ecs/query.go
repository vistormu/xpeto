package ecs

// =======
// filters
// =======
type Filter = func(w *World, e Entity) bool

func With[T any]() Filter {
	return func(w *World, e Entity) bool {
		_, ok := GetComponent[T](w, e)
		return ok
	}
}

func Without[T any]() Filter {
	return func(w *World, e Entity) bool {
		_, ok := GetComponent[T](w, e)
		return !ok
	}
}

func Or(filters ...Filter) Filter {
	return func(w *World, e Entity) bool {
		for _, f := range filters {
			if f(w, e) {
				return true
			}
		}
		return false
	}
}

// =======
// bundles
// =======
type bundle1[A any] struct {
	e Entity
	a *A
}

func (b bundle1[A]) Entity() Entity { return b.e }
func (b bundle1[A]) A() *A          { return b.a }

type bundle2[A, B any] struct {
	e  Entity
	a  *A
	bb *B
}

func (b bundle2[A, B]) Entity() Entity { return b.e }
func (b bundle2[A, B]) A() *A          { return b.a }
func (b bundle2[A, B]) B() *B          { return b.bb }

type bundle3[A, B, C any] struct {
	e  Entity
	a  *A
	bb *B
	c  *C
}

func (b bundle3[A, B, C]) Entity() Entity { return b.e }
func (b bundle3[A, B, C]) A() *A          { return b.a }
func (b bundle3[A, B, C]) B() *B          { return b.bb }
func (b bundle3[A, B, C]) C() *C          { return b.c }

type bundle4[A, B, C, D any] struct {
	e  Entity
	a  *A
	bb *B
	c  *C
	d  *D
}

func (b bundle4[A, B, C, D]) Entity() Entity { return b.e }
func (b bundle4[A, B, C, D]) A() *A          { return b.a }
func (b bundle4[A, B, C, D]) B() *B          { return b.bb }
func (b bundle4[A, B, C, D]) C() *C          { return b.c }
func (b bundle4[A, B, C, D]) D() *D          { return b.d }

// =======
// helpers
// =======
type denseProvider func() []Entity

func pickSmallestDense(providers ...denseProvider) denseProvider {
	if len(providers) == 0 {
		return func() []Entity { return nil }
	}
	best := providers[0]
	bestLen := len(best())
	for i := 1; i < len(providers); i++ {
		cur := providers[i]
		if l := len(cur()); l < bestLen {
			best, bestLen = cur, l
		}
	}
	return best
}

func passFilters(w *World, e Entity, filters []Filter) bool {
	for _, flt := range filters {
		if !flt(w, e) {
			return false
		}
	}
	return true
}

// =====
// query
// =====

// 1
type Query1[A any] struct {
	storeA  *store[A]
	filters []Filter
	w       *World
}

func (q *Query1[A]) Iter() func(func(i int, b bundle1[A]) bool) {
	return func(yield func(i int, b bundle1[A]) bool) {
		w := q.w
		i := 0
		for rowA, e := range q.storeA.dense {
			if !w.population.has(e) {
				continue
			}
			if !passFilters(w, e, q.filters) {
				continue
			}
			if !yield(i, bundle1[A]{
				e: e,
				a: &q.storeA.values[rowA],
			}) {
				return
			}
			i++
		}
	}
}

func (q *Query1[A]) Get() []bundle1[A] {
	return collect(q.Iter())
}

func NewQuery1[A any](w *World, filters ...Filter) *Query1[A] {
	return &Query1[A]{
		storeA:  getStore[A](w.registry),
		filters: filters,
		w:       w,
	}
}

// 2
type Query2[A, B any] struct {
	storeA  *store[A]
	storeB  *store[B]
	filters []Filter
	w       *World
	dense   denseProvider
}

func (q *Query2[A, B]) Iter() func(func(i int, b bundle2[A, B]) bool) {
	return func(yield func(i int, b bundle2[A, B]) bool) {
		w := q.w
		i := 0
		for _, e := range q.dense() {
			if !w.population.has(e) {
				continue
			}
			rowA, ok := q.storeA.location[e]
			if !ok {
				continue
			}
			rowB, ok := q.storeB.location[e]
			if !ok {
				continue
			}
			if !passFilters(w, e, q.filters) {
				continue
			}
			if !yield(i, bundle2[A, B]{
				e:  e,
				a:  &q.storeA.values[rowA],
				bb: &q.storeB.values[rowB],
			}) {
				return
			}
			i++
		}
	}
}

func (q *Query2[A, B]) Get() []bundle2[A, B] {
	return collect(q.Iter())
}

func NewQuery2[A, B any](w *World, filters ...Filter) *Query2[A, B] {
	a := getStore[A](w.registry)
	b := getStore[B](w.registry)

	return &Query2[A, B]{
		storeA:  a,
		storeB:  b,
		filters: filters,
		w:       w,
		dense: pickSmallestDense(
			func() []Entity { return a.dense },
			func() []Entity { return b.dense },
		),
	}
}

// 3
type Query3[A, B, C any] struct {
	storeA  *store[A]
	storeB  *store[B]
	storeC  *store[C]
	filters []Filter
	w       *World
	dense   denseProvider
}

func (q *Query3[A, B, C]) Iter() func(func(i int, b bundle3[A, B, C]) bool) {
	return func(yield func(i int, b bundle3[A, B, C]) bool) {
		w := q.w
		i := 0
		for _, e := range q.dense() {
			if !w.population.has(e) {
				continue
			}
			rowA, ok := q.storeA.location[e]
			if !ok {
				continue
			}
			rowB, ok := q.storeB.location[e]
			if !ok {
				continue
			}
			rowC, ok := q.storeC.location[e]
			if !ok {
				continue
			}
			if !passFilters(w, e, q.filters) {
				continue
			}
			if !yield(i, bundle3[A, B, C]{
				e:  e,
				a:  &q.storeA.values[rowA],
				bb: &q.storeB.values[rowB],
				c:  &q.storeC.values[rowC],
			}) {
				return
			}
			i++
		}
	}
}

func (q *Query3[A, B, C]) Get() []bundle3[A, B, C] {
	return collect(q.Iter())
}

func NewQuery3[A, B, C any](w *World, filters ...Filter) *Query3[A, B, C] {
	a := getStore[A](w.registry)
	b := getStore[B](w.registry)
	c := getStore[C](w.registry)

	return &Query3[A, B, C]{
		storeA:  a,
		storeB:  b,
		storeC:  c,
		filters: filters,
		w:       w,
		dense: pickSmallestDense(
			func() []Entity { return a.dense },
			func() []Entity { return b.dense },
			func() []Entity { return c.dense },
		),
	}
}

// 4
type Query4[A, B, C, D any] struct {
	storeA  *store[A]
	storeB  *store[B]
	storeC  *store[C]
	storeD  *store[D]
	filters []Filter
	w       *World
	dense   denseProvider
}

func (q *Query4[A, B, C, D]) Iter() func(func(i int, b bundle4[A, B, C, D]) bool) {
	return func(yield func(i int, b bundle4[A, B, C, D]) bool) {
		w := q.w
		i := 0
		for _, e := range q.dense() {
			if !w.population.has(e) {
				continue
			}
			rowA, ok := q.storeA.location[e]
			if !ok {
				continue
			}
			rowB, ok := q.storeB.location[e]
			if !ok {
				continue
			}
			rowC, ok := q.storeC.location[e]
			if !ok {
				continue
			}
			rowD, ok := q.storeD.location[e]
			if !ok {
				continue
			}
			if !passFilters(w, e, q.filters) {
				continue
			}
			if !yield(i, bundle4[A, B, C, D]{
				e:  e,
				a:  &q.storeA.values[rowA],
				bb: &q.storeB.values[rowB],
				c:  &q.storeC.values[rowC],
				d:  &q.storeD.values[rowD],
			}) {
				return
			}
			i++
		}
	}
}

func (q *Query4[A, B, C, D]) Get() []bundle4[A, B, C, D] {
	return collect(q.Iter())
}

func NewQuery4[A, B, C, D any](w *World, filters ...Filter) *Query4[A, B, C, D] {
	a := getStore[A](w.registry)
	b := getStore[B](w.registry)
	c := getStore[C](w.registry)
	d := getStore[D](w.registry)

	return &Query4[A, B, C, D]{
		storeA:  a,
		storeB:  b,
		storeC:  c,
		storeD:  d,
		filters: filters,
		w:       w,
		dense: pickSmallestDense(
			func() []Entity { return a.dense },
			func() []Entity { return b.dense },
			func() []Entity { return c.dense },
			func() []Entity { return d.dense },
		),
	}
}

// =======
// helpers
// =======
func collect[T any](iter func(func(i int, t T) bool)) []T {
	out := make([]T, 0)
	iter(func(_ int, t T) bool {
		out = append(out, t)
		return true
	})
	return out
}

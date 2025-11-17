package ecs

import (
	"testing"
)

func TestEntities(t *testing.T) {
	w := NewWorld()

	// add first entity
	e := AddEntity(w)

	if e.index() != 0 {
		t.Fatal("wrong index")
	}

	if e.gen() != 0 {
		t.Fatal("wrong gen")
	}

	if w.population.alive != 1 {
		t.Fatal("wrong number of alive entities")
	}

	if w.population.free.Length() != 0 {
		t.Fatal("wrong free length")
	}

	if len(w.population.gens) != 1 {
		t.Fatal("wrong gens len")
	}

	// remove first entity
	RemoveEntity(w, e)

	if w.population.alive != 0 {
		t.Fatal("wrong number of alive entities")
	}

	if w.population.free.Length() != 1 {
		t.Fatal("wrong free length")
	}

	if len(w.population.gens) != 1 {
		t.Fatal("wrong gens len")
	}

	// add second entity
	e = AddEntity(w)

	if e.index() != 0 {
		t.Fatal("wrong index")
	}

	if e.gen() != 1 {
		t.Fatal("wrong gen")
	}

	if w.population.alive != 1 {
		t.Fatal("wrong number of alive entities")
	}

	if w.population.free.Length() != 0 {
		t.Fatal("wrong free length")
	}

	if len(w.population.gens) != 1 {
		t.Fatal("wrong gens len")
	}

	// add third entity
	e = AddEntity(w)

	if e.index() != 1 {
		t.Fatal("wrong index")
	}

	if e.gen() != 0 {
		t.Fatal("wrong gen")
	}

	if w.population.alive != 2 {
		t.Fatal("wrong number of alive entities")
	}

	if w.population.free.Length() != 0 {
		t.Fatal("wrong free length")
	}

	if len(w.population.gens) != 2 {
		t.Fatal("wrong gens len")
	}
}

type mockComponent struct {
	value int
}

func TestComponent(t *testing.T) {
	w := NewWorld()

	e := AddEntity(w)
	AddComponent(w, e, mockComponent{value: 1})

	if w.registry.Len() != 1 {
		t.Fatal("wrong length of stores")
	}

	c, ok := GetComponent[mockComponent](w, e)
	if !ok {
		t.Fatal("component not found")
	}

	if c.value != 1 {
		t.Fatal("wrong component value")
	}

	c.value = 10

	c2, _ := GetComponent[mockComponent](w, e)

	if c2.value != 10 {
		t.Fatal("wrong component value")
	}

	ok = RemoveComponent[mockComponent](w, e)
	if !ok {
		t.Fatal("error removing component")
	}

	if w.registry.Len() != 1 {
		t.Fatal("wrong number of stores")
	}
}

type mockResource struct {
	value int
}

func TestResources(t *testing.T) {
	w := NewWorld()

	r := mockResource{value: 10}

	AddResource(w, r)

	result, ok := GetResource[mockResource](w)
	if !ok {
		t.Fatal("did not find resource")
	}
	if result.value != 10 {
		t.Fatal("wrong value")
	}

	r.value = 20

	_, ok = GetResource[*mockResource](w)
	if ok {
		t.Fatal("shouldn't work")
	}

	ok = RemoveResource[mockResource](w)
	if !ok {
		t.Fatal("couldn't remove resource")
	}

	_, ok = GetResource[mockResource](w)
	if ok {
		t.Fatal("resource was not removed")
	}

	AddResource(w, &mockResource{value: 5})
	result, ok = GetResource[mockResource](w)
	if !ok {
		t.Fatal("mock resource should be inside the world")
	}
	if result.value != 5 {
		t.Fatal("wrong value")
	}
}

type mockComponent2 struct{}

func TestQuery(t *testing.T) {
	w := NewWorld()

	e1 := AddEntity(w)
	ok := AddComponent(w, e1, mockComponent{value: 1})
	if !ok {
		t.Fatal("could not add component")
	}

	e2 := AddEntity(w)
	ok1 := AddComponent(w, e2, mockComponent{value: 1})
	ok2 := AddComponent(w, e2, mockComponent2{})
	if !ok1 || !ok2 {
		t.Fatal("could not load ocmponents")
	}

	// queries
	q := NewQuery1[mockComponent](w)

	count := 0
	for i, b := range q.Iter() {
		mc := b.Components()
		if mc.value != 1 {
			t.Fatal("wrong value")
		}

		mc.value = 10

		if i != count {
			t.Fatal("iterator fAiled")
		}

		count++
	}
	if count != 2 {
		t.Fatalf("wrong count: %d", count)
	}

	if c, ok := GetComponent[mockComponent](w, e1); ok {
		if c.value != 10 {
			t.Fatal("wrong value")
		}
	}

	q2 := NewQuery2[mockComponent, mockComponent2](w)

	count = 0
	for _, b := range q2.Iter() {
		mc1, _ := b.Components()

		if mc1.value != 10 {
			t.Fatal("wrong value")
		}
		count++
	}
	if count != 1 {
		t.Fatal("wrong count")
	}

	q3 := NewQuery1[mockComponent](w, Without[mockComponent2]())

	count = 0
	for _, b := range q3.Iter() {
		mc := b.Components()

		if mc.value != 10 {
			t.Fatal("wrong value")
		}

		count++
	}

	if count != 1 {
		t.Fatal("wrong count")
	}
}

func TestWrongComponentAdding(t *testing.T) {
	w := NewWorld()

	e1 := AddEntity(w)
	ok := AddComponent(w, e1, &mockComponent{value: 1})
	if ok {
		t.Fatal("could add component")
	}
}

func TestQueryFirstQuery1(t *testing.T) {
	w := NewWorld()

	// no components yet → First should return ok=false
	q := NewQuery1[mockComponent](w)
	if _, ok := q.First(); ok {
		t.Fatal("expected First to return ok=false when there are no matches")
	}

	// add two entities with the same component
	e1 := AddEntity(w)
	if !AddComponent(w, e1, mockComponent{value: 1}) {
		t.Fatal("could not add component to e1")
	}
	e2 := AddEntity(w)
	if !AddComponent(w, e2, mockComponent{value: 2}) {
		t.Fatal("could not add component to e2")
	}

	b, ok := q.First()
	if !ok {
		t.Fatal("expected First to find at least one entity")
	}

	if b.Entity() != e1 {
		t.Fatalf("expected First to return e1, got %v", b.Entity())
	}

	mc := b.Components()
	if mc.value != 1 {
		t.Fatalf("expected First to return component value 1, got %d", mc.value)
	}
}

func TestQuerySingleQuery1(t *testing.T) {
	w := NewWorld()

	q := NewQuery1[mockComponent](w)

	// zero matches
	if _, ok := q.Single(); ok {
		t.Fatal("expected Single to return ok=false with zero matches")
	}

	// exactly one match
	e1 := AddEntity(w)
	if !AddComponent(w, e1, mockComponent{value: 10}) {
		t.Fatal("could not add component to e1")
	}

	b, ok := q.Single()
	if !ok {
		t.Fatal("expected Single to return ok=true with one match")
	}
	if b.Entity() != e1 {
		t.Fatalf("expected Single to return e1, got %v", b.Entity())
	}
	if mc := b.Components(); mc.value != 10 {
		t.Fatalf("expected component value 10, got %d", mc.value)
	}

	// two matches → Single must fail
	e2 := AddEntity(w)
	if !AddComponent(w, e2, mockComponent{value: 20}) {
		t.Fatal("could not add component to e2")
	}

	if _, ok := q.Single(); ok {
		t.Fatal("expected Single to return ok=false with more than one match")
	}
}

func TestQueryFirstAndSingleQuery2(t *testing.T) {
	w := NewWorld()

	// prepare world
	e1 := AddEntity(w)
	if !AddComponent(w, e1, mockComponent{value: 1}) {
		t.Fatal("could not add mockComponent to e1")
	}
	if !AddComponent(w, e1, mockComponent2{}) {
		t.Fatal("could not add mockComponent2 to e1")
	}

	q := NewQuery2[mockComponent, mockComponent2](w)

	// Single with exactly one match
	b, ok := q.Single()
	if !ok {
		t.Fatal("expected Single on Query2 to succeed with one match")
	}

	a, _ := b.Components()
	if a.value != 1 {
		t.Fatalf("expected component value 1, got %d", a.value)
	}
	if b.Entity() != e1 {
		t.Fatalf("expected Single to return e1, got %v", b.Entity())
	}

	// First should return the same entity
	b2, ok := q.First()
	if !ok {
		t.Fatal("expected First on Query2 to succeed with one match")
	}
	if b2.Entity() != e1 {
		t.Fatalf("expected First to return e1, got %v", b2.Entity())
	}

	// add a second matching entity
	e2 := AddEntity(w)
	if !AddComponent(w, e2, mockComponent{value: 2}) {
		t.Fatal("could not add mockComponent to e2")
	}
	if !AddComponent(w, e2, mockComponent2{}) {
		t.Fatal("could not add mockComponent2 to e2")
	}

	// Single must now fail because there are two matches
	if _, ok := q.Single(); ok {
		t.Fatal("expected Single on Query2 to fail with more than one match")
	}

	// First must still succeed and return some match (the first one)
	b3, ok := q.First()
	if !ok {
		t.Fatal("expected First on Query2 to still succeed with multiple matches")
	}
	if b3.Entity() != e1 {
		t.Fatalf("expected First on Query2 to still return e1, got %v", b3.Entity())
	}
}

func TestBundlesMatchesIterQuery1(t *testing.T) {
	w := NewWorld()

	e1 := AddEntity(w)
	e2 := AddEntity(w)
	AddComponent(w, e1, mockComponent{value: 1})
	AddComponent(w, e2, mockComponent{value: 2})

	q := NewQuery1[mockComponent](w)

	// collect via Bundles
	bs := q.Bundles()
	if len(bs) != 2 {
		t.Fatalf("expected Bundles length 2, got %d", len(bs))
	}

	// collect via Iter
	var fromIter []Entity
	for _, b := range q.Iter() {
		fromIter = append(fromIter, b.Entity())
	}

	if len(fromIter) != len(bs) {
		t.Fatalf("mismatch between Iter and Bundles lengths: %d vs %d", len(fromIter), len(bs))
	}

	for i := range bs {
		if bs[i].Entity() != fromIter[i] {
			t.Fatalf("entity mismatch at %d: Bundles=%v Iter=%v", i, bs[i].Entity(), fromIter[i])
		}
	}
}

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

	if len(w.registry.stores) != 1 {
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

	if len(w.registry.stores) != 1 {
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
		if b.A().value != 1 {
			t.Fatal("wrong value")
		}

		b.A().value = 10

		t.Log(i)

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
		if b.A().value != 10 {
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
		if b.A().value != 10 {
			t.Fatal("wrong value")
		}

		count++
	}

	if count != 1 {
		t.Fatal("wrong count")
	}
}

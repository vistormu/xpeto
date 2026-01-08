package ecs

import "testing"

type mockComponent struct{ value int }
type mockComponent2 struct{ value int }
type mockComponent3 struct{ value int }
type mockComponent4 struct{ value int }

type mockResource struct{ value int }

func TestEntities_GenerationStartsAtOne(t *testing.T) {
	w := NewWorld()

	e := AddEntity(w)

	if e.index() != 0 {
		t.Fatal("wrong index")
	}
	if e.gen() != 1 {
		t.Fatalf("wrong generation: got %d want %d", e.gen(), 1)
	}
	if !HasEntity(w, e) {
		t.Fatal("expected entity to be alive")
	}
}

func TestEntities_RemoveAndReuseIncrementsGeneration(t *testing.T) {
	w := NewWorld()

	e1 := AddEntity(w)
	if e1.gen() != 1 {
		t.Fatalf("wrong generation: got %d want %d", e1.gen(), 1)
	}

	if !RemoveEntity(w, e1) {
		t.Fatal("expected RemoveEntity to succeed")
	}
	if HasEntity(w, e1) {
		t.Fatal("expected entity to be dead after removal")
	}

	e2 := AddEntity(w)

	// if index is reused, generation must increase
	if e2.index() == e1.index() {
		if e2.gen() != e1.gen()+1 {
			t.Fatalf("expected generation increment on reuse: got %d want %d", e2.gen(), e1.gen()+1)
		}
	}
}

func TestComponents_AddGetMutateRemove(t *testing.T) {
	w := NewWorld()
	e := AddEntity(w)

	if !AddComponent(w, e, mockComponent{value: 5}) {
		t.Fatal("error adding component")
	}

	c, ok := GetComponent[mockComponent](w, e)
	if !ok {
		t.Fatal("error getting component")
	}
	if c.value != 5 {
		t.Fatal("wrong component value")
	}

	// mutate through pointer
	c.value = 10

	c2, ok := GetComponent[mockComponent](w, e)
	if !ok {
		t.Fatal("error getting component after mutation")
	}
	if c2.value != 10 {
		t.Fatal("wrong component value after mutation")
	}

	if !RemoveComponent[mockComponent](w, e) {
		t.Fatal("error removing component")
	}

	if _, ok := GetComponent[mockComponent](w, e); ok {
		t.Fatal("expected missing component after removal")
	}
}

func TestComponents_AddTwiceOverwritesValue(t *testing.T) {
	w := NewWorld()
	e := AddEntity(w)

	if !AddComponent(w, e, mockComponent{value: 1}) {
		t.Fatal("error adding component")
	}
	if !AddComponent(w, e, mockComponent{value: 7}) {
		t.Fatal("error adding component second time")
	}

	c, ok := GetComponent[mockComponent](w, e)
	if !ok || c.value != 7 {
		t.Fatal("expected overwrite on second add")
	}
}

func TestComponents_PointerTypesAreRejected(t *testing.T) {
	w := NewWorld()
	e := AddEntity(w)

	v := &mockComponent{value: 3}
	if AddComponent(w, e, v) {
		t.Fatal("expected pointer component type to be rejected")
	}

	if _, ok := GetComponent[*mockComponent](w, e); ok {
		t.Fatal("expected pointer component type to be missing")
	}
}

func TestEntityOps_FailOnDeadEntity(t *testing.T) {
	w := NewWorld()
	e := AddEntity(w)

	if !RemoveEntity(w, e) {
		t.Fatal("expected RemoveEntity to succeed")
	}

	if AddComponent(w, e, mockComponent{value: 1}) {
		t.Fatal("expected AddComponent to fail on dead entity")
	}
	if _, ok := GetComponent[mockComponent](w, e); ok {
		t.Fatal("expected GetComponent to fail on dead entity")
	}
	if RemoveComponent[mockComponent](w, e) {
		t.Fatal("expected RemoveComponent to fail on dead entity")
	}
}

func TestRemoveEntity_RemovesAllComponents(t *testing.T) {
	w := NewWorld()
	e := AddEntity(w)

	AddComponent(w, e, mockComponent{value: 1})
	AddComponent(w, e, mockComponent2{value: 2})
	AddComponent(w, e, mockComponent3{value: 3})

	if !RemoveEntity(w, e) {
		t.Fatal("expected RemoveEntity to succeed")
	}

	if _, ok := GetComponent[mockComponent](w, e); ok {
		t.Fatal("expected component to be gone after entity removal")
	}
	if _, ok := GetComponent[mockComponent2](w, e); ok {
		t.Fatal("expected component2 to be gone after entity removal")
	}
	if _, ok := GetComponent[mockComponent3](w, e); ok {
		t.Fatal("expected component3 to be gone after entity removal")
	}
}

// func TestReserveComponents_DoesNotBreakStorage(t *testing.T) {
// 	w := NewWorld()

// 	ReserveComponents[mockComponent](w, 64)

// 	es := make([]Entity, 0, 32)
// 	for i := 0; i < 32; i++ {
// 		e := AddEntity(w)
// 		es = append(es, e)
// 		if !AddComponent(w, e, mockComponent{value: i}) {
// 			t.Fatal("error adding component after reserve")
// 		}
// 	}

// 	for i, e := range es {
// 		c, ok := GetComponent[mockComponent](w, e)
// 		if !ok || c.value != i {
// 			t.Fatal("wrong component value after reserve")
// 		}
// 	}
// }

func TestResources_AddGetOverwriteRemove(t *testing.T) {
	w := NewWorld()

	if !AddResource(w, mockResource{value: 1}) {
		t.Fatal("error adding resource")
	}

	r, ok := GetResource[mockResource](w)
	if !ok {
		t.Fatal("error getting resource")
	}
	if r.value != 1 {
		t.Fatal("wrong resource value")
	}

	// overwrite
	if !AddResource(w, mockResource{value: 9}) {
		t.Fatal("error overwriting resource")
	}

	r2, ok := GetResource[mockResource](w)
	if !ok || r2.value != 9 {
		t.Fatal("expected overwritten resource value")
	}

	if !RemoveResource[mockResource](w) {
		t.Fatal("error removing resource")
	}

	if _, ok := GetResource[mockResource](w); ok {
		t.Fatal("expected missing resource after removal")
	}
}

func TestResources_PointerTypesAreRejected(t *testing.T) {
	w := NewWorld()

	r := &mockResource{value: 3}
	if AddResource(w, r) {
		t.Fatal("expected pointer resource type to be rejected")
	}

	if _, ok := GetResource[*mockResource](w); ok {
		t.Fatal("expected GetResource[*T] to fail")
	}

	if RemoveResource[*mockResource](w) {
		t.Fatal("expected RemoveResource[*T] to fail")
	}
}

func TestQuery1_BundlesIterFirstSingle(t *testing.T) {
	w := NewWorld()

	e1 := AddEntity(w)
	e2 := AddEntity(w)
	e3 := AddEntity(w)

	AddComponent(w, e1, mockComponent{value: 1})
	AddComponent(w, e2, mockComponent{value: 2})
	AddComponent(w, e3, mockComponent{value: 3})

	q := NewQuery1[mockComponent](w)

	bs := q.Bundles()
	if len(bs) != 3 {
		t.Fatalf("expected 3 bundles, got %d", len(bs))
	}

	seen := 0
	for _, b := range q.Iter() {
		a := b.Components()
		if a == nil {
			t.Fatal("expected non nil component pointer")
		}
		_ = b.Entity()
		seen++
	}
	if seen != 3 {
		t.Fatalf("expected 3 iter results, got %d", seen)
	}

	first, ok := q.First()
	if !ok {
		t.Fatal("expected First to succeed")
	}
	if first.Components() == nil {
		t.Fatal("expected First component pointer")
	}

	if _, ok := q.Single(); ok {
		t.Fatal("expected Single to fail on multiple results")
	}
}

func TestQuery1_Filters_WithWithout(t *testing.T) {
	w := NewWorld()

	e1 := AddEntity(w)
	e2 := AddEntity(w)
	e3 := AddEntity(w)

	AddComponent(w, e1, mockComponent{value: 1})
	AddComponent(w, e2, mockComponent{value: 2})
	AddComponent(w, e3, mockComponent{value: 3})
	AddComponent(w, e3, mockComponent2{value: 20})

	q := NewQuery1[mockComponent](w, Without[mockComponent2]())

	for _, b := range q.Iter() {
		if b.Entity() == e3 {
			t.Fatal("expected Without filter to exclude entity")
		}
	}
}

func TestQuery2_Query3_Query4_BasicIteration(t *testing.T) {
	w := NewWorld()

	e1 := AddEntity(w)
	e2 := AddEntity(w)

	AddComponent(w, e1, mockComponent{value: 1})
	AddComponent(w, e1, mockComponent2{value: 2})
	AddComponent(w, e1, mockComponent3{value: 3})
	AddComponent(w, e1, mockComponent4{value: 4})

	AddComponent(w, e2, mockComponent{value: 10})
	AddComponent(w, e2, mockComponent2{value: 20})
	AddComponent(w, e2, mockComponent3{value: 30})
	AddComponent(w, e2, mockComponent4{value: 40})

	q2 := NewQuery2[mockComponent, mockComponent2](w)
	n2 := 0
	for _, b := range q2.Iter() {
		a, bb := b.Components()
		if a == nil || bb == nil {
			t.Fatal("expected non nil component pointers in Query2")
		}
		n2++
	}
	if n2 != 2 {
		t.Fatalf("expected 2 results in Query2, got %d", n2)
	}

	q3 := NewQuery3[mockComponent, mockComponent2, mockComponent3](w)
	n3 := 0
	for _, b := range q3.Iter() {
		a, bb, c := b.Components()
		if a == nil || bb == nil || c == nil {
			t.Fatal("expected non nil component pointers in Query3")
		}
		n3++
	}
	if n3 != 2 {
		t.Fatalf("expected 2 results in Query3, got %d", n3)
	}

	q4 := NewQuery4[mockComponent, mockComponent2, mockComponent3, mockComponent4](w)
	n4 := 0
	for _, b := range q4.Iter() {
		a, bb, c, d := b.Components()
		if a == nil || bb == nil || c == nil || d == nil {
			t.Fatal("expected non nil component pointers in Query4")
		}
		n4++
	}
	if n4 != 2 {
		t.Fatalf("expected 2 results in Query4, got %d", n4)
	}
}

func TestQueryOr_SemanticsAreBoundToDrivingStore(t *testing.T) {
	w := NewWorld()

	eA := AddEntity(w)
	AddComponent(w, eA, mockComponent{value: 1})

	eB := AddEntity(w)
	AddComponent(w, eB, mockComponent2{value: 2})

	// query is driven by mockComponent store, so it cannot return eB
	q := NewQuery1[mockComponent](w, Or(With[mockComponent](), With[mockComponent2]()))

	count := 0
	for _, b := range q.Iter() {
		if b.Entity() == eB {
			t.Fatal("expected Query1[mockComponent] to never return entities without mockComponent")
		}
		count++
	}
	if count != 1 {
		t.Fatalf("expected 1 result, got %d", count)
	}
}

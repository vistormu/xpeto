package input

import (
	"github.com/vistormu/go-dsa/set"
)

type ButtonInput[T comparable] struct {
	pressed        *set.HashSet[T]
	justPressed    *set.HashSet[T]
	justReleased   *set.HashSet[T]
	pressDurations map[T]int
}

func newButtonInput[T comparable]() ButtonInput[T] {
	return ButtonInput[T]{
		pressed:        set.NewHashSet[T](),
		justPressed:    set.NewHashSet[T](),
		justReleased:   set.NewHashSet[T](),
		pressDurations: make(map[T]int),
	}
}

func (bi *ButtonInput[T]) begin() {
	bi.justPressed.Clear()
	bi.justReleased.Clear()
}

func (bi *ButtonInput[T]) compute() {
	for b := range bi.pressed.Iter() {
		if bi.pressDurations[b] < 0 {
			bi.pressDurations[b] = 0
		}
		bi.pressDurations[b]++
	}

	for b := range bi.pressDurations {
		if !bi.pressed.Contains(b) {
			delete(bi.pressDurations, b)
		}
	}
}

func (bi *ButtonInput[T]) end() {}

func (bi *ButtonInput[T]) press(button T) {
	if bi.pressed.Contains(button) {
		return
	}
	bi.pressed.Add(button)
	bi.justPressed.Add(button)
}

func (bi *ButtonInput[T]) release(button T) {
	if !bi.pressed.Contains(button) {
		return
	}
	bi.pressed.Remove(button)
	bi.justReleased.Add(button)
	delete(bi.pressDurations, button)
}

func (bi *ButtonInput[T]) reset() {
	bi.pressed.Clear()
	bi.justPressed.Clear()
	bi.justReleased.Clear()
	clear(bi.pressDurations)
}

// ========
// user API
// ========
func (bi *ButtonInput[T]) IsPressed(button T) bool {
	return bi.pressed.Contains(button)
}

func (bi *ButtonInput[T]) IsJustPressed(button T) bool {
	return bi.justPressed.Contains(button)
}

func (bi *ButtonInput[T]) IsJustReleased(button T) bool {
	return bi.justReleased.Contains(button)
}

func (bi *ButtonInput[T]) Duration(button T) int {
	d, ok := bi.pressDurations[button]
	if !ok || d < 0 {
		return 0
	}
	return d
}

func (bi *ButtonInput[T]) Pressed() []T {
	return bi.pressed.ToSlice()
}

func (bi *ButtonInput[T]) JustPressed() []T {
	return bi.justPressed.ToSlice()
}

func (bi *ButtonInput[T]) JustReleased() []T {
	return bi.justReleased.ToSlice()
}

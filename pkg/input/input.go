package input

import (
	"github.com/vistormu/go-dsa/set"
)

// ======
// button
// ======
type ButtonInput[T comparable] struct {
	pressed        *set.HashSet[T]
	pressDurations map[T]int
	justPressed    *set.HashSet[T]
	justReleased   *set.HashSet[T]
}

func newButtonInput[T comparable]() *ButtonInput[T] {
	return &ButtonInput[T]{
		pressed:        set.NewHashSet[T](),
		pressDurations: make(map[T]int, 0),
		justPressed:    set.NewHashSet[T](),
		justReleased:   set.NewHashSet[T](),
	}
}

func (bi *ButtonInput[T]) Clear() {
	bi.justPressed.Clear()
	bi.justReleased.Clear()
}

func (bi *ButtonInput[T]) SetDuration(button T, duration int) {
	bi.pressDurations[button] = duration
}

func (bi *ButtonInput[T]) Press(button T) {
	if !bi.pressed.Contains(button) {
		bi.pressed.Add(button)
		bi.justPressed.Add(button)
	}
}

func (bi *ButtonInput[T]) Release(button T) {
	if bi.pressed.Contains(button) {
		bi.pressed.Remove(button)
		bi.justReleased.Add(button)
		delete(bi.pressDurations, button)
	}
}

// ===
// API
// ===
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
	duration, ok := bi.pressDurations[button]
	if !ok {
		return 0
	}

	return duration
}

func (bi *ButtonInput[T]) Pressed() []T {
	return bi.pressed.Values()
}

// ======
// cursor
// ======
type CursorInput struct {
	X, Y         int
	Dx, Dy       int
	PrevX, PrevY int
}

func (ci *CursorInput) Position() (int, int) {
	return ci.X, ci.Y
}

func (ci *CursorInput) Delta() (int, int) {
	return ci.Dx, ci.Dy
}

// =====
// wheel
// =====
type WheelInput struct {
	X, Y float64
}

// ====
// axis
// ====
type AxisInput struct {
	Value    float64
	Previous float64
	Delta    float64
}

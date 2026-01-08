package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/event"
	"github.com/vistormu/xpeto/pkg/input"
)

// =====
// table
// =====
var keyMap = func() []uint16 {
	m := make([]uint16, int(ebiten.KeyMax)+1)

	set := func(ek ebiten.Key, k input.Key) { m[int(ek)] = uint16(k) + 1 }

	// Letters
	set(ebiten.KeyA, input.KeyA)
	set(ebiten.KeyB, input.KeyB)
	set(ebiten.KeyC, input.KeyC)
	set(ebiten.KeyD, input.KeyD)
	set(ebiten.KeyE, input.KeyE)
	set(ebiten.KeyF, input.KeyF)
	set(ebiten.KeyG, input.KeyG)
	set(ebiten.KeyH, input.KeyH)
	set(ebiten.KeyI, input.KeyI)
	set(ebiten.KeyJ, input.KeyJ)
	set(ebiten.KeyK, input.KeyK)
	set(ebiten.KeyL, input.KeyL)
	set(ebiten.KeyM, input.KeyM)
	set(ebiten.KeyN, input.KeyN)
	set(ebiten.KeyO, input.KeyO)
	set(ebiten.KeyP, input.KeyP)
	set(ebiten.KeyQ, input.KeyQ)
	set(ebiten.KeyR, input.KeyR)
	set(ebiten.KeyS, input.KeyS)
	set(ebiten.KeyT, input.KeyT)
	set(ebiten.KeyU, input.KeyU)
	set(ebiten.KeyV, input.KeyV)
	set(ebiten.KeyW, input.KeyW)
	set(ebiten.KeyX, input.KeyX)
	set(ebiten.KeyY, input.KeyY)
	set(ebiten.KeyZ, input.KeyZ)

	// Digits
	set(ebiten.Key0, input.Key0)
	set(ebiten.Key1, input.Key1)
	set(ebiten.Key2, input.Key2)
	set(ebiten.Key3, input.Key3)
	set(ebiten.Key4, input.Key4)
	set(ebiten.Key5, input.Key5)
	set(ebiten.Key6, input.Key6)
	set(ebiten.Key7, input.Key7)
	set(ebiten.Key8, input.Key8)
	set(ebiten.Key9, input.Key9)

	// Numpad
	set(ebiten.KeyNumpad0, input.KeyNumpad0)
	set(ebiten.KeyNumpad1, input.KeyNumpad1)
	set(ebiten.KeyNumpad2, input.KeyNumpad2)
	set(ebiten.KeyNumpad3, input.KeyNumpad3)
	set(ebiten.KeyNumpad4, input.KeyNumpad4)
	set(ebiten.KeyNumpad5, input.KeyNumpad5)
	set(ebiten.KeyNumpad6, input.KeyNumpad6)
	set(ebiten.KeyNumpad7, input.KeyNumpad7)
	set(ebiten.KeyNumpad8, input.KeyNumpad8)
	set(ebiten.KeyNumpad9, input.KeyNumpad9)
	set(ebiten.KeyNumpadAdd, input.KeyNumpadAdd)
	set(ebiten.KeyNumpadDecimal, input.KeyNumpadDecimal)
	set(ebiten.KeyNumpadDivide, input.KeyNumpadDivide)
	set(ebiten.KeyNumpadEnter, input.KeyNumpadEnter)
	set(ebiten.KeyNumpadEqual, input.KeyNumpadEqual)
	set(ebiten.KeyNumpadMultiply, input.KeyNumpadMultiply)
	set(ebiten.KeyNumpadSubtract, input.KeyNumpadSubtract)

	// Punctuation and symbols
	set(ebiten.KeyBracketLeft, input.KeyBracketLeft)
	set(ebiten.KeyBracketRight, input.KeyBracketRight)
	set(ebiten.KeyComma, input.KeyComma)
	set(ebiten.KeyBackslash, input.KeyBackslash)
	set(ebiten.KeyEqual, input.KeyEqual)
	set(ebiten.KeyBackquote, input.KeyBackquote)
	set(ebiten.KeyIntlBackslash, input.KeyIntlBackslash)
	set(ebiten.KeyMinus, input.KeyMinus)
	set(ebiten.KeyPeriod, input.KeyPeriod)
	set(ebiten.KeyQuote, input.KeyQuote)
	set(ebiten.KeySemicolon, input.KeySemicolon)
	set(ebiten.KeySlash, input.KeySlash)

	// Common controls
	set(ebiten.KeySpace, input.KeySpace)
	set(ebiten.KeyEnter, input.KeyEnter)
	set(ebiten.KeyTab, input.KeyTab)
	set(ebiten.KeyEscape, input.KeyEscape)
	set(ebiten.KeyBackspace, input.KeyBackspace)
	set(ebiten.KeyDelete, input.KeyDelete)
	set(ebiten.KeyInsert, input.KeyInsert)

	// Arrows
	set(ebiten.KeyUp, input.KeyArrowUp)
	set(ebiten.KeyDown, input.KeyArrowDown)
	set(ebiten.KeyLeft, input.KeyArrowLeft)
	set(ebiten.KeyRight, input.KeyArrowRight)

	// Modifiers (left/right)
	set(ebiten.KeyAltLeft, input.KeyAltLeft)
	set(ebiten.KeyAltRight, input.KeyAltRight)
	set(ebiten.KeyControlLeft, input.KeyControlLeft)
	set(ebiten.KeyControlRight, input.KeyControlRight)
	set(ebiten.KeyMetaLeft, input.KeyMetaLeft)
	set(ebiten.KeyMetaRight, input.KeyMetaRight)
	set(ebiten.KeyShiftLeft, input.KeyShiftLeft)
	set(ebiten.KeyShiftRight, input.KeyShiftRight)

	// Modifiers (generic)
	set(ebiten.KeyAlt, input.KeyAlt)
	set(ebiten.KeyControl, input.KeyControl)
	set(ebiten.KeyShift, input.KeyShift)
	set(ebiten.KeyMeta, input.KeyMeta)

	// Navigation and system keys
	set(ebiten.KeyCapsLock, input.KeyCapsLock)
	set(ebiten.KeyContextMenu, input.KeyContextMenu)
	set(ebiten.KeyEnd, input.KeyEnd)
	set(ebiten.KeyHome, input.KeyHome)
	set(ebiten.KeyNumLock, input.KeyNumLock)
	set(ebiten.KeyPageDown, input.KeyPageDown)
	set(ebiten.KeyPageUp, input.KeyPageUp)
	set(ebiten.KeyPause, input.KeyPause)
	set(ebiten.KeyPrintScreen, input.KeyPrintScreen)
	set(ebiten.KeyScrollLock, input.KeyScrollLock)

	// Function keys
	set(ebiten.KeyF1, input.KeyF1)
	set(ebiten.KeyF2, input.KeyF2)
	set(ebiten.KeyF3, input.KeyF3)
	set(ebiten.KeyF4, input.KeyF4)
	set(ebiten.KeyF5, input.KeyF5)
	set(ebiten.KeyF6, input.KeyF6)
	set(ebiten.KeyF7, input.KeyF7)
	set(ebiten.KeyF8, input.KeyF8)
	set(ebiten.KeyF9, input.KeyF9)
	set(ebiten.KeyF10, input.KeyF10)
	set(ebiten.KeyF11, input.KeyF11)
	set(ebiten.KeyF12, input.KeyF12)
	set(ebiten.KeyF13, input.KeyF13)
	set(ebiten.KeyF14, input.KeyF14)
	set(ebiten.KeyF15, input.KeyF15)
	set(ebiten.KeyF16, input.KeyF16)
	set(ebiten.KeyF17, input.KeyF17)
	set(ebiten.KeyF18, input.KeyF18)
	set(ebiten.KeyF19, input.KeyF19)
	set(ebiten.KeyF20, input.KeyF20)
	set(ebiten.KeyF21, input.KeyF21)
	set(ebiten.KeyF22, input.KeyF22)
	set(ebiten.KeyF23, input.KeyF23)
	set(ebiten.KeyF24, input.KeyF24)

	return m
}()

func mapKey(ek ebiten.Key) (input.Key, bool) {
	if ek < 0 || int(ek) >= len(keyMap) {
		return 0, false
	}
	v := keyMap[int(ek)]
	if v == 0 {
		return 0, false
	}
	return input.Key(v - 1), true
}

// ======
// buffer
// ======
type keyboardState struct {
	justPressedKeys  []ebiten.Key
	justReleasedKeys []ebiten.Key
}

func newKeyboardState() keyboardState {
	return keyboardState{
		justPressedKeys:  make([]ebiten.Key, 0),
		justReleasedKeys: make([]ebiten.Key, 0),
	}
}

// =====
// event
// =====
func emitKeyboard(w *ecs.World) {
	st := ecs.EnsureResource(w, newKeyboardState)

	st.justPressedKeys = inpututil.AppendJustPressedKeys(st.justPressedKeys[:0])
	for _, ek := range st.justPressedKeys {
		if k, ok := mapKey(ek); ok {
			event.AddEvent(w, input.KeyEvent{Key: k, Pressed: true})
		}
	}

	st.justReleasedKeys = inpututil.AppendJustReleasedKeys(st.justReleasedKeys[:0])
	for _, ek := range st.justReleasedKeys {
		if k, ok := mapKey(ek); ok {
			event.AddEvent(w, input.KeyEvent{Key: k, Pressed: false})
		}
	}
}

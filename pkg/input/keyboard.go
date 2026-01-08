package input

import (
	"fmt"
)

// ===
// key
// ===
type Key int

const (
	KeyA Key = iota
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ

	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	KeyNumpad0
	KeyNumpad1
	KeyNumpad2
	KeyNumpad3
	KeyNumpad4
	KeyNumpad5
	KeyNumpad6
	KeyNumpad7
	KeyNumpad8
	KeyNumpad9
	KeyNumpadAdd
	KeyNumpadDecimal
	KeyNumpadDivide
	KeyNumpadEnter
	KeyNumpadEqual
	KeyNumpadMultiply
	KeyNumpadSubtract

	KeyBracketLeft
	KeyBracketRight
	KeyComma
	KeyBackspace
	KeyBackslash
	KeyEqual
	KeyBackquote
	KeyIntlBackslash
	KeyMinus
	KeyPeriod
	KeyQuote
	KeySemicolon
	KeySlash
	KeySpace
	KeyTab

	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
	KeyArrowUp

	KeyAltLeft
	KeyAltRight
	KeyControlLeft
	KeyControlRight
	KeyMetaLeft
	KeyMetaRight
	KeyShiftLeft
	KeyShiftRight
	KeyAlt
	KeyControl
	KeyShift
	KeyMeta

	KeyCapsLock
	KeyContextMenu
	KeyDelete
	KeyEnd
	KeyEnter
	KeyEscape
	KeyHome
	KeyInsert
	KeyNumLock
	KeyPageDown
	KeyPageUp
	KeyPause
	KeyPrintScreen
	KeyScrollLock

	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyF13
	KeyF14
	KeyF15
	KeyF16
	KeyF17
	KeyF18
	KeyF19
	KeyF20
	KeyF21
	KeyF22
	KeyF23
	KeyF24

	KeyMax = KeyF24
)

func (k Key) String() string {
	switch k {
	case KeyA:
		return "KeyA"
	case KeyB:
		return "KeyB"
	case KeyC:
		return "KeyC"
	case KeyD:
		return "KeyD"
	case KeyE:
		return "KeyE"
	case KeyF:
		return "KeyF"
	case KeyG:
		return "KeyG"
	case KeyH:
		return "KeyH"
	case KeyI:
		return "KeyI"
	case KeyJ:
		return "KeyJ"
	case KeyK:
		return "KeyK"
	case KeyL:
		return "KeyL"
	case KeyM:
		return "KeyM"
	case KeyN:
		return "KeyN"
	case KeyO:
		return "KeyO"
	case KeyP:
		return "KeyP"
	case KeyQ:
		return "KeyQ"
	case KeyR:
		return "KeyR"
	case KeyS:
		return "KeyS"
	case KeyT:
		return "KeyT"
	case KeyU:
		return "KeyU"
	case KeyV:
		return "KeyV"
	case KeyW:
		return "KeyW"
	case KeyX:
		return "KeyX"
	case KeyY:
		return "KeyY"
	case KeyZ:
		return "KeyZ"
	case KeyAltLeft:
		return "KeyAltLeft"
	case KeyAltRight:
		return "KeyAltRight"
	case KeyArrowDown:
		return "KeyArrowDown"
	case KeyArrowLeft:
		return "KeyArrowLeft"
	case KeyArrowRight:
		return "KeyArrowRight"
	case KeyArrowUp:
		return "KeyArrowUp"
	case KeyBackquote:
		return "KeyBackquote"
	case KeyBackslash:
		return "KeyBackslash"
	case KeyBackspace:
		return "KeyBackspace"
	case KeyBracketLeft:
		return "KeyBracketLeft"
	case KeyBracketRight:
		return "KeyBracketRight"
	case KeyCapsLock:
		return "KeyCapsLock"
	case KeyComma:
		return "KeyComma"
	case KeyContextMenu:
		return "KeyContextMenu"
	case KeyControlLeft:
		return "KeyControlLeft"
	case KeyControlRight:
		return "KeyControlRight"
	case KeyDelete:
		return "KeyDelete"
	case Key0:
		return "Key0"
	case Key1:
		return "Key1"
	case Key2:
		return "Key2"
	case Key3:
		return "Key3"
	case Key4:
		return "Key4"
	case Key5:
		return "Key5"
	case Key6:
		return "Key6"
	case Key7:
		return "Key7"
	case Key8:
		return "Key8"
	case Key9:
		return "Key9"
	case KeyEnd:
		return "KeyEnd"
	case KeyEnter:
		return "KeyEnter"
	case KeyEqual:
		return "KeyEqual"
	case KeyEscape:
		return "KeyEscape"
	case KeyF1:
		return "KeyF1"
	case KeyF2:
		return "KeyF2"
	case KeyF3:
		return "KeyF3"
	case KeyF4:
		return "KeyF4"
	case KeyF5:
		return "KeyF5"
	case KeyF6:
		return "KeyF6"
	case KeyF7:
		return "KeyF7"
	case KeyF8:
		return "KeyF8"
	case KeyF9:
		return "KeyF9"
	case KeyF10:
		return "KeyF10"
	case KeyF11:
		return "KeyF11"
	case KeyF12:
		return "KeyF12"
	case KeyF13:
		return "KeyF13"
	case KeyF14:
		return "KeyF14"
	case KeyF15:
		return "KeyF15"
	case KeyF16:
		return "KeyF16"
	case KeyF17:
		return "KeyF17"
	case KeyF18:
		return "KeyF18"
	case KeyF19:
		return "KeyF19"
	case KeyF20:
		return "KeyF20"
	case KeyF21:
		return "KeyF21"
	case KeyF22:
		return "KeyF22"
	case KeyF23:
		return "KeyF23"
	case KeyF24:
		return "KeyF24"
	case KeyHome:
		return "KeyHome"
	case KeyInsert:
		return "KeyInsert"
	case KeyIntlBackslash:
		return "KeyIntlBackslash"
	case KeyMetaLeft:
		return "KeyMetaLeft"
	case KeyMetaRight:
		return "KeyMetaRight"
	case KeyMinus:
		return "KeyMinus"
	case KeyNumLock:
		return "KeyNumLock"
	case KeyNumpad0:
		return "KeyNumpad0"
	case KeyNumpad1:
		return "KeyNumpad1"
	case KeyNumpad2:
		return "KeyNumpad2"
	case KeyNumpad3:
		return "KeyNumpad3"
	case KeyNumpad4:
		return "KeyNumpad4"
	case KeyNumpad5:
		return "KeyNumpad5"
	case KeyNumpad6:
		return "KeyNumpad6"
	case KeyNumpad7:
		return "KeyNumpad7"
	case KeyNumpad8:
		return "KeyNumpad8"
	case KeyNumpad9:
		return "KeyNumpad9"
	case KeyNumpadAdd:
		return "KeyNumpadAdd"
	case KeyNumpadDecimal:
		return "KeyNumpadDecimal"
	case KeyNumpadDivide:
		return "KeyNumpadDivide"
	case KeyNumpadEnter:
		return "KeyNumpadEnter"
	case KeyNumpadEqual:
		return "KeyNumpadEqual"
	case KeyNumpadMultiply:
		return "KeyNumpadMultiply"
	case KeyNumpadSubtract:
		return "KeyNumpadSubtract"
	case KeyPageDown:
		return "KeyPageDown"
	case KeyPageUp:
		return "KeyPageUp"
	case KeyPause:
		return "KeyPause"
	case KeyPeriod:
		return "KeyPeriod"
	case KeyPrintScreen:
		return "KeyPrintScreen"
	case KeyQuote:
		return "KeyQuote"
	case KeyScrollLock:
		return "KeyScrollLock"
	case KeySemicolon:
		return "KeySemicolon"
	case KeyShiftLeft:
		return "KeyShiftLeft"
	case KeyShiftRight:
		return "KeyShiftRight"
	case KeySlash:
		return "KeySlash"
	case KeySpace:
		return "KeySpace"
	case KeyTab:
		return "KeyTab"
	}
	return fmt.Sprintf("Key(%d)", k)
}

// ========
// keyboard
// ========
type Keyboard = ButtonInput[Key]

func newKeyboard() Keyboard {
	return newButtonInput[Key]()
}

// ======
// events
// ======
type KeyEvent struct {
	Key     Key
	Pressed bool
}

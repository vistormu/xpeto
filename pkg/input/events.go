package input

type KeyJustPressed struct {
	Key Key
}

type KeyJustReleased struct {
	Key Key
}

type MouseButtonJustPressed struct {
	Button MouseButton
}

type MouseButtonJustReleased struct {
	Button MouseButton
}

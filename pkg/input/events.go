package input

type EventKeyJustPressed struct {
	Key Key
}

type EventKeyJustReleased struct {
	Key Key
}

type EventMouseButtonJustPressed struct {
	Button MouseButton
}

type EventMouseButtonJustReleased struct {
	Button MouseButton
}

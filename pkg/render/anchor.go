package render

type Anchor uint8

const (
	AnchorCenter Anchor = iota
	AnchorTopLeft
	AnchorTop
	AnchorTopRight
	AnchorLeft
	AnchorRight
	AnchorBottomLeft
	AnchorBottom
	AnchorBottomRight
)

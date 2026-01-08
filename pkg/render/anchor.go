package render

type AnchorType uint8

const (
	AnchorCenter AnchorType = iota
	AnchorTopLeft
	AnchorTop
	AnchorTopRight
	AnchorLeft
	AnchorRight
	AnchorBottomLeft
	AnchorBottom
	AnchorBottomRight
)

type Anchor struct {
	Type AnchorType
}

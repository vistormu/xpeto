package text

import (
	"image/color"
	"math"

	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/render"
)

// =====
// align
// =====
type Align uint8

const (
	AlignStart Align = iota
	AlignCenter
	AlignEnd
)

// ====
// wrap
// ====
type WrapMode uint8

const (
	WrapNone WrapMode = iota
	WrapWord
	WrapRune
)

// ====
// text
// ====
type Text struct {
	Font asset.Asset

	Content string
	Color   color.Color

	Size     float64
	MaxWidth float64

	Align Align
	Wrap  WrapMode

	OrderKey render.OrderKey
}

func NewText(content string, opts ...option) Text {
	t := Text{
		Font:     asset.Asset(0),
		Content:  content,
		Color:    nil,
		Size:     18,
		Align:    AlignStart,
		MaxWidth: 0,
		Wrap:     WrapNone,
		OrderKey: render.NewOrderKey(0, 0, 0),
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&t)
		}
	}

	if t.Color == nil {
		t.Color = color.White
	}

	if math.IsNaN(t.Size) || math.IsInf(t.Size, 0) || t.Size <= 0 {
		t.Size = 18
	}

	if math.IsNaN(t.MaxWidth) || math.IsInf(t.MaxWidth, 0) || t.MaxWidth < 0 {
		t.MaxWidth = 0
	}

	return t
}

// =======
// options
// =======
type option = func(*Text)

type textOpt struct{}

var TextOpt textOpt

func (textOpt) Font(font asset.Asset) option {
	return func(t *Text) { t.Font = font }
}

func (textOpt) Color(c color.Color) option {
	return func(t *Text) { t.Color = c }
}

func (textOpt) Size(size float64) option {
	return func(t *Text) { t.Size = size }
}

func (textOpt) Align(a Align) option {
	return func(t *Text) { t.Align = a }
}

func (textOpt) MaxWidth(w float64) option {
	return func(t *Text) { t.MaxWidth = w }
}

func (textOpt) Wrap(m WrapMode) option {
	return func(t *Text) { t.Wrap = m }
}

func (textOpt) Order(layer uint16, order uint16, tie uint32) option {
	return func(t *Text) { t.OrderKey = render.NewOrderKey(layer, order, tie) }
}

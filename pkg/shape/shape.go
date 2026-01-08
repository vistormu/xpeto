package shape

import (
	"image/color"

	"github.com/vistormu/go-dsa/geometry"
	"github.com/vistormu/xpeto/pkg/render"
)

type ShapeKind uint8

const (
	None ShapeKind = iota
	Arrow
	Capsule
	Ellipse
	Line
	Path
	Polygon
	Ray
	Rect
	Segment
)

// =====
// shape
// =====
type Shape struct {
	Kind ShapeKind

	Arrow   geometry.Arrow[float32]
	Capsule geometry.Capsule[float32]
	Ellipse geometry.Ellipse[float32]
	Line    geometry.Line[float32]
	Path    geometry.Path[float32]
	Polygon geometry.Polygon[float32]
	Ray     geometry.Ray[float32]
	Rect    geometry.Rect[float32]
	Segment geometry.Segment[float32]

	Fills    []Fill
	Strokes  []Stroke
	OrderKey render.OrderKey
}

func NewShape(opts ...option) Shape {
	s := Shape{
		Kind:     None,
		Fills:    make([]Fill, 0),
		Strokes:  make([]Stroke, 0),
		OrderKey: render.NewOrderKey(0, 0, 0),
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&s)
		}
	}

	return s
}

// =======
// options
// =======
type option func(*Shape)

type shapeOpt struct{}

var ShapeOpt shapeOpt

func (s shapeOpt) FillSolid(c color.Color) option {
	return func(shape *Shape) {
		shape.Fills = append(shape.Fills, Fill{
			Type:  FillSolid,
			Color: c,
		})
	}
}

func (shapeOpt) Stroke(c color.Color, w float32) option {
	return func(s *Shape) {
		if w < 0 {
			w = 0
		}
		s.Strokes = append(s.Strokes, Stroke{
			Color: c,
			Width: w,
		})
	}
}

func (shapeOpt) Order(layer uint16, order uint16, tie uint32) option {
	return func(s *Shape) {
		s.OrderKey = render.NewOrderKey(layer, order, tie)
	}
}

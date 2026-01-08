package shape

import (
	"github.com/vistormu/go-dsa/geometry"
)

// =====
// arrow
// =====
func (s shapeOpt) Arrow(x1, y1, x2, y2, headWidth, headLength float32) option {
	return func(s *Shape) {
		s.Kind = Arrow
		s.Arrow = geometry.NewArrow(
			geometry.NewVector(x1, y1),
			geometry.NewVector(x2, y2),
			headLength,
			headWidth,
		)
	}
}

// =======
// capsule
// =======
func (s shapeOpt) Capsule(x1, y1, x2, y2, r float32) option {
	return func(s *Shape) {
		s.Kind = Capsule
		s.Capsule = geometry.NewCapsule(
			geometry.NewSegment(
				geometry.NewVector(x1, y1),
				geometry.NewVector(x2, y2),
			),
			r,
		)
	}
}

// =======
// ellipse
// =======
func (s shapeOpt) Ellipse(rx, ry float32) option {
	return func(s *Shape) {
		s.Kind = Ellipse
		s.Ellipse = geometry.NewEllipse(rx, ry)
	}
}

func (s shapeOpt) Circle(r float32) option {
	return s.Ellipse(r, r)
}

func (s shapeOpt) CircleD(d float32) option {
	return s.Ellipse(d/2, d/2)
}

// ====
// line
// ====
func (s shapeOpt) Line(x1, y1, x2, y2 float32) option {
	return func(s *Shape) {
		s.Kind = Line
		s.Line = geometry.NewLine(
			geometry.NewVector(x1, y1),
			geometry.NewVector(x2, y2),
		)
	}
}

// ====
// path
// ====
func (s shapeOpt) Path(xy ...float32) option {
	return func(s *Shape) {
		s.Kind = Path
		s.Path = geometry.NewPath[float32]()

		if len(xy)%2 != 0 {
			return
		}

		for i := 0; i < len(xy); i += 2 {
			s.Path.AddXY(xy[i], xy[i+1])
		}
	}
}

// =======
// polygon
// =======
func (s shapeOpt) Polygon(xy ...float32) option {
	return func(s *Shape) {
		s.Kind = Polygon
		s.Polygon = geometry.NewPolygon[float32]()

		if len(xy)%2 != 0 {
			return
		}

		for i := 0; i < len(xy); i += 2 {
			s.Polygon.Add(geometry.NewVector(xy[i], xy[i+1]))
		}
	}
}

// ===
// ray
// ===
func (s shapeOpt) Ray(x1, y1, x2, y2 float32) option {
	return func(s *Shape) {
		s.Kind = Ray
		s.Ray = geometry.NewRay(
			geometry.NewVector(x1, y1),
			geometry.NewVector(x2, y2),
		)
	}
}

// ====
// rect
// ====
func (s shapeOpt) Rect(w, h float32) option {
	return func(s *Shape) {
		s.Kind = Rect
		s.Rect = geometry.NewRect(w, h)
	}
}

func (s shapeOpt) Square(v float32) option {
	return s.Rect(v, v)
}

// =======
// segment
// =======
func (s shapeOpt) Segment(x1, y1, x2, y2 float32) option {
	return func(s *Shape) {
		s.Kind = Segment
		s.Segment = geometry.NewSegment(
			geometry.NewVector(x1, y1),
			geometry.NewVector(x2, y2),
		)
	}
}

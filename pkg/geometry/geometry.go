package geometry

import "github.com/vistormu/go-dsa/constraints"

type GeometryKind uint8

const (
	GeometryNone GeometryKind = iota
	GeometryEllipse
	GeometryRect
	GeometryCapsule
	GeometryPolygon
	GeometrySegment
	GeometryRay
	GeometryPath
)

type Geometry[T constraints.Number] struct {
	Kind    GeometryKind
	Ellipse Ellipse[T]
	Rect    Rect[T]
	Capsule Capsule[T]
	Polygon Polygon[T]
	Segment Segment[T]
	Ray     Ray[T]
	Path    Path[T]
}

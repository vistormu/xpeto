package physics

type BodyType uint8

const (
	Static BodyType = iota
	Kinematic
	Dynamic
)

type RigidBody struct {
	Type BodyType
	Mass float64
}

type Aabb struct {
	HalfWidth  float64
	HalfHeight float64
}

type Velocity struct {
	X, Y float64
}

type GravityScale struct {
	Value float64
}

package physics

type BodyType uint8

const (
	Static BodyType = iota
	Kinematic
	Dynamic
)

type RigidBody struct {
	Type        BodyType
	Mass        float64
	Restitution float64
	Friction    float64
}

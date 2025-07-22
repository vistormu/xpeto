package platformer

type Playable struct {
	MaxSpeed           float32
	GroundAcceleration float32
	AirAcceleration    float32
	GroundFriction     float32
	AirFriction        float32
	JumpSpeed          float32
	ShortHopFactor     float32
	CoyoteTime         float32
	BufferTime         float32
}

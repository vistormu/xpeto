package physics

import "github.com/vistormu/xpeto/internal/ecs"

type EventCollisionType uint8

const (
	CollisionStarted EventCollisionType = iota
	CollisionEnded
	CollisionStay
)

type EventCollision struct {
	A, B         ecs.Entity
	Type         EventCollisionType
	ContactCount int
}

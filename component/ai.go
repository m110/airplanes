package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

const (
	AITypeConstantVelocity AIType = iota
	AITypeFollowPath
)

type AIType int

type AIData struct {
	Type AIType

	Speed float64

	StartedMoving bool

	Path       []math.Vec2
	PathLoops  bool
	NextTarget int
}

var AI = donburi.NewComponentType[AIData]()

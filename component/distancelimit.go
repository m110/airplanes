package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type DistanceLimitData struct {
	MaxDistance float64

	DistanceTraveled float64
	PreviousPosition math.Vec2
	Initialized      bool
}

var DistanceLimit = donburi.NewComponentType[DistanceLimitData]()

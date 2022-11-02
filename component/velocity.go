package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type VelocityData struct {
	Velocity         math.Vec2
	RotationVelocity float64
}

var Velocity = donburi.NewComponentType[VelocityData]()

func GetVelocity(entry *donburi.Entry) *VelocityData {
	return donburi.Get[VelocityData](entry, Velocity)
}

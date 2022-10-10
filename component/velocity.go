package component

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/engine"
)

type VelocityData struct {
	Velocity engine.Vector
}

var Velocity = donburi.NewComponentType[VelocityData]()

func GetVelocity(entry *donburi.Entry) *VelocityData {
	return donburi.Get[VelocityData](entry, Velocity)
}

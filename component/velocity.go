package component

import "github.com/yohamta/donburi"

type VelocityData struct {
	X        float64
	Y        float64
	Rotation float64
}

var Velocity = donburi.NewComponentType[VelocityData]()

func GetVelocity(entry *donburi.Entry) *VelocityData {
	return donburi.Get[VelocityData](entry, Velocity)
}

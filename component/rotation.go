package component

import "github.com/yohamta/donburi"

type RotationData struct {
	Angle float64
}

var Rotation = donburi.NewComponentType[RotationData]()

func GetRotation(entry *donburi.Entry) *RotationData {
	return donburi.Get[RotationData](entry, Rotation)
}

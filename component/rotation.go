package component

import "github.com/yohamta/donburi"

type RotationData struct {
	Angle float64
	// The original rotation of the sprite
	// "Facing right" is considered 0 degrees
	OriginalAngle float64
}

var Rotation = donburi.NewComponentType[RotationData]()

func GetRotation(entry *donburi.Entry) *RotationData {
	return donburi.Get[RotationData](entry, Rotation)
}

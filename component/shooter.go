package component

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/engine"
)

type ShooterType int

const (
	ShooterTypeRocket ShooterType = iota
	ShooterTypeMissiles
	ShooterTypeBeam
)

type ShooterData struct {
	Type       ShooterType
	ShootTimer *engine.Timer
}

var Shooter = donburi.NewComponentType[ShooterData]()

func GetShooter(entry *donburi.Entry) *ShooterData {
	return donburi.Get[ShooterData](entry, Shooter)
}

package component

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/engine"
)

type ShooterType int

const (
	ShooterTypeBullet ShooterType = iota
	ShooterTypeMissile
	ShooterTypeBeam
)

type ShooterData struct {
	Type       ShooterType
	ShootTimer *engine.Timer
}

var Shooter = donburi.NewComponentType[ShooterData]()

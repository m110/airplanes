package component

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/engine"
)

type AltitudeData struct {
	// Altitude is the level above the ground, in percent.
	// 0.0 is ground level, 1.0 is the highest point.
	Altitude float64
	Velocity float64

	// TODO: Not sure if this fits this component
	Falling bool
}

func (a *AltitudeData) Update() {
	a.Altitude = engine.Clamp(a.Altitude+a.Velocity, 0, 1)
}

var Altitude = donburi.NewComponentType[AltitudeData]()

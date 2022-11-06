package component

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/engine"
)

type CameraData struct {
	Moving    bool
	MoveTimer *engine.Timer
}

var Camera = donburi.NewComponentType[CameraData]()

package component

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/engine"
)

type PlayerAirplaneData struct {
	PlayerNumber int

	Invulnerable      bool
	InvulnerableTimer *engine.Timer
}

var PlayerAirplane = donburi.NewComponentType[PlayerAirplaneData]()

func GetPlayerAirplane(entry *donburi.Entry) *PlayerAirplaneData {
	return donburi.Get[PlayerAirplaneData](entry, PlayerAirplane)
}

package component

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/engine"
)

type PlayerData struct {
	PlayerNumber int

	Lives        int
	Respawning   bool
	RespawnTimer *engine.Timer
}

func (d *PlayerData) Damage() {
	if d.Respawning {
		return
	}

	d.Lives--

	if d.Lives > 0 {
		d.Respawning = true
		d.RespawnTimer.Reset()
	}
}

var Player = donburi.NewComponentType[PlayerData]()

func GetPlayer(entry *donburi.Entry) *PlayerData {
	return donburi.Get[PlayerData](entry, Player)
}

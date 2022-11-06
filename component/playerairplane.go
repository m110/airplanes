package component

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/engine"
)

type PlayerAirplaneData struct {
	PlayerNumber int
	// TODO Duplicated across PlayerAirplane and Player?
	Faction PlayerFaction

	Invulnerable          bool
	InvulnerableTimer     *engine.Timer
	InvulnerableIndicator *SpriteData
}

func (d *PlayerAirplaneData) StartInvulnerability() {
	d.Invulnerable = true
	d.InvulnerableTimer.Reset()
	d.InvulnerableIndicator.Hidden = false
}

func (d *PlayerAirplaneData) StopInvulnerability() {
	d.Invulnerable = false
	d.InvulnerableIndicator.Hidden = true
}

var PlayerAirplane = donburi.NewComponentType[PlayerAirplaneData]()

package component

import "github.com/yohamta/donburi"

type PlayerSelectData struct {
	Index   int
	Faction PlayerFaction

	Selected     bool
	Ready        bool
	PlayerNumber int
}

func (p *PlayerSelectData) Select(playerNumber int) {
	p.Selected = true
	p.PlayerNumber = playerNumber
}

func (p *PlayerSelectData) Unselect() {
	p.Selected = false
	p.PlayerNumber = 0
}

func (p *PlayerSelectData) LockIn() {
	p.Ready = true
}

func (p *PlayerSelectData) Release() {
	p.Ready = false
}

var PlayerSelect = donburi.NewComponentType[PlayerSelectData]()

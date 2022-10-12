package component

import "github.com/yohamta/donburi"

type PlayerSelectData struct {
	Index   int
	Faction PlayerFaction

	Selected     bool
	Ready        bool
	PlayerNumber int
}

var PlayerSelect = donburi.NewComponentType[PlayerSelectData]()

func GetPlayerSelect(entry *donburi.Entry) *PlayerSelectData {
	return donburi.Get[PlayerSelectData](entry, PlayerSelect)
}

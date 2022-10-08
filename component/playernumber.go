package component

import "github.com/yohamta/donburi"

type PlayerNumberData struct {
	Number int
}

var PlayerNumber = donburi.NewComponentType[PlayerNumberData]()

func GetPlayerNumber(entry *donburi.Entry) *PlayerNumberData {
	return donburi.Get[PlayerNumberData](entry, PlayerNumber)
}

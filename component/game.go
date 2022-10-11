package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type GameData struct {
	Score    int
	Settings Settings
}

func (d *GameData) AddScore(score int) {
	d.Score += score
}

type Settings struct {
	ScreenWidth  int
	ScreenHeight int
}

var Game = donburi.NewComponentType[GameData]()

func GetGame(entry *donburi.Entry) *GameData {
	return donburi.Get[GameData](entry, Game)
}

func MustFindGame(w donburi.World) *GameData {
	game, ok := query.NewQuery(filter.Contains(Game)).FirstEntity(w)
	if !ok {
		panic("game not found")
	}
	return GetGame(game)
}

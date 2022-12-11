package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type GameData struct {
	Score    int
	Paused   bool
	GameOver bool
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

func MustFindGame(w donburi.World) *GameData {
	game, ok := query.NewQuery(filter.Contains(Game)).First(w)
	if !ok {
		panic("game not found")
	}
	return Game.Get(game)
}

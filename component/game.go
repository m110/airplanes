package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type GameStatus struct {
	Score    int
	Paused   bool
	GameOver bool
	Settings Settings
}

func (d *GameStatus) AddScore(score int) {
	d.Score += score
}

type Settings struct {
	ScreenWidth  int
	ScreenHeight int
}

var Game = donburi.NewComponentType[GameStatus]()

func GetGame(entry *donburi.Entry) *GameStatus {
	return donburi.Get[GameStatus](entry, Game)
}

func MustFindGame(w donburi.World) *GameStatus {
	game, ok := query.NewQuery(filter.Contains(Game)).FirstEntity(w)
	if !ok {
		panic("game not found")
	}
	return GetGame(game)
}

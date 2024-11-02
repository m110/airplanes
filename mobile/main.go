package mobile

import (
	"github.com/hajimehoshi/ebiten/v2/mobile"
	"github.com/m110/airplanes/game"
)

func init() {
	mobile.SetGame(game.NewGame())
}

func Dummy() {}

package mobile

import (
	"github.com/hajimehoshi/ebiten/v2/mobile"

	"github.com/m110/airplanes/game"
)

func init() {
	mobile.SetGame(game.NewGame(game.Config{
		Quick:        true,
		ScreenWidth:  480,
		ScreenHeight: 1040,
	}))
}

func Dummy() {}

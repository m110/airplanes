package scene

import (
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/m110/airplanes/assets"
)

type Title struct {
	screenWidth     int
	screenHeight    int
	newGameCallback func()
}

func NewTitle(screenWidth int, screenHeight int, newGameCallback func()) *Title {
	return &Title{
		screenWidth:     screenWidth,
		screenHeight:    screenHeight,
		newGameCallback: newGameCallback,
	}
}

func (t *Title) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeySpace) {
		t.newGameCallback()
		return
	}

	touchIDs := inpututil.AppendJustPressedTouchIDs(nil)
	if len(touchIDs) > 0 {
		t.newGameCallback()
		return
	}
}

func (t *Title) Draw(screen *ebiten.Image) {
	text.Draw(screen, "m110's Airplanes", assets.NarrowFont, t.screenWidth/4, 100, color.White)
	text.Draw(screen, "Player 1: WASD + Space", assets.NarrowFont, t.screenWidth/6, 250, color.White)
	text.Draw(screen, "Player 2: Arrows + Enter", assets.NarrowFont, t.screenWidth/6, 350, color.White)
	text.Draw(screen, "Press space to start", assets.NarrowFont, t.screenWidth/5, 500, color.White)
}

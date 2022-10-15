package scene

import (
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
}

func (t *Title) Draw(screen *ebiten.Image) {
	text.Draw(screen, "m110's Airplanes", assets.NarrowFont, t.screenWidth/4, 100, color.White)
}
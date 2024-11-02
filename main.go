package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/scene"
	"github.com/m110/airplanes/system"
)

var (
	screenWidth  = 480
	screenHeight = 640
)

type Scene interface {
	Update()
	Draw(screen *ebiten.Image)
}

type Game struct {
	scene Scene
}

func NewGame() *Game {
	assets.MustLoadAssets()

	g := &Game{}
	g.switchToTitle()
	return g
}

func (g *Game) switchToTitle() {
	g.scene = scene.NewTitle(screenWidth, screenHeight, g.switchToAirbase)
}

func (g *Game) switchToAirbase() {
	g.scene = scene.NewAirbase(screenWidth, screenHeight, g.switchToGame, g.switchToTitle)
}

func (g *Game) switchToGame(players []system.ChosenPlayer) {
	g.scene = scene.NewGame(players, screenWidth, screenHeight)
}

func (g *Game) Update() error {
	g.scene.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.scene.Draw(screen)
}

func (g *Game) Layout(width, height int) (int, int) {
	return width, height
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)

	err := ebiten.RunGame(NewGame())
	if err != nil {
		log.Fatal(err)
	}
}

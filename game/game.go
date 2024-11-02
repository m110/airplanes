package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/scene"
	"github.com/m110/airplanes/system"
)

type Scene interface {
	Update()
	Draw(screen *ebiten.Image)
}

type Game struct {
	scene        Scene
	screenWidth  int
	screenHeight int
}

func NewGame() *Game {
	assets.MustLoadAssets()

	return &Game{}
}

func (g *Game) switchToTitle() {
	g.scene = scene.NewTitle(g.screenWidth, g.screenHeight, g.switchToAirbase)
}

func (g *Game) switchToAirbase() {
	g.scene = scene.NewAirbase(g.screenWidth, g.screenHeight, g.switchToGame, g.switchToTitle)
}

func (g *Game) switchToGame(players []system.ChosenPlayer) {
	g.scene = scene.NewGame(players, g.screenWidth, g.screenHeight)
}

func (g *Game) Update() error {
	if g.scene == nil {
		g.switchToTitle()
	}
	g.scene.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.scene.Draw(screen)
}

func (g *Game) Layout(width, height int) (int, int) {
	g.screenWidth = width
	g.screenHeight = height
	return width, height
}

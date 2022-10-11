package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/scenes"
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
	g.scene = scenes.NewTitle(screenWidth, screenHeight, g.switchToGame)
}

func (g *Game) switchToGame() {
	g.scene = scenes.NewGame(screenWidth, screenHeight)
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
	rand.Seed(time.Now().UTC().UnixNano())

	err := ebiten.RunGame(NewGame())
	if err != nil {
		log.Fatal(err)
	}
}

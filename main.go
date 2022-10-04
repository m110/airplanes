package main

import (
	"github.com/m110/airplanes/archetypes"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/system"
)

type System interface {
	Update(w donburi.World)
}

type Drawable interface {
	Draw(w donburi.World, screen *ebiten.Image)
}

type Game struct {
	world     donburi.World
	systems   []System
	drawables []Drawable
}

func NewGame() *Game {
	assets.LoadAssets()

	g := &Game{
		world: createWorld(),
	}

	g.systems = []System{
		system.NewVelocity(),
		system.NewControls(),
	}
	g.drawables = []Drawable{
		system.NewRenderer(),
	}

	return g
}

func createWorld() donburi.World {
	world := donburi.NewWorld()

	archetypes.NewPlayerOne(world)
	archetypes.NewPlayerTwo(world)

	return world
}

func (g *Game) Update() error {
	for _, s := range g.systems {
		s.Update(g.world)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	for _, s := range g.drawables {
		s.Draw(g.world, screen)
	}
}

func (g *Game) Layout(width, height int) (int, int) {
	return width, height
}

func main() {
	ebiten.SetWindowSize(800, 600)
	rand.Seed(time.Now().UTC().UnixNano())

	err := ebiten.RunGame(NewGame())
	if err != nil {
		log.Fatal(err)
	}
}

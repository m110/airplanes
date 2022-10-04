package main

import (
	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/system"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
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

func NewGame() (*Game, error) {
	err := assets.LoadAssets()
	if err != nil {
		return nil, err
	}

	g := &Game{
		world: createWorld(),
	}

	g.systems = []System{}
	g.drawables = []Drawable{
		system.NewRenderer(),
	}

	return g, nil
}

func createWorld() donburi.World {
	world := donburi.NewWorld()

	airplaneEntity := world.Create(component.Position, component.Velocity, component.Sprite)
	airplane := world.Entry(airplaneEntity)
	donburi.SetValue(airplane, component.Position, component.PositionData{X: 100, Y: 100})
	donburi.SetValue(airplane, component.Velocity, component.VelocityData{X: 1, Y: 1})
	donburi.SetValue(airplane, component.Sprite, component.SpriteData{Image: assets.ShipYellowSmall})

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

	g, err := NewGame()
	if err != nil {
		log.Fatal(err)
	}

	err = ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/archetypes"
	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/system"
)

var (
	screenWidth  = 480
	screenHeight = 600
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
		system.NewControls(),
		system.NewVelocity(),
		system.NewBounds(screenWidth, screenHeight),
	}
	g.drawables = []Drawable{
		system.NewRenderer(),
	}

	return g
}

func createWorld() donburi.World {
	world := donburi.NewWorld()

	archetypes.NewPlayerOne(world, component.PositionData(assets.Level1.Player1Spawn))
	archetypes.NewPlayerTwo(world, component.PositionData(assets.Level1.Player2Spawn))

	camera := world.Create(component.CameraTag, component.Position, component.Velocity)
	cameraEntry := world.Entry(camera)
	donburi.SetValue(cameraEntry, component.Position, component.PositionData{
		X: 0,
		Y: float64(assets.Level1.Background.Bounds().Dy() - screenHeight),
	})
	donburi.SetValue(cameraEntry, component.Velocity, component.VelocityData{
		Y: -0.5,
	})

	level := world.Create(component.Position, component.Sprite)
	levelEntry := world.Entry(level)
	donburi.SetValue(levelEntry, component.Sprite, component.SpriteData{
		Image: assets.Level1.Background,
		Layer: component.SpriteLayerBackground,
	})

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
	ebiten.SetWindowSize(screenWidth, screenHeight)
	rand.Seed(time.Now().UTC().UnixNano())

	err := ebiten.RunGame(NewGame())
	if err != nil {
		log.Fatal(err)
	}
}

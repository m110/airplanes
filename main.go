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
	"github.com/m110/airplanes/engine"
	"github.com/m110/airplanes/system"
)

var (
	screenWidth  = 480
	screenHeight = 640
)

type System interface {
	Update(w donburi.World)
}

type Drawable interface {
	Draw(w donburi.World, screen *ebiten.Image)
}

type Game struct {
	level     int
	world     donburi.World
	systems   []System
	drawables []Drawable
}

func NewGame() *Game {
	assets.LoadAssets()

	g := &Game{
		level: 0,
	}

	debug := system.NewDebug()

	g.systems = []System{
		system.NewControls(),
		system.NewVelocity(),
		system.NewBounds(screenWidth, screenHeight),
		system.NewCameraBounds(),
		system.NewAI(),
		system.NewDespawn(screenWidth, screenHeight),
		system.NewProgression(g.nextLevel),
		debug,
	}

	g.drawables = []Drawable{
		system.NewRenderer(),
		debug,
	}

	g.loadLevel()

	return g
}

func (g *Game) nextLevel() {
	if g.level == len(assets.Levels)-1 {
		// TODO all levels done
		return
	}
	g.level++
	g.loadLevel()
}

func (g *Game) loadLevel() {
	g.world = createWorld(g.level)
}

func createWorld(levelIndex int) donburi.World {
	levelAsset := assets.Levels[levelIndex]

	world := donburi.NewWorld()

	level := world.Entry(world.Create(component.Level))
	component.GetLevel(level).ProgressionTimer = engine.NewTimer(time.Second * 3)

	archetypes.NewPlayerOne(world, component.PositionData(levelAsset.Player1Spawn))
	archetypes.NewPlayerTwo(world, component.PositionData(levelAsset.Player2Spawn))

	archetypes.NewCamera(world, component.PositionData{
		X: 0,
		Y: float64(levelAsset.Background.Bounds().Dy() - screenHeight),
	})

	levelEntity := world.Create(component.Position, component.Sprite)
	levelEntry := world.Entry(levelEntity)
	donburi.SetValue(levelEntry, component.Sprite, component.SpriteData{
		Image: levelAsset.Background,
		Layer: component.SpriteLayerBackground,
	})

	for _, enemy := range levelAsset.Enemies {
		archetypes.NewEnemy(
			world,
			component.PositionData(enemy.Position),
			enemy.Rotation,
			enemy.Speed,
			enemy.Path,
		)
	}

	world.Create(component.Debug)

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

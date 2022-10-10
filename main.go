package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

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
	assets.MustLoadAssets()

	g := &Game{
		level: 0,
	}

	render := system.NewRenderer()
	debug := system.NewDebug()

	g.systems = []System{
		system.NewControls(),
		system.NewVelocity(),
		system.NewBounds(),
		system.NewCameraBounds(),
		system.NewAI(),
		system.NewDespawn(),
		system.NewCollision(),
		system.NewProgression(g.nextLevel),
		system.NewHealth(),
		system.NewRespawn(),
		system.NewInvulnerable(),
		render,
		debug,
	}

	g.drawables = []Drawable{
		render,
		debug,
		system.NewHUD(),
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
	// TODO Customizable number of players
	g.world = g.createWorld(g.level, 2)
}

func (g *Game) createWorld(levelIndex int, players int) donburi.World {
	levelAsset := assets.Levels[levelIndex]

	world := donburi.NewWorld()

	settings := world.Entry(world.Create(component.Settings))
	donburi.SetValue(settings, component.Settings, component.SettingsData{
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
	})

	level := world.Entry(world.Create(component.Level))
	component.GetLevel(level).ProgressionTimer = engine.NewTimer(time.Second * 3)

	archetypes.NewCamera(world, engine.Vector{
		X: 0,
		Y: float64(levelAsset.Background.Bounds().Dy() - screenHeight),
	})

	levelEntity := world.Create(component.Transform, component.Sprite)
	levelEntry := world.Entry(levelEntity)
	donburi.SetValue(levelEntry, component.Sprite, component.SpriteData{
		Image: levelAsset.Background,
		Layer: component.SpriteLayerBackground,
		Pivot: component.SpritePivotTopLeft,
	})

	for _, enemy := range levelAsset.Enemies {
		switch enemy.Class {
		case assets.EnemyClassAirplane:
			archetypes.NewEnemyAirplane(
				world,
				enemy.Position,
				enemy.Rotation,
				enemy.Speed,
				enemy.Path,
			)
		case assets.EnemyClassTank:
			archetypes.NewEnemyTank(
				world,
				enemy.Position,
				enemy.Rotation,
				enemy.Speed,
				enemy.Path,
			)
		default:
			panic("unknown enemy class: " + enemy.Class)
		}
	}

	if g.world == nil {
		// Spawn new players
		for i := 1; i <= players; i++ {
			player := archetypes.NewPlayer(world, i)
			archetypes.NewPlayerAirplane(world, *component.GetPlayer(player))
		}
	} else {
		// Transfer existing players from the previous level
		query.NewQuery(filter.Contains(component.Player)).EachEntity(g.world, func(entry *donburi.Entry) {
			player := component.GetPlayer(entry)
			// In case the level ends while the player's respawning
			player.Respawning = false

			archetypes.NewPlayerFromPlayerData(world, *player)
			if player.Lives > 0 {
				archetypes.NewPlayerAirplane(world, *player)
			}
		})
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

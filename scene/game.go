package scene

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
	"github.com/m110/airplanes/system"
)

type System interface {
	Update(w donburi.World)
}

type Drawable interface {
	Draw(w donburi.World, screen *ebiten.Image)
}

type Game struct {
	players   []system.ChosenPlayer
	level     int
	world     donburi.World
	systems   []System
	drawables []Drawable

	screenWidth  int
	screenHeight int
}

func NewGame(players []system.ChosenPlayer, screenWidth int, screenHeight int) *Game {
	g := &Game{
		players:      players,
		level:        0,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
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
	render := system.NewRenderer()
	debug := system.NewDebug(g.loadLevel)

	g.systems = []System{
		system.NewControls(),
		system.NewVelocity(),
		system.NewBounds(),
		system.NewCameraBounds(),
		system.NewSpawn(),
		system.NewAI(),
		system.NewDespawn(),
		system.NewCollision(),
		system.NewProgression(g.nextLevel),
		system.NewHealth(),
		system.NewRespawn(g.restart),
		system.NewInvulnerable(),
		system.NewCamera(),
		system.NewObserver(),
		system.NewShooter(),
		system.NewEvolution(),
		system.NewAltitude(),
		system.NewEvents(),
		render,
		debug,
	}

	g.drawables = []Drawable{
		render,
		debug,
		system.NewHUD(),
	}

	g.world = g.createWorld(g.level)
}

func (g *Game) createWorld(levelIndex int) donburi.World {
	levelAsset := assets.Levels[levelIndex]

	world := donburi.NewWorld()

	level := world.Entry(world.Create(component.Level))
	component.Level.Get(level).ProgressionTimer = engine.NewTimer(time.Second * 3)

	archetype.NewCamera(world, math.Vec2{
		X: 0,
		Y: float64(levelAsset.Background.Bounds().Dy() - g.screenHeight),
	})

	levelEntry := world.Entry(
		world.Create(transform.Transform, component.Sprite),
	)
	component.Sprite.SetValue(levelEntry, component.SpriteData{
		Image: levelAsset.Background,
		Layer: component.SpriteLayerBackground,
		Pivot: component.SpritePivotTopLeft,
	})

	for i := range levelAsset.Enemies {
		enemy := levelAsset.Enemies[i]
		pos := enemy.Position
		// TODO Sprite offset could be based on the sprite
		pos.Y += 32

		switch enemy.Class {
		case assets.EnemyClassAirplane:
			archetype.NewEnemySpawn(world, pos, func(w donburi.World) {
				archetype.NewEnemyAirplane(
					w,
					enemy.Position,
					enemy.Rotation,
					enemy.Speed,
					enemy.Path,
				)
			})
		case assets.EnemyClassTank:
			archetype.NewEnemySpawn(world, pos, func(w donburi.World) {
				archetype.NewEnemyTank(
					w,
					enemy.Position,
					enemy.Rotation,
					enemy.Speed,
					enemy.Path,
				)
			})
		default:
			panic("unknown enemy class: " + enemy.Class)
		}
	}

	if g.world == nil {
		game := world.Entry(world.Create(component.Game))
		component.Game.SetValue(game, component.GameData{
			Score: 0,
			Settings: component.Settings{
				ScreenWidth:  g.screenWidth,
				ScreenHeight: g.screenHeight,
			},
		})

		// Spawn new players
		for _, p := range g.players {
			player := archetype.NewPlayer(world, p.PlayerNumber, p.Faction)
			archetype.NewPlayerAirplane(world, *component.Player.Get(player), p.Faction, 0)
		}
	} else {
		// Keep the same game data across levels
		gameData := component.MustFindGame(g.world)
		newGameData := world.Entry(world.Create(component.Game))
		component.Game.Set(newGameData, gameData)

		// Transfer existing players from the previous level
		query.NewQuery(filter.Contains(component.Player)).EachEntity(g.world, func(entry *donburi.Entry) {
			player := component.Player.Get(entry)
			// In case the level ends while the player's respawning
			player.Respawning = false

			archetype.NewPlayerFromPlayerData(world, *player)
			if player.Lives > 0 {
				archetype.NewPlayerAirplane(world, *player, player.PlayerFaction, player.EvolutionLevel())
			}
		})
	}

	world.Create(component.Debug)

	system.SetupEvents(world)

	return world
}

func (g *Game) restart() {
	g.world = nil
	g.level = 0
	g.loadLevel()
}

func (g *Game) Update() {
	gameData := component.MustFindGame(g.world)
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		gameData.Paused = !gameData.Paused
	}

	if gameData.Paused {
		return
	}

	for _, s := range g.systems {
		s.Update(g.world)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	for _, s := range g.drawables {
		s.Draw(g.world, screen)
	}
}

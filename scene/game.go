package scene

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
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
	level     int
	world     donburi.World
	systems   []System
	drawables []Drawable

	screenWidth  int
	screenHeight int
}

func NewGame(screenWidth int, screenHeight int) *Game {
	g := &Game{
		level:        0,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}

	render := system.NewRenderer()
	debug := system.NewDebug(g.loadLevel)

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
		system.NewRespawn(g.restart),
		system.NewInvulnerable(),
		system.NewCamera(),
		system.NewObserver(),
		system.NewShooter(),
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

	level := world.Entry(world.Create(component.Level))
	component.GetLevel(level).ProgressionTimer = engine.NewTimer(time.Second * 3)

	archetype.NewCamera(world, engine.Vector{
		X: 0,
		Y: float64(levelAsset.Background.Bounds().Dy() - g.screenHeight),
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
			archetype.NewEnemyAirplane(
				world,
				enemy.Position,
				enemy.Rotation,
				enemy.Speed,
				enemy.Path,
			)
		case assets.EnemyClassTank:
			archetype.NewEnemyTank(
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
		game := world.Entry(world.Create(component.Game))
		donburi.SetValue(game, component.Game, component.GameData{
			Score: 0,
			Settings: component.Settings{
				ScreenWidth:  g.screenWidth,
				ScreenHeight: g.screenHeight,
			},
		})

		// Spawn new players
		for i := 1; i <= players; i++ {
			player := archetype.NewPlayer(world, i)
			archetype.NewPlayerAirplane(world, *component.GetPlayer(player))
		}
	} else {
		// Keep the same game data across levels
		gameData := component.MustFindGame(g.world)
		newGameData := world.Entry(world.Create(component.Game))
		donburi.Set(newGameData, component.Game, gameData)

		// Transfer existing players from the previous level
		query.NewQuery(filter.Contains(component.Player)).EachEntity(g.world, func(entry *donburi.Entry) {
			player := component.GetPlayer(entry)
			// In case the level ends while the player's respawning
			player.Respawning = false

			archetype.NewPlayerFromPlayerData(world, *player)
			if player.Lives > 0 {
				archetype.NewPlayerAirplane(world, *player)
			}
		})
	}

	world.Create(component.Debug)

	return world
}

func (g *Game) restart() {
	// TODO: Definitely a hack. Needed because GameData is cached in systems.
	// Consider a different approach to GameData, perhaps not as a component?
	component.MustFindGame(g.world).Score = 0
	component.MustFindGame(g.world).GameOver = false

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

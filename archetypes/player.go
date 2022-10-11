package archetypes

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

type playerInputs struct {
	Up    ebiten.Key
	Right ebiten.Key
	Down  ebiten.Key
	Left  ebiten.Key
	Shoot ebiten.Key
}

type playerSettings struct {
	Image  func() *ebiten.Image
	Inputs playerInputs
}

var players = map[int]playerSettings{
	1: {
		Image: func() *ebiten.Image { return assets.AirplaneYellowSmall },
		Inputs: playerInputs{
			Up:    ebiten.KeyW,
			Right: ebiten.KeyD,
			Down:  ebiten.KeyS,
			Left:  ebiten.KeyA,
			Shoot: ebiten.KeySpace,
		},
	},
	2: {
		Image: func() *ebiten.Image { return assets.AirplaneGreenSmall },
		Inputs: playerInputs{
			Up:    ebiten.KeyUp,
			Right: ebiten.KeyRight,
			Down:  ebiten.KeyDown,
			Left:  ebiten.KeyLeft,
			Shoot: ebiten.KeyK,
		},
	},
}

func playerSpawn(w donburi.World, playerNumber int) engine.Vector {
	game := component.MustFindGame(w)
	cameraPos := component.GetTransform(MustFindCamera(w)).LocalPosition

	switch playerNumber {
	case 1:
		return engine.Vector{
			X: float64(game.Settings.ScreenWidth) * 0.25,
			Y: cameraPos.Y + float64(game.Settings.ScreenHeight)*0.9,
		}
	case 2:
		return engine.Vector{
			X: float64(game.Settings.ScreenWidth) * 0.75,
			Y: cameraPos.Y + float64(game.Settings.ScreenHeight)*0.9,
		}
	default:
		panic(fmt.Sprintf("unknown player number: %v", playerNumber))
	}
}

func NewPlayer(w donburi.World, playerNumber int) *donburi.Entry {
	_, ok := players[playerNumber]
	if !ok {
		panic(fmt.Sprintf("unknown player number: %v", playerNumber))
	}

	player := component.PlayerData{
		PlayerNumber: playerNumber,
		Lives:        3,
		RespawnTimer: engine.NewTimer(time.Second * 3),
		WeaponLevel:  component.WeaponLevelSingle,
	}

	// TODO It looks like a constructor would fit here
	player.ShootTimer = engine.NewTimer(player.WeaponCooldown())

	return NewPlayerFromPlayerData(w, player)
}

func NewPlayerFromPlayerData(w donburi.World, playerData component.PlayerData) *donburi.Entry {
	player := w.Entry(w.Create(component.Player))
	donburi.SetValue(player, component.Player, playerData)
	return player
}

func NewPlayerAirplane(w donburi.World, player component.PlayerData) {
	settings, ok := players[player.PlayerNumber]
	if !ok {
		panic(fmt.Sprintf("unknown player number: %v", player.PlayerNumber))
	}

	airplane := w.Entry(
		w.Create(
			component.PlayerAirplane,
			component.Transform,
			component.Velocity,
			component.Sprite,
			component.Input,
			component.Bounds,
			component.Collider,
		),
	)

	donburi.SetValue(airplane, component.PlayerAirplane, component.PlayerAirplaneData{
		PlayerNumber:      player.PlayerNumber,
		InvulnerableTimer: engine.NewTimer(time.Second * 3),
		Invulnerable:      true,
	})

	originalRotation := -90.0

	pos := playerSpawn(w, player.PlayerNumber)
	donburi.SetValue(airplane, component.Transform, component.TransformData{
		LocalPosition: pos,
		LocalRotation: originalRotation,
	})

	image := settings.Image()

	donburi.SetValue(airplane, component.Sprite, component.SpriteData{
		Image:            image,
		Layer:            component.SpriteLayerAirUnits,
		Pivot:            component.SpritePivotCenter,
		OriginalRotation: originalRotation,
	})

	width, height := image.Size()
	donburi.SetValue(airplane, component.Collider, component.ColliderData{
		Width:  float64(width),
		Height: float64(height),
		Layer:  component.CollisionLayerPlayers,
	})

	inputs := settings.Inputs
	donburi.SetValue(airplane, component.Input, component.InputData{
		MoveUpKey:    inputs.Up,
		MoveRightKey: inputs.Right,
		MoveDownKey:  inputs.Down,
		MoveLeftKey:  inputs.Left,
		MoveSpeed:    3.5,
		ShootKey:     inputs.Shoot,
	})

	NewShadow(w, airplane)
}

func MustFindPlayerByNumber(w donburi.World, playerNumber int) *component.PlayerData {
	var foundPlayer *component.PlayerData
	query.NewQuery(filter.Contains(component.Player)).EachEntity(w, func(e *donburi.Entry) {
		player := component.GetPlayer(e)
		if player.PlayerNumber == playerNumber {
			foundPlayer = player
		}
	})

	if foundPlayer == nil {
		panic(fmt.Sprintf("player not found: %v", playerNumber))
	}

	return foundPlayer
}

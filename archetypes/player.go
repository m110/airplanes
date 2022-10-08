package archetypes

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

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

func playerSpawn(w donburi.World, playerNumber int) component.PositionData {
	settings := component.MustFindSettings(w)
	cameraPos := component.GetPosition(MustFindCamera(w))

	switch playerNumber {
	case 1:
		return component.PositionData{
			X: float64(settings.ScreenWidth) * 0.25,
			Y: cameraPos.Y + float64(settings.ScreenHeight)*0.9,
		}
	case 2:
		return component.PositionData{
			X: float64(settings.ScreenWidth) * 0.75,
			Y: cameraPos.Y + float64(settings.ScreenHeight)*0.9,
		}
	default:
		panic(fmt.Sprintf("unknown player number: %v", playerNumber))
	}
}

func NewPlayer(w donburi.World, playerNumber int) {
	_, ok := players[playerNumber]
	if !ok {
		panic(fmt.Sprintf("unknown player number: %v", playerNumber))
	}

	player := w.Entry(w.Create(component.Player))
	donburi.SetValue(player, component.Player, component.PlayerData{
		PlayerNumber: playerNumber,
		Lives:        3,
		RespawnTimer: engine.NewTimer(time.Second * 3),
	})
}

func NewPlayerAirplane(w donburi.World, playerNumber int) {
	settings, ok := players[playerNumber]
	if !ok {
		panic(fmt.Sprintf("unknown player number: %v", playerNumber))
	}

	airplane := w.Entry(
		w.Create(
			component.PlayerAirplane,
			component.Position,
			component.Velocity,
			component.Sprite,
			component.Input,
			component.Bounds,
			component.Collider,
		),
	)

	donburi.SetValue(airplane, component.PlayerAirplane, component.PlayerAirplaneData{
		PlayerNumber:      playerNumber,
		InvulnerableTimer: engine.NewTimer(time.Second * 3),
		Invulnerable:      true,
	})

	donburi.SetValue(airplane, component.Position, playerSpawn(w, playerNumber))

	image := settings.Image()

	donburi.SetValue(airplane, component.Sprite, component.SpriteData{
		Image: image,
		Layer: component.SpriteLayerUnits,
		Pivot: component.SpritePivotCenter,
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
		ShootTimer:   engine.NewTimer(time.Millisecond * 400),
	})
}

package archetype

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

type PlayerInputs struct {
	Up    ebiten.Key
	Right ebiten.Key
	Down  ebiten.Key
	Left  ebiten.Key
	Shoot ebiten.Key
}

type PlayerSettings struct {
	Inputs PlayerInputs
}

func AirplaneImageByFaction(faction component.PlayerFaction, level int) *ebiten.Image {
	switch faction {
	case component.PlayerFactionBlue:
		return assets.AirplanesBlue[level]
	case component.PlayerFactionRed:
		return assets.AirplanesRed[level]
	case component.PlayerFactionGreen:
		return assets.AirplanesGreen[level]
	case component.PlayerFactionYellow:
		return assets.AirplanesYellow[level]
	default:
		panic(fmt.Sprintf("unknown player faction: %v", faction))
	}
}

var Players = map[int]PlayerSettings{
	1: {
		Inputs: PlayerInputs{
			Up:    ebiten.KeyW,
			Right: ebiten.KeyD,
			Down:  ebiten.KeyS,
			Left:  ebiten.KeyA,
			Shoot: ebiten.KeySpace,
		},
	},
	2: {
		Inputs: PlayerInputs{
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

func NewPlayer(w donburi.World, playerNumber int, faction component.PlayerFaction) *donburi.Entry {
	_, ok := Players[playerNumber]
	if !ok {
		panic(fmt.Sprintf("unknown player number: %v", playerNumber))
	}

	player := component.PlayerData{
		PlayerNumber:  playerNumber,
		PlayerFaction: faction,
		Lives:         3,
		RespawnTimer:  engine.NewTimer(time.Second * 3),
		WeaponLevel:   component.WeaponLevelSingle,
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

func NewPlayerAirplane(w donburi.World, player component.PlayerData, faction component.PlayerFaction, evolutionLevel int) {
	settings, ok := Players[player.PlayerNumber]
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
			component.Evolution,
		),
	)

	shield := w.Entry(
		w.Create(
			component.Transform,
			component.Sprite,
		),
	)

	donburi.SetValue(shield, component.Sprite, component.SpriteData{
		Image:            assets.AirplaneShield,
		Layer:            component.SpriteLayerIndicators,
		Pivot:            component.SpritePivotCenter,
		OriginalRotation: -90.0,
	})

	donburi.SetValue(airplane, component.PlayerAirplane, component.PlayerAirplaneData{
		PlayerNumber:          player.PlayerNumber,
		Faction:               faction,
		InvulnerableTimer:     engine.NewTimer(time.Second * 3),
		InvulnerableIndicator: component.GetSprite(shield),
	})

	component.GetPlayerAirplane(airplane).StartInvulnerability()

	originalRotation := -90.0

	pos := playerSpawn(w, player.PlayerNumber)
	donburi.SetValue(airplane, component.Transform, component.TransformData{
		LocalPosition: pos,
		LocalRotation: originalRotation,
	})

	component.GetTransform(airplane).AppendChild(airplane, shield, false)

	image := AirplaneImageByFaction(faction, evolutionLevel)

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

	donburi.SetValue(airplane, component.Evolution, component.EvolutionData{
		Level:       evolutionLevel,
		GrowTimer:   engine.NewTimer(time.Second * 1),
		ShrinkTimer: engine.NewTimer(time.Second * 1),
	})

	NewShadow(w, airplane)

	evolution := w.Entry(
		w.Create(
			component.Transform,
			component.Sprite,
			EvolutionTag,
		),
	)
	component.GetTransform(airplane).AppendChild(airplane, evolution, false)
	donburi.SetValue(evolution, component.Sprite, component.SpriteData{
		Layer:            component.SpriteLayerAirUnits,
		Pivot:            component.SpritePivotCenter,
		OriginalRotation: originalRotation,
		Hidden:           true,
	})
}

var EvolutionTag = donburi.NewTag()

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

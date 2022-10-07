package archetypes

import (
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

func NewPlayerOne(w donburi.World, position component.PositionData) *donburi.Entry {
	inputs := playerInputs{
		Up:    ebiten.KeyW,
		Right: ebiten.KeyD,
		Down:  ebiten.KeyS,
		Left:  ebiten.KeyA,
		Shoot: ebiten.KeySpace,
	}

	return newPlayer(w, position, assets.ShipYellowSmall, inputs)
}

func NewPlayerTwo(w donburi.World, position component.PositionData) *donburi.Entry {
	inputs := playerInputs{
		Up:    ebiten.KeyUp,
		Right: ebiten.KeyRight,
		Down:  ebiten.KeyDown,
		Left:  ebiten.KeyLeft,
		Shoot: ebiten.KeyK,
	}

	return newPlayer(w, position, assets.ShipGreenSmall, inputs)
}

func newPlayer(
	w donburi.World,
	position component.PositionData,
	image *ebiten.Image,
	inputs playerInputs,
) *donburi.Entry {
	player := w.Entry(
		w.Create(
			component.PlayerTag,
			component.Position,
			component.Velocity,
			component.Sprite,
			component.Input,
			component.Bounds,
			component.Collider,
		),
	)

	donburi.SetValue(player, component.Position, position)

	donburi.SetValue(player, component.Sprite, component.SpriteData{
		Image: image,
		Layer: component.SpriteLayerUnits,
	})

	width, height := image.Size()
	donburi.SetValue(player, component.Collider, component.ColliderData{
		Width:  width,
		Height: height,
		Layer:  component.CollisionLayerPlayers,
	})

	donburi.SetValue(player, component.Input, component.InputData{
		MoveUpKey:    inputs.Up,
		MoveRightKey: inputs.Right,
		MoveDownKey:  inputs.Down,
		MoveLeftKey:  inputs.Left,
		MoveSpeed:    3.5,
		ShootKey:     inputs.Shoot,
		ShootTimer:   engine.NewTimer(time.Millisecond * 400),
	})

	return player
}

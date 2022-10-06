package archetypes

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

func NewPlayerOne(w donburi.World, position component.PositionData) *donburi.Entry {
	player := newPlayer(w, position)

	donburi.SetValue(player, component.Sprite, component.SpriteData{
		Image: assets.ShipYellowSmall,
		Layer: component.SpriteLayerUnits,
	})
	donburi.SetValue(player, component.Input, component.InputData{
		MoveUpKey:    ebiten.KeyW,
		MoveRightKey: ebiten.KeyD,
		MoveDownKey:  ebiten.KeyS,
		MoveLeftKey:  ebiten.KeyA,
		MoveSpeed:    3.5,
		ShootKey:     ebiten.KeySpace,
		ShootTimer:   engine.NewTimer(time.Millisecond * 400),
	})

	return player
}

func NewPlayerTwo(w donburi.World, position component.PositionData) *donburi.Entry {
	player := newPlayer(w, position)

	donburi.SetValue(player, component.Sprite, component.SpriteData{
		Image: assets.ShipGreenSmall,
		Layer: component.SpriteLayerUnits,
	})
	donburi.SetValue(player, component.Input, component.InputData{
		MoveUpKey:    ebiten.KeyUp,
		MoveRightKey: ebiten.KeyRight,
		MoveDownKey:  ebiten.KeyDown,
		MoveLeftKey:  ebiten.KeyLeft,
		MoveSpeed:    3.5,
		ShootKey:     ebiten.KeyK,
		ShootTimer:   engine.NewTimer(time.Millisecond * 400),
	})

	return player
}

func newPlayer(w donburi.World, position component.PositionData) *donburi.Entry {
	player := w.Entry(
		w.Create(
			component.PlayerTag,
			component.Position,
			component.Velocity,
			component.Sprite,
			component.Input,
			component.Bounds,
		),
	)
	donburi.SetValue(player, component.Position, position)

	return player
}

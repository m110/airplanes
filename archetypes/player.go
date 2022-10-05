package archetypes

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

func NewPlayerOne(w donburi.World) *donburi.Entry {
	player := newPlayer(w)

	donburi.SetValue(player, component.Position, component.PositionData{X: 100, Y: 500})
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
		ShootTimer:   engine.NewTimer(time.Millisecond * 300),
	})

	return player
}

func NewPlayerTwo(w donburi.World) *donburi.Entry {
	player := newPlayer(w)

	donburi.SetValue(player, component.Position, component.PositionData{X: 500, Y: 500})
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
		ShootTimer:   engine.NewTimer(time.Millisecond * 300),
	})

	return player
}

func newPlayer(w donburi.World) *donburi.Entry {
	entity := w.Create(
		component.PlayerTag,
		component.Position,
		component.Velocity,
		component.Sprite,
		component.Input,
	)
	player := w.Entry(entity)

	return player
}

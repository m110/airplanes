package archetypes

import (
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

func NewPlayerOne(w donburi.World, position component.PositionData) {
	inputs := playerInputs{
		Up:    ebiten.KeyW,
		Right: ebiten.KeyD,
		Down:  ebiten.KeyS,
		Left:  ebiten.KeyA,
		Shoot: ebiten.KeySpace,
	}

	newPlayerShip(w, position, assets.ShipYellowSmall, inputs, 1)
	newPlayerIfNotExists(w, 1)
}

func NewPlayerTwo(w donburi.World, position component.PositionData) {
	inputs := playerInputs{
		Up:    ebiten.KeyUp,
		Right: ebiten.KeyRight,
		Down:  ebiten.KeyDown,
		Left:  ebiten.KeyLeft,
		Shoot: ebiten.KeyK,
	}

	newPlayerShip(w, position, assets.ShipGreenSmall, inputs, 2)
	newPlayerIfNotExists(w, 2)
}

func newPlayerIfNotExists(
	w donburi.World,
	number int,
) {
	exists := false
	query.NewQuery(filter.Contains(component.Player)).EachEntity(w, func(entry *donburi.Entry) {
		if component.GetPlayer(entry).PlayerNumber == number {
			exists = true
		}
	})

	if exists {
		return
	}

	player := w.Entry(w.Create(component.Player))
	donburi.SetValue(player, component.Player, component.PlayerData{
		PlayerNumber: number,
		Lives:        3,
		RespawnTimer: engine.NewTimer(time.Second * 3),
	})
}

func newPlayerShip(
	w donburi.World,
	position component.PositionData,
	image *ebiten.Image,
	inputs playerInputs,
	number int,
) {
	ship := w.Entry(
		w.Create(
			component.PlayerNumber,
			component.Position,
			component.Velocity,
			component.Sprite,
			component.Input,
			component.Bounds,
			component.Collider,
		),
	)

	donburi.SetValue(ship, component.PlayerNumber, component.PlayerNumberData{
		Number: number,
	})

	donburi.SetValue(ship, component.Position, position)

	donburi.SetValue(ship, component.Sprite, component.SpriteData{
		Image: image,
		Layer: component.SpriteLayerUnits,
		Pivot: component.SpritePivotCenter,
	})

	width, height := image.Size()
	donburi.SetValue(ship, component.Collider, component.ColliderData{
		Width:  float64(width),
		Height: float64(height),
		Layer:  component.CollisionLayerPlayers,
	})

	donburi.SetValue(ship, component.Input, component.InputData{
		MoveUpKey:    inputs.Up,
		MoveRightKey: inputs.Right,
		MoveDownKey:  inputs.Down,
		MoveLeftKey:  inputs.Left,
		MoveSpeed:    3.5,
		ShootKey:     inputs.Shoot,
		ShootTimer:   engine.NewTimer(time.Millisecond * 400),
	})
}

package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/m110/airplanes/archetypes"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/component"
)

type Controls struct {
	query *query.Query
}

func NewControls() *Controls {
	return &Controls{
		query: query.NewQuery(
			filter.Contains(
				component.Position,
				component.Input,
				component.Velocity,
				component.Sprite,
			),
		),
	}
}

func (i *Controls) Update(w donburi.World) {
	i.query.EachEntity(w, func(entry *donburi.Entry) {
		input := component.GetInput(entry)
		velocity := component.GetVelocity(entry)

		velocity.X = 0
		velocity.Y = 0

		if ebiten.IsKeyPressed(input.MoveUpKey) {
			velocity.Y = -input.MoveSpeed
		} else if ebiten.IsKeyPressed(input.MoveDownKey) {
			velocity.Y = input.MoveSpeed
		}

		if ebiten.IsKeyPressed(input.MoveRightKey) {
			velocity.X = input.MoveSpeed
		}
		if ebiten.IsKeyPressed(input.MoveLeftKey) {
			velocity.X = -input.MoveSpeed
		}

		input.ShootTimer.Update()
		if ebiten.IsKeyPressed(input.ShootKey) && input.ShootTimer.IsReady() {
			bullet := archetypes.NewBullet(w)

			bulletPosition := component.GetPosition(bullet)
			bulletSprite := component.GetSprite(bullet)

			position := component.GetPosition(entry)
			sprite := component.GetSprite(entry)

			bulletPosition.X = position.X + float64(sprite.Image.Bounds().Dx())/2.0 - float64(bulletSprite.Image.Bounds().Dx())/2.0
			bulletPosition.Y = position.Y - float64(bulletSprite.Image.Bounds().Dy())

			input.ShootTimer.Reset()
		}
	})
}

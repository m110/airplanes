package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

type Bounds struct {
	query                     *query.Query
	screenWidth, screenHeight int
}

func NewBounds(screenWidth, screenHeight int) *Bounds {
	return &Bounds{
		query: query.NewQuery(filter.Contains(
			component.PlayerTag,
			component.Position,
			component.Sprite,
		)),
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
}

func (b *Bounds) Update(w donburi.World) {
	b.query.EachEntity(w, func(entry *donburi.Entry) {
		position := component.GetPosition(entry)
		sprite := component.GetSprite(entry)

		position.X = engine.Clamp(position.X, 0, float64(b.screenWidth-sprite.Image.Bounds().Dx()))
		position.Y = engine.Clamp(position.Y, 0, float64(b.screenHeight-sprite.Image.Bounds().Dy()))
	})
}

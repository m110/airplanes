package system

import (
	"github.com/m110/airplanes/archetypes"
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
			component.Bounds,
		)),
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
}

func (b *Bounds) Update(w donburi.World) {
	camera := archetypes.MustFindCamera(w)
	cameraPos := component.GetPosition(camera)

	b.query.EachEntity(w, func(entry *donburi.Entry) {
		bounds := component.GetBounds(entry)
		if bounds.Disabled {
			return
		}

		position := component.GetPosition(entry)
		sprite := component.GetSprite(entry)

		minX := cameraPos.X
		maxX := cameraPos.X + float64(b.screenWidth-sprite.Image.Bounds().Dx())

		minY := cameraPos.Y
		maxY := cameraPos.Y + float64(b.screenHeight-sprite.Image.Bounds().Dy())

		position.X = engine.Clamp(position.X, minX, maxX)
		position.Y = engine.Clamp(position.Y, minY, maxY)
	})
}

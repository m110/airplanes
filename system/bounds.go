package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/archetypes"
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
			component.PlayerNumber,
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

		w, h := sprite.Image.Size()
		width, height := float64(w), float64(h)

		var minX, maxX, minY, maxY float64

		switch sprite.Pivot {
		case component.SpritePivotTopLeft:
			minX = cameraPos.X
			maxX = cameraPos.X + float64(b.screenWidth) - width

			minY = cameraPos.Y
			maxY = cameraPos.Y + float64(b.screenHeight) - height
		case component.SpritePivotCenter:
			minX = cameraPos.X + width/2
			maxX = cameraPos.X + float64(b.screenWidth) - width/2

			minY = cameraPos.Y + height/2
			maxY = cameraPos.Y + float64(b.screenHeight) - height/2
		}

		position.X = engine.Clamp(position.X, minX, maxX)
		position.Y = engine.Clamp(position.Y, minY, maxY)
	})
}

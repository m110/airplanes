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
	query    *query.Query
	settings *component.SettingsData
}

func NewBounds() *Bounds {
	return &Bounds{
		query: query.NewQuery(filter.Contains(
			component.PlayerAirplane,
			component.Transform,
			component.Sprite,
			component.Bounds,
		)),
	}
}

func (b *Bounds) Update(w donburi.World) {
	if b.settings == nil {
		b.settings = component.MustFindSettings(w)
		if b.settings == nil {
			return
		}
	}

	camera := archetypes.MustFindCamera(w)
	cameraPos := component.GetTransform(camera).LocalPosition

	b.query.EachEntity(w, func(entry *donburi.Entry) {
		bounds := component.GetBounds(entry)
		if bounds.Disabled {
			return
		}

		transform := component.GetTransform(entry)
		sprite := component.GetSprite(entry)

		w, h := sprite.Image.Size()
		width, height := float64(w), float64(h)

		var minX, maxX, minY, maxY float64

		switch sprite.Pivot {
		case component.SpritePivotTopLeft:
			minX = cameraPos.X
			maxX = cameraPos.X + float64(b.settings.ScreenWidth) - width

			minY = cameraPos.Y
			maxY = cameraPos.Y + float64(b.settings.ScreenHeight) - height
		case component.SpritePivotCenter:
			minX = cameraPos.X + width/2
			maxX = cameraPos.X + float64(b.settings.ScreenWidth) - width/2

			minY = cameraPos.Y + height/2
			maxY = cameraPos.Y + float64(b.settings.ScreenHeight) - height/2
		}

		transform.LocalPosition.X = engine.Clamp(transform.LocalPosition.X, minX, maxX)
		transform.LocalPosition.Y = engine.Clamp(transform.LocalPosition.Y, minY, maxY)
	})
}

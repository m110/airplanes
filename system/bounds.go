package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

type Bounds struct {
	query *donburi.Query
	game  *component.GameData
}

func NewBounds() *Bounds {
	return &Bounds{
		query: donburi.NewQuery(filter.Contains(
			component.PlayerAirplane,
			transform.Transform,
			component.Sprite,
			component.Bounds,
		)),
	}
}

func (b *Bounds) Update(w donburi.World) {
	if b.game == nil {
		b.game = component.MustFindGame(w)
		if b.game == nil {
			return
		}
	}

	camera := archetype.MustFindCamera(w)
	cameraPos := transform.Transform.Get(camera).LocalPosition

	b.query.Each(w, func(entry *donburi.Entry) {
		bounds := component.Bounds.Get(entry)
		if bounds.Disabled {
			return
		}

		t := transform.Transform.Get(entry)
		sprite := component.Sprite.Get(entry)

		w, h := sprite.Image.Size()
		width, height := float64(w), float64(h)

		var minX, maxX, minY, maxY float64

		switch sprite.Pivot {
		case component.SpritePivotTopLeft:
			minX = cameraPos.X
			maxX = cameraPos.X + float64(b.game.Settings.ScreenWidth) - width

			minY = cameraPos.Y
			maxY = cameraPos.Y + float64(b.game.Settings.ScreenHeight) - height
		case component.SpritePivotCenter:
			minX = cameraPos.X + width/2
			maxX = cameraPos.X + float64(b.game.Settings.ScreenWidth) - width/2

			minY = cameraPos.Y + height/2
			maxY = cameraPos.Y + float64(b.game.Settings.ScreenHeight) - height/2
		}

		t.LocalPosition.X = engine.Clamp(t.LocalPosition.X, minX, maxX)
		t.LocalPosition.Y = engine.Clamp(t.LocalPosition.Y, minY, maxY)
	})
}

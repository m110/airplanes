package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/archetypes"
	"github.com/m110/airplanes/component"
)

type Despawn struct {
	query *query.Query
	game  *component.GameData
}

func NewDespawn() *Despawn {
	return &Despawn{
		query: query.NewQuery(filter.Contains(component.Despawnable)),
	}
}

func (d *Despawn) Update(w donburi.World) {
	if d.game == nil {
		d.game = component.MustFindGame(w)
		if d.game == nil {
			return
		}
	}

	cameraPos := component.GetTransform(archetypes.MustFindCamera(w)).LocalPosition

	d.query.EachEntity(w, func(entry *donburi.Entry) {
		position := component.GetTransform(entry).LocalPosition
		sprite := component.GetSprite(entry)
		despawnable := component.GetDespawnable(entry)

		maxX := position.X + float64(sprite.Image.Bounds().Dx())
		maxY := position.Y + float64(sprite.Image.Bounds().Dy())

		cameraMaxY := cameraPos.Y + float64(d.game.Settings.ScreenHeight)
		cameraMaxX := cameraPos.X + float64(d.game.Settings.ScreenWidth)

		if !despawnable.Spawned {
			if position.Y > cameraPos.Y && maxY < cameraMaxY &&
				position.X > cameraPos.X && maxX < cameraMaxX {
				despawnable.Spawned = true
			}

			return
		}

		if maxY < cameraPos.Y || position.Y > cameraMaxY ||
			maxX < cameraPos.X || position.X > cameraMaxX {
			Destroy(w, entry)
		}
	})
}

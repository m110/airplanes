package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/hierarchy"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/archetype"
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

	cameraPos := transform.Transform.Get(archetype.MustFindCamera(w)).LocalPosition

	d.query.EachEntity(w, func(entry *donburi.Entry) {
		position := transform.Transform.Get(entry).LocalPosition
		sprite := component.Sprite.Get(entry)
		despawnable := component.Despawnable.Get(entry)

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
			hierarchy.RemoveRecursive(entry)
		}
	})
}

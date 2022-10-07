package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/archetypes"
	"github.com/m110/airplanes/component"
)

type Despawn struct {
	query        *query.Query
	screenWidth  int
	screenHeight int
}

func NewDespawn(screenWidth int, screenHeight int) *Despawn {
	return &Despawn{
		query: query.NewQuery(filter.Contains(component.Despawnable)),
		// TODO Move these out to "settings"?
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
}

func (d *Despawn) Update(w donburi.World) {
	cameraPos := component.GetPosition(archetypes.MustFindCamera(w))

	d.query.EachEntity(w, func(entry *donburi.Entry) {
		position := component.GetPosition(entry)
		sprite := component.GetSprite(entry)
		despawnable := component.GetDespawnable(entry)

		maxX := position.X + float64(sprite.Image.Bounds().Dx())
		maxY := position.Y + float64(sprite.Image.Bounds().Dy())

		cameraMaxY := cameraPos.Y + float64(d.screenHeight)
		cameraMaxX := cameraPos.X + float64(d.screenWidth)

		if !despawnable.Spawned {
			if position.Y > cameraPos.Y && maxY < cameraMaxY &&
				position.X > cameraPos.X && maxX < cameraMaxX {
				despawnable.Spawned = true
			}

			return
		}

		if maxY < cameraPos.Y || position.Y > cameraMaxY ||
			maxX < cameraPos.X || position.X > cameraMaxX {
			w.Remove(entry.Entity())
		}
	})
}

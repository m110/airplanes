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
		despawnable := component.GetDespawnable(entry)
		if !despawnable.Spawned {
			return
		}

		position := component.GetPosition(entry)
		sprite := component.GetSprite(entry)

		if position.Y+float64(sprite.Image.Bounds().Dy()) < cameraPos.Y ||
			position.Y > cameraPos.Y+float64(d.screenHeight) ||
			position.X+float64(sprite.Image.Bounds().Dx()) < cameraPos.X ||
			position.X > cameraPos.X+float64(d.screenWidth) {
			w.Remove(entry.Entity())
		}
	})
}

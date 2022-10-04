package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/component"
)

type Render struct {
	query *query.Query
}

func NewRenderer() *Render {
	return &Render{
		query: query.NewQuery(
			filter.Contains(component.Position, component.Sprite),
		),
	}
}

func (r *Render) Draw(w donburi.World, screen *ebiten.Image) {
	r.query.EachEntity(w, func(entry *donburi.Entry) {
		position := component.GetPosition(entry)
		sprite := component.GetSprite(entry)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(position.X, position.Y)
		screen.DrawImage(sprite.Image, op)
	})
}

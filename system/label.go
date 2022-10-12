package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
	"golang.org/x/image/colornames"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
)

type Label struct {
	query *query.Query
}

func NewLabel() *Label {
	return &Label{
		query: query.NewQuery(
			filter.Contains(component.Transform, component.Label),
		),
	}
}

func (l *Label) Draw(w donburi.World, screen *ebiten.Image) {
	l.query.EachEntity(w, func(entry *donburi.Entry) {
		label := component.GetLabel(entry)
		if label.Hidden {
			return
		}

		transform := component.GetTransform(entry)

		// TODO Rotation, Scale, customizable font and color
		text.Draw(
			screen,
			label.Text,
			assets.SmallFont,
			int(transform.WorldPosition().X),
			int(transform.WorldPosition().Y),
			colornames.White,
		)
	})
}

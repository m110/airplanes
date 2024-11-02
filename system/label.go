package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"golang.org/x/image/colornames"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
)

type Label struct {
	query *donburi.Query
}

func NewLabel() *Label {
	return &Label{
		query: donburi.NewQuery(
			filter.Contains(transform.Transform, component.Label),
		),
	}
}

func (l *Label) Draw(w donburi.World, screen *ebiten.Image) {
	l.query.Each(w, func(entry *donburi.Entry) {
		label := component.Label.Get(entry)
		if label.Hidden {
			return
		}

		pos := transform.WorldPosition(entry)

		// TODO Rotation, Scale, customizable font and color
		text.Draw(
			screen,
			label.Text,
			assets.SmallFont,
			int(pos.X),
			int(pos.Y),
			colornames.White,
		)
	})
}

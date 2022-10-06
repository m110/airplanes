package system

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/lo"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/archetypes"
	"github.com/m110/airplanes/component"
)

type Render struct {
	query     *query.Query
	offscreen *ebiten.Image
}

func NewRenderer() *Render {
	return &Render{
		query: query.NewQuery(
			filter.Contains(component.Position, component.Sprite),
		),
		// TODO figure out the proper size
		offscreen: ebiten.NewImage(1000, 2000),
	}
}

func (r *Render) Draw(w donburi.World, screen *ebiten.Image) {
	camera := archetypes.MustFindCamera(w)
	cam := component.GetCamera(camera)

	if !cam.Moving {
		cam.MoveTimer.Update()
		if cam.MoveTimer.IsReady() {
			cam.Moving = true
			component.GetVelocity(camera).Y = -0.5
		}
	}

	cameraPos := component.GetPosition(camera)

	r.offscreen.Clear()

	var entries []*donburi.Entry
	r.query.EachEntity(w, func(entry *donburi.Entry) {
		entries = append(entries, entry)
	})

	byLayer := lo.GroupBy(entries, func(entry *donburi.Entry) int {
		return component.GetSprite(entry).Layer
	})
	layers := lo.Keys(byLayer)
	sort.Ints(layers)

	for _, layer := range layers {
		for _, entry := range byLayer[layer] {
			position := component.GetPosition(entry)
			sprite := component.GetSprite(entry)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(position.X, position.Y)
			r.offscreen.DrawImage(sprite.Image, op)
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-cameraPos.X, -cameraPos.Y)
	screen.DrawImage(r.offscreen, op)
}

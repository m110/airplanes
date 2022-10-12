package system

import (
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/lo"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
)

type Render struct {
	query     *query.Query
	offscreen *ebiten.Image
	debug     *component.DebugData
}

func NewRenderer() *Render {
	return &Render{
		query: query.NewQuery(
			filter.Contains(component.Transform, component.Sprite),
		),
		// TODO figure out the proper size
		offscreen: ebiten.NewImage(3000, 3000),
	}
}

func (r *Render) Update(w donburi.World) {
	if r.debug == nil {
		debug, ok := query.NewQuery(filter.Contains(component.Debug)).FirstEntity(w)
		if !ok {
			return
		}

		r.debug = component.GetDebug(debug)
	}
}

func (r *Render) Draw(w donburi.World, screen *ebiten.Image) {
	camera := archetype.MustFindCamera(w)
	cameraPos := component.GetTransform(camera).LocalPosition

	r.offscreen.Clear()

	var entries []*donburi.Entry
	r.query.EachEntity(w, func(entry *donburi.Entry) {
		entries = append(entries, entry)
	})

	byLayer := lo.GroupBy(entries, func(entry *donburi.Entry) int {
		return int(component.GetSprite(entry).Layer)
	})
	layers := lo.Keys(byLayer)
	sort.Ints(layers)

	for _, layer := range layers {
		for _, entry := range byLayer[layer] {
			sprite := component.GetSprite(entry)

			if sprite.Hidden {
				continue
			}

			transform := component.GetTransform(entry)

			w, h := sprite.Image.Size()
			halfW, halfH := float64(w)/2, float64(h)/2

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(-halfW, -halfH)
			op.GeoM.Rotate(float64(int(transform.WorldRotation()-sprite.OriginalRotation)%360) * 2 * math.Pi / 360)
			op.GeoM.Translate(halfW, halfH)

			position := transform.WorldPosition()

			x := position.X
			y := position.Y

			switch sprite.Pivot {
			case component.SpritePivotCenter:
				x -= halfW
				y -= halfH
			}

			if transform.LocalScale.X != 0 || transform.LocalScale.Y != 0 {
				op.GeoM.Translate(-halfW, -halfH)
				op.GeoM.Scale(transform.LocalScale.X, transform.LocalScale.Y)
				op.GeoM.Translate(halfW, halfH)
			}
			op.GeoM.Translate(x, y)
			r.offscreen.DrawImage(sprite.Image, op)
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-cameraPos.X, -cameraPos.Y)
	screen.DrawImage(r.offscreen, op)
}

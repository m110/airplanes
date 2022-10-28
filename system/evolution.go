package system

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
)

type Evolution struct {
	query *query.Query
}

func NewEvolution() *Evolution {
	return &Evolution{
		query: query.NewQuery(filter.Contains(component.Evolution)),
	}
}

func (s *Evolution) Update(w donburi.World) {
	s.query.EachEntity(w, func(entry *donburi.Entry) {
		evolution := component.GetEvolution(entry)
		if !evolution.Evolving {
			return
		}

		t := transform.GetTransform(entry)
		evolutionChild, ok := transform.FindChildWithComponent(entry, component.EvolutionTag)
		if !ok {
			panic("no evolution found")
		}

		evolutionChildTransform := transform.GetTransform(evolutionChild)
		evolutionChildTransform.LocalScale.Set(
			evolution.GrowTimer.PercentDone(),
			evolution.GrowTimer.PercentDone(),
		)

		// TODO Change just once
		newSprite := component.GetSprite(evolutionChild)
		newSprite.Image = archetype.AirplaneImageByFaction(component.GetPlayerAirplane(entry).Faction, evolution.Level)
		newSprite.Show()
		newSprite.ColorOverride = &component.ColorOverride{
			R: 1,
			G: 1,
			B: 1,
			A: 1,
		}
		component.GetSprite(entry).ColorOverride = &component.ColorOverride{
			R: 1,
			G: 1,
			B: 1,
			A: 1,
		}

		shadow, _ := transform.FindChildWithComponent(entry, component.ShadowTag)
		shadowSprite := component.GetSprite(shadow)

		shadowSprite.Image.Clear()

		op := &ebiten.DrawImageOptions{}
		w, h := newSprite.Image.Size()
		halfW, halfH := float64(w)/2, float64(h)/2
		op.GeoM.Translate(-halfW, -halfH)
		op.GeoM.Rotate(float64(int(transform.WorldRotation(evolutionChild)-newSprite.OriginalRotation)%360) * 2 * math.Pi / 360)
		op.GeoM.Translate(halfW, halfH)
		op.GeoM.Translate(-halfW, -halfH)
		op.GeoM.Scale(evolutionChildTransform.LocalScale.X, evolutionChildTransform.LocalScale.Y)
		op.GeoM.Translate(halfW, halfH)
		shadowSprite.Image.DrawImage(newSprite.Image, op)

		evolution.GrowTimer.Update()
		if evolution.GrowTimer.IsReady() {
			evolution.ShrinkTimer.Update()

			transform.GetTransform(entry).LocalScale.Set(
				1.0-evolution.ShrinkTimer.PercentDone(),
				1.0-evolution.ShrinkTimer.PercentDone(),
			)
		}

		// TODO To function
		op = &ebiten.DrawImageOptions{}
		currentSprite := component.GetSprite(entry)
		w, h = currentSprite.Image.Size()
		halfW, halfH = float64(w)/2, float64(h)/2
		op.GeoM.Translate(-halfW, -halfH)
		op.GeoM.Rotate(float64(int(transform.WorldRotation(entry)-currentSprite.OriginalRotation)%360) * 2 * math.Pi / 360)
		op.GeoM.Translate(halfW, halfH)
		if t.LocalScale.X > 0 || t.LocalScale.Y > 0 {
			op.GeoM.Translate(-halfW, -halfH)
			op.GeoM.Scale(t.LocalScale.X, t.LocalScale.Y)
			op.GeoM.Translate(halfW, halfH)
		}
		shadowSprite.Image.DrawImage(currentSprite.Image, op)

		if evolution.ShrinkTimer.IsReady() {
			evolution.StopEvolving()
			sprite := component.GetSprite(entry)
			sprite.Image = newSprite.Image
			sprite.ColorOverride = nil

			component.GetSprite(shadow).Image = ebiten.NewImageFromImage(newSprite.Image)
			newSprite.Hide()
		}
	})
}

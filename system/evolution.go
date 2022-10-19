package system

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
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

		transform := component.GetTransform(entry)
		evolutionChild := transform.FindChildWithComponent(component.EvolutionTag)
		if evolutionChild == nil {
			panic("no evolution found")
		}

		evolutionChildTransform := component.GetTransform(evolutionChild)
		evolutionChildTransform.LocalScale = engine.Vector{
			X: evolution.GrowTimer.PercentDone(),
			Y: evolution.GrowTimer.PercentDone(),
		}

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

		shadow := component.GetTransform(entry).FindChildWithComponent(component.ShadowTag)
		shadowSprite := component.GetSprite(shadow)

		shadowSprite.Image.Clear()

		op := &ebiten.DrawImageOptions{}
		w, h := newSprite.Image.Size()
		halfW, halfH := float64(w)/2, float64(h)/2
		op.GeoM.Translate(-halfW, -halfH)
		op.GeoM.Rotate(float64(int(evolutionChildTransform.WorldRotation()-newSprite.OriginalRotation)%360) * 2 * math.Pi / 360)
		op.GeoM.Translate(halfW, halfH)
		op.GeoM.Translate(-halfW, -halfH)
		op.GeoM.Scale(evolutionChildTransform.LocalScale.X, evolutionChildTransform.LocalScale.Y)
		op.GeoM.Translate(halfW, halfH)
		shadowSprite.Image.DrawImage(newSprite.Image, op)

		evolution.GrowTimer.Update()
		if evolution.GrowTimer.IsReady() {
			evolution.ShrinkTimer.Update()

			component.GetTransform(entry).LocalScale = engine.Vector{
				X: 1.0 - evolution.ShrinkTimer.PercentDone(),
				Y: 1.0 - evolution.ShrinkTimer.PercentDone(),
			}
		}

		// TODO To function
		op = &ebiten.DrawImageOptions{}
		currentSprite := component.GetSprite(entry)
		w, h = currentSprite.Image.Size()
		halfW, halfH = float64(w)/2, float64(h)/2
		op.GeoM.Translate(-halfW, -halfH)
		op.GeoM.Rotate(float64(int(transform.WorldRotation()-currentSprite.OriginalRotation)%360) * 2 * math.Pi / 360)
		op.GeoM.Translate(halfW, halfH)
		if transform.LocalScale.X > 0 || transform.LocalScale.Y > 0 {
			op.GeoM.Translate(-halfW, -halfH)
			op.GeoM.Scale(transform.LocalScale.X, transform.LocalScale.Y)
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

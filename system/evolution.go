package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
)

type Evolution struct {
	query *donburi.Query

	shadowBuffer *ebiten.Image
}

func NewEvolution() *Evolution {
	return &Evolution{
		query: donburi.NewQuery(filter.Contains(component.Evolution)),
		// TODO Not that universal in terms of size
		shadowBuffer: ebiten.NewImage(assets.AirplanesBlue[0].Size()),
	}
}

func (s *Evolution) Update(w donburi.World) {
	// TODO Handle player evolving while already evolving (queue evolutions)
	s.query.Each(w, func(entry *donburi.Entry) {
		evolution := component.Evolution.Get(entry)
		if !evolution.Evolving {
			return
		}

		currentEvolution, _ := transform.FindChildWithComponent(entry, component.CurrentEvolutionTag)
		nextEvolution, _ := transform.FindChildWithComponent(entry, component.NextEvolutionTag)

		currentEvolutionSprite := component.Sprite.Get(currentEvolution)
		nextEvolutionSprite := component.Sprite.Get(nextEvolution)

		currentEvolutionTransform := transform.Transform.Get(currentEvolution)
		nextEvolutionTransform := transform.Transform.Get(nextEvolution)

		shadow, _ := transform.FindChildWithComponent(entry, component.ShadowTag)
		shadowSprite := component.Sprite.Get(shadow)

		sprite := component.Sprite.Get(entry)

		if !evolution.StartedEvolving {
			// Hide sprite
			sprite.Hide()

			// Show evolutions instead
			currentEvolutionTransform.LocalScale = dmath.NewVec2(1, 1)
			currentEvolutionSprite.Image = whiteImageFromImage(sprite.Image)
			currentEvolutionSprite.Show()

			nextEvolutionTransform.LocalScale = dmath.NewVec2(0, 0)
			nextEvolutionSprite.Image = whiteImageFromImage(archetype.AirplaneImageByFaction(component.PlayerAirplane.Get(entry).Faction, evolution.Level))
			nextEvolutionSprite.Show()

			evolution.StartedEvolving = true
		}

		evolution.GrowTimer.Update()

		nextEvolutionTransform.LocalScale = dmath.NewVec2(
			evolution.GrowTimer.PercentDone(),
			evolution.GrowTimer.PercentDone(),
		)

		if evolution.GrowTimer.IsReady() {
			evolution.ShrinkTimer.Update()

			currentEvolutionTransform.LocalScale = dmath.NewVec2(
				1.0-evolution.ShrinkTimer.PercentDone(),
				1.0-evolution.ShrinkTimer.PercentDone(),
			)
		}

		w, h := sprite.Image.Size()
		halfW, halfH := float64(w)/2, float64(h)/2

		s.shadowBuffer.Clear()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-halfW, -halfH)
		op.GeoM.Scale(currentEvolutionTransform.LocalScale.X, currentEvolutionTransform.LocalScale.Y)
		op.GeoM.Translate(halfW, halfH)
		s.shadowBuffer.DrawImage(currentEvolutionSprite.Image, op)

		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-halfW, -halfH)
		op.GeoM.Scale(nextEvolutionTransform.LocalScale.X, nextEvolutionTransform.LocalScale.Y)
		op.GeoM.Translate(halfW, halfH)
		s.shadowBuffer.DrawImage(nextEvolutionSprite.Image, op)

		shadowSprite.Image.Clear()
		op = &ebiten.DrawImageOptions{}
		archetype.ShadowDrawOptions(op)
		shadowSprite.Image.DrawImage(s.shadowBuffer, op)

		if evolution.ShrinkTimer.IsReady() {
			evolution.StopEvolving()

			// Hide evolutions
			currentEvolutionSprite.Hide()
			nextEvolutionSprite.Hide()

			// Show sprite and shadow
			sprite.Image = archetype.AirplaneImageByFaction(component.PlayerAirplane.Get(entry).Faction, evolution.Level)
			sprite.Show()
			shadowSprite.Image = archetype.ShadowImage(sprite.Image)
		}
	})
}

func whiteImageFromImage(src *ebiten.Image) *ebiten.Image {
	img := ebiten.NewImage(src.Size())
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(0, 0, 0, 1)
	op.ColorM.Translate(1, 1, 1, 0)
	img.DrawImage(src, op)
	return img
}

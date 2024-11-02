package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
)

type Altitude struct {
	query *donburi.Query
}

func NewAltitude() *Altitude {
	return &Altitude{
		query: donburi.NewQuery(filter.Contains(
			transform.Transform,
			component.Altitude,
		)),
	}
}

func (a *Altitude) Update(w donburi.World) {
	a.query.Each(w, func(entry *donburi.Entry) {
		altitude := component.Altitude.Get(entry)
		altitude.Update()

		scale := 0.8 + 0.2*altitude.Altitude
		t := transform.Transform.Get(entry)
		t.LocalScale.X = scale
		t.LocalScale.Y = scale

		shadow, ok := transform.FindChildWithComponent(entry, component.ShadowTag)
		if ok {
			shadowTransform := transform.Transform.Get(shadow)
			shadowTransform.LocalPosition.X = -archetype.MaxShadowPosition * altitude.Altitude
			shadowTransform.LocalPosition.Y = archetype.MaxShadowPosition * altitude.Altitude
		}

		// Grounded units don't move
		if altitude.Falling && altitude.Altitude == 0 {
			if entry.HasComponent(component.Velocity) {
				velocity := component.Velocity.Get(entry)
				velocity.Velocity = math.Vec2{}
				velocity.RotationVelocity = 0
			}
			sprite := component.Sprite.Get(entry)
			sprite.Layer = component.SpriteLayerDebris
		}
	})
}

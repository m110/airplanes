package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
)

type Altitude struct {
	query *query.Query
}

func NewAltitude() *Altitude {
	return &Altitude{
		query: query.NewQuery(filter.Contains(
			component.Transform,
			component.Altitude,
		)),
	}
}

func (a *Altitude) Update(w donburi.World) {
	a.query.EachEntity(w, func(entry *donburi.Entry) {
		altitude := component.GetAltitude(entry)
		altitude.Update()

		transform := component.GetTransform(entry)
		// TODO Add scaling once render takes WorldScale into account
		// scale := 0.5 + 0.2*altitude.Altitude
		// transform.LocalScale.X = scale
		// transform.LocalScale.Y = scale

		shadow, ok := transform.FindChildWithComponent(component.ShadowTag)
		if ok {
			shadowTransform := component.GetTransform(shadow)
			shadowTransform.LocalPosition.X = -archetype.MaxShadowPosition * altitude.Altitude
			shadowTransform.LocalPosition.Y = archetype.MaxShadowPosition * altitude.Altitude
		}
	})
}

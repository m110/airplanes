package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
)

type Spawn struct {
	query *donburi.Query
}

func NewSpawn() *Spawn {
	return &Spawn{
		query: donburi.NewQuery(filter.Contains(component.Spawnable)),
	}
}

func (s *Spawn) Update(w donburi.World) {
	cameraPos := transform.WorldPosition(archetype.MustFindCamera(w))

	s.query.Each(w, func(entry *donburi.Entry) {
		t := transform.Transform.Get(entry)

		if cameraPos.Y <= t.LocalPosition.Y {
			spawnable := component.Spawnable.Get(entry)
			spawnable.SpawnFunc(w)
			component.Destroy(entry)
		}
	})
}

package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/hierarchy"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
)

type Spawn struct {
	query *query.Query
}

func NewSpawn() *Spawn {
	return &Spawn{
		query: query.NewQuery(filter.Contains(component.Spawnable)),
	}
}

func (s *Spawn) Update(w donburi.World) {
	cameraPos := transform.WorldPosition(archetype.MustFindCamera(w))

	s.query.EachEntity(w, func(entry *donburi.Entry) {
		t := transform.GetTransform(entry)

		if cameraPos.Y <= t.LocalPosition.Y {
			spawnable := component.GetSpawnable(entry)
			spawnable.SpawnFunc(w)
			hierarchy.RemoveRecursive(entry)
		}
	})
}

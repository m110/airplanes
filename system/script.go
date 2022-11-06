package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/component"
)

type Script struct {
	query *query.Query
}

func NewScript() *Script {
	return &Script{
		query: query.NewQuery(filter.Contains(component.Script)),
	}
}

func (s *Script) Update(w donburi.World) {
	s.query.EachEntity(w, func(entry *donburi.Entry) {
		script := component.Script.Get(entry)
		script.Update(w)
	})
}

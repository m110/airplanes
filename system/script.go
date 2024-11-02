package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/airplanes/component"
)

type Script struct {
	query *donburi.Query
}

func NewScript() *Script {
	return &Script{
		query: donburi.NewQuery(filter.Contains(component.Script)),
	}
}

func (s *Script) Update(w donburi.World) {
	s.query.Each(w, func(entry *donburi.Entry) {
		script := component.Script.Get(entry)
		script.Update(w)
	})
}

package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/query"
)

type ObserverData struct {
	LookFor *query.Query
	Target  *TransformData
}

var Observer = donburi.NewComponentType[ObserverData]()

func GetObserver(entry *donburi.Entry) *ObserverData {
	return donburi.Get[ObserverData](entry, Observer)
}

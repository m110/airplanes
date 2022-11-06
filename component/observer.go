package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/query"
)

type ObserverData struct {
	LookFor *query.Query
	Target  *donburi.Entry
}

var Observer = donburi.NewComponentType[ObserverData]()

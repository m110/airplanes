package component

import "github.com/yohamta/donburi"

type DestroyedData struct{}

var Destroyed = donburi.NewComponentType[DestroyedData]()

func Destroy(e *donburi.Entry) {
	if !e.Valid() {
		return
	}
	if !e.HasComponent(Destroyed) {
		e.AddComponent(Destroyed)
	}
}

package system

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/component"
)

func Destroy(w donburi.World, entry *donburi.Entry) {
	if entry.HasComponent(component.Transform) {
		parent := component.GetTransform(entry)
		for _, child := range parent.Children {
			Destroy(w, child)
		}
	}

	w.Remove(entry.Entity())
}

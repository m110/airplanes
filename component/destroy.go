package component

import (
	"github.com/yohamta/donburi"
)

// TODO It feels like it should be somewhere in the engine package but it causes an import cycle
// Perhaps the "parent" relation shouldn't be a component?
func Destroy(w donburi.World, entry *donburi.Entry) {
	if entry.HasComponent(Parent) {
		parent := GetParent(entry)
		for _, child := range parent.Children {
			Destroy(w, child)
		}
	}

	w.Remove(entry.Entity())
}

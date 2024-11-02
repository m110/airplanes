package engine

import (
	"fmt"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/component"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

func MustGetParent(entry *donburi.Entry) *donburi.Entry {
	parent, ok := transform.GetParent(entry)
	if !ok {
		panic("parent not found")
	}
	return parent
}

func MustFindChildWithComponent(parent *donburi.Entry, componentType component.IComponentType) *donburi.Entry {
	entry, ok := transform.FindChildWithComponent(parent, componentType)
	if !ok {
		panic(fmt.Sprintf("entry not found with component %T", componentType))
	}
	return entry
}

func FindWithComponent(w donburi.World, componentType component.IComponentType) (*donburi.Entry, bool) {
	return donburi.NewQuery(filter.Contains(componentType)).First(w)
}

func MustFindWithComponent(w donburi.World, componentType component.IComponentType) *donburi.Entry {
	entry, ok := FindWithComponent(w, componentType)
	if !ok {
		panic(fmt.Sprintf("entry not found with component %T", componentType))
	}
	return entry
}

type Component[T any] interface {
	donburi.IComponentType
	Get(entry *donburi.Entry) *T
}

func MustFindComponent[T any](w donburi.World, c Component[T]) *T {
	entry, ok := donburi.NewQuery(filter.Contains(c)).First(w)
	if !ok {
		panic("MustFindComponent: entry not found")
	}

	return c.Get(entry)
}

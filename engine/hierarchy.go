package engine

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/component"
	"github.com/yohamta/donburi/features/transform"
)

func FindChildrenWithComponent(e *donburi.Entry, c component.IComponentType) []*donburi.Entry {
	if !e.Valid() {
		return nil
	}

	children, ok := transform.GetChildren(e)
	if !ok {
		return nil
	}

	var result []*donburi.Entry
	for _, child := range children {
		if !child.Valid() {
			continue
		}

		if child.HasComponent(c) {
			result = append(result, child)
		}

		result = append(result, FindChildrenWithComponent(child, c)...)
	}

	return result
}

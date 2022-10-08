package system

import (
	"fmt"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

type Collision struct {
	query *query.Query
}

func NewCollision() *Collision {
	return &Collision{
		query: query.NewQuery(filter.Contains(component.Collider)),
	}
}

type collisionEffect func(w donburi.World, entry *donburi.Entry, other *donburi.Entry)

var collisionEffects = map[component.ColliderLayer]map[component.ColliderLayer]collisionEffect{
	component.CollisionLayerBullets: {
		component.CollisionLayerEnemies: func(w donburi.World, entry *donburi.Entry, other *donburi.Entry) {
			w.Remove(entry.Entity())
			component.GetHealth(other).Damage()
		},
	},
	component.CollisionLayerEnemies: {
		component.CollisionLayerPlayers: func(w donburi.World, entry *donburi.Entry, other *donburi.Entry) {
			// TODO damage player
		},
	},
}

func (c *Collision) Update(w donburi.World) {
	var entries []*donburi.Entry
	c.query.EachEntity(w, func(entry *donburi.Entry) {
		// Skip entities not spawned yet
		if entry.HasComponent(component.Despawnable) {
			if !component.GetDespawnable(entry).Spawned {
				return
			}
		}
		entries = append(entries, entry)
	})

	for _, entry := range entries {
		collider := component.GetCollider(entry)

		for _, other := range entries {
			if entry.Entity().Id() == other.Entity().Id() {
				continue
			}

			// One of the entities could already be removed from the world due to collision effect
			if !w.Valid(entry.Entity()) || !w.Valid(other.Entity()) {
				continue
			}

			otherCollider := component.GetCollider(other)

			effects, ok := collisionEffects[collider.Layer]
			if !ok {
				continue
			}

			effect, ok := effects[otherCollider.Layer]
			if !ok {
				continue
			}

			if !entry.HasComponent(component.Position) {
				panic(fmt.Sprintf("%#v missing position\n", entry.Entity().Id()))
			}
			pos := component.GetPosition(entry)
			otherPos := component.GetPosition(other)

			// TODO The current approach doesn't take rotation into account
			rect := engine.NewRect(pos.X, pos.Y, collider.Width, collider.Height)
			otherRect := engine.NewRect(otherPos.X, otherPos.Y, otherCollider.Width, otherCollider.Height)

			if rect.Intersects(otherRect) {
				effect(w, entry, other)
			}
		}
	}
}

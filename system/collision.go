package system

import (
	"fmt"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/archetypes"
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
	component.CollisionLayerPlayerBullets: {
		component.CollisionLayerAirEnemies: func(w donburi.World, entry *donburi.Entry, other *donburi.Entry) {
			Destroy(w, entry)
			component.GetHealth(other).Damage()
		},
		component.CollisionLayerGroundEnemies: func(w donburi.World, entry *donburi.Entry, other *donburi.Entry) {
			Destroy(w, entry)
			component.GetHealth(other).Damage()
		},
	},
	component.CollisionLayerPlayers: {
		component.CollisionLayerCollectibles: func(w donburi.World, entry *donburi.Entry, other *donburi.Entry) {
			airplane := component.GetPlayerAirplane(entry)
			player := archetypes.MustFindPlayerByNumber(w, airplane.PlayerNumber)

			// TODO Is this the best place to do this?
			switch component.GetCollectible(other).Type {
			case component.CollectibleTypeWeaponUpgrade:
				player.UpgradeWeapon()
			case component.CollectibleTypeShield:
				airplane.StartInvulnerability()
			case component.CollectibleTypeHealth:
				player.AddLive()
			}

			Destroy(w, other)
		},
	},
	component.CollisionLayerAirEnemies: {
		component.CollisionLayerPlayers: func(w donburi.World, entry *donburi.Entry, other *donburi.Entry) {
			component.GetHealth(entry).Destroy()
			damagePlayer(w, other)
		},
	},
	component.CollisionLayerEnemyBullets: {
		component.CollisionLayerPlayers: func(w donburi.World, entry *donburi.Entry, other *donburi.Entry) {
			Destroy(w, entry)
			damagePlayer(w, other)
		},
	},
}

func damagePlayer(w donburi.World, entry *donburi.Entry) {
	if component.GetPlayerAirplane(entry).Invulnerable {
		return
	}

	playerNumber := component.GetPlayerAirplane(entry).PlayerNumber
	Destroy(w, entry)

	player := archetypes.MustFindPlayerByNumber(w, playerNumber)
	player.Damage()
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
		if !w.Valid(entry.Entity()) {
			continue
		}

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

			if !entry.HasComponent(component.Transform) {
				panic(fmt.Sprintf("%#v missing position\n", entry.Entity().Id()))
			}
			pos := component.GetTransform(entry).LocalPosition
			otherPos := component.GetTransform(other).LocalPosition

			// TODO The current approach doesn't take rotation into account
			rect := engine.NewRect(pos.X, pos.Y, collider.Width, collider.Height)
			otherRect := engine.NewRect(otherPos.X, otherPos.Y, otherCollider.Width, otherCollider.Height)

			if rect.Intersects(otherRect) {
				effect(w, entry, other)
			}
		}
	}
}

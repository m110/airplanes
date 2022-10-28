package system

import (
	"fmt"

	"github.com/yohamta/donburi/features/hierarchy"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/archetype"
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
			hierarchy.RemoveRecursive(entry)
			component.GetHealth(other).Damage()
		},
		component.CollisionLayerGroundEnemies: func(w donburi.World, entry *donburi.Entry, other *donburi.Entry) {
			hierarchy.RemoveRecursive(entry)
			component.GetHealth(other).Damage()
		},
	},
	component.CollisionLayerPlayers: {
		component.CollisionLayerCollectibles: func(w donburi.World, entry *donburi.Entry, other *donburi.Entry) {
			airplane := component.GetPlayerAirplane(entry)
			player := archetype.MustFindPlayerByNumber(w, airplane.PlayerNumber)

			// TODO Is this the best place to do this?
			switch component.GetCollectible(other).Type {
			case component.CollectibleTypeWeaponUpgrade:
				player.UpgradeWeapon()

				evolution := component.GetEvolution(entry)
				if player.EvolutionLevel() > evolution.Level {
					evolution.Evolve()
				}
			case component.CollectibleTypeShield:
				airplane.StartInvulnerability()
			case component.CollectibleTypeHealth:
				player.AddLive()
			}

			hierarchy.RemoveRecursive(other)
		},
	},
	component.CollisionLayerAirEnemies: {
		component.CollisionLayerPlayers: func(w donburi.World, entry *donburi.Entry, other *donburi.Entry) {
			// TODO deduplicate
			component.MustFindGame(w).AddScore(1)

			if entry.HasComponent(component.Wreckable) {
				archetype.NewAirplaneWreck(w, entry, component.GetSprite(entry))
			}

			hierarchy.RemoveRecursive(entry)

			damagePlayer(w, other)
		},
	},
	component.CollisionLayerEnemyBullets: {
		component.CollisionLayerPlayers: func(w donburi.World, entry *donburi.Entry, other *donburi.Entry) {
			hierarchy.RemoveRecursive(entry)
			damagePlayer(w, other)
		},
	},
}

func damagePlayer(w donburi.World, entry *donburi.Entry) {
	if component.GetPlayerAirplane(entry).Invulnerable {
		return
	}

	playerNumber := component.GetPlayerAirplane(entry).PlayerNumber

	if entry.HasComponent(component.Wreckable) {
		archetype.NewAirplaneWreck(w, entry, component.GetSprite(entry))
	}
	hierarchy.RemoveRecursive(entry)

	player := archetype.MustFindPlayerByNumber(w, playerNumber)
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
		if !entry.Valid() {
			continue
		}

		collider := component.GetCollider(entry)

		for _, other := range entries {
			if entry.Entity().Id() == other.Entity().Id() {
				continue
			}

			// One of the entities could already be removed from the world due to collision effect
			if !entry.Valid() || !other.Valid() {
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

			if !entry.HasComponent(transform.Transform) {
				panic(fmt.Sprintf("%#v missing position\n", entry.Entity().Id()))
			}
			pos := transform.GetTransform(entry).LocalPosition
			otherPos := transform.GetTransform(other).LocalPosition

			// TODO The current approach doesn't take rotation into account
			// TODO The current approach doesn't take scale into account
			rect := engine.NewRect(pos.X, pos.Y, collider.Width, collider.Height)
			otherRect := engine.NewRect(otherPos.X, otherPos.Y, otherCollider.Width, otherCollider.Height)

			if rect.Intersects(otherRect) {
				effect(w, entry, other)
			}
		}
	}
}

package system

import (
	"fmt"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

type Collision struct {
	query *donburi.Query
}

func NewCollision() *Collision {
	return &Collision{
		query: donburi.NewQuery(filter.Contains(component.Collider)),
	}
}

type collisionEffect func(w donburi.World, entry *donburi.Entry, other *donburi.Entry)

var collisionEffects = map[component.ColliderLayer]map[component.ColliderLayer]collisionEffect{
	component.CollisionLayerPlayerBullets: {
		component.CollisionLayerAirEnemies: func(w donburi.World, entry *donburi.Entry, other *donburi.Entry) {
			component.Destroy(entry)
			component.Health.Get(other).Damage()
		},
		component.CollisionLayerGroundEnemies: func(w donburi.World, entry *donburi.Entry, other *donburi.Entry) {
			component.Destroy(entry)
			component.Health.Get(other).Damage()
		},
	},
	component.CollisionLayerPlayers: {
		component.CollisionLayerCollectibles: func(w donburi.World, entry *donburi.Entry, other *donburi.Entry) {
			airplane := component.PlayerAirplane.Get(entry)
			player := archetype.MustFindPlayerByNumber(w, airplane.PlayerNumber)

			// TODO Is this the best place to do this?
			switch component.Collectible.Get(other).Type {
			case component.CollectibleTypeWeaponUpgrade:
				player.UpgradeWeapon()

				evolution := component.Evolution.Get(entry)
				if player.EvolutionLevel() > evolution.Level {
					evolution.Evolve()
				}
			case component.CollectibleTypeShield:
				airplane.StartInvulnerability()
			case component.CollectibleTypeHealth:
				player.AddLive()
			}

			component.Destroy(other)
		},
	},
	component.CollisionLayerAirEnemies: {
		component.CollisionLayerPlayers: func(w donburi.World, entry *donburi.Entry, other *donburi.Entry) {
			EnemyKilledEvent.Publish(w, EnemyKilled{
				Enemy: entry,
			})

			damagePlayer(w, other)
		},
	},
	component.CollisionLayerEnemyBullets: {
		component.CollisionLayerPlayers: func(w donburi.World, entry *donburi.Entry, other *donburi.Entry) {
			component.Destroy(entry)
			damagePlayer(w, other)
		},
	},
}

func damagePlayer(w donburi.World, entry *donburi.Entry) {
	if component.PlayerAirplane.Get(entry).Invulnerable {
		return
	}

	playerNumber := component.PlayerAirplane.Get(entry).PlayerNumber

	if entry.HasComponent(component.Wreckable) {
		archetype.NewAirplaneWreck(w, entry, component.Sprite.Get(entry))
	}
	component.Destroy(entry)

	player := archetype.MustFindPlayerByNumber(w, playerNumber)
	player.Damage()
}

func (c *Collision) Update(w donburi.World) {
	var entries []*donburi.Entry
	c.query.Each(w, func(entry *donburi.Entry) {
		// Skip entities not spawned yet
		if entry.HasComponent(component.Despawnable) {
			if !component.Despawnable.Get(entry).Spawned {
				return
			}
		}
		entries = append(entries, entry)
	})

	for _, entry := range entries {
		if !entry.Valid() {
			continue
		}

		collider := component.Collider.Get(entry)

		for _, other := range entries {
			if entry.Entity().Id() == other.Entity().Id() {
				continue
			}

			otherCollider := component.Collider.Get(other)

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
			pos := transform.Transform.Get(entry).LocalPosition
			otherPos := transform.Transform.Get(other).LocalPosition

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

package system

import (
	"math/rand"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/hierarchy"
	"github.com/yohamta/donburi/features/transform"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

type EnemyKilled struct {
	Enemy *donburi.Entry
}

func OnEnemyKilled(w donburi.World, event engine.Event) {
	e := event.(EnemyKilled)

	// TODO deduplicate
	component.MustFindGame(w).AddScore(1)

	if e.Enemy.HasComponent(component.Wreckable) {
		archetype.NewAirplaneWreck(w, e.Enemy, component.GetSprite(e.Enemy))
	}

	// TODO A temporary high chance for test purposes
	r := rand.Intn(10)
	if r < 7 {
		archetype.NewRandomCollectible(w, transform.GetTransform(e.Enemy).LocalPosition)
	}

	hierarchy.RemoveRecursive(e.Enemy)
}

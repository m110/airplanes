package system

import (
	"math/rand"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
	"github.com/yohamta/donburi/features/transform"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
)

type EnemyKilled struct {
	Enemy *donburi.Entry
}

var EnemyKilledEvent = events.NewEventType[EnemyKilled]()

func OnEnemyKilledWreck(w donburi.World, event EnemyKilled) {
	if event.Enemy.Valid() && event.Enemy.HasComponent(component.Wreckable) {
		archetype.NewAirplaneWreck(w, event.Enemy, component.Sprite.Get(event.Enemy))
	}
}

func OnEnemyKilledAddScore(w donburi.World, event EnemyKilled) {
	component.MustFindGame(w).AddScore(1)
}

func OnEnemyKilledSpawnCollectible(w donburi.World, event EnemyKilled) {
	if !event.Enemy.Valid() {
		return
	}

	r := rand.Intn(10)
	if r < 2 {
		archetype.NewRandomCollectible(w, transform.Transform.Get(event.Enemy).LocalPosition)
	}
}

func OnEnemyKilledDestroyEnemy(w donburi.World, event EnemyKilled) {
	component.Destroy(event.Enemy)
}

func SetupEvents(w donburi.World) {
	EnemyKilledEvent.Subscribe(w, OnEnemyKilledWreck)
	EnemyKilledEvent.Subscribe(w, OnEnemyKilledAddScore)
	EnemyKilledEvent.Subscribe(w, OnEnemyKilledSpawnCollectible)
	EnemyKilledEvent.Subscribe(w, OnEnemyKilledDestroyEnemy)
}

type Events struct{}

func NewEvents() *Events {
	return &Events{}
}

func (e *Events) Update(w donburi.World) {
	EnemyKilledEvent.ProcessEvents(w)
}

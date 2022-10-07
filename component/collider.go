package component

import (
	"github.com/yohamta/donburi"
)

const (
	CollisionLayerBullets ColliderLayer = iota
	CollisionLayerEnemies
	CollisionLayerPlayers
)

type ColliderLayer int

type ColliderData struct {
	Width  float64
	Height float64
	Layer  ColliderLayer
}

var Collider = donburi.NewComponentType[ColliderData]()

func GetCollider(entry *donburi.Entry) *ColliderData {
	return donburi.Get[ColliderData](entry, Collider)
}

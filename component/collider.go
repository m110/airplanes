package component

import (
	"github.com/yohamta/donburi"
)

const (
	CollisionLayerBullets = iota
	CollisionLayerEnemies
	CollisionLayerPlayers
)

type ColliderData struct {
	Width  int
	Height int
	Layer  int
}

var Collider = donburi.NewComponentType[ColliderData]()

func GetCollider(entry *donburi.Entry) *ColliderData {
	return donburi.Get[ColliderData](entry, Collider)
}

package archetypes

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

func NewRandomCollectible(w donburi.World, position engine.Vector) {
	collectible := w.Entry(w.Create(
		component.Transform,
		component.Sprite,
		component.Collider,
		component.Collectible,
		component.Despawnable,
	))

	var image *ebiten.Image
	var collectibleType component.CollectibleType
	switch rand.Intn(3) {
	case 0:
		image = assets.WeaponUpgrade
		collectibleType = component.CollectibleTypeWeaponUpgrade
	case 1:
		image = assets.Shield
		collectibleType = component.CollectibleTypeShield
	case 2:
		image = assets.Health
		collectibleType = component.CollectibleTypeHealth
	}

	donburi.SetValue(collectible, component.Transform, component.TransformData{
		Position: position,
	})

	donburi.SetValue(collectible, component.Sprite, component.SpriteData{
		Image: image,
		Layer: component.SpriteLayerCollectibles,
		Pivot: component.SpritePivotCenter,
	})

	width, height := image.Size()
	donburi.SetValue(collectible, component.Collider, component.ColliderData{
		Width:  float64(width),
		Height: float64(height),
		Layer:  component.CollisionLayerCollectibles,
	})

	donburi.SetValue(collectible, component.Collectible, component.CollectibleData{
		Type: collectibleType,
	})
}

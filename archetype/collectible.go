package archetype

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
)

func NewRandomCollectible(w donburi.World, position math.Vec2) {
	collectible := w.Entry(w.Create(
		transform.Transform,
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

	transform.Transform.Get(collectible).LocalPosition = position

	component.Sprite.SetValue(collectible, component.SpriteData{
		Image: image,
		Layer: component.SpriteLayerCollectibles,
		Pivot: component.SpritePivotCenter,
	})

	width, height := image.Size()
	component.Collider.SetValue(collectible, component.ColliderData{
		Width:  float64(width),
		Height: float64(height),
		Layer:  component.CollisionLayerCollectibles,
	})

	component.Collectible.SetValue(collectible, component.CollectibleData{
		Type: collectibleType,
	})
}

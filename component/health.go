package component

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/engine"
)

type HealthData struct {
	Health               int
	JustDamaged          bool
	DamageIndicatorTimer *engine.Timer
}

func (d *HealthData) Damage() {
	if d.Health <= 0 {
		return
	}

	d.Health--
	d.JustDamaged = true
	d.DamageIndicatorTimer.Reset()
}

var Health = donburi.NewComponentType[HealthData]()

func GetHealth(entry *donburi.Entry) *HealthData {
	return donburi.Get[HealthData](entry, Health)
}

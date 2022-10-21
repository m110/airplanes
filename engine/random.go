package engine

import "math/rand"

func RandomRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

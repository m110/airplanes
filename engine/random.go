package engine

import "math/rand"

func RandomFloatRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func RandomIntRange(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func RandomFrom[T comparable](list []T) T {
	index := rand.Intn(len(list))
	return list[index]
}

func RandomFromOrEmpty[T comparable](list []T) *T {
	index := rand.Intn(len(list) + 1)
	if index == len(list) {
		return nil
	}
	return &list[index]
}

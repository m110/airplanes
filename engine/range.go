package engine

import "time"

type IntRange struct {
	Min int
	Max int
}

func (r IntRange) Random() int {
	return RandomIntRange(r.Min, r.Max)
}

type FloatRange struct {
	Min float64
	Max float64
}

func (r FloatRange) Random() float64 {
	return RandomFloatRange(r.Min, r.Max)
}

type DurationRange struct {
	Min time.Duration
	Max time.Duration
}

func (r DurationRange) Random() time.Duration {
	return time.Duration(RandomIntRange(int(r.Min), int(r.Max)))
}

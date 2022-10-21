package engine

import "math"

type Vector struct {
	X float64
	Y float64
}

func (v Vector) IsZero() bool {
	return v.X == 0 && v.Y == 0
}

func (v Vector) Add(other Vector) Vector {
	return Vector{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

func (v Vector) Sub(other Vector) Vector {
	return Vector{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

func (v Vector) Mul(other Vector) Vector {
	return Vector{
		X: v.X * other.X,
		Y: v.Y * other.Y,
	}
}

func (v Vector) AddScalar(value float64) Vector {
	return Vector{
		X: v.X + value,
		Y: v.Y + value,
	}
}

func (v Vector) MulScalar(value float64) Vector {
	return Vector{
		X: v.X * value,
		Y: v.Y * value,
	}
}

func (v Vector) Distance(other Vector) float64 {
	return math.Sqrt(math.Pow(v.X-other.X, 2) + math.Pow(v.Y-other.Y, 2))
}

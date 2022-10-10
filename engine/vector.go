package engine

type Vector struct {
	X float64
	Y float64
}

func (v Vector) Add(other Vector) Vector {
	return Vector{
		X: v.X + other.X,
		Y: v.Y + other.Y,
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

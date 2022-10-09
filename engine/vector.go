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

package helpers

import "github.com/deeean/go-vector/vector2"

type MyVect2 struct {
	vector2.Vector2
}

func NewMyVect2(X float64, Y float64) *MyVect2 {
	return &MyVect2{Vector2: vector2.Vector2{X: X, Y: Y}}
}

func (v *MyVect2) Add(other *MyVect2) *MyVect2 {
	return &MyVect2{*v.Vector2.Add(&other.Vector2)}
}

func (v MyVect2) Subtract(other *MyVect2) *MyVect2 {
	return v.Add(other.MulScalar(-1))
}

func (v *MyVect2) Distance(other *MyVect2) float64 {
	return v.Vector2.Distance(&other.Vector2)
}

func (v *MyVect2) MulScalar(scalar float64) *MyVect2 {
	return &MyVect2{*v.Vector2.MulScalar(scalar)}
}

func (v *MyVect2) Dot(other *MyVect2) float64 {
	return v.Vector2.Dot(&other.Vector2)
}

type Copiable interface {
	Copy() Copiable
}

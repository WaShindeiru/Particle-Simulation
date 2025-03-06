package circle

import (
	"collision/helpers"
	"github.com/deeean/go-vector/vector2"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math"
	"math/rand"
)

type circleFactory struct {
	index int
}

func (factory *circleFactory) GetId() int {
	temp := factory.index
	factory.index++
	return temp
}

func (fact *circleFactory) GetCircleRandomRadius(x_max int, y_max int) *BasicCircle {
	radius := rand.Intn(8) + 8
	return fact.GetCircleRandom(radius, x_max, y_max)
}

func (fact *circleFactory) GetCircleRandom(r int, x_max int, y_max int) *BasicCircle {
	x_ := float64(rand.Intn(x_max-6*r) + 2*r)
	y_ := float64(rand.Intn(y_max-6*r) + 2*r)

	vel_x := float64(rand.Intn(6) - 3)
	vel_y := float64(rand.Intn(6) - 3)

	return fact.GetCircle(r, x_max, y_max, x_, y_, vel_x, vel_y)
}

func (fact *circleFactory) GetCircle(r int, x_max int, y_max int, x, y, vel_x, vel_y float64) *BasicCircle {
	temp := color.RGBA{R: uint8(rand.Intn(255)), G: uint8(rand.Intn(255)), B: uint8(rand.Intn(255)), A: 255}
	img_temp := ebiten.NewImage(2*r, 2*r)

	for y := 0; y < 2*r; y++ {
		for x := 0; x < 2*r; x++ {
			var dy float64 = float64(r) - float64(y)
			var dx float64 = float64(r) - float64(x)

			if math.Sqrt(dx*dx+dy*dy) < float64(r) {
				img_temp.Set(x, y, temp)
			}
		}
	}

	return &BasicCircle{
		pos:      helpers.MyVect2{vector2.Vector2{X: x, Y: y}},
		vel:      helpers.MyVect2{vector2.Vector2{X: vel_x, Y: vel_y}},
		midpoint: *helpers.NewMyVect2(x+float64(r), y+float64(r)),
		mass:     math.Pow(float64(r)/radius_avg, 2) * mass_avg,
		max_x:    x_max,
		max_y:    y_max,
		radius:   r,
		img:      img_temp,
		id:       CircleId(fact.GetId()),
	}
}

var lock chan bool = make(chan bool, 1)

var circleFactorySingleton *circleFactory

func GetInstance() *circleFactory {
	if circleFactorySingleton == nil {
		lock <- true

		if circleFactorySingleton == nil {
			circleFactorySingleton = &circleFactory{index: 1}
		}

		<-lock
	}

	return circleFactorySingleton
}

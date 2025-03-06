package circle

import (
	"collision/helpers"
	"github.com/deeean/go-vector/vector2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"math"
	"math/rand"
)

const (
	Image_Width  = 64
	Image_Height = 64
)

type SimpleCircleFactory struct {
	index int
}

func (factory *SimpleCircleFactory) GetId() int {
	temp := factory.index
	factory.index++
	return temp
}

func (fact *SimpleCircleFactory) GetCircleRandomRadius(x_max int, y_max int) *BasicCircle {
	radius := rand.Intn(8) + 8
	return fact.GetCircleRandom(radius, x_max, y_max)
}

func (fact *SimpleCircleFactory) GetCircleRandom(r int, x_max int, y_max int) *BasicCircle {
	x_ := float64(rand.Intn(x_max-6*r) + 2*r)
	y_ := float64(rand.Intn(y_max-6*r) + 2*r)

	vel_x := float64(rand.Intn(6) - 3)
	vel_y := float64(rand.Intn(6) - 3)

	return fact.GetCircle(r, x_max, y_max, x_, y_, vel_x, vel_y)
}

func (fact *SimpleCircleFactory) GetCircle(r int, x_max int, y_max int, x, y, vel_x, vel_y float64) *BasicCircle {
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
		img: img_temp,
		CircleGeneric: CircleGeneric{
			pos:      helpers.MyVect2{vector2.Vector2{X: x, Y: y}},
			vel:      helpers.MyVect2{vector2.Vector2{X: vel_x, Y: vel_y}},
			midpoint: *helpers.NewMyVect2(x+float64(r), y+float64(r)),
			mass:     math.Pow(float64(r)/radius_avg, 2) * mass_avg,
			max_x:    x_max,
			max_y:    y_max,
			radius:   r,
			id:       CircleId(fact.GetId()),
		},
	}
}

var lock chan bool = make(chan bool, 1)

var circleFactorySingleton *SimpleCircleFactory

func GetInstance() *SimpleCircleFactory {
	if circleFactorySingleton == nil {
		lock <- true

		if circleFactorySingleton == nil {
			circleFactorySingleton = &SimpleCircleFactory{index: 1}
		}

		<-lock
	}

	return circleFactorySingleton
}

type GameCircleFactory struct {
	index int
	image []*ebiten.Image
}

func (factory *GameCircleFactory) GetId() int {
	temp := factory.index
	factory.index++
	return temp
}

func (fact *GameCircleFactory) GetCircleRandomRadius(x_max int, y_max int) *GameCircle {
	radius := rand.Intn(8) + 27
	return fact.GetCircleRandom(radius, x_max, y_max)
}

func (fact *GameCircleFactory) GetCircleRandom(r int, x_max int, y_max int) *GameCircle {
	x_ := float64(rand.Intn(x_max-6*r) + 2*r)
	y_ := float64(rand.Intn(y_max-6*r) + 2*r)

	vel_x := float64(rand.Intn(6) - 3)
	vel_y := float64(rand.Intn(6) - 3)

	return fact.GetCircle(r, x_max, y_max, x_, y_, vel_x, vel_y)
}

func (fact *GameCircleFactory) GetCircle(r int, x_max int, y_max int, x, y, vel_x, vel_y float64) *GameCircle {
	type_ := Type(rand.Intn(3))

	return &GameCircle{
		img:      fact.image,
		circType: type_,
		CircleGeneric: CircleGeneric{
			pos:      helpers.MyVect2{vector2.Vector2{X: x, Y: y}},
			vel:      helpers.MyVect2{vector2.Vector2{X: vel_x, Y: vel_y}},
			midpoint: *helpers.NewMyVect2(x+float64(r), y+float64(r)),
			mass:     math.Pow(float64(r)/radius_avg, 2) * mass_avg,
			max_x:    x_max,
			max_y:    y_max,
			radius:   r,
			id:       CircleId(fact.GetId()),
		},
	}
}

func newGameFactory() *GameCircleFactory {

	gameFactory := &GameCircleFactory{
		index: 1,
		image: make([]*ebiten.Image, 3),
	}

	var img *ebiten.Image
	var err error

	img, _, err = ebitenutil.NewImageFromFile("/home/washindeiru/significant/go/collision_temp_new/collision/image/paper.png")
	helpers.Assert(err == nil, "not ok")
	gameFactory.image[0] = img

	img, _, err = ebitenutil.NewImageFromFile("/home/washindeiru/significant/go/collision_temp_new/collision/image/rock.png")
	helpers.Assert(err == nil, "not ok")
	gameFactory.image[1] = img

	img, _, err = ebitenutil.NewImageFromFile("/home/washindeiru/significant/go/collision_temp_new/collision/image/scissors.png")
	helpers.Assert(err == nil, "not ok")
	gameFactory.image[2] = img

	//var temp color.RGBA
	//for i := 0; i < 3; i++ {
	//	gameFactory.image[i] = ebiten.NewImage(Image_Width, Image_Height)
	//	temp = color.RGBA{R: uint8(rand.Intn(255)), G: uint8(rand.Intn(255)), B: uint8(rand.Intn(255)), A: 255}
	//
	//	for y := 0; y < Image_Height; y++ {
	//		for x := 0; x < Image_Width; x++ {
	//			var dy float64 = float64(Image_Height/2) - float64(y)
	//			var dx float64 = float64(Image_Width/2) - float64(x)
	//
	//			if math.Sqrt(dx*dx+dy*dy) < float64(Image_Width/2) {
	//				gameFactory.image[i].Set(x, y, temp)
	//			}
	//		}
	//	}
	//}

	return gameFactory
}

var lockGame chan bool = make(chan bool, 1)

var gameFactorySingleton *GameCircleFactory

func GetGameFactoryInstance() *GameCircleFactory {
	if gameFactorySingleton == nil {
		lockGame <- true

		if gameFactorySingleton == nil {
			gameFactorySingleton = newGameFactory()
		}

		<-lockGame
	}

	return gameFactorySingleton
}

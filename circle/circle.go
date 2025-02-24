package circle

import (
	"collision/helpers"
	"fmt"
	"github.com/deeean/go-vector/vector2"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math"
	"math/rand"
)

const (
	radius_avg = 16
	mass_avg   = 10
)

type Circle struct {
	img *ebiten.Image
	pos helpers.MyVect2

	midpoint helpers.MyVect2
	//x_center float64
	//y_center float64

	vel helpers.MyVect2

	max_x int
	max_y int

	mass float64

	radius int
}

func (c Circle) GetPos() helpers.MyVect2 {
	return c.pos
}

func (c Circle) GetMidpoint() helpers.MyVect2 {
	return c.midpoint
}

func (c *Circle) SetPos(v helpers.MyVect2) {
	c.pos = v
	c.midpoint = helpers.MyVect2{Vector2: vector2.Vector2{X: c.pos.X + float64(c.radius), Y: c.pos.Y + float64(c.radius)}}
}

func (c *Circle) SetMidpoint(v helpers.MyVect2) {
	c.midpoint = v
	c.pos = *helpers.NewMyVect2(c.midpoint.X-float64(c.radius), c.midpoint.Y-float64(c.radius))
}

func New_random_radius(x_max int, y_max int) Circle {
	radius := rand.Intn(4) + 16
	return New(radius, x_max, y_max)
}

func New(r int, x_max int, y_max int) Circle {
	x_ := float64(rand.Intn(x_max-6*r) + 2*r)
	y_ := float64(rand.Intn(y_max-6*r) + 2*r)

	vel_x := float64(rand.Intn(6) - 3)
	vel_y := float64(rand.Intn(6) - 3)

	return New_(r, x_max, y_max, x_, y_, vel_x, vel_y)
}

func New_(r int, x_max int, y_max int, x, y, vel_x, vel_y float64) Circle {
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

	return Circle{
		pos:      helpers.MyVect2{vector2.Vector2{X: x, Y: y}},
		vel:      helpers.MyVect2{vector2.Vector2{X: vel_x, Y: vel_y}},
		midpoint: *helpers.NewMyVect2(x+float64(r), y+float64(r)),
		mass:     math.Pow(float64(r/radius_avg), 2) * mass_avg,
		max_x:    x_max,
		max_y:    y_max,
		radius:   r,
		img:      img_temp,
	}
}

func (c *Circle) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.pos.X), float64(c.pos.Y))
	screen.DrawImage(c.img, op)
}

func (c *Circle) Update() {
	var prev Circle = *c

	vx := c.vel.X
	vy := c.vel.Y

	pos := c.GetPos()
	c.SetPos(*pos.Add(helpers.NewMyVect2(vx, vy)))

	*c = handleCollisionWithBorder(prev, *c)
}

// continuous collision
// tunneling prevention
func handleCollisionWithBorder(prev Circle, curr Circle) Circle {
	if curr.pos.Y < 0 {
		timeOfCollision := prev.pos.Y / (prev.pos.Y - curr.pos.Y)
		correct_y := prev.pos.Y + curr.vel.Y*timeOfCollision
		curr.vel = helpers.MyVect2{vector2.Vector2{X: curr.vel.X, Y: -curr.vel.Y}}
		correct_y += curr.vel.Y * (1 - timeOfCollision)

		curr.SetPos(*helpers.NewMyVect2(curr.GetPos().X, correct_y))

	} else if my := float64(curr.max_y - 2*curr.radius); my <= curr.pos.Y {
		timeOfCollision := (my - prev.pos.Y) / (curr.pos.Y - prev.pos.Y)
		correct_y := prev.pos.Y + curr.vel.Y*timeOfCollision
		curr.vel = helpers.MyVect2{vector2.Vector2{X: curr.vel.X, Y: -curr.vel.Y}}
		correct_y += curr.vel.Y * (1 - timeOfCollision)

		curr.SetPos(*helpers.NewMyVect2(curr.GetPos().X, correct_y))

	} else if curr.pos.X < 0 {
		timeOfCollision := prev.pos.X / (prev.pos.X - curr.pos.X)
		correct_x := prev.pos.X + curr.vel.X*timeOfCollision
		curr.vel = helpers.MyVect2{vector2.Vector2{X: -curr.vel.X, Y: curr.vel.Y}}
		correct_x += curr.vel.X * (1 - timeOfCollision)

		curr.SetPos(*helpers.NewMyVect2(correct_x, curr.GetPos().Y))

	} else if mx := float64(curr.max_x - 2*curr.radius); mx <= curr.pos.X {
		timeOfCollision := (mx - prev.pos.X) / (curr.pos.X - prev.pos.X)
		correct_x := prev.pos.X + curr.vel.X*timeOfCollision
		curr.vel = helpers.MyVect2{vector2.Vector2{X: -curr.vel.X, Y: curr.vel.Y}}
		correct_x += curr.vel.X * (1 - timeOfCollision)

		curr.SetPos(*helpers.NewMyVect2(correct_x, curr.GetPos().Y))
	}

	return curr
}

func (c Circle) CheckCollision(other Circle) bool {
	dx := c.GetMidpoint().X - other.GetMidpoint().X
	dy := c.GetMidpoint().Y - other.GetMidpoint().Y
	if dx*dx+dy*dy < math.Pow(float64(c.radius+other.radius), 2) {
		return true
	}
	return false
}

func (c *Circle) ElasticCollision(other *Circle) {
	first_result, second_result := moveOverappedCircles(*c, *other)
	tmp := first_result.CheckCollision(second_result)
	fmt.Print(tmp)
	//first_result, second_result := *c, *other

	orgPos := first_result.pos
	otherPos := second_result.pos

	orgVel := c.vel
	otherVel := other.vel

	distance := orgPos.Distance(&otherPos)
	distance_squared := distance * distance

	c.vel = *orgVel.Add(
		orgPos.Add(otherPos.MulScalar(-1.0)).MulScalar(
			orgVel.Add(otherVel.MulScalar(-1.0)).Dot(orgPos.Add(otherPos.MulScalar(-1.0)))).
			MulScalar(1.0 / distance_squared).
			MulScalar(2 * other.mass / (c.mass + other.mass)).
			MulScalar(-1.0))

	other.vel = *otherVel.Add(
		otherPos.Add(orgPos.MulScalar(-1.0)).MulScalar(
			otherVel.Add(orgVel.MulScalar(-1.0)).Dot(otherPos.Add(orgPos.MulScalar(-1.0)))).
			MulScalar(1.0 / distance_squared).
			MulScalar(2 * c.mass / (c.mass + other.mass)).
			MulScalar(-1.0))

	c.SetPos(first_result.pos)
	other.SetPos(second_result.pos)
}

func moveOverappedCircles(first Circle, second Circle) (Circle, Circle) {
	first_result := first
	second_result := second

	const threshold = 1

	first_center := first.GetMidpoint()
	second_center := second.GetMidpoint()
	midpoint := first_center.Add(&second_center).MulScalar(0.5)

	var dist_vec helpers.MyVect2
	dist_vec = *first.midpoint.Subtract(&second.midpoint)
	temp := dist_vec.MulScalar(1.0 / dist_vec.Magnitude())
	check := temp.Magnitude()
	fmt.Println(check)
	first_result.SetMidpoint(*midpoint.Add(temp.MulScalar(float64(first.radius) + threshold)))
	//first_result.pos = *midpoint.Add(temp.MulScalar(float64(first.radius)))

	dist_vec = *second.midpoint.Subtract(&first.midpoint)
	temp2 := dist_vec.MulScalar(1.0 / dist_vec.Magnitude())
	check_v2 := temp2.Magnitude()
	fmt.Println(check_v2)
	second_result.SetMidpoint(*midpoint.Add(temp2.MulScalar(float64(second.radius) + threshold)))
	//second_result.pos = *midpoint.Add(temp2.MulScalar(float64(second.radius)))

	return first_result, second_result
}

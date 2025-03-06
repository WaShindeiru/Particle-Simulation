package circle

import (
	"collision/helpers"
	"fmt"
	"github.com/deeean/go-vector/vector2"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

const (
	radius_avg = 32
	mass_avg   = 4
)

type CircleId uint

type BasicCircle struct {
	img      *ebiten.Image
	pos      helpers.MyVect2
	midpoint helpers.MyVect2
	vel      helpers.MyVect2
	max_x    int
	max_y    int
	mass     float64
	radius   int
	id       CircleId
}

func (c *BasicCircle) GetId() CircleId { return c.id }

func (c *BasicCircle) String() string {
	return fmt.Sprintf("BasicCircle{pos: %v, vel: %v, radius: %d}", c.GetPos(), c.vel, c.radius)
}

func (c BasicCircle) GetPos() helpers.MyVect2 {
	return c.pos
}

func (c BasicCircle) GetMidpoint() *helpers.MyVect2 {
	temp := c.midpoint
	return &temp
}

func (c BasicCircle) GetMax() helpers.MyVect2 {
	return *c.pos.Add(helpers.NewMyVect2(2.0*float64(c.radius), 2.0*float64(c.radius)))
}

func (c *BasicCircle) SetPos(v helpers.MyVect2) {
	c.pos = v
	c.midpoint = helpers.MyVect2{Vector2: vector2.Vector2{X: c.pos.X + float64(c.radius), Y: c.pos.Y + float64(c.radius)}}
}

func (c *BasicCircle) SetMidpoint(v helpers.MyVect2) {
	c.midpoint = v
	c.pos = *helpers.NewMyVect2(c.midpoint.X-float64(c.radius), c.midpoint.Y-float64(c.radius))
}

func (c *BasicCircle) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.pos.X), float64(c.pos.Y))
	screen.DrawImage(c.img, op)
}

func (c *BasicCircle) Update() {
	var prev BasicCircle = *c

	vx := c.vel.X
	vy := c.vel.Y

	pos := c.GetPos()
	c.SetPos(*pos.Add(helpers.NewMyVect2(vx, vy)))

	*c = handleCollisionWithBorder(prev, *c)
}

// continuous collision
// tunneling prevention
func handleCollisionWithBorder(prev BasicCircle, curr BasicCircle) BasicCircle {
	if curr.pos.Y < 0 {
		timeOfCollision := prev.GetPos().Y / (prev.GetPos().Y - curr.GetPos().Y)
		correct_y := prev.GetPos().Y + curr.vel.Y*timeOfCollision
		curr.vel = helpers.MyVect2{vector2.Vector2{X: curr.vel.X, Y: -curr.vel.Y}}
		correct_y += curr.vel.Y * (1 - timeOfCollision)

		curr.SetPos(*helpers.NewMyVect2(curr.GetPos().X, correct_y))

	} else if my := float64(curr.max_y - 2*curr.radius); my <= curr.GetPos().Y {
		timeOfCollision := (my - prev.GetPos().Y) / (curr.GetPos().Y - prev.GetPos().Y)
		correct_y := prev.GetPos().Y + curr.vel.Y*timeOfCollision
		curr.vel = helpers.MyVect2{vector2.Vector2{X: curr.vel.X, Y: -curr.vel.Y}}
		correct_y += curr.vel.Y * (1 - timeOfCollision)

		curr.SetPos(*helpers.NewMyVect2(curr.GetPos().X, correct_y))

	} else if curr.GetPos().X < 0 {
		timeOfCollision := prev.GetPos().X / (prev.GetPos().X - curr.GetPos().X)
		correct_x := prev.GetPos().X + curr.vel.X*timeOfCollision
		curr.vel = helpers.MyVect2{vector2.Vector2{X: -curr.vel.X, Y: curr.vel.Y}}
		correct_x += curr.vel.X * (1 - timeOfCollision)

		curr.SetPos(*helpers.NewMyVect2(correct_x, curr.GetPos().Y))

	} else if mx := float64(curr.max_x - 2*curr.radius); mx <= curr.GetPos().X {
		timeOfCollision := (mx - prev.GetPos().X) / (curr.GetPos().X - prev.GetPos().X)
		correct_x := prev.GetPos().X + curr.vel.X*timeOfCollision
		curr.vel = helpers.MyVect2{vector2.Vector2{X: -curr.vel.X, Y: curr.vel.Y}}
		correct_x += curr.vel.X * (1 - timeOfCollision)

		curr.SetPos(*helpers.NewMyVect2(correct_x, curr.GetPos().Y))
	}

	return curr
}

func (c *BasicCircle) GetRadius() float64 { return float64(c.radius) }

func (c *BasicCircle) CheckCollision(other Circle) bool {
	dx := c.GetMidpoint().X - other.GetMidpoint().X
	dy := c.GetMidpoint().Y - other.GetMidpoint().Y
	if dx*dx+dy*dy < math.Pow(float64(c.radius)+other.GetRadius(), 2) {
		return true
	}
	return false
}

func (c *BasicCircle) GetVel() *helpers.MyVect2 {
	temp := c.vel
	return &temp
}

func (c *BasicCircle) GetMass() float64 { return c.mass }

func (c *BasicCircle) setVel(vel_ helpers.MyVect2) {
	c.vel = vel_
}

func ElasticCollision(this, other Circle) {
	first_result, second_result := moveOverappedCircles(this, other)
	//first_result, second_result := *c, *other

	orgPos := first_result.GetPos()
	otherPos := second_result.GetPos()

	orgVel := this.GetVel()
	otherVel := other.GetVel()

	distance := orgPos.Distance(&otherPos)
	distance_squared := distance * distance

	this.setVel(*orgVel.Subtract(
		orgPos.Subtract(&otherPos).MulScalar(
			orgVel.Subtract(otherVel).Dot(orgPos.Subtract(&otherPos))).
			MulScalar(1.0 / distance_squared).
			MulScalar(2 * other.GetMass() / (this.GetMass() + other.GetMass()))))

	other.setVel(*otherVel.Subtract(
		otherPos.Subtract(&orgPos).MulScalar(
			otherVel.Subtract(orgVel).Dot(otherPos.Subtract(&orgPos))).
			MulScalar(1.0 / distance_squared).
			MulScalar(2 * this.GetMass() / (this.GetMass() + other.GetMass()))))

	this.SetPos(first_result.GetPos())
	other.SetPos(second_result.GetPos())
}

func moveOverappedCircles(first Circle, second Circle) (Circle, Circle) {
	first_result := first.Copy().(Circle)
	second_result := second.Copy().(Circle)

	const threshold = 0

	first_center := first.GetMidpoint()
	second_center := second.GetMidpoint()
	midpoint := first_center.Add(second_center).MulScalar(0.5)

	var dist_vec helpers.MyVect2
	dist_vec = *first.GetMidpoint().Subtract(second.GetMidpoint())
	temp := dist_vec.MulScalar(1.0 / dist_vec.Magnitude())
	first_result.SetMidpoint(*midpoint.Add(temp.MulScalar(float64(first.GetRadius()) + threshold)))
	//first_result.GetPos() = *midpoint.Add(temp.MulScalar(float64(first.radius)))

	dist_vec = *second.GetMidpoint().Subtract(first.GetMidpoint())
	temp2 := dist_vec.MulScalar(1.0 / dist_vec.Magnitude())
	second_result.SetMidpoint(*midpoint.Add(temp2.MulScalar(float64(second.GetRadius()) + threshold)))
	//second_result.GetPos() = *midpoint.Add(temp2.MulScalar(float64(second.radius)))

	return first_result, second_result
}

func (c *BasicCircle) Copy() helpers.Copiable {
	return &BasicCircle{
		img:      c.img,
		pos:      c.pos,
		midpoint: c.midpoint,
		vel:      c.vel,
		max_x:    c.max_x,
		max_y:    c.max_y,
		mass:     c.mass,
		radius:   c.radius,
		id:       c.id,
	}
}

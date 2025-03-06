package circle

import (
	"collision/helpers"
	"fmt"
	"github.com/deeean/go-vector/vector2"
	"math"
)

const (
	radius_avg = 32
	mass_avg   = 4
)

type CircleId uint

type CircleGeneric struct {
	pos      helpers.MyVect2
	midpoint helpers.MyVect2
	vel      helpers.MyVect2
	max_x    int
	max_y    int
	mass     float64
	radius   int
	id       CircleId
}

func (c *CircleGeneric) GetId() CircleId { return c.id }

func (c *CircleGeneric) String() string {
	return fmt.Sprintf("BasicCircle{pos: %v, vel: %v, radius: %d}", c.GetPos(), c.vel, c.radius)
}

func (c *CircleGeneric) GetPos() helpers.MyVect2 {
	return c.pos
}

func (c *CircleGeneric) GetMidpoint() *helpers.MyVect2 {
	temp := c.midpoint
	return &temp
}

func (c *CircleGeneric) GetMax() helpers.MyVect2 {
	return *c.pos.Add(helpers.NewMyVect2(2.0*float64(c.radius), 2.0*float64(c.radius)))
}

func (c *CircleGeneric) SetPos(v helpers.MyVect2) {
	c.pos = v
	c.midpoint = helpers.MyVect2{Vector2: vector2.Vector2{X: c.pos.X + float64(c.radius), Y: c.pos.Y + float64(c.radius)}}
}

func (c *CircleGeneric) SetMidpoint(v helpers.MyVect2) {
	c.midpoint = v
	c.pos = *helpers.NewMyVect2(c.midpoint.X-float64(c.radius), c.midpoint.Y-float64(c.radius))
}

func (c *CircleGeneric) GetRadius() float64 { return float64(c.radius) }

func (c *CircleGeneric) CheckCollision(other Circle) bool {
	dx := c.GetMidpoint().X - other.GetMidpoint().X
	dy := c.GetMidpoint().Y - other.GetMidpoint().Y
	if dx*dx+dy*dy < math.Pow(float64(c.radius)+other.GetRadius(), 2) {
		return true
	}
	return false
}

func (c *CircleGeneric) GetVel() *helpers.MyVect2 {
	temp := c.vel
	return &temp
}

func (c *CircleGeneric) GetMass() float64 { return c.mass }

func (c *CircleGeneric) setVel(vel_ helpers.MyVect2) {
	c.vel = vel_
}

func (c *CircleGeneric) GetLimits() *helpers.MyVect2 {
	return helpers.NewMyVect2(float64(c.max_x), float64(c.max_y))
}

// continuous collision
// tunneling prevention
func handleCollisionWithBorder(prev Circle, curr Circle) Circle {
	max_x, max_y := curr.GetLimits().X, curr.GetLimits().Y
	radius := curr.GetRadius()

	if curr.GetPos().Y < 0 {
		timeOfCollision := prev.GetPos().Y / (prev.GetPos().Y - curr.GetPos().Y)
		correct_y := prev.GetPos().Y + curr.GetVel().Y*timeOfCollision
		curr.setVel(helpers.MyVect2{vector2.Vector2{X: curr.GetVel().X, Y: -curr.GetVel().Y}})
		correct_y += curr.GetVel().Y * (1 - timeOfCollision)

		curr.SetPos(*helpers.NewMyVect2(curr.GetPos().X, correct_y))

	} else if my := float64(max_y - 2*radius); my <= curr.GetPos().Y {
		timeOfCollision := (my - prev.GetPos().Y) / (curr.GetPos().Y - prev.GetPos().Y)
		correct_y := prev.GetPos().Y + curr.GetVel().Y*timeOfCollision
		curr.setVel(helpers.MyVect2{vector2.Vector2{X: curr.GetVel().X, Y: -curr.GetVel().Y}})
		correct_y += curr.GetVel().Y * (1 - timeOfCollision)

		curr.SetPos(*helpers.NewMyVect2(curr.GetPos().X, correct_y))

	} else if curr.GetPos().X < 0 {
		timeOfCollision := prev.GetPos().X / (prev.GetPos().X - curr.GetPos().X)
		correct_x := prev.GetPos().X + curr.GetVel().X*timeOfCollision
		curr.setVel(helpers.MyVect2{vector2.Vector2{X: -curr.GetVel().X, Y: curr.GetVel().Y}})
		correct_x += curr.GetVel().X * (1 - timeOfCollision)

		curr.SetPos(*helpers.NewMyVect2(correct_x, curr.GetPos().Y))

	} else if mx := float64(max_x - 2*radius); mx <= curr.GetPos().X {
		timeOfCollision := (mx - prev.GetPos().X) / (curr.GetPos().X - prev.GetPos().X)
		correct_x := prev.GetPos().X + curr.GetVel().X*timeOfCollision
		curr.setVel(helpers.MyVect2{vector2.Vector2{X: -curr.GetVel().X, Y: curr.GetVel().Y}})
		correct_x += curr.GetVel().X * (1 - timeOfCollision)

		curr.SetPos(*helpers.NewMyVect2(correct_x, curr.GetPos().Y))
	}

	return curr
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

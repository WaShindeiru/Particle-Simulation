package circle

import (
	"collision/helpers"
	"github.com/hajimehoshi/ebiten/v2"
)

type Circle interface {
	GetId() CircleId
	GetPos() helpers.MyVect2
	GetMidpoint() *helpers.MyVect2
	GetMax() helpers.MyVect2
	SetPos(helpers.MyVect2)
	SetMidpoint(helpers.MyVect2)
	Draw(image *ebiten.Image)
	Update()
	//handleCollisionWithBorder(Circle, Circle) Circle
	CheckCollision(Circle) bool
	GetRadius() float64
	GetVel() *helpers.MyVect2
	GetMass() float64
	//ElasticCollision(Circle)
	helpers.Copiable
	setVel(helpers.MyVect2)
	GetLimits() *helpers.MyVect2
}

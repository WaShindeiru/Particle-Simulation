package circle

import (
	"collision/helpers"
	"github.com/hajimehoshi/ebiten/v2"
)

type BasicCircle struct {
	CircleGeneric
	img *ebiten.Image
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

	circle, ok := handleCollisionWithBorder(&prev, c).(*BasicCircle)
	helpers.Assert(ok, "error")

	*c = *circle
}

func (c *BasicCircle) Copy() helpers.Copiable {
	return &BasicCircle{
		img: c.img,
		CircleGeneric: CircleGeneric{
			pos:      c.pos,
			midpoint: c.midpoint,
			vel:      c.vel,
			max_x:    c.max_x,
			max_y:    c.max_y,
			mass:     c.mass,
			radius:   c.radius,
			id:       c.id,
		},
	}
}

type Type int

const (
	Paper Type = iota
	Rock
	Scissors
)

type GameCircle struct {
	CircleGeneric
	img      []*ebiten.Image
	circType Type
}

var LookUpTable = [3][3]Type{
	{Paper, Paper, Scissors},
	{Paper, Rock, Rock},
	{Scissors, Rock, Scissors},
}

func GetNewType(type1 Type, type2 Type) Type {
	return LookUpTable[type1][type2]
}

func (c *GameCircle) Update() {
	var prev = *c

	vx := c.vel.X
	vy := c.vel.Y

	pos := c.GetPos()
	c.SetPos(*pos.Add(helpers.NewMyVect2(vx, vy)))

	circle, ok := handleCollisionWithBorder(&prev, c).(*GameCircle)
	helpers.Assert(ok, "error")

	*c = *circle
}

func (c *GameCircle) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(c.radius*2)/float64(Image_Height), float64(c.radius*2)/float64(Image_Width))
	op.GeoM.Translate(float64(c.pos.X), float64(c.pos.Y))
	screen.DrawImage(c.img[c.circType], op)
}

func (c *GameCircle) Copy() helpers.Copiable {
	return &GameCircle{
		img:      c.img,
		circType: c.circType,
		CircleGeneric: CircleGeneric{
			pos:      c.pos,
			midpoint: c.midpoint,
			vel:      c.vel,
			max_x:    c.max_x,
			max_y:    c.max_y,
			mass:     c.mass,
			radius:   c.radius,
			id:       c.id,
		},
	}
}

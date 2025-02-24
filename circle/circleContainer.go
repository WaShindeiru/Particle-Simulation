package circle

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type CircleContainer struct {
	Circles []Circle
	Size    int
}

func (this *CircleContainer) Add(circle *Circle) {
	this.Circles = append(this.Circles, *circle)
	this.Size++
}

func NewCircleContainer(size_ int, screenWidth int, screenHeight int) *CircleContainer {
	temp := CircleContainer{
		Circles: make([]Circle, size_),
		Size:    size_,
	}

	for i := 0; i < size_; i++ {
		temp.Circles[i] = New_random_radius(screenWidth, screenHeight)
	}

	return &temp
}

func (c *CircleContainer) Draw(screen *ebiten.Image) {
	for i := 0; i < c.Size; i++ {
		c.Circles[i].Draw(screen)
	}
}

func (c *CircleContainer) Update() {
	for i := 0; i < c.Size; i++ {
		c.Circles[i].Update()
	}

	for i := 0; i < c.Size; i++ {
		for j := i + 1; j < c.Size; j++ {
			if c.Circles[i].CheckCollision(c.Circles[j]) {
				c.Circles[i].ElasticCollision(&c.Circles[j])
			}
		}
	}
}

func (c CircleContainer) KineticEnergy() float64 {
	var energy float64
	for _, circle := range c.Circles {
		energy += circle.mass * 0.5 * math.Pow(circle.vel.Magnitude(), 2)
	}

	return energy
}

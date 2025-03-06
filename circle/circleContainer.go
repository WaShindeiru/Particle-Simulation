package circle

import (
	"collision/helpers"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"sort"
	"sync"
)

type CircleContainer[C Circle] struct {
	Circles []C
	Size    int
	queue   chan *circleTask[C]
	context *SimulationContext[C]
}

func (this *CircleContainer[C]) Add(circle C) {
	this.Circles = append(this.Circles, circle)
	this.Size++
}

//func NewCircleContainer(size_ int, screenWidth int, screenHeight int) *CircleContainer {
//	temp := CircleContainer{
//		Circles: make([]*BasicCircle, size_),
//		Size:    size_,
//	}
//
//	factory := GetInstance()
//
//	for i := 0; i < size_; i++ {
//		temp.Circles[i] = factory.GetCircleRandomRadius(screenWidth, screenHeight)
//	}
//
//	return &temp
//}

func (c *CircleContainer[T]) Draw(screen *ebiten.Image) {
	for i := 0; i < c.Size; i++ {
		c.Circles[i].Draw(screen)
	}
}

func copySlice[T helpers.Copiable](cont []T) []T {
	tempCont := make([]T, len(cont))

	for i, val := range cont {
		tempCont[i] = val.Copy().(T)
	}

	return tempCont
}

func (c *CircleContainer[T]) Update() {
	cont := copySlice(c.Circles)

	for i := 0; i < c.Size; i++ {
		cont[i].Update()
	}

	sort.Slice(cont, func(i, j int) bool {
		return cont[i].GetPos().X < cont[j].GetPos().X
	})

	for i := 0; i < c.Size; i++ {
		circle1 := cont[i]
		for j := i + 1; j < c.Size; j++ {
			circle2 := cont[j]

			if circle2.GetPos().X > circle1.GetMax().X {
				break
			}

			if cont[i].CheckCollision(cont[j]) {
				ElasticCollision(cont[i], cont[j])
			}
		}
	}

	c.context.SetUpdatedCircles(cont)
}

func (c *CircleContainer[T]) UpdateConcurrent() {
	go Serve(c.queue)

	cont := copySlice(c.Circles)

	for i := 0; i < c.Size; i++ {
		cont[i].Update()
	}

	//for _, val := range cont {
	//	c.context.SetCircle(val)
	//}

	var wg sync.WaitGroup

	sort.Slice(cont, func(i, j int) bool {
		return cont[i].GetPos().X < cont[j].GetPos().X
	})

	c.context.SetUpdatedCircles(cont)

	for i := 0; i < c.Size; i++ {
		circle1 := cont[i]
		for j := i + 1; j < c.Size; j++ {
			circle2 := cont[j]

			if circle2.GetPos().X > circle1.GetMax().X {
				//if j == i+1 {
				//	wg.Add(1)
				//	go func() {
				//		defer wg.Done()
				//		var temp BasicCircle
				//		temp = *circle1
				//		c.context.SetCircle(&temp)
				//		temp = *circle2
				//		c.context.SetCircle(&temp)
				//	}()
				//}

				break
			}

			wg.Add(1)
			c.queue <- &circleTask[T]{
				iPos:      i,
				jPos:      j,
				container: cont,
				wg:        &wg,
				context:   c.context}
		}
	}

	close(c.queue)
	wg.Wait()
}

func (c CircleContainer[C]) KineticEnergy() float64 {
	var energy float64
	for _, circle := range c.Circles {
		energy += circle.GetMass() * 0.5 * math.Pow(circle.GetVel().Magnitude(), 2)
	}

	return energy
}

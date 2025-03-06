package circle

import "sort"

type GameCircleContainer struct {
	CircleContainer[*GameCircle]
}

func (c *GameCircleContainer) Update() {
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
				newType := GetNewType(cont[i].circType, cont[j].circType)
				cont[i].circType = newType
				cont[j].circType = newType
			}
		}
	}

	c.context.SetUpdatedCircles(cont)
}

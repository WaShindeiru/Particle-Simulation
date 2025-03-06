package circle

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GameSimulationContext struct {
	*SimulationContext[*GameCircle]
	paperCount    int
	rockCount     int
	scissorsCount int
	counter       int
}

func (context *GameSimulationContext) GetCurrentIteration() *GameCircleContainer {
	size := len(context.State)
	temp := make([]*GameCircle, 0, size)

	for _, circ := range context.State {
		temp = append(temp, circ)
	}

	return &GameCircleContainer{
		CircleContainer: CircleContainer[*GameCircle]{
			Circles: temp,
			Size:    size,
			queue:   make(chan *circleTask[*GameCircle], 1000),
			context: context.SimulationContext,
		},
	}
}

func (context *GameSimulationContext) UpdateSingle() {
	context.iteration++
	context.buffer = make(StateMap[*GameCircle])
	cont := context.GetCurrentIteration()
	cont.Update()

	for _, circ := range context.updatedCircles {
		id := circ.GetId()
		if _, exists := context.buffer[id]; !exists {
			context.buffer[id] = circ
		}
	}

	context.State = context.buffer
	context.buffer = nil

	context.counter = (context.counter + 1) % 10
	if context.counter == 0 {
		var paperTemp, rockTemp, scissorsTemp int = 0, 0, 0

		for _, val := range context.State {
			if val.circType == Paper {
				paperTemp++
			} else if val.circType == Rock {
				rockTemp++
			} else if val.circType == Scissors {
				scissorsTemp++
			}
		}

		context.paperCount = paperTemp
		context.rockCount = rockTemp
		context.scissorsCount = scissorsTemp
	}
}

func (context *GameSimulationContext) Draw(screen *ebiten.Image) {
	cont := context.GetCurrentIteration()
	cont.Draw(screen)
	//
	//fontBytes, err := os.ReadFile("path/to/your/font.ttf")
	//if err != nil {
	//	log.Fatalf("Failed to read font file: %v", err)
	//}

	//face, err := font.
	msg := fmt.Sprintf(
		"TPS: %0.2f\n"+
			"FPS: %0.2f\n"+
			"Energy: %f\n"+
			"Paper Count: %d\n"+
			"Rock Count: %d\n"+
			"Scissor Count: %d\n",
		ebiten.ActualTPS(), ebiten.ActualFPS(), cont.KineticEnergy(),
		context.paperCount, context.rockCount, context.scissorsCount)

	//op := &text.DrawOptions{}
	//
	//text.Draw(screen, msg, face.UnsafeInternal(), op)

	ebitenutil.DebugPrint(screen, msg)
}

func (context *GameSimulationContext) Update() error {
	context.UpdateSingle()
	//context.iteration++
	//context.buffer = make(StateMap[A])
	//cont := context.GetCurrentIteration()
	//cont.UpdateConcurrent()
	//
	//for _, circ := range context.updatedCircles {
	//	id := circ.GetId()
	//	if _, exists := context.buffer[id]; !exists {
	//		context.buffer[id] = circ
	//	}
	//}
	//
	//context.State = context.buffer
	//context.buffer = nil

	return nil
}

func NewContextGame(size_ int, screenWidth int, screenHeight int) *GameSimulationContext {
	state_ := make(StateMap[*GameCircle])

	factory := GetGameFactoryInstance()

	for i := 0; i < size_; i++ {
		circle := factory.GetCircleRandomRadius(screenWidth, screenHeight)
		state_[circle.GetId()] = circle
	}

	return &GameSimulationContext{
		SimulationContext: &SimulationContext[*GameCircle]{
			State:     state_,
			iteration: 1,
		},
	}
}

func NewContextEmpty(screenWidth int, screenHeight int) *SimulationContext[*BasicCircle] {
	state_ := make(StateMap[*BasicCircle])
	return &SimulationContext[*BasicCircle]{
		State:     state_,
		iteration: 1,
	}
}

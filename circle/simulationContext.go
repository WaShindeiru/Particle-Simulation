package circle

import (
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"sync"
)

type StateMap[T Circle] map[CircleId]T

type SimulationContext[A Circle] struct {
	State          StateMap[A]
	buffer         StateMap[A]
	updatedCircles []A
	lock           sync.RWMutex
	iteration      uint
}

//func (context *SimulationContext) BeginNewIteration() {
//	context.State[0] = context.State[1]
//	context.State[1] = make(StateMap)
//	context.iteration++
//}

func (context *SimulationContext[A]) SetUpdatedCircles(updated []A) {
	context.updatedCircles = updated
}

func (context *SimulationContext[A]) Draw(screen *ebiten.Image) {
	cont := context.GetCurrentIteration()
	cont.Draw(screen)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\nEnergy: %f\n", ebiten.ActualTPS(),
		ebiten.ActualFPS(), cont.KineticEnergy()))
}

func (context *SimulationContext[A]) UpdateSingle() {
	context.iteration++
	context.buffer = make(StateMap[A])
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
}

func (context *SimulationContext[A]) Update() error {
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

func (context *SimulationContext[A]) Exists(id CircleId) bool {
	context.lock.RLock()
	defer context.lock.RUnlock()
	_, exists := context.buffer[id]
	return exists
}

func (context *SimulationContext[A]) GetCircle(id CircleId) (A, error) {
	context.lock.RLock()
	defer context.lock.RUnlock()
	value, exists := context.buffer[id]

	if exists {
		return value, nil
	}

	var zero A
	return zero, errors.New("invalid key")
}

//func (context *SimulationContext) SetCircle(circ *BasicCircle) {
//	context.lock.Lock()
//	defer context.lock.Unlock()
//	id := circ.GetId()
//	context.buffer[id] = circ
//}

func (context *SimulationContext[T]) SetCircle(circ T) {
	context.lock.Lock()
	defer context.lock.Unlock()
	id := circ.GetId()

	found, exists := context.buffer[id]

	if exists {
		vel := found.GetVel()
		to_add := circ.GetVel()
		found.setVel(*vel.Add(to_add))
	} else {
		context.buffer[id] = circ
	}
}

func (context *SimulationContext[T]) GetCurrentIteration() *CircleContainer[T] {
	size := len(context.State)
	temp := make([]T, 0, size)

	for _, circ := range context.State {
		temp = append(temp, circ)
	}

	return &CircleContainer[T]{
		Circles: temp,
		Size:    size,
		queue:   make(chan *circleTask[T], 1000),
		context: context,
	}
}

//func New() *SimulationContext {
//	state_ := make([]StateMap, 2)
//	state_[0] = make(StateMap)
//	state_[1] = make(StateMap)
//
//	return &SimulationContext{State: state_, iteration: 1}
//}

func NewContext(size_ int, screenWidth int, screenHeight int) *SimulationContext[*BasicCircle] {
	state_ := make(StateMap[*BasicCircle])

	factory := GetInstance()

	for i := 0; i < size_; i++ {
		circle := factory.GetCircleRandomRadius(screenWidth, screenHeight)
		state_[circle.GetId()] = circle
	}

	return &SimulationContext[*BasicCircle]{
		State:     state_,
		iteration: 1,
	}
}

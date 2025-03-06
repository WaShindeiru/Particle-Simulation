package circle

import (
	"sync"
)

type circleTask[C Circle] struct {
	iPos      int
	jPos      int
	container []C
	wg        *sync.WaitGroup
	context   *SimulationContext[C]
}

func work[C Circle](sem <-chan int, task *circleTask[C]) {
	defer task.wg.Done()

	circ1, circ2 := task.container[task.iPos], task.container[task.jPos]
	circ1_, _ := circ1.Copy().(C)
	circ2_, _ := circ2.Copy().(C)

	if circ1_.CheckCollision(circ2_) {
		ElasticCollision(circ1_, circ2_)

		task.context.SetCircle(circ1_)
		task.context.SetCircle(circ2_)
	}

	<-sem
}

func Serve[C Circle](queue <-chan *circleTask[C]) {
	//numCores := runtime.NumCPU()
	numCores := 1
	var sem = make(chan int, numCores)

	for task := range queue {
		sem <- 1
		work(sem, task)
	}

	//for {
	//	sem <- 1
	//	task := <-queue
	//	work(sem, task)
	//}
}

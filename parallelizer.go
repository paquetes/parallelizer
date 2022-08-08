package parallelizer

import (
	"runtime"
	"sync"

	"github.com/paquetes/disaster"
)

// Group goroutine pool
type Group struct {
	queue chan int
	wg    *sync.WaitGroup
}

// Add task
func (g *Group) Add(delta int) {
	for i := 0; i < delta; i++ {
		g.queue <- 1
	}

	for i := 0; i > delta; i-- {
		<-g.queue
	}

	g.wg.Add(delta)
}

// Done
func (g *Group) Done() {
	<-g.queue
	g.wg.Done()
}

// Wait task complete
func (g *Group) Wait() {
	g.wg.Wait()
}

// Close
func (g *Group) Close() {
	close(g.queue)
}

// New goroutine pool
func New(size int) *Group {
	if size <= 0 {
		size = 1
	}

	return &Group{
		queue: make(chan int, size),
		wg:    &sync.WaitGroup{},
	}
}

// Parallelize
func Parallelize(maxWorkers int, callback func(pg *Group)) {
	defer disaster.Catch()

	if maxWorkers == 0 {
		maxWorkers = runtime.NumCPU()
	}

	pg := New(maxWorkers)

	callback(pg)

	pg.Wait()
}

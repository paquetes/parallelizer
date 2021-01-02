package parallelizer

import (
	"runtime"
	"sync"

	"github.com/paquetes/disaster"
)

// Group goroutine池
type Group struct {
	queue chan int
	wg    *sync.WaitGroup
}

// Add 添加任务
func (g *Group) Add(delta int) {
	for i := 0; i < delta; i++ {
		g.queue <- 1
	}

	for i := 0; i > delta; i-- {
		<-g.queue
	}

	g.wg.Add(delta)
}

// Done 任务完成
func (g *Group) Done() {
	<-g.queue
	g.wg.Done()
}

// Wait 等待任务完成
func (g *Group) Wait() {
	g.wg.Wait()
}

// Close 关闭队列
func (g *Group) Close() {
	close(g.queue)
}

// New 创建goroutine池
func New(size int) *Group {
	if size <= 0 {
		size = 1
	}

	return &Group{
		queue: make(chan int, size),
		wg:    &sync.WaitGroup{},
	}
}

// Parallelize 并发执行给定任务
func Parallelize(maxWorkers int, callback func(pg *Group)) {
	defer disaster.Catch()

	if maxWorkers == 0 {
		maxWorkers = runtime.NumCPU()
	}

	pg := New(maxWorkers)

	callback(pg)

	pg.Wait()
}

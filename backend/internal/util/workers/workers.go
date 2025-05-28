package workers

import (
	"context"
	"log"
	"sync"
)

type WorkerPool[T any] struct {
	Wg            *sync.WaitGroup
	Channel       chan T
	WorkerHandler func(task T)
}

func NewPool[T any](channelSize int, withWaitGroup bool, handler func(task T)) *WorkerPool[T] {
	var wg *sync.WaitGroup
	if withWaitGroup {
		wg = &sync.WaitGroup{}
	}
	return &WorkerPool[T]{
		Wg:            wg,
		Channel:       make(chan T, channelSize),
		WorkerHandler: handler,
	}
}

func (wp *WorkerPool[T]) Launch(workerCount int, ctx context.Context) {
	for i := range workerCount {
		go func(id int) {
			log.Printf("👷 Worker %d started", id)
			for {
				select {
				case <-ctx.Done():
					log.Printf("👷 Worker %d shutting down", id)
					return
				case task := <-wp.Channel:
					wp.WorkerHandler(task)
					if wp.Wg != nil {
						wp.Wg.Done()
					}
				}
			}
		}(i)
	}
}

func (wp *WorkerPool[T]) Dispach(task T) {
	if wp.Wg != nil {
		wp.Wg.Add(1)
	}
	wp.Channel <- task
}

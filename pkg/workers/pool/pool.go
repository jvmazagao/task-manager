package pool

import (
	"context"
	"log"
	registry "src/pkg/processors/registry"
	"src/pkg/queue"
	"src/pkg/workers"
	"sync"
)

type WorkerPool struct {
	workers  []*workers.Worker
	queue    *queue.TaskQueue
	registry *registry.Registry
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewWorkerPool(queue *queue.TaskQueue, registry *registry.Registry, size int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workers:  make([]*workers.Worker, size),
		queue:    queue,
		registry: registry,
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (wp *WorkerPool) Start() {
	log.Println("Starting pool")
	for i := 0; i < len(wp.workers); i++ {
		worker := workers.NewWorker(i, wp.queue, wp.registry, wp.ctx)
		wp.wg.Add(1)
		wp.workers[i] = &worker

		go func(w *workers.Worker) {
			defer wp.wg.Done()
			w.Start()
		}(wp.workers[i])
	}
	log.Printf("Started %d workers\n", len(wp.workers))
}

func (wp *WorkerPool) Stop() {
	log.Println("Initiating graceful shutdown...")
	wp.cancel()
	wp.wg.Wait()
	log.Println("All workers stopped")
}

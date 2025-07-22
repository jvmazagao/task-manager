package pool

import (
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
}

func NewWorkerPool(queue *queue.TaskQueue, registry *registry.Registry, size int) *WorkerPool {
	return &WorkerPool{
		workers:  make([]*workers.Worker, size),
		queue:    queue,
		registry: registry,
	}
}

func (wp *WorkerPool) Start() {
	log.Println("Starting pool")
	for i := 0; i < len(wp.workers); i++ {
		worker := workers.NewWorker(i, wp.queue, wp.registry)
		wp.wg.Add(1)
		wp.workers[i] = &worker

		go func(w *workers.Worker) {
			defer wp.wg.Done()
			w.Start()
		}(&worker)
	}
	log.Printf("Started %d workers\n", len(wp.workers))
}

func (wp *WorkerPool) GracefulShutdown() {
	log.Println("Initiating graceful shutdown...")

	wp.queue.Close()

	for _, worker := range wp.workers {
		worker.Stop()
	}

	wp.wg.Wait()

	log.Println("All workers stopped")
}

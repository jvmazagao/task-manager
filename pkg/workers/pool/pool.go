package pool

import (
	"fmt"
	registry "src/pkg/processors/registry"
	"src/pkg/queue"
	"src/pkg/workers"
	"sync"
	"sync/atomic"
)

type WorkerPool struct {
	workers   []*workers.Worker
	queue     *queue.Queue
	registry  *registry.Registry
	workerWG  sync.WaitGroup
	startupWG sync.WaitGroup

	closeCh   chan struct{}
	isRunning atomic.Bool
}

func NewWorkerPool(q *queue.Queue, reg *registry.Registry, numWorkers int) WorkerPool {
	var ws []*workers.Worker
	for i := range numWorkers {
		work := workers.NewWorker(i, q, reg)
		ws = append(ws, &work)
	}

	return WorkerPool{
		workers:   ws,
		queue:     q,
		registry:  reg,
		workerWG:  sync.WaitGroup{},
		startupWG: sync.WaitGroup{},
		closeCh:   make(chan struct{}),
		isRunning: atomic.Bool{},
	}
}

func (wp *WorkerPool) Start() error {
	wp.isRunning.Store(true)

	fmt.Printf("Starting %d workers...\n", len(wp.workers))
	wp.startupWG.Add(len(wp.workers))

	for _, worker := range wp.workers {
		wp.workerWG.Add(1)
		go func(w *workers.Worker) {
			defer wp.workerWG.Done()

			fmt.Printf("Worker goroutine started, signaling ready...\n")
			wp.startupWG.Done()

			w.Start()
		}(worker)
	}

	fmt.Println("Waiting for all workers to be ready...")
	wp.startupWG.Wait()
	fmt.Println("All workers ready!")

	return nil
}

func (wp *WorkerPool) GracefulShutdown() error {
	wp.queue.WaitForTasks()

	if !wp.isRunning.Load() {
		return fmt.Errorf("WorkerPool not running")
	}

	wp.queue.StartShutdown()

	for _, worker := range wp.workers {
		worker.Stop()
	}

	wp.workerWG.Wait()

	wp.isRunning.Store(false)
	fmt.Println("Graceful shutdown complete!")
	return nil
}

func (wp *WorkerPool) IsRunning() bool {
	return wp.isRunning.Load()
}

func (wp *WorkerPool) WorkerCount() int {
	return len(wp.workers)
}

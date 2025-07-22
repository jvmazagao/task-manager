package workers

import (
	"context"
	"log"

	"src/pkg/model"
	registry "src/pkg/processors/registry"
	taskQueue "src/pkg/queue"
	"time"
)

type Worker struct {
	id       int
	queue    *taskQueue.TaskQueue
	registry *registry.Registry
	ctx      context.Context
}

func NewWorker(id int, queue *taskQueue.TaskQueue, registry *registry.Registry, ctx context.Context) Worker {
	return Worker{
		id:       id,
		queue:    queue,
		registry: registry,
		ctx:      ctx,
	}
}
func (w *Worker) Start() {
	w.run()
}

func (w *Worker) run() {
	log.Printf("Worker %d: run() started\n", w.id)

	for {
		select {
		case <-w.ctx.Done():
			log.Printf("Worker %d: context cancelled\n", w.id)
			return
		default:
			log.Printf("Worker %d: waiting for task\n", w.id)
			task, err := w.queue.GetTask()
			if err != nil {
				log.Printf("Worker %d: error getting task: %v\n", w.id, err)
				return
			}
			log.Printf("Worker %d: got task %s\n", w.id, task.GetID())
			w.processTask(*task)
		}
	}
}

func (w *Worker) processTask(task model.PriorityTask) {
	log.Printf("Worker %d processing task %s (priority: %d)\n",
		w.id, task.ID, task.GetPriority())

	processor, err := w.registry.GetProcessor(task.GetTaskType())
	if err != nil {
		log.Printf("Worker %d: No processor for task type %s\n", w.id, task.GetTaskType())
		return
	}

	done := make(chan error, 1)
	go func() {
		done <- processor.ProcessTask(task)
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Printf("Worker %d: Task %s failed: %v\n", w.id, task.GetID(), err)
		} else {
			log.Printf("Worker %d: Task %s completed\n", w.id, task.GetID())
		}
	case <-time.After(5 * time.Second):
		log.Printf("Worker %d: Task %s timed out\n", w.id, task.GetID())
	}
}

func (w *Worker) Stop() {
	w.ctx.Done()
}

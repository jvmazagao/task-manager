package workers

import (
	"fmt"

	"src/pkg/model"
	registry "src/pkg/processors/registry"
	taskQueue "src/pkg/queue"
	"time"
)

type Worker struct {
	id       int
	queue    *taskQueue.Queue
	registry *registry.Registry
	stopCh   chan struct{}
}

func NewWorker(id int, queue *taskQueue.Queue, registry *registry.Registry) Worker {
	return Worker{
		id:       id,
		queue:    queue,
		registry: registry,
		stopCh:   make(chan struct{}),
	}
}

func (w *Worker) Start() {
	fmt.Printf("Starting worker %d...\n", w.id)

	for {
		select {
		case <-w.stopCh:
			fmt.Printf("Stop worker %d...\n", w.id)
			return
		default:
			if task, hasTask := w.queue.GetTask(); hasTask {
				w.processTask(task)
			} else {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

func (w *Worker) Stop() {
	close(w.stopCh)
}

func (w *Worker) processTask(t any) {
	defer w.queue.TaskCompleted()

	task, ok := t.(model.TaskData)

	if !ok {
		fmt.Printf("Worker: %d. Failed to process %T\n", w.id, t)
		return
	}

	taskType := task.GetTaskType()
	processor, err := w.registry.GetProcessor(taskType)

	if err != nil {
		fmt.Printf("Worker: %d. No processor for the task type %s\n", w.id, taskType)
		return
	}

	if err := processor.ProcessTask(task); err != nil {
		fmt.Printf("Worker %d. Failed to process task: %v\n", w.id, err)
	}
}

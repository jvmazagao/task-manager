package queue

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Queue struct {
	tasks          chan any
	activeTasks    sync.WaitGroup
	isShuttingDown atomic.Bool
	lock           sync.Mutex
}

func NewTaskQueue(bufferSize int) Queue {
	return Queue{tasks: make(chan any, bufferSize)}
}

func (q *Queue) AddTask(task any) error {
	if q.isShuttingDown.Load() {
		return fmt.Errorf("Queue is shutting down. No accepting more tasks")
	}
	select {
	case q.tasks <- task:
		q.activeTasks.Add(1)
		return nil
	default:
		return fmt.Errorf("queue is full")
	}
}

func (q *Queue) GetTask() (any, bool) {
	select {
	case task := <-q.tasks:
		return task, true
	default:
		return nil, false
	}
}

func (q *Queue) TaskCompleted() {
	q.activeTasks.Done()
}

func (q *Queue) StartShutdown() {
	q.isShuttingDown.Store(true)
}

func (q *Queue) WaitForTasks() {
	q.activeTasks.Wait()
}

func (q *Queue) IsEmpty() bool {
	return len(q.tasks) == 0
}

func (q *Queue) RemainingTasks() int {
	return len(q.tasks)
}

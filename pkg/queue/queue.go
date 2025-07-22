package queue

import (
	"errors"
	"src/pkg/model"
)

type TaskQueue struct {
	pq       *SafePriorityQueue
	capacity int
}

func New(capacity int) TaskQueue {
	return TaskQueue{
		pq:       NewSafePriorityQueue(),
		capacity: 10,
	}
}

func (tq *TaskQueue) Close() {
	tq.pq.Close()
}

func (tq *TaskQueue) AddTask(task interface{}) error {

	if tq.pq.items.Len() >= tq.capacity {
		return errors.New("Task queue on max capacity")
	}

	priorityTask, ok := any(task).(*model.PriorityTask)
	if !ok {
		taskVal, ok := any(task).(*model.Task)
		if !ok {
			return errors.New("Task not supported")
		}
		priorityTask = &model.PriorityTask{
			Task:     *taskVal,
			Priority: 1000,
		}
	}

	tq.pq.Push(priorityTask)
	return nil
}

func (tq *TaskQueue) GetTask() (*model.PriorityTask, error) {
	return tq.pq.Pop()
}

package model

import "time"

type Task struct {
	ID        string
	Type      string
	CreatedAt time.Time
}

type PriorityTask struct {
	Task
	Priority int
	Index    int
}

type TaskData interface {
	GetID() string
	GetTaskType() string
	GetPriority() int
}

func (t *PriorityTask) GetID() string {
	return t.ID
}

func (t *PriorityTask) GetPriority() int {
	return t.Priority
}

func (t *PriorityTask) GetTaskType() string {
	return t.Type
}

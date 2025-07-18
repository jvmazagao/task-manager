package model

import "time"

type Task struct {
	ID        string
	Type      string
	Priority  int
	CreatedAt time.Time
}

type TaskData interface {
	GetTaskType() string
}

package processors

import (
	"fmt"
	"src/pkg/model"
)

type TaskProcessor interface {
	CanProcess(taskType string) bool
	ProcessTask(t any) error
}

type DefaulProcessor struct {
}

func (dp DefaulProcessor) ProcessTask(t any) error {
	task, ok := t.(model.PriorityTask)

	if !ok {
		return fmt.Errorf("expected email task, got %T", t)
	}

	fmt.Printf("Processed email task: ID=%v\n", task.ID)

	return nil
}

func (dp DefaulProcessor) CanProcess(taskType string) bool {
	return true
}

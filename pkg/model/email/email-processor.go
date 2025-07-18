package email

import (
	"fmt"
)

type EmailProcessor struct {
}

func (ep EmailProcessor) ProcessTask(t any) error {
	task, ok := t.(EmailTask)

	if !ok {
		return fmt.Errorf("expected email task, got %T", t)
	}

	if !ep.CanProcess(task.GetTaskType()) {
		return fmt.Errorf("send_mail processor expected, got %T", task.GetTaskType())
	}

	fmt.Printf("Processed email task: ID=%v, Receiver=%v, Sender=%v, Subject=%v\n", task.ID, task.Receiver, task.Sender, task.Subject)

	return nil
}

func (ep EmailProcessor) CanProcess(taskType string) bool {
	if taskType == "send_email" {
		return true
	}

	return false
}

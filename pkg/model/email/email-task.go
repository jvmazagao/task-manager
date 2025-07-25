package email

import "src/pkg/model"

type EmailTask struct {
	model.PriorityTask
	Receiver string
	Sender   string
	Subject  string
	Content  map[string]string
}

func (t EmailTask) GetTaskType() string {
	return "send_email"
}

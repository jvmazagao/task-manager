package processors

type TaskProcessor interface {
	CanProcess(taskType string) bool
	ProcessTask(t any) error
}

package main

import (
	"src/pkg/model"
	"src/pkg/model/email"
	registry "src/pkg/processors/registry"
	taskQueue "src/pkg/queue"
	"src/pkg/workers/pool"

	"fmt"
)

func main() {
	queue := taskQueue.NewTaskQueue(10)
	registry := registry.NewRegistry()
	registry.AddProcessor("send_email", &email.EmailProcessor{})

	wp := pool.NewWorkerPool(&queue, &registry, 3)

	wp.Start()

	for i := 0; i < 10; i++ {
		task := email.EmailTask{
			Task: model.Task{
				ID:   fmt.Sprintf("email-%d", i),
				Type: "send_email",
			},
			Receiver: fmt.Sprintf("user%d@example.com", i),
		}
		queue.AddTask(task)
	}

	wp.GracefulShutdown()
}

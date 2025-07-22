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
	queue := taskQueue.New(10)
	registry := registry.NewRegistry()
	registry.AddProcessor("send_email", &email.EmailProcessor{})

	wp := pool.NewWorkerPool(&queue, &registry, 3)

	wp.Start()

	for i := 0; i < 10; i++ {
		task := &model.PriorityTask{
			Task: model.Task{
				ID:   fmt.Sprintf("email-%d", i),
				Type: "any",
			},
			Priority: i % 3,
		}
		queue.AddTask(task)
	}

	wp.Stop()
}

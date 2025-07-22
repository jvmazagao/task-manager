package processors

import (
	"fmt"
	"src/pkg/processors"
)

type Registry struct {
	processors map[string]processors.TaskProcessor
}

func NewRegistry() Registry {
	return Registry{processors: make(map[string]processors.TaskProcessor)}
}

func (r *Registry) AddProcessor(processorType string, p processors.TaskProcessor) {
	r.processors[processorType] = p
}

func (r *Registry) GetProcessor(processorType string) (processors.TaskProcessor, error) {

	if processor, ok := r.processors[processorType]; ok {
		return processor, nil
	}

	fmt.Println("Processing with default processor")

	return processors.DefaulProcessor{}, nil
}

package queue

import (
	"container/heap"
	"errors"
	"log"
	"src/pkg/model"
	"sync"
)

type SafePriorityQueue struct {
	mu     sync.Mutex
	items  PriorityQueue
	cond   *sync.Cond
	closed bool
}

func NewSafePriorityQueue() *SafePriorityQueue {
	spq := SafePriorityQueue{
		items: make(PriorityQueue, 0),
	}
	spq.cond = sync.NewCond(&spq.mu)
	heap.Init(&spq.items)
	return &spq
}

func (spq *SafePriorityQueue) Push(task *model.PriorityTask) {
	spq.mu.Lock()
	defer spq.mu.Unlock()

	heap.Push(&spq.items, task)

	spq.cond.Signal()
}

func (spq *SafePriorityQueue) Pop() (*model.PriorityTask, error) {
	spq.mu.Lock()
	defer spq.mu.Unlock()

	for spq.items.Len() == 0 && !spq.closed {
		log.Println("Pop: waiting for items...")
		spq.cond.Wait()
	}

	if spq.closed && spq.items.Len() == 0 {
		log.Println("Pop: queue closed and empty")
		return nil, errors.New("queue closed")
	}

	item := heap.Pop(&spq.items).(*model.PriorityTask)
	return item, nil
}

func (spq *SafePriorityQueue) Close() {
	spq.mu.Lock()
	defer spq.mu.Unlock()
	spq.closed = true
	spq.cond.Broadcast() // Wake all waiting workers
}

func (spq *SafePriorityQueue) Size() int {
	spq.mu.Lock()
	defer spq.mu.Unlock()

	return spq.items.Len()
}

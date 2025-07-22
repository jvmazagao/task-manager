package queue

import "src/pkg/model"

type PriorityQueue []*model.PriorityTask

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]

	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x any) {
	item := x.(*model.PriorityTask)

	n := len(*pq)

	item.Index = n

	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)

	item := old[n-1]

	old[n-1] = nil
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

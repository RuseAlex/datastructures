package priorityQueue

import (
	"datastructures/queue"
)

type PriorityQueue[T any] struct {
	q    []queue.Queue[T] //slice of queue
	size int
}

func New[T any](numPriorities int) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		q: make([]queue.Queue[T], numPriorities),
	}
}

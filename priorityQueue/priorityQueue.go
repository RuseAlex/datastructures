package priorityQueue

import (
	"github.com/RuseAlex/datastructures/queue"
)

type PriorityQueue[T any] struct {
	q    []queue.Queue[T] //slice of queue
	size int
}

// New create a new priority queue based on the number of priorities
func New[T any](numPriorities int) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		q: make([]queue.Queue[T], numPriorities),
	}
}

// Push adds a new queue to the priority queue
func (pq *PriorityQueue[T]) Push(elem T, priority int) {
	pq.q[priority-1].Push(elem)
	pq.size++
}

// Remove return an item from the queue before removing it
func (pq *PriorityQueue[T]) Remove() *T {
	pq.size--
	for i := 0; i < len(pq.q); i++ {
		if pq.q[i].Len() > 0 {
			return pq.q[i].Remove()
		}
	}

	return nil
}

// First returns the first item in the queue
func (pq *PriorityQueue[T]) First() *T {
	for _, que := range pq.q {
		if que.Len() > 0 {
			return que.First()
		}
	}

	return nil
}

// IsEmpty checks to see if the priority queue is empty
func (pq *PriorityQueue[T]) IsEmpty() bool {
	result := true
	for _, que := range pq.q {
		if que.Len() > 0 {
			result = false
			break
		}
	}

	return result
}

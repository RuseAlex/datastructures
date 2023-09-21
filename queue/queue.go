package queue

type Queue[T any] struct {
	items []T
}

// New create a new queue
func New[T any]() *Queue[T] {
	return &Queue[T]{}
}

// Push add a new element to the queue
func (q *Queue[T]) Push(elem T) {
	q.items = append(q.items, elem)
}

// Remove returns the first item in the queue before removing it
func (q *Queue[T]) Remove() *T {
	if len(q.items) <= 0 {
		return nil
	}

	removedVal := q.items[0]
	q.items = q.items[1:]
	return &removedVal
}

// Len returns the length of the queue
func (q *Queue[T]) Len() int {
	return len(q.items)
}

// Iterator returns a pointer to an iterator we can use to traverse the queue
func (q *Queue[T]) Iterator() *Iterator[T] {
	return &Iterator[T]{
		curPos: 0,
		items:  q.items,
	}
}

// Iterator is used to iterate over the queue's items
type Iterator[T any] struct {
	curPos int
	items  []T
}

// HasNext checks to see if there is an item to move forward to
func (i *Iterator[T]) HasNext() bool {
	if i.curPos < (len(i.items) - 1) {
		return true
	}

	return false
}

// Next returns the current item and then moves forward to the next
func (i *Iterator[T]) Next() *T {
	curItem := i.items[i.curPos]
	i.curPos++
	return &curItem
}

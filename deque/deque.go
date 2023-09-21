package deque

type Deque[T any] struct {
	items []T
}

// New creates a new Dequeue
func New[T any]() *Deque[T] {
	return &Deque[T]{}
}

// InsertFront inserts a new element to the front of the dequeue
func (d *Deque[T]) InsertFront(elem T) {
	d.items = append(d.items, elem)
	for i := len(d.items) - 1; i > 0; i-- {
		d.items[i] = d.items[i-1]
	}

	d.items[0] = elem
}

// InsertBack inserts a new element in the back of the dequeue
func (d *Deque[T]) InsertBack(elem T) {
	d.items = append(d.items, elem)
}

// First returns a pointer to the value of the first item in the dequeue
func (d *Deque[T]) First() *T {
	if len(d.items) <= 0 {
		return nil
	}

	return &d.items[0]
}

// Last returns a pointer to the value of the last item in the dequeue
func (d *Deque[T]) Last() *T {
	if len(d.items) <= 0 {
		return nil
	}

	return &d.items[len(d.items)-1]
}

// RemoveFirst removes the first item in the dequeue
func (d *Deque[T]) RemoveFirst() *T {
	if len(d.items) <= 0 {
		return nil
	}

	returnItem := d.items[0]
	d.items = d.items[1:]
	return &returnItem
}

// RemoveLast removes the last item in the deque
func (d *Deque[T]) RemoveLast() *T {
	length := len(d.items)
	if length <= 0 {
		return nil
	}

	returnItem := d.items[length-1]
	d.items = d.items[:(length - 1)]
	return &returnItem
}

// IsEmpty checks to see if the dequeue is empty
func (d *Deque[T]) IsEmpty() bool {
	return len(d.items) == 0
}

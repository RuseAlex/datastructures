package deque

type Deque[T any] struct {
	items []T
}

func New[T any]() *Deque[T] {
	return &Deque[T]{}
}

func (d *Deque[T]) InsertFront(elem T) {
	d.items = append(d.items, elem)
	for i := len(d.items) - 1; i > 0; i-- {
		d.items[i] = d.items[i-1]
	}

	d.items[0] = elem
}

func (d *Deque[T]) InsertBack(elem T) {
	d.items = append(d.items, elem)
}

func (d *Deque[T]) First() *T {
	if len(d.items) <= 0 {
		return nil
	}

	return &d.items[0]
}

func (d *Deque[T]) Last() *T {
	if len(d.items) <= 0 {
		return nil
	}

	return &d.items[len(d.items)-1]
}

func (d *Deque[T]) RemoveFirst() *T {
	if len(d.items) <= 0 {
		return nil
	}

	returnItem := d.items[0]
	d.items = d.items[1:]
	return &returnItem
}

func (d *Deque[T]) RemoveLast() *T {
	length := len(d.items)
	if length <= 0 {
		return nil
	}

	returnItem := d.items[length-1]
	d.items = d.items[:(length - 1)]
	return &returnItem
}

func (d *Deque[T]) Empty() bool {
	return len(d.items) == 0
}

package listof

type ListOf[T any] struct {
	items []T
}

// New creates a new empty list
func New[T any]() *ListOf[T] {
	return &ListOf[T]{}
}

// NewFromSlice creates a new list from a slice
func NewFromSlice[T any](elems ...T) *ListOf[T] {
	items := make([]T, len(elems))
	items = append(items, elems...)
	return &ListOf[T]{
		items: items,
	}
}

// NewWithSize creates a new list of a certain size
func NewWithSize[T any](size int) *ListOf[T] {
	if size < 0 {
		return nil
	}

	return &ListOf[T]{
		items: make([]T, size),
	}
}

// Get returns the item at index idx
func (l *ListOf[T]) Get(idx int) *T {
	if idx < 0 || idx > len(l.items) {
		return nil
	}

	return &l.items[idx]
}

// Set the value of an item in the list at index idx
func (l *ListOf[T]) Set(idx int, val T) {
	if idx < 0 || idx > len(l.items) {
		return
	}

	l.items[idx] = val
}

// Remove an item from the list
func (l *ListOf[T]) Remove(idx int) *T {
	if idx < 0 || idx > len(l.items) {
		return nil
	}

	popped := l.items[idx]
	l.items = append(l.items[:idx], l.items[idx+1:]...)
	return &popped
}

// Add a new element to the list
func (l *ListOf[T]) Add(elem T) {
	l.items = append(l.items, elem)
}

// AddSlice to the list
func (l *ListOf[T]) AddSlice(elems ...T) {
	l.items = append(l.items, elems...)
}

// AddList adds another list to the original one
func (l *ListOf[T]) AddList(list ListOf[T]) {
	l.items = append(l.items, list.items...)
}

// Clear return all elements in the list to their default/zero value
func (l *ListOf[T]) Clear() {
	l.items = make([]T, len(l.items))
}

// Len returns the length of the list
func (l *ListOf[T]) Len() int {
	return len(l.items)
}

// Copy returns a copy of the list
func (l *ListOf[T]) Copy() *ListOf[T] {
	newStack := ListOf[T]{
		items: l.items,
	}

	return &newStack
}

// Contains takes a value of T and equal func to check if there is an element in the list that is equal to val
func (l *ListOf[T]) Contains(val T, eq func(T, T) bool) bool {
	for _, elem := range l.items {
		if eq(val, elem) == true {
			return true
		}
	}

	return false
}

// Map takes a function that maps each element of the list to a new one
func (l *ListOf[T]) Map(mapFun func(T) T) *ListOf[T] {
	newList := ListOf[T]{
		items: make([]T, len(l.items)),
	}

	for idx, elem := range l.items {
		newList.items[idx] = mapFun(elem)
	}

	return &newList
}

// Filter applies a filter function on the list
func (l *ListOf[T]) Filter(filter func(T) bool) *ListOf[T] {
	newList := ListOf[T]{
		items: make([]T, len(l.items)),
	}

	for idx, elem := range l.items {
		if filter(elem) == true {
			newList.items[idx] = elem
		}
	}

	return &newList
}

// Reduce applies a reduce function on the list
func (l *ListOf[T]) Reduce(iniVal T, filter func(T, T) T) T {
	res := iniVal
	for _, elem := range l.items {
		res = filter(res, elem)
	}

	return res
}

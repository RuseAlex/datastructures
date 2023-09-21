package stack

type Stack[T any] struct {
	items []T
}

// New create a new stack and returns a pointer to it
func New[T any]() *Stack[T] {
	return &Stack[T]{}
}

// Push adds a new item on the top of the stack
func (s *Stack[T]) Push(elem T) {
	s.items = append(s.items, elem)
}

// Pop return the item on the top of the stack before removing it
func (s *Stack[T]) Pop() *T {
	length := len(s.items)
	if length <= 0 {
		return nil
	}

	popped := s.items[length-1]
	s.items = s.items[:(length - 1)]
	return &popped
}

// Top return a pointer to the item on the top of the stack
func (s *Stack[T]) Top() *T {
	if len(s.items) <= 0 {
		return nil
	}

	return &s.items[len(s.items)-1]
}

// Len returns the length of the stack
func (s *Stack[T]) Len() int {
	return len(s.items)
}

// Copy return a pointer to a copy of the stack
func (s *Stack[T]) Copy() *Stack[T] {
	newStack := Stack[T]{
		items: s.items,
	}

	return &newStack
}

package stack

type Stack[T any] struct {
	items []T
}

func New[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Push(elem T) {
	s.items = append(s.items, elem)
}

func (s *Stack[T]) Pop() *T {
	length := len(s.items)
	if length <= 0 {
		return nil
	}

	popped := s.items[length-1]
	s.items = s.items[:(length - 1)]
	return &popped
}

func (s *Stack[T]) Top() *T {
	if len(s.items) <= 0 {
		return nil
	}

	return &s.items[len(s.items)-1]
}

func (s *Stack[T]) Len() int {
	return len(s.items)
}

func (s *Stack[T]) Copy() *Stack[T] {
	newStack := Stack[T]{
		items: s.items,
	}

	return &newStack
}

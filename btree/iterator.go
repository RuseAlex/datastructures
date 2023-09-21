package btree

type Iterator[K, V any] struct {
	btree *BTree[K, V]
	node  *Node[K, V]
	entry *Entry[K, V]
	pos   int
}

const (
	START = iota
	BETWEEN
	END
)

func (b *BTree[K, V]) Iterator() Iterator[K, V] {
	return Iterator[K, V]{btree: b, node: nil, pos: START}
}

// Next moves the iterator to the next element and returns true
// if there is a next element in the tree
func (i *Iterator[K, V]) Next() bool {
	if i.pos == END {
		goto end
	}

	if i.pos == START {
		left := i.btree.Left()
		if left == nil {
			goto end
		}
		i.node = left
		i.entry = left.Entries[0]
		goto between
	}
	{
		e, _ := i.btree.searchNode(i.node, i.entry.Key)
		if e+1 < len(i.node.Children) {
			i.node = i.node.Children[e+1]
			for len(i.node.Children) > 0 {
				i.node = i.node.Children[0]
			}

			i.entry = i.node.Entries[0]
			goto between
		}

		if e+1 < len(i.node.Entries) {
			i.entry = i.node.Entries[e+1]
			goto between
		}
	}

	for i.node.Parent != nil {
		i.node = i.node.Parent
		e, _ := i.btree.searchNode(i.node, i.entry.Key)
		if e < len(i.node.Entries) {
			i.entry = i.node.Entries[e]
			goto between
		}
	}

end:
	i.End()
	return false

between:
	i.pos = BETWEEN
	return true
}

// Prev moves the iterator to the previous element and returns true
// if there was a previous element that we passed
func (i *Iterator[K, V]) Prev() bool {
	if i.pos == START {
		goto start
	}
	if i.pos == END {
		right := i.btree.Right()
		if right == nil {
			goto start
		}
		i.node = right
		i.entry = right.Entries[len(right.Entries)-1]
		goto between
	}
	{
		e, _ := i.btree.searchNode(i.node, i.entry.Key)
		if e < len(i.node.Children) {
			i.node = i.node.Children[e]
			for len(i.node.Children) > 0 {
				i.node = i.node.Children[len(i.node.Children)-1]
			}
			i.entry = i.node.Entries[len(i.node.Entries)-1]
			goto between
		}
		if e-1 >= 0 {
			i.entry = i.node.Entries[e-1]
			goto between
		}
	}
	for i.node.Parent != nil {
		i.node = i.node.Parent
		e, _ := i.btree.searchNode(i.node, i.entry.Key)
		if e-1 >= 0 {
			i.entry = i.node.Entries[e-1]
			goto between
		}
	}

start:
	i.Begin()
	return false

between:
	i.pos = BETWEEN
	return true
}

// Value returns the current element's value
func (i *Iterator[K, V]) Value() V {
	return i.entry.Value
}

// Key return the current element's key
func (i *Iterator[K, V]) Key() K {
	return i.entry.Key
}

// Node returns the current element's node
func (i *Iterator[K, V]) Node() *Node[K, V] {
	return i.node
}

// Begin resets the iterator to its initial state
func (i *Iterator[K, V]) Begin() {
	i.node = nil
	i.pos = START
	i.entry = nil
}

// End moves the iterator past the last element
func (i *Iterator[K, V]) End() {
	i.node = nil
	i.pos = END
	i.entry = nil
}

// First moves the iterator to the first element and returns true if there aws
// a first element we passed by.
func (i *Iterator[K, V]) First() bool {
	i.Begin()
	return i.Next()
}

func (i *Iterator[K, V]) Last() bool {
	i.End()
	return i.Prev()
}

func (i *Iterator[K, V]) NextTo(f func(key K, val V) bool) bool {
	for i.Next() {
		key, val := i.Key(), i.Value()
		if f(key, val) {
			return true
		}
	}
	return false
}

func (i *Iterator[K, V]) PrevTo(f func(key K, val V) bool) bool {
	for i.Prev() {
		key, value := i.Key(), i.Value()
		if f(key, value) {
			return true
		}
	}

	return false
}

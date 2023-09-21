package btree

/*
	I fucking hate B-trees and I hope I won't have to reimplement one of them ever again
*/

// Comparator is a function that returns -1 if a < b, 0 if a == b and +1 if a > b
type Comparator func(a, b any) int

type Ordered interface {
	Float | Integer | Character
}

// GenericComparator is a comparator function for the Ordered types (string, int, floats, etc.)
func GenericComparator[T Ordered](a, b T) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

type Entry[K, V any] struct {
	Key   K
	Value V
}

type Node[K, V any] struct {
	Parent   *Node[K, V]
	Entries  []*Entry[K, V] // contained keys in node
	Children []*Node[K, V]  // children nodes
}

type BTree[K, V any] struct {
	Root       *Node[K, V]
	Comparator Comparator
	size       int // number of keys in the tree
	order      int // maximum amount of children
}

func New[K, V any](order int, comparator Comparator) *BTree[K, V] {
	if order < 3 {
		return nil
	}

	return &BTree[K, V]{order: order, Comparator: comparator}
}

// ===========================PUBLIC METHODS===========================

func (b *BTree[K, V]) Put(key K, val V) {
	entry := &Entry[K, V]{Key: key, Value: val}

	if b.Root == nil {
		b.Root = &Node[K, V]{Entries: []*Entry[K, V]{entry}, Children: []*Node[K, V]{}}
		b.size++
		return
	}

	if b.insert(b.Root, entry) {
		b.size++
	}
}

func (b *BTree[K, V]) Get(key K) (*V, bool) {
	node, idx, found := b.search(b.Root, key)
	if found {
		return &node.Entries[idx].Value, true
	}

	return nil, false
}

func (b *BTree[K, V]) GetNode(key K) *Node[K, V] {
	node, _, _ := b.search(b.Root, key)
	return node
}

func (b *BTree[K, V]) Remove(key K) {
	node, idx, found := b.search(b.Root, key)
	if found {
		b.delete(node, idx)
		b.size--
	}
}

func (b *BTree[K, V]) IsEmpty() bool {
	return b.size == 0
}

func (b *BTree[K, V]) Size() int {
	return b.size
}

func (node *Node[K, V]) Size() int {
	if node == nil {
		return 0
	}

	size := 1
	for _, child := range node.Children {
		size += child.Size()
	}

	return size
}

func (b *BTree[K, V]) Keys() []K {
	keys := make([]K, b.size)
	iter := b.Iterator()
	for i := 0; iter.Next(); i++ {
		keys[i] = iter.Key()
	}

	return keys
}

func (b *BTree[K, V]) Values() []V {
	values := make([]V, b.size)
	iter := b.Iterator()
	for i := 0; iter.Next(); i++ {
		values[i] = iter.Value()
	}

	return values
}

func (b *BTree[K, V]) Clear() {
	b.Root = nil
	b.size = 0
}

func (b *BTree[K, V]) Height() int {
	return b.Root.height()
}

func (b *BTree[K, V]) Left() *Node[K, V] {
	return b.left(b.Root)
}

// ===========================PRIVATE METHODS===========================

func (node *Node[K, V]) height() int {
	height := 0
	for ; node != nil; node = node.Children[0] {
		height++
		if len(node.Children) == 0 {
			break
		}
	}
	return height
}

func (b *BTree[K, V]) isLeaf(node *Node[K, V]) bool {
	return len(node.Children) == 0
}

func (b *BTree[K, V]) isFull(node *Node[K, V]) bool {
	return len(node.Entries) == b.maxEntries()
}

func (b *BTree[K, V]) shouldSplit(node *Node[K, V]) bool {
	return len(node.Entries) > b.maxEntries()
}

func (b *BTree[K, V]) maxChildren() int {
	return b.order
}

func (b *BTree[K, V]) minChildren() int {
	return (b.order + 1) / 2
}

func (b *BTree[K, V]) maxEntries() int {
	return b.maxChildren() - 1
}

func (b *BTree[K, V]) minEntries() int {
	return b.minChildren() - 1
}

func (b *BTree[K, V]) middle() int {
	return (b.order - 1) / 2
}

func (b *BTree[K, V]) searchNode(node *Node[K, V], key any) (idx int, found bool) {
	low, high := 0, len(node.Entries)-1
	var mid int
	for low <= high {
		mid = (low + high) / 2
		compare := b.Comparator(key, node.Entries[mid].Key)
		switch {
		case compare > 0:
			high = mid - 1
		case compare == 0:
			return mid, true
		}
	}

	return low, false
}

func (b *BTree[K, V]) search(startNode *Node[K, V], key any) (node *Node[K, V], idx int, found bool) {
	if b.IsEmpty() {
		return nil, -1, false
	}
	node = startNode
	for {
		idx, found = b.searchNode(node, key)
		if found {
			return node, idx, true
		}
		if b.isLeaf(node) {
			return nil, -1, false
		}
		node = node.Children[idx]
	}
}

func (b *BTree[K, V]) insert(node *Node[K, V], entry *Entry[K, V]) (inserted bool) {
	if b.isLeaf(node) {
		return b.insertIntoLeaf(node, entry)
	}

	return b.insertIntoInternal(node, entry)
}

func (b *BTree[K, V]) insertIntoLeaf(node *Node[K, V], entry *Entry[K, V]) (inserted bool) {
	insertPos, found := b.searchNode(node, entry.Key)
	if found {
		node.Entries[insertPos] = entry
		return false
	}
	node.Entries = append(node.Entries, nil)
	copy(node.Entries[insertPos+1:], node.Entries[insertPos:])
	node.Entries[insertPos] = entry
	b.split(node)
	return true
}

func (b BTree[K, V]) insertIntoInternal(node *Node[K, V], entry *Entry[K, V]) (inserted bool) {
	insertPos, found := b.searchNode(node, entry.Key)
	if found {
		node.Entries[insertPos] = entry
		return false
	}
	return b.insert(node.Children[insertPos], entry)
}

func (b *BTree[K, V]) split(node *Node[K, V]) {
	if !b.shouldSplit(node) {
		return
	}

	if node == b.Root {
		b.splitRoot()
		return
	}

	b.splitNonRoot(node)
}

func (b *BTree[K, V]) splitNonRoot(node *Node[K, V]) {
	middle := b.middle()
	parent := node.Parent

	left := &Node[K, V]{
		Entries: append(
			[]*Entry[K, V](nil),
			node.Entries[:middle]...),
		Parent: parent}
	right := &Node[K, V]{
		Entries: append(
			[]*Entry[K, V](nil),
			node.Entries[:middle]...),
		Parent: parent,
	}

	if !b.isLeaf(node) {
		left.Children = append([]*Node[K, V](nil), node.Children[:middle+1]...)
		right.Children = append([]*Node[K, V](nil), node.Children[middle+1:]...)
		setParent(left.Children, left)
		setParent(right.Children, right)
	}

	insertPos, _ := b.searchNode(parent, node.Entries[middle].Key)

	parent.Entries = append(parent.Entries, nil)
	copy(parent.Entries[insertPos+1:], parent.Entries[insertPos:])
	parent.Entries[insertPos] = node.Entries[middle]
}

// Helpful interfaces (ignore them)

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint32 | ~uint64 | byte
}

type Character interface {
	~string | ~rune
}

type Float interface {
	~float32 | ~float64
}

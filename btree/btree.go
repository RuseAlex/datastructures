package btree

// A compare function that returns -1 if a < b, 0 if a == b and +1 if a > b
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

type Entry struct {
	Key   any
	Value any
}

type Node struct {
	Parent   *Node
	Entries  []*Entry // contained keys in node
	Children []*Node  // children nodes
}

type BTree struct {
	Root       *Node
	Comparator Comparator
	size       int // number of keys in the tree
	order      int // maximum amount of children
}

func New(order int, comparator Comparator) *BTree {
	if order < 3 {
		return nil
	}

	return &BTree{order: order, Comparator: comparator}
}

func (b *BTree) Put(key any, val any) {
	entry := &Entry{Key: key, Value: val}

	if b.Root == nil {
		b.Root = &Node{Entries: []*Entry{entry}, Children: []*Node{}}
		b.size++
		return
	}

	if b.insert(b.Root, entry) {
		b.size++
	}
}

func (b *BTree) Get(key any) (val any, found bool) {
	node, idx, found := b.searchRecursively(b.Root, key)
	if found {
		return node.Entries[idx].Value, true
	}

	return nil, false
}

func (b *BTree) GetNode(key any) *Node {
	node, _, _ := b.searchRecursively(b.Root, key)
	return node
}

func (b *BTree) Remove(key any) {
	node, idx, found := b.searchRecursively(b.Root, key)
	if found {
		b.delete(node, idx)
		b.size--
	}
}

func (b *BTree) Empty() bool {
	return b.size == 0
}

func (b *BTree) Size() int {
	return b.size
}

func (node *Node) Size() int {
	if node == nil {
		return 0
	}

	size := 1
	for _, child := range node.Children {
		size += child.Size()
	}

	return size
}

func (b *BTree) Keys() []any {
	keys := make([]any, b.size)
	iter := b.Iterator()
	for i := 0; iter.Next(); i++ {
		keys[i] = iter.Key()
	}

	return keys
}

func (b *BTree) Values() []any {
	values := make([]any, b.size)
	iter := b.Iterator()
	for i := 0; iter.Next(); i++ {
		values[i] = iter.Value()
	}

	return values
}

func (b *BTree) Clear() {
	b.Root = nil
	b.size = 0
}

func (b *BTree) Height() int {
	return b.Root.height()
}

func (b *BTree) Left() *Node {
	return b.left(b.Root)
}

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint32 | ~uint64 | byte
}

type Character interface {
	~string | ~rune
}

type Float interface {
	~float32 | ~float64
}

package btree

/*
	I fucking hate B-trees and I hope I won't have to reimplement one of them ever again
*/

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

// Put inserts key-value pair node into the B-Tree
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

// Get searches the node in the tree by key and returns its value or nil if key is not found in btree
func (b *BTree[K, V]) Get(key K) (*V, bool) {
	node, idx, found := b.search(b.Root, key)
	if found {
		return &node.Entries[idx].Value, true
	}

	return nil, false
}

// GetNode searches the node in the tree by key and returns it
func (b *BTree[K, V]) GetNode(key K) *Node[K, V] {
	node, _, _ := b.search(b.Root, key)
	return node
}

// Remove removes the node from the tree by key
func (b *BTree[K, V]) Remove(key K) {
	node, idx, found := b.search(b.Root, key)
	if found {
		b.delete(node, idx)
		b.size--
	}
}

// IsEmpty checks if the btree is empty
func (b *BTree[K, V]) IsEmpty() bool {
	return b.size == 0
}

// Size returns the number of elements stored in the btree
func (b *BTree[K, V]) Size() int {
	return b.size
}

// Size returns the number of elements stored in the subtree
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

// Keys returns all keys in-order
func (b *BTree[K, V]) Keys() []K {
	keys := make([]K, b.size)
	iter := b.Iterator()
	for i := 0; iter.Next(); i++ {
		keys[i] = iter.Key()
	}

	return keys
}

// Values returns all values in-order based on the key
func (b *BTree[K, V]) Values() []V {
	values := make([]V, b.size)
	iter := b.Iterator()
	for i := 0; iter.Next(); i++ {
		values[i] = iter.Value()
	}

	return values
}

// Clear removes all nodes from the tree
func (b *BTree[K, V]) Clear() {
	b.Root = nil
	b.size = 0
}

// Height returns the height of the tree
func (b *BTree[K, V]) Height() int {
	return b.Root.height()
}

// Left returns the left-most (i.e. min) node or nil if tree is empty
func (b *BTree[K, V]) Left() *Node[K, V] {
	return b.left(b.Root)
}

// LeftKey returns the left-most (i.e. min) key or nil
func (b *BTree[K, V]) LeftKey() *K {
	if left := b.Left(); left != nil {
		return &left.Entries[0].Key
	}

	return nil
}

// LeftValue returns the left-most value or nil if tree is empty.
func (b *BTree[K, V]) LeftValue() *V {
	if left := b.Left(); left != nil {
		return &left.Entries[0].Value
	}

	return nil
}

// Right returns the right-most (max) node or nil if tree is empty.
func (b *BTree[K, V]) Right() *Node[K, V] {
	return b.right(b.Root)
}

// RightKey returns the right-most (max) key or nil if tree is empty.
func (b *BTree[K, V]) RightKey() *K {
	if right := b.Right(); right != nil {
		return &right.Entries[len(right.Entries)-1].Key
	}

	return nil
}

// RightValue returns the right-most value or nil if tree is empty.
func (b *BTree[K, V]) RightValue() *V {
	if right := b.Right(); right != nil {
		return &right.Entries[len(right.Entries)-1].Value
	}

	return nil
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
		setParent[K, V](left.Children, left)
		setParent[K, V](right.Children, right)
	}

	insertPos, _ := b.searchNode(parent, node.Entries[middle].Key)

	parent.Entries = append(parent.Entries, nil)
	copy(parent.Entries[insertPos+1:], parent.Entries[insertPos:])
	parent.Entries[insertPos] = node.Entries[middle]

	parent.Children[insertPos] = left

	parent.Children = append(parent.Children, nil)
	copy(parent.Children[insertPos+2:], parent.Children[insertPos+1:])
	parent.Children[insertPos+1] = right

	b.split(parent)
}

func (b *BTree[K, V]) splitRoot() {
	middle := b.middle()

	left := &Node[K, V]{Entries: append(
		[]*Entry[K, V](nil),
		b.Root.Entries[:middle]...)}
	right := &Node[K, V]{Entries: append(
		[]*Entry[K, V](nil),
		b.Root.Entries[middle+1:]...)}

	if !b.isLeaf(b.Root) {
		left.Children = append([]*Node[K, V](nil), b.Root.Children[:middle+1]...)
		right.Children = append([]*Node[K, V](nil), b.Root.Children[middle+1:]...)
		setParent[K, V](left.Children, left)
		setParent[K, V](right.Children, right)
	}

	newRoot := &Node[K, V]{
		Entries:  []*Entry[K, V]{b.Root.Entries[middle]},
		Children: []*Node[K, V]{left, right},
	}

	left.Parent = newRoot
	right.Parent = newRoot
	b.Root = newRoot
}

func setParent[K, V any](nodes []*Node[K, V], parent *Node[K, V]) {
	for _, node := range nodes {
		node.Parent = parent
	}
}

func (b *BTree[K, V]) left(node *Node[K, V]) *Node[K, V] {
	if b.IsEmpty() {
		return nil
	}
	curNode := node
	for {
		if b.isLeaf(curNode) {
			return curNode
		}

		curNode = curNode.Children[0]
	}
}

func (b *BTree[K, V]) right(node *Node[K, V]) *Node[K, V] {
	if b.IsEmpty() {
		return nil
	}
	curNode := node
	for {
		if b.isLeaf(curNode) {
			return curNode
		}
		curNode = curNode.Children[len(curNode.Children)-1]
	}
}

func (b *BTree[K, V]) leftSibling(node *Node[K, V], key K) (*Node[K, V], int) {
	if node.Parent != nil {
		idx, _ := b.searchNode(node.Parent, key)
		idx--
		if idx >= 0 && idx < len(node.Parent.Children) {
			return node.Parent.Children[idx], idx
		}
	}

	return nil, -1
}

func (b *BTree[K, V]) rightSibling(node *Node[K, V], key K) (*Node[K, V], int) {
	if node.Parent != nil {
		idx, _ := b.searchNode(node.Parent, key)
		idx++
		if idx < len(node.Parent.Children) {
			return node.Parent.Children[idx], idx
		}
	}

	return nil, -1
}

func (b *BTree[K, V]) delete(node *Node[K, V], idx int) {
	if b.isLeaf(node) {
		deletedKey := node.Entries[idx].Key
		b.deleteEntry(node, idx)
		b.rebalance(node, deletedKey)
		if len(b.Root.Entries) == 0 {
			b.Root = nil
		}
		return
	}

	// deleting from an internal node
	leftLargestNode := b.right(node.Children[idx]) //largest node in the left subtree (we assume it exists)
	leftLargestEntryIndex := len(leftLargestNode.Entries) - 1
	node.Entries[idx] = leftLargestNode.Entries[leftLargestEntryIndex]
	deletedKey := leftLargestNode.Entries[leftLargestEntryIndex].Key
	b.deleteEntry(leftLargestNode, leftLargestEntryIndex)
	b.rebalance(leftLargestNode, deletedKey)
}

func (b *BTree[K, V]) rebalance(node *Node[K, V], deletedKey K) {
	// check if we need to re-balance it in the first place
	if node == nil || len(node.Entries) >= b.minEntries() {
		return
	}

	// we try to borrow from the left sibling
	leftSibling, leftSiblingIndex := b.leftSibling(node, deletedKey)
	if leftSibling != nil && len(leftSibling.Entries) > b.minEntries() {
		// rotate right
		node.Entries = append([]*Entry[K, V]{node.Parent.Entries[leftSiblingIndex]}, node.Entries...)
		node.Parent.Entries[leftSiblingIndex] = leftSibling.Entries[len(leftSibling.Entries)-1]
		b.deleteEntry(leftSibling, len(leftSibling.Entries)-1)
		if !b.isLeaf(leftSibling) {
			leftSiblingRightMostChild := leftSibling.Children[len(leftSibling.Children)-1]
			leftSiblingRightMostChild.Parent = node
			node.Children = append([]*Node[K, V]{leftSiblingRightMostChild}, node.Children...)
			b.deleteChild(leftSibling, len(leftSibling.Children)-1)
		}
		return
	}

	// we try to borrow now from the right sibling
	rightSibling, rightSiblingIndex := b.rightSibling(node, deletedKey)
	if rightSibling != nil && len(rightSibling.Entries) > b.minEntries() {
		// rotate left
		node.Entries = append(node.Entries, node.Parent.Entries[rightSiblingIndex-1])
		node.Parent.Entries[rightSiblingIndex-1] = rightSibling.Entries[0]
		b.deleteEntry(rightSibling, 0)
		if !b.isLeaf(rightSibling) {
			rightSiblingLeftMostChild := rightSibling.Children[0]
			rightSiblingLeftMostChild.Parent = node
			node.Children = append(node.Children, rightSiblingLeftMostChild)
			b.deleteChild(rightSibling, 0)
		}
		return
	}

	// merge with siblings
	if rightSibling != nil {
		// merge with the right sibling
		node.Entries = append(node.Entries, node.Parent.Entries[rightSiblingIndex-1])
		node.Entries = append(node.Entries, rightSibling.Entries...)
		deletedKey = node.Parent.Entries[rightSiblingIndex-1].Key
		b.deleteEntry(node.Parent, rightSiblingIndex-1)
		b.appendChildren(node.Parent.Children[rightSiblingIndex], node)
		b.deleteChild(node.Parent, rightSiblingIndex)
	} else if leftSibling != nil {
		// merge with the left sibling
		entries := append([]*Entry[K, V](nil), leftSibling.Entries...)
		entries = append(entries, node.Parent.Entries[leftSiblingIndex])
		node.Entries = append(entries, node.Entries...)
		deletedKey = node.Parent.Entries[leftSiblingIndex].Key
		b.deleteEntry(node.Parent, leftSiblingIndex)
		b.prependChildren(node.Parent.Children[leftSiblingIndex], node)
		b.deleteChild(node.Parent, leftSiblingIndex)
	}

	// make the merged node the root if its parent was the root and the root are empty
	if node.Parent == b.Root && len(b.Root.Entries) == 0 {
		b.Root = node
		node.Parent = nil
		return
	}

	// parent node may underflow, so we might need to re-balance it
	b.rebalance(node.Parent, deletedKey)
}

func (b *BTree[K, V]) prependChildren(fromNode *Node[K, V], toNode *Node[K, V]) {
	children := append([]*Node[K, V](nil), fromNode.Children...)
	toNode.Children = append(children, toNode.Children...)
	setParent(fromNode.Children, toNode)
}

func (b *BTree[K, V]) appendChildren(fromNode *Node[K, V], toNode *Node[K, V]) {
	toNode.Children = append(toNode.Children, fromNode.Children...)
	setParent(fromNode.Children, toNode)
}

func (b *BTree[K, V]) deleteEntry(node *Node[K, V], idx int) {
	copy(node.Entries[idx:], node.Entries[idx+1:])
	node.Entries[len(node.Entries)-1] = nil
	node.Entries = node.Entries[:len(node.Entries)-1]
}

func (b *BTree[K, V]) deleteChild(node *Node[K, V], idx int) {
	if idx >= len(node.Children) {
		return
	}
	copy(node.Children[idx:], node.Children[idx+1:])
	node.Children[len(node.Children)-1] = nil
	node.Children = node.Children[:len(node.Children)-1]
}

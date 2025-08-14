package listext

// Node is an element of the doubly linked list.
type Node[V any] struct {
	next  *Node[V]
	prev  *Node[V]
	Value V
}

// Next returns the nodes next Value or nil if it is at the tail.
func (n *Node[V]) Next() *Node[V] {
	return n.next
}

// Prev returns the nodes previous Value or nil if it is at the head.
func (n *Node[V]) Prev() *Node[V] {
	return n.prev
}

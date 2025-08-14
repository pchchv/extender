package listext

// Node is an element of the doubly linked list.
type Node[V any] struct {
	next  *Node[V]
	prev  *Node[V]
	Value V
}

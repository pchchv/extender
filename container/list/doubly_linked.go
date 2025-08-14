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

// DoublyLinkedList is a doubly linked list
type DoublyLinkedList[V any] struct {
	head, tail *Node[V]
	len        int
}

// NewDoublyLinked creates a DoublyLinkedList for use.
func NewDoublyLinked[V any]() *DoublyLinkedList[V] {
	return new(DoublyLinkedList[V])
}

// PushBack appends an element to the back of a list.
func (d *DoublyLinkedList[V]) PushBack(v V) *Node[V] {
	node := &Node[V]{
		Value: v,
	}

	d.pushBack(node)
	return d.tail
}

// PushFront adds an element first in the list.
func (d *DoublyLinkedList[V]) PushFront(v V) *Node[V] {
	node := &Node[V]{
		Value: v,
	}

	d.pushFront(node)
	return d.head
}

// IsEmpty returns true if the list is empty.
func (d *DoublyLinkedList[V]) IsEmpty() bool {
	return d.len == 0
}

// Len returns length of the Linked List.
func (d *DoublyLinkedList[V]) Len() int {
	return d.len
}

func (d *DoublyLinkedList[V]) pushBack(node *Node[V]) {
	node.prev = d.tail
	node.next = nil
	if d.tail == nil {
		d.head = node
	} else {
		d.tail.next = node
	}

	d.tail = node
	d.len++
}

func (d *DoublyLinkedList[V]) pushFront(node *Node[V]) {
	node.next = d.head
	node.prev = nil
	if d.head == nil {
		d.tail = node
	} else {
		d.head.prev = node
	}

	d.head = node
	d.len++
}

// PopBack removes the last element from a list and returns it or nil.
func (d *DoublyLinkedList[V]) PopBack() *Node[V] {
	if d.IsEmpty() {
		return nil
	}

	node := d.tail
	d.tail = node.prev
	if d.tail == nil {
		d.head = nil
	} else {
		d.tail.next = nil
	}

	d.len--
	// ensure no leakage
	node.next, node.prev = nil, nil
	return node
}

// PopFront removes the first element and returns it or nil.
func (d *DoublyLinkedList[V]) PopFront() *Node[V] {
	if d.IsEmpty() {
		return nil
	}

	node := d.head
	d.head = node.next
	if d.head == nil {
		d.tail = nil
	} else {
		d.head.prev = nil
	}

	d.len--
	// ensure no leakage
	node.next, node.prev = nil, nil
	return node
}

// Back returns the end/tail element for use without removing it or nil list is empty.
func (d *DoublyLinkedList[V]) Back() *Node[V] {
	return d.tail
}

// Front returns the front/head element for use without removing it or nil list is empty.
func (d *DoublyLinkedList[V]) Front() *Node[V] {
	return d.head
}

package listext

import (
	"testing"

	. "github.com/pchchv/go-assert"
)

func TestSingleEntryPopBack(t *testing.T) {
	l := NewDoublyLinked[int]()
	Equal(t, l.IsEmpty(), true)
	Equal(t, l.Len(), 0)

	// push some data and then re-check
	zeroNode := l.PushFront(0)
	Equal(t, zeroNode.Value, 0)
	Equal(t, l.IsEmpty(), false)
	Equal(t, l.Len(), 1)
	Equal(t, zeroNode.Prev(), nil)
	Equal(t, zeroNode.Next(), nil)

	// test popping where one node is both head and tail
	back := l.PopBack()
	Equal(t, back.Value, 0)
	Equal(t, back.Next(), nil)
	Equal(t, back.Prev(), nil)
	Equal(t, l.IsEmpty(), true)
	Equal(t, l.Len(), 0)

	front := l.PopFront()
	Equal(t, front, nil)
}

func TestSingleEntryPopFront(t *testing.T) {
	l := NewDoublyLinked[int]()
	Equal(t, l.IsEmpty(), true)
	Equal(t, l.Len(), 0)

	// push some data and then re-check
	zeroNode := l.PushFront(0)
	Equal(t, zeroNode.Value, 0)
	Equal(t, l.IsEmpty(), false)
	Equal(t, l.Len(), 1)
	Equal(t, zeroNode.Prev(), nil)
	Equal(t, zeroNode.Next(), nil)

	// test popping where one node is both head and tail
	front := l.PopFront()
	Equal(t, front.Value, 0)
	Equal(t, front.Next(), nil)
	Equal(t, front.Prev(), nil)
	Equal(t, l.IsEmpty(), true)
	Equal(t, l.Len(), 0)

	back := l.PopBack()
	Equal(t, back, nil)
}

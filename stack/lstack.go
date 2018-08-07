package stack

import (
	"sync"
)

type node struct {
	data interface{}
	next *node
}

// LinkedStack represents a LIFO stack, which using singly linked list as its underlying implemetation.
type LinkedStack struct {
	mux  *sync.Mutex
	head *node
	size int
}

// LStack is alias for LinkedStack.
type LStack = LinkedStack

// NewLStack returns a new instance of LinkedStack.
func NewLStack() *LStack {
	return &LStack{mux: &sync.Mutex{}}
}

// Push pushes the data into the LinkedStack.
func (ls *LStack) Push(data ...interface{}) {
	ls.mux.Lock()
	defer ls.mux.Unlock()

	for _, v := range data {
		newNode := new(node)
		newNode.data = v
		newNode.next = ls.head
		ls.head = newNode
		ls.size++
	}
}

// Pop returns the data popped from the LinkedStack.
func (ls *LStack) Pop() interface{} {
	ls.mux.Lock()
	defer ls.mux.Unlock()

	if ls.size > 0 {
		ret := ls.head
		ls.head = ret.next
		ls.size--
		return ret.data
	}

	return nil
}

// Size return the size of the LinkedStack.
func (ls *LStack) Size() int {
	ls.mux.Lock()
	defer ls.mux.Unlock()

	return ls.size
}

// Clear drops all of the data in the LinkedStack.
func (ls *LStack) Clear() {
	ls.mux.Lock()
	defer ls.mux.Unlock()

	ls.head = nil
	ls.size = 0
}

// Empty returns true if the LinkedStack has no data, otherwise false.
func (ls *LStack) Empty() bool {
	ls.mux.Lock()
	defer ls.mux.Unlock()

	return ls.size == 0
}

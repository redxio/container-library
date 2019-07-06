package stack

type node struct {
	data interface{}
	next *node
}

// LinkedStack represents a LIFO stack, which using singly linked list as its underlying implemetation.
type LinkedStack struct {
	head *node
	size int
}

// LStack is alias for LinkedStack.
type LStack = LinkedStack

// NewLStack returns a new instance of LinkedStack.
func NewLStack() *LStack {
	return &LStack{}
}

// Push pushes the data into the LinkedStack.
func (ls *LStack) Push(data ...interface{}) {
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
	return ls.size
}

// Clear drops all of the data in the LinkedStack.
func (ls *LStack) Clear() {
	ls.head = nil
	ls.size = 0
}

// Empty returns true if the LinkedStack has no data, otherwise false.
func (ls *LStack) Empty() bool {
	return ls.size == 0
}

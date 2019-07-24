// Package stack implements a LIFO stack.
package stack

// Stack represents a LIFO stack implemented using dynamic array.
type Stack struct {
	data []interface{}
}

// NewStack returns a new instance of Stack.
func NewStack() *Stack {
	return &Stack{}
}

// NewStackSize returns a stack with initial size.
func NewStackSize(size int) *Stack {
	return &Stack{make([]interface{}, size)}
}

// Push pushed data into Stack.
func (s *Stack) Push(data ...interface{}) {
	s.data = append(s.data, data...)
}

// Pop returns the data popped from the Stack.
func (s *Stack) Pop() interface{} {
	if len(s.data) > 0 {
		ret := s.data[len(s.data)-1]
		s.data = s.data[:len(s.data)-1]
		return ret
	}

	return nil
}

// Size returns the size of the Stack.
func (s *Stack) Size() int {
	return len(s.data)
}

// Clear drops all of data in the Stack.
func (s *Stack) Clear() {
	s.data = nil
}

// Empty reports whether the Stack is empty.
func (s *Stack) Empty() bool {
	return len(s.data) == 0
}

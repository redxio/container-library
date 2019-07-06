// Package stack implements a library for stack with thread safety.
package stack

// Stack represents a LIFO stack implemented using dynamic array.
type Stack struct {
	data []interface{}
}

// NewStack returns a new instance of Stack.
func NewStack() *Stack {
	return &Stack{}
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

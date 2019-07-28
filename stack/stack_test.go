package stack_test

import (
	"testing"

	"github.com/NzKSO/container/stack"
	"github.com/NzKSO/container/testdata"
)

func createStack() *stack.Stack {
	s := stack.NewStack()
	for _, v := range testdata.TestCases {
		s.Push(v)
	}

	return s
}

func TestStackPush(t *testing.T) {
	s := createStack()

	if s.Size() != len(testdata.TestCases) {
		t.Errorf("%v != %v", s.Size(), len(testdata.TestCases))
	}
}

func TestStackPop(t *testing.T) {
	s := createStack()

	for i := len(testdata.TestCases) - 1; i >= 0; i-- {
		v := s.Pop()
		if testdata.TestCases[i] != v {
			t.Errorf("%v != %v", testdata.TestCases[i], v)
		}
	}

	v := s.Pop()
	if v != nil {
		t.Errorf("%v != nil", v)
	}
}

func TestStackReset(t *testing.T) {
	s := createStack()

	s.Reset()
	v := s.Pop()
	if v != nil {
		t.Errorf("%v != nil", v)
	}
	if s.Size() != 0 {
		t.Errorf("%v != 0", s.Size())
	}
	if !s.Empty() {
		t.Errorf("Stack is empty? %v", s.Empty())
	}
}

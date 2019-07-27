package stack_test

import (
	"testing"

	"github.com/NzKSO/container/stack"
)

func createLStack() *stack.LStack {
	ls := stack.NewLStack()
	for _, v := range testCase {
		ls.Push(v)
	}

	return ls
}

func TestLStackPush(t *testing.T) {
	ls := createLStack()

	if ls.Size() != len(testCase) {
		t.Errorf("%v != %v", ls.Size(), len(testCase))
	}
}

func TestLStackPop(t *testing.T) {
	ls := createLStack()

	for i := len(testCase) - 1; i >= 0; i-- {
		v := ls.Pop()
		if v != testCase[i] {
			t.Errorf("%v != %v", v, testCase[i])
		}
	}

	v := ls.Pop()
	if v != nil {
		t.Errorf("%v != nil", v)
	}
}

func TestLStackReset(t *testing.T) {
	ls := createLStack()

	ls.Reset()
	v := ls.Pop()
	if v != nil {
		t.Errorf("%v != nil", v)
	}
	if ls.Size() != 0 {
		t.Errorf("%v != 0", ls.Size())
	}
	if !ls.Empty() {
		t.Errorf("LStack is empty? %v", ls.Empty())
	}
}

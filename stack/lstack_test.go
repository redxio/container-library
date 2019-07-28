package stack_test

import (
	"testing"

	"github.com/NzKSO/container/stack"
	"github.com/NzKSO/container/testdata"
)

func createLStack() *stack.LStack {
	ls := stack.NewLStack()
	for _, v := range testdata.TestCases {
		ls.Push(v)
	}

	return ls
}

func TestLStackPush(t *testing.T) {
	ls := createLStack()

	if ls.Size() != len(testdata.TestCases) {
		t.Errorf("%v != %v", ls.Size(), len(testdata.TestCases))
	}
}

func TestLStackPop(t *testing.T) {
	ls := createLStack()

	for i := len(testdata.TestCases) - 1; i >= 0; i-- {
		v := ls.Pop()
		if v != testdata.TestCases[i] {
			t.Errorf("%v != %v", v, testdata.TestCases[i])
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

package stack_test

import (
	"testing"

	"github.com/NzKSO/container/stack"
)

type corp struct {
	id      int
	company string
}

var testCase = []corp{
	{22, "General Electric"},
	{19, "Siemens"},
	{2, "Alphabet"},
	{6, "Softbank"},
	{30, "Wells Fargo"},
	{7, "Intel"},
	{33, "Morgan Stanley"},
	{23, "Amazon"},
	{17, "Goldman Sachs Group"},
	{29, "Dell"},
	{0, "Facebook"},
	{24, "Coca-Cola"},
	{12, "Qualcomm"},
	{32, "Verizon Communications"},
	{28, "Apple"},
	{3, "IBM"},
	{11, "Microsoft"},
	{18, "Samsung Electronics"},
	{8, "Booking"},
	{1, "Twitter"},
	{13, "Berkshire Hathaway"},
	{35, "JPMorgan Chase"},
	{31, "Bank of America"},
	{14, "Netflix"},
	{5, "Toyota Motor"},
	{21, "AT&T"},
	{27, "Citigroup"},
	{9, "Wal-Mart Stores"},
	{16, "Uber"},
	{20, "Royal Dutch Shell"},
	{34, "Comcast"},
	{4, "Cisco Systems"},
	{26, "Walt Disney"},
	{15, "Oracle"},
	{10, "Boeing"},
	{25, "HP"},
}

func createStack() *stack.Stack {
	s := stack.NewStack()
	for _, v := range testCase {
		s.Push(v)
	}

	return s
}

func TestStackPush(t *testing.T) {
	s := createStack()

	if s.Size() != len(testCase) {
		t.Errorf("%v != %v", s.Size(), len(testCase))
	}
}

func TestStackPop(t *testing.T) {
	s := createStack()

	for i := len(testCase) - 1; i >= 0; i-- {
		v := s.Pop()
		if testCase[i] != v {
			t.Errorf("%v != %v", testCase[i], v)
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

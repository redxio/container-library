package queue_test

import (
	"testing"

	"github.com/NzKSO/container/queue"
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

func createQueue() *queue.Queue {
	q := queue.NewQueue()
	for _, v := range testCase {
		q.EnQueue(v)
	}

	return q
}

func TestEnQueue(t *testing.T) {
	q := createQueue()

	if q.Size() != len(testCase) {
		t.Errorf("%v != %v", q.Size(), len(testCase))
	}
}

func TestLeQueue(t *testing.T) {
	q := createQueue()

	for _, v := range testCase {
		lv := q.LeQueue()
		if v != lv {
			t.Errorf("%v != %v", v, lv)
		}
	}
}

func TestClear(t *testing.T) {
	q := createQueue()

	q.Clear()
	lv := q.LeQueue()
	if lv != nil {
		t.Errorf("%v != nil", lv)
	}
	if q.Size() != 0 {
		t.Errorf("%v != 0", q.Size())
	}
	if !q.Empty() {
		t.Errorf("LQueue is empty? %v", q.Empty())
	}
}

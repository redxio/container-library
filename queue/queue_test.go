package queue_test

import (
	"testing"

	"github.com/NzKSO/container/queue"
	"github.com/NzKSO/container/testdata"
)

func createQueue() *queue.Queue {
	q := queue.NewQueue()
	for _, v := range testdata.TestCases {
		q.EnQueue(v)
	}

	return q
}

func TestEnQueue(t *testing.T) {
	q := createQueue()

	if q.Size() != len(testdata.TestCases) {
		t.Errorf("%v != %v", q.Size(), len(testdata.TestCases))
	}
}

func TestLeQueue(t *testing.T) {
	q := createQueue()

	for _, v := range testdata.TestCases {
		lv := q.LeQueue()
		if v != lv {
			t.Errorf("%v != %v", v, lv)
		}
	}
}

func TestReset(t *testing.T) {
	q := createQueue()

	q.Reset()
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

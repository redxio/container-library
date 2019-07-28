package queue_test

import (
	"testing"

	"github.com/NzKSO/container/queue"
	"github.com/NzKSO/container/testdata"
)

func createLQueue() *queue.LQueue {
	lq := queue.NewLQueue()

	for _, v := range testdata.TestCases {
		lq.EnQueue(v)
	}

	return lq
}

func TestLEnQueue(t *testing.T) {
	lq := createLQueue()

	if lq.Empty() {
		t.Errorf("Queue is empty? %v", lq.Empty())
	}
	if len(testdata.TestCases) != lq.Size() {
		t.Errorf("%v != %v", len(testdata.TestCases), lq.Size())
	}
}

func TestLLeQueue(t *testing.T) {
	lq := createLQueue()

	for _, v := range testdata.TestCases {
		lv := lq.LeQueue().(testdata.Corp)
		if v != lv {
			t.Errorf("%v != %v", v, lv)
		}
	}

	if lq.Size() != 0 {
		t.Errorf("%v != 0", lq.Size())
	}
	if !lq.Empty() {
		t.Errorf("LQueue is empty? %v", lq.Empty())
	}
}

func TestLReset(t *testing.T) {
	lq := createLQueue()

	lq.Reset()

	lv := lq.LeQueue()
	if lv != nil {
		t.Errorf("%v != nil", lv)
	}
	if lq.Size() != 0 {
		t.Errorf("%v != 0", lq.Size())
	}
	if !lq.Empty() {
		t.Errorf("LQueue is empty? %v", lq.Empty())
	}
}

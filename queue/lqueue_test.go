package queue_test

import (
	"testing"

	"github.com/NzKSO/container/queue"
)

func createLQueue() *queue.LQueue {
	lq := queue.NewLQueue()

	for _, v := range testCase {
		lq.EnQueue(v)
	}

	return lq
}

func TestLEnQueue(t *testing.T) {
	lq := createLQueue()

	if lq.Empty() {
		t.Errorf("Queue is empty? %v", lq.Empty())
	}
	if len(testCase) != lq.Size() {
		t.Errorf("%v != %v", len(testCase), lq.Size())
	}
}

func TestLLeQueue(t *testing.T) {
	lq := createLQueue()

	for _, v := range testCase {
		lv := lq.LeQueue().(corp)
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

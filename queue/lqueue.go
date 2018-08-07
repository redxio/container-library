package queue

import (
	"sync"
)

type qnode struct {
	prev *qnode
	next *qnode
	data interface{}
}

// LQueue represents a FIFO queue implemented using doubly linked list.
type LQueue struct {
	mux  *sync.Mutex
	tail *qnode
	head *qnode
	size int
}

// NewLQueue returns a new instance of LQueue.
func NewLQueue() *LQueue {
	return &LQueue{mux: &sync.Mutex{}}
}

// EnQueue enters data into the tail of the LQueue.
func (lq *LQueue) EnQueue(data ...interface{}) {
	lq.mux.Lock()
	defer lq.mux.Unlock()

	for _, val := range data {
		if lq.head == nil && lq.tail == nil {
			newQNode := new(qnode)
			newQNode.data = val
			lq.head = newQNode
			lq.tail = newQNode
			lq.size++
			continue
		}

		newQNode := new(qnode)
		newQNode.data = val
		newQNode.prev = lq.tail
		lq.tail.next = newQNode
		lq.tail = newQNode
		lq.size++
	}
}

// LeQueue let data leave from the head of the LQueue.
func (lq *LQueue) LeQueue() interface{} {
	lq.mux.Lock()
	defer lq.mux.Unlock()

	if lq.size > 0 {
		ret := lq.head.data
		lq.head = lq.head.next
		if lq.size == 1 {
			lq.tail = nil
		}
		lq.size--
		return ret
	}

	return nil
}

// Size returns the size of the LQueue.
func (lq *LQueue) Size() int {
	lq.mux.Lock()
	defer lq.mux.Unlock()

	return lq.size
}

// Clear clears the LQueue, it will drop all data.
func (lq *LQueue) Clear() {
	lq.mux.Lock()
	defer lq.mux.Unlock()

	lq.head = nil
	lq.tail = nil
	lq.size = 0
}

// Empty returns true if th LQueue is empty, otherwise false.
func (lq *LQueue) Empty() bool {
	lq.mux.Lock()
	defer lq.mux.Unlock()

	return lq.size == 0 && lq.head == nil && lq.tail == nil
}
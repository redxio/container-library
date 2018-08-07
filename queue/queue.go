// Package queue implements queue-related operations with thread safety.
package queue

import (
	"sync"
)

// Queue represents a FIFO queue implemented using dynamic array.
type Queue struct {
	mux  *sync.Mutex
	data []interface{}
}

// NewQueue returns a new instance of LQueue.
func NewQueue() *Queue {
	return &Queue{mux: &sync.Mutex{}}
}

// EnQueue enters data to the tail of the Queue.
func (q *Queue) EnQueue(data ...interface{}) {
	q.mux.Lock()
	defer q.mux.Unlock()

	q.data = append(q.data, data...)
}

// LeQueue leaves from the head of the Queue.
func (q *Queue) LeQueue() interface{} {
	q.mux.Lock()
	defer q.mux.Unlock()

	if len(q.data) > 0 {
		ret := q.data[0]
		q.data = q.data[1:]
		return ret
	}

	return nil
}

// Size returns the size of the Queue.
func (q *Queue) Size() int {
	q.mux.Lock()
	defer q.mux.Unlock()

	return len(q.data)
}

// Clear clears all of data in the Queue.
func (q *Queue) Clear() {
	q.mux.Lock()
	defer q.mux.Unlock()

	q.data = nil
}

// Empty returns true if the Queue is empty.
func (q *Queue) Empty() bool {
	q.mux.Lock()
	defer q.mux.Unlock()

	return len(q.data) == 0
}

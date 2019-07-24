// Package queue implements queue-related operations with thread safety.
package queue

// Queue represents a FIFO queue implemented using dynamic array.
type Queue struct {
	data []interface{}
}

// NewQueue returns a new instance of LQueue.
func NewQueue() *Queue {
	return &Queue{}
}

// NewQueueSize returns a queue with initial size.
func NewQueueSize(size int) *Queue {
	return &Queue{make([]interface{}, size)}
}

// EnQueue enters data to the tail of the Queue.
func (q *Queue) EnQueue(data ...interface{}) {
	q.data = append(q.data, data...)
}

// LeQueue leaves from the head of the Queue.
func (q *Queue) LeQueue() interface{} {
	if len(q.data) > 0 {
		ret := q.data[0]
		q.data = q.data[1:]
		return ret
	}

	return nil
}

// Size returns the size of the Queue.
func (q *Queue) Size() int {
	return len(q.data)
}

// Clear clears all of data in the Queue.
func (q *Queue) Clear() {
	q.data = nil
}

// Empty returns true if the Queue is empty.
func (q *Queue) Empty() bool {
	return len(q.data) == 0
}

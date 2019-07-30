package queue

import (
	"container/list"
	"sync"
	"time"
)

type dqItem struct {
	data   interface{}
	expire time.Time
}

// DelayQueue represents a delay queue
type DelayQueue struct {
	rw        sync.RWMutex
	queue     *list.List
	size      int
	semaphore chan bool
	delay     chan interface{}
}

// TravFunc is used for traversing delay queue, the argument of TravFunc is data item in queue
type TravFunc func(interface{})

func (dq *DelayQueue) delayService() {
	for {
		dq.rw.Lock()
		if dq.queue.Len() == 0 {
			dq.rw.Unlock()
			<-dq.semaphore
			continue
		}

		elem := dq.queue.Front()
		item := elem.Value.(*dqItem)
		now := time.Now()

		if now.Before(item.expire) {
			dq.rw.Unlock()
			time.Sleep(item.expire.Sub(now))
			continue
		}

		if dq.delay != nil {
			dq.delay <- item.data
		}
		dq.queue.Remove(elem)

		dq.rw.Unlock()
	}
}

// NewDelayQueue returns a initialized delay queue
func NewDelayQueue() *DelayQueue {
	dq := &DelayQueue{queue: list.New(), semaphore: make(chan bool)}
	go dq.delayService()
	return dq
}

// EnQueue enters delay queue, stay in queue for delay milliseconds then leave immediately, it will leave immediately if delay less or equal than 0.
func (dq *DelayQueue) EnQueue(data interface{}, delay int64) {
	expireTime := time.Now().Add(time.Millisecond * time.Duration(delay))

	if delay <= 0 {
		if dq.delay != nil {
			dq.delay <- data
		}
		return
	}

	dq.rw.Lock()

	elem := dq.queue.Back()
	var mark *list.Element

	for elem != nil && expireTime.Before(elem.Value.(*dqItem).expire) {
		mark = elem
		elem = elem.Prev()
	}

	if mark == nil {
		dq.queue.PushBack(&dqItem{data, expireTime})
	} else {
		dq.queue.InsertBefore(&dqItem{data, expireTime}, mark)
	}

	dq.rw.Unlock()

	if dq.queue.Len() == 1 {
		dq.semaphore <- true
	}
}

// Receive returns a received only channel, which can be used for receiving something left from queue
func (dq *DelayQueue) Receive() <-chan interface{} {
	if dq.delay == nil {
		dq.delay = make(chan interface{})
	}
	return dq.delay
}

// Trav returns a received only channel, which can be used to receive traversing result
func (dq *DelayQueue) Trav() <-chan interface{} {
	ch := make(chan interface{})

	go func() {
		dq.rw.RLock()
		defer dq.rw.RUnlock()
		defer close(ch)

		for elem := dq.queue.Front(); elem != nil; elem = elem.Next() {
			ch <- elem.Value.(*dqItem).data
		}

	}()

	return ch
}

// TravWithFunc traverses delay queue dq with function f
func (dq *DelayQueue) TravWithFunc(f TravFunc) {
	dq.rw.RLock()
	defer dq.rw.RUnlock()

	if f == nil {
		return
	}

	for elem := dq.queue.Front(); elem != nil; elem = elem.Next() {
		f(elem.Value.(*dqItem).data)
	}
}

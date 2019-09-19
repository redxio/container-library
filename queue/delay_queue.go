package queue

import (
	"container/list"
	"sync"
	"time"
)

type dqItem struct {
	value  interface{}
	expire time.Time
}

type semaphore struct{}

var greenLight, yellowLight semaphore

// DelayQueue represents a delay queue
type DelayQueue struct {
	rw            sync.RWMutex
	queue         *list.List
	size          int
	worker        chan semaphore
	reconsumption chan semaphore
	releaseRLock  bool
	delay         chan interface{}
}

// TravFunc is used for traversing delay queue, the argument of TravFunc is data item in queue
type TravFunc func(interface{})

func (dq *DelayQueue) delayService() {
	var (
		elem *list.Element
		item *dqItem
		now  time.Time
	)

	for {
		dq.rw.RLock()

		if dq.releaseRLock {
			dq.releaseRLock = false
		}

		if dq.queue.Len() == 0 {
			dq.rw.RUnlock()
			dq.releaseRLock = true
			<-dq.worker
			continue
		}

		if len(dq.reconsumption) > 0 {
			<-dq.reconsumption
		}

		elem = dq.queue.Front()
		item = elem.Value.(*dqItem)
		now = time.Now()

		if now.Before(item.expire) {
			dq.rw.RUnlock()
			dq.releaseRLock = true
			select {
			case <-dq.reconsumption:
				continue
			case <-time.After(item.expire.Sub(now)):
			}
		}

		if !dq.releaseRLock {
			dq.rw.RUnlock()
		}

		dq.rw.Lock()

		if len(dq.reconsumption) > 0 {
			<-dq.reconsumption
			dq.rw.Unlock()
			continue
		}

		if dq.delay != nil {
			dq.delay <- item.value
		}
		dq.queue.Remove(elem)

		dq.rw.Unlock()
	}
}

// NewDelayQueue returns a initialized delay queue
func NewDelayQueue() *DelayQueue {
	dq := &DelayQueue{
		queue:         list.New(),
		worker:        make(chan semaphore, 1),
		reconsumption: make(chan semaphore, 1),
	}
	go dq.delayService()
	return dq
}

// EnQueue enters delay queue, stay in queue for delay milliseconds then leave immediately, it will leave immediately if delay less or equal than 0.
func (dq *DelayQueue) EnQueue(value interface{}, delay int64) {
	expireTime := time.Now().Add(time.Millisecond * time.Duration(delay))

	if delay <= 0 {
		if dq.delay != nil {
			dq.delay <- value
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

	item := &dqItem{value, expireTime}

	if elem == nil && mark == nil {
		dq.queue.PushFront(item)
	} else if mark == nil {
		dq.queue.PushBack(item)
	} else if mark == dq.queue.Front() && elem == nil {
		dq.queue.PushFront(item)
		dq.reconsumption <- yellowLight
	} else {
		dq.queue.InsertBefore(item, mark)
	}

	if dq.queue.Len() == 1 {
		dq.worker <- greenLight
	}

	dq.rw.Unlock()
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
			ch <- elem.Value.(*dqItem).value
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
		f(elem.Value.(*dqItem).value)
	}
}

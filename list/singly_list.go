// Package list implements singly linked list, which partially
// supports concurrent operations on list.
package list

import (
	"sync"

	"github.com/NzKSO/container"
)

// Node represents a node of singly linked list.
type Node struct {
	data interface{}
	next *Node
}

// Next returns the next node of current node.
func (pnode *Node) Next() *Node {
	return pnode.next
}

// GetData returns the data of current node.
func (pnode *Node) GetData() interface{} {
	return pnode.data
}

// SinglyList represents a singly linked list.
type SinglyList struct {
	sync.RWMutex
	head            *Node
	size            int
	NumPerGoroutine int // specify every how many nodes of list start a goroutine
}

type findResult struct {
	prev, find *Node
}

type splitResult struct {
	prev, head, tail *Node
}

// SortFunc represents the type of sorting method implemented by the user.
type SortFunc func(head *Node, size ...int)

// NewSinglyList returns a pointer to linked list.
func NewSinglyList() *SinglyList {
	return &SinglyList{}
}

// Insert inserts data into linked list.
func (ll *SinglyList) Insert(data container.Interface) {
	newNode := new(Node)
	newNode.data = data
	newNode.next = ll.head
	ll.head = newNode
	ll.size++
}

func (ll *SinglyList) splitList() <-chan *splitResult {
	ch := make(chan *splitResult)

	go func() {
		defer close(ch)

		if ll.NumPerGoroutine <= 0 || ll.NumPerGoroutine >= ll.size || ((ll.size - ll.NumPerGoroutine) < (ll.NumPerGoroutine / 3)) {
			ch <- &splitResult{nil, ll.head, nil}
			return
		}

		move, head := ll.head, ll.head

		var (
			prev *Node
			ct   int
		)

		for move != nil {
			result := splitResult{}

			for size := 1; size < ll.NumPerGoroutine && move.next != nil; size++ {
				move = move.next
				if (ll.size - ct) < (ll.NumPerGoroutine / 3) {
					for move.next != nil {
						move = move.next
					}
				}
			}
			ct += ll.NumPerGoroutine

			result.head = head
			result.tail = move.next
			result.prev = prev

			prev = move
			move = move.next
			head = move

			ch <- &result
		}
	}()

	return ch
}

// Delete deletes data specified by key from linked list.
func (ll *SinglyList) Delete(key interface{}) error {
	if ll.size == 0 && ll.head == nil {
		return container.ErrEmptyList
	}

	res := ll.multiGoroutinesFind(ll.splitList(), key)

	if res != nil {
		if res.prev == nil {
			ll.head = res.find.next
			ll.size--
			return nil
		}
		ll.Lock()
		res.prev.next = res.find.next
		ll.Unlock()
		ll.size--
		return nil
	}
	return container.ErrNotExist
}

// Search searches data associated with key by lanuching multiple goroutines
func (ll *SinglyList) Search(key interface{}) (interface{}, error) {
	if ll.head == nil && ll.size == 0 {
		return nil, container.ErrEmptyList
	}

	res := ll.multiGoroutinesFind(ll.splitList(), key)

	if res != nil {
		return res.find.data, nil
	}
	return nil, container.ErrNotExist
}

func (ll *SinglyList) multiGoroutinesFind(splitCh <-chan *splitResult, key interface{}) *findResult {
	findResCh := make(chan *findResult)
	termGoroutineCh := make(chan struct{}, 1)

	var wg sync.WaitGroup

	for split := range splitCh {
		wg.Add(1)
		go func(split *splitResult) {
			defer wg.Done()

			walk := split.head
			ll.RLock()
			end := split.tail
			ll.RUnlock()

			var (
				prev *Node
				itf  container.Interface
			)

			for walk != end {
				itf = walk.data.(container.Interface)
				if itf.Find(key) {
					if walk == split.head {
						findResCh <- &findResult{split.prev, walk}
						return
					}
					findResCh <- &findResult{prev, walk}
					return
				}
				select {
				case <-termGoroutineCh:
					return
				default:
					prev = walk
					ll.RLock()
					walk = walk.next
					ll.RUnlock()
				}
			}
		}(split)
	}

	go func() {
		defer close(findResCh)
		wg.Wait()
	}()

	for {
		res, ok := <-findResCh
		if ok {
			termGoroutineCh <- struct{}{}
			return res
		}
		return res
	}
}

// Update updates data associated with key in linked list.
func (ll *SinglyList) Update(key interface{}, val interface{}) error {
	if ll.head == nil && ll.size == 0 {
		return container.ErrEmptyList
	}

	res := ll.multiGoroutinesFind(ll.splitList(), key)
	if res.find != nil {
		itf := res.find.data.(container.Interface)
		itf.Set(val)
		return nil
	}

	return container.ErrNotExist
}

// Traversal returns a received only channel, which can be used to receive results
// that returned by traversing linked list.
func (ll *SinglyList) Traversal() <-chan interface{} {
	ch := make(chan interface{}, ll.size)
	go func() {
		defer close(ch)

		if ll.head == nil && ll.size == 0 {
			return
		}

		walk := ll.head
		for walk != nil {
			ch <- walk.data
			walk = walk.next
		}
	}()

	return ch
}

func (ll *SinglyList) reverse(split *splitResult, wg *sync.WaitGroup) {
	defer wg.Done()

	move := split.head
	prev := split.prev
	var temp *Node

	for move != split.tail {
		temp = move.next
		ll.Lock()
		move.next = prev
		ll.Unlock()
		prev = move
		move = temp
	}

	if split.tail == nil {
		ll.head = prev
	}
}

// Reverse reverses the list concurrently.
func (ll *SinglyList) Reverse() {
	if ll.head == nil || ll.head.next == nil {
		return
	}

	var wg sync.WaitGroup
	for split := range ll.splitList() {
		wg.Add(1)
		go ll.reverse(split, &wg)
	}
	wg.Wait()
}

// Empty returns true if the list is empty, otherwise false.
func (ll *SinglyList) Empty() bool {
	return ll.head == nil && ll.size == 0
}

// Size returns the size of list ll.
func (ll *SinglyList) Size() int {
	return ll.size
}

// Reset resets ll to its initial state, it will drop all of data.
func (ll *SinglyList) Reset() {
	ll.head = nil
	ll.size = 0
}

// BubbleSort represents Bubble sorting, which can be used as parameter to method SortWith.
func BubbleSort(head *Node) {
	var end *Node
	var start = head

	for start != end {
		for start.next != end {
			itf := start.data.(container.Lesser)
			ret := itf.Less(start.next.data)
			if !ret {
				start.data, start.next.data = start.next.data, start.data
			}
			start = start.next
		}
		end = start
		start = head
	}

	/*var slow = head
	var fast *Node
	for slow != fast {
		fast = slow.next
		for fast != nil {
			itf := slow.data.(Lesser)
			ret := itf.Less(fast.data)
			if !ret {
				slow.data, fast.data = fast.data, slow.data
			}
			fast = fast.next
		}
		slow = slow.next
	}*/
}

// Sort sorts the list using merge sorting by default
func (ll *SinglyList) Sort() {
	if ll.head == nil || ll.head.next == nil {
		return
	}
	mergeSort(&ll.head)
}

func getMiddleNode(head *Node) *Node {
	slow, fast := head, head

	for fast.next != nil && fast.next.next != nil {
		slow = slow.next
		fast = fast.next.next
	}

	return slow
}

func mergeSort(phead **Node) {
	if *phead == nil || (*phead).next == nil {
		return
	}
	var front, back *Node

	middle := getMiddleNode(*phead)
	front = *phead
	back = middle.next
	middle.next = nil

	mergeSort(&front)
	mergeSort(&back)

	*phead = mergeList(front, back)
}

func mergeList(front, back *Node) *Node {
	var head *Node

	if front == nil {
		return back
	} else if back == nil {
		return front
	}

	itf := front.data.(container.Lesser)
	ret := itf.Less(back.data)
	if ret {
		head = front
		head.next = mergeList(front.next, back)
	} else {
		head = back
		head.next = mergeList(front, back.next)
	}

	return head
}

// InsertionSort represents insertion sorting, which can be used as parameter to method SortWith.
func InsertionSort(phead **Node) {
	var (
		sorted, next *Node
		current      = *phead
	)

	for current != nil {
		next = current.next
		sortedInsert(&sorted, current)
		current = next
	}
	*phead = sorted
}

func sortedInsert(phead **Node, newNode *Node) {
	itf := newNode.data.(container.Lesser)
	if *phead == nil || itf.Less((*phead).data) {
		newNode.next = *phead
		*phead = newNode
	} else {
		current := *phead
		for current.next != nil && !itf.Less(current.next.data) {
			current = current.next
		}
		newNode.next = current.next
		current.next = newNode
	}
}

// SortWith sorts the list using user defined sorting method.
func (ll *SinglyList) SortWith(sort SortFunc) {
	if ll.head == nil || ll.head.next == nil {
		return
	}
	sort(ll.head, ll.size)
}

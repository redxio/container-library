// Package list implements singly linked list, which partially
// supports concurrent operations on list.
package list

import (
	"runtime"
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

// LinkedList represents a singly linked list.
type LinkedList struct {
	head *Node
	size int
}

type findResult struct {
	prev, find *Node
}

type splitResult struct {
	beforeHead, head, tail *Node
}

// LList is alias for LinkedList.
type LList = LinkedList

// SortFunc represents the type of sorting method implemented by the user.
type SortFunc func(head *Node, size ...int)

var (
	// SizOfSublst represents that every how much size of list start a goroutine. -1 is default values.
	SizOfSublst = -1

	// SetSizOfSublst is used to set the size of the child linked list according to the length of the linked list
	SetSizOfSublst func(size int) int

	nCPUs = runtime.NumCPU()

	rw sync.RWMutex
)

// NewLinkedList returns a pointer to linked list.
func NewLinkedList() *LinkedList {
	return &LinkedList{}
}

// Insert inserts data into linked list.
func (ll *LinkedList) Insert(data container.Interface) {
	newNode := new(Node)
	newNode.data = data
	newNode.next = ll.head
	ll.head = newNode
	ll.size++
}

func getSizOfSublst(size int) int {
	var sizofsublst int
	if SetSizOfSublst != nil {
		if sizofsublst = SetSizOfSublst(size); sizofsublst < 0 {
			panic("Negative number")
		} else if sizofsublst == 0 {
			sizofsublst = size
		}
	} else if SizOfSublst > 0 {
		sizofsublst = SizOfSublst
	} else if SizOfSublst == 0 {
		sizofsublst = size
	} else if sizofsublst = size / nCPUs; size <= nCPUs {
		sizofsublst = size
	}
	return sizofsublst
}

func splitList(ll *LinkedList, sizofsublst int) <-chan *splitResult {
	nGoroutines := ll.size / sizofsublst
	if ll.size < sizofsublst {
		nGoroutines++
	} else {
		rest := ll.size % sizofsublst
		if rest >= 6 {
			nGoroutines++
		}
	}
	ch := make(chan *splitResult, nGoroutines)

	go func() {
		defer close(ch)
		var beforeHead *Node
		head := ll.head
		tail := ll.head
		ct := 0

		for i := 0; i < nGoroutines; i++ {
			for ct < sizofsublst-1 && tail.next != nil {
				tail = tail.next
				ct++
			}
			if i == nGoroutines-1 {
				for tail.next != nil {
					tail = tail.next
				}
				ch <- &splitResult{beforeHead, head, tail}
				break
			}
			ch <- &splitResult{beforeHead, head, tail}
			beforeHead = tail
			rw.RLock()
			head = tail.next
			rw.RUnlock()
			tail = head
			ct = 0
		}
	}()
	return ch
}

// Delete deletes data specified by key from linked list.
func (ll *LinkedList) Delete(key interface{}) error {
	if ll.size == 0 && ll.head == nil {
		return container.ErrEmptyList
	}

	splitCh := splitList(ll, getSizOfSublst(ll.size))
	res := multiGoroutinesFind(ll.head, splitCh, key)

	if res != nil {
		if res.prev == nil {
			ll.head = res.find.next
			ll.size--
			return nil
		}
		rw.Lock()
		res.prev.next = res.find.next
		rw.Unlock()
		ll.size--
		return nil
	}
	return container.ErrNotExist
}

// Search searches data associated with key by lanuching multiple goroutines
func (ll *LinkedList) Search(key interface{}) (interface{}, error) {
	if ll.head == nil && ll.size == 0 {
		return nil, container.ErrEmptyList
	}

	splitCh := splitList(ll, getSizOfSublst(ll.size))
	res := multiGoroutinesFind(ll.head, splitCh, key)

	if res != nil {
		return res.find.data, nil
	}
	return nil, container.ErrNotExist
}

func multiGoroutinesFind(head *Node, splitCh <-chan *splitResult, key interface{}) *findResult {
	findResCh := make(chan *findResult)
	termGoroutineCh := make(chan struct{}, 1)
	var wg sync.WaitGroup

	for split := range splitCh {
		wg.Add(1)
		go func(split *splitResult) {
			defer wg.Done()
			walk := split.head
			var prev *Node
			rw.RLock()
			end := split.tail.next
			rw.RUnlock()

			for walk != end {
				itf := walk.data.(container.Interface)
				if itf.Find(key) {
					if walk == split.head {
						findResCh <- &findResult{split.beforeHead, walk}
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
					rw.RLock()
					walk = walk.next
					rw.RUnlock()
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

func findSubList(split *splitResult, key interface{}, findResCh chan *findResult, terminatedCh chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	walk := split.head
	end := split.tail.next
	var prev *Node

	for walk != end {
		itf := walk.data.(container.Interface)
		if itf.Find(key) {
			if walk == split.head {
				findResCh <- &findResult{split.beforeHead, walk}
				return
			}
			findResCh <- &findResult{prev, walk}
			return
		}
		select {
		case <-terminatedCh:
			return
		default:
			prev = walk
			walk = walk.next
		}
	}
}

// Update updates data associated with key in linked list.
func (ll *LinkedList) Update(key interface{}, val interface{}) error {
	if ll.head == nil && ll.size == 0 {
		return container.ErrEmptyList
	}

	splitCh := splitList(ll, getSizOfSublst(ll.size))
	res := multiGoroutinesFind(ll.head, splitCh, key)
	if res.find != nil {
		itf := res.find.data.(container.Interface)
		itf.Set(val)
		return nil
	}

	return container.ErrNotExist
}

// Traversal returns a received only channel which can be used to receive results
// that returned by traversing linked list.
func (ll *LinkedList) Traversal() <-chan interface{} {
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

func reverse(split *splitResult, ln **Node, wg *sync.WaitGroup) {
	defer wg.Done()
	move := split.head
	prev := split.beforeHead
	last := split.tail.next

	if last == nil {
		*ln = split.tail
	}
	for move != last {
		temp := move.next
		rw.Lock()
		move.next = prev
		rw.Unlock()
		prev = move
		move = temp
	}
}

// Reverse reverses the linked list concurrently.
func (ll *LinkedList) Reverse() {
	if ll.head == nil || ll.head.next == nil {
		return
	}

	var wg sync.WaitGroup
	splitSize := getSizOfSublst(ll.size)
	splitCh := splitList(ll, splitSize)
	for split := range splitCh {
		wg.Add(1)
		go reverse(split, &ll.head, &wg)
	}
	wg.Wait()
}

// Empty returns true if the linked list is empty, otherwise false.
func (ll *LinkedList) Empty() bool {
	return ll.head == nil && ll.size == 0
}

// Size returns the size of list ll.
func (ll *LinkedList) Size() int {
	return ll.size
}

// Reset resets list ll to its initial state, it will drop all of data.
func (ll *LinkedList) Reset() {
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

// Sort sorts the linked list, using merge sorting by default
func (ll *LinkedList) Sort() {
	if ll.head == nil || ll.head.next == nil {
		return
	}
	mergeSort(&ll.head)
}

func getMiddleNode(head *Node) *Node {
	slow := head
	fast := head

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
	var sorted *Node
	var current = *phead

	for current != nil {
		next := current.next
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

// SortWith sorts the linked list using user defined sorting method.
func (ll *LinkedList) SortWith(sort SortFunc) {
	if ll.head == nil || ll.head.next == nil {
		return
	}
	sort(ll.head, ll.size)
}

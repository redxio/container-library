// Package tree implements binary search tree related operations with thread safety
package tree

import (
	"math"
	"reflect"
	"sync"

	"github.com/NzKSO/container"
	"github.com/NzKSO/container/queue"
)

// Tnode represents a node in a binary tree.
type Tnode struct {
	lightChild *Tnode
	rightChild *Tnode
	data       interface{}
}

// GetLchild returns member lightChild pointed by tn.
func (tn *Tnode) GetLchild() *Tnode {
	return tn.lightChild
}

// GetRchild returns member rightChild pointed by tn.
func (tn *Tnode) GetRchild() *Tnode {
	return tn.rightChild
}

// GetData returns member data pointed by tn.
func (tn *Tnode) GetData() interface{} {
	return tn.data
}

// BSTree represents an binary tree.
type BSTree struct {
	mux  *sync.Mutex
	root *Tnode
	size int
}

// 0, 1, 2, 3 respectively denotes traversing in Inorder, Preorder,
// Postorder, levelorder.
const (
	InorderTrav = iota
	PreorderTrav
	PostorderTrav
	LevelTrav
)

// TravFunc is used for traversing the binary search tree, which accepts a pointer to root node
// and a channel as used to send traversing sequence.
type TravFunc func(root *Tnode, ch chan<- interface{})

// NewBSTree returns an empty binary tree.
func NewBSTree() *BSTree {
	return &BSTree{mux: &sync.Mutex{}}
}

func insert(tn *Tnode, data container.Interface) (*Tnode, error) {
	if tn == nil {
		return &Tnode{data: data}, nil
	}

	itf := tn.data.(container.Interface)
	if itf.Find(data) {
		return tn, container.ErrDataExists
	}

	var err error
	toRight := itf.Less(data)
	if toRight {
		tn.rightChild, err = insert(tn.rightChild, data)
	} else {
		tn.lightChild, err = insert(tn.lightChild, data)
	}

	return tn, err
}

// Insert inserts data to binary tree.
func (bt *BSTree) Insert(data container.Interface) error {
	bt.mux.Lock()
	defer bt.mux.Unlock()

	var err error
	bt.root, err = insert(bt.root, data)
	if err != nil {
		return err
	}
	bt.size++

	return nil
}

func lookup(find *Tnode, parent *Tnode, key interface{}) (*Tnode, *Tnode) {
	if find == nil {
		return nil, nil
	}

	v := find.data.(container.Interface)
	if v.Find(key) {
		return find, parent
	}

	toRight := v.Less(key)
	if toRight {
		find, parent = lookup(find.rightChild, find, key)
	} else {
		find, parent = lookup(find.lightChild, find, key)
	}

	return find, parent

}

// Search search tree to found data associated with the provided key. If the tree is empty, returns ErrEmptyTree.
// If not found, returns ErrNotExist.
func (bt *BSTree) Search(key interface{}) (interface{}, error) {
	bt.mux.Lock()
	defer bt.mux.Unlock()
	if bt.root == nil && bt.size == 0 {
		return nil, container.ErrEmptyTree
	}

	tn, _ := lookup(bt.root, nil, key)
	if tn == nil {
		return nil, container.ErrNotExist
	}

	return tn.data, nil
}

// Delete deletes the data found by key. if tree is empty, Delete returns ErrEmptyTree,
// if the data doesn't exist, it will return ErrNotExist.
func (bt *BSTree) Delete(key interface{}) error {
	bt.mux.Lock()
	defer bt.mux.Unlock()

	if bt.root == nil && bt.size == 0 {
		return container.ErrEmptyTree
	}

	find, parent := lookup(bt.root, nil, key)
	if find == nil {
		return container.ErrNotExist
	}

	if find.lightChild == nil && find.rightChild == nil { // Node to be removed has 0 child node
		if parent == nil {
			bt.root = nil
		} else {
			if parent.lightChild == find {
				parent.lightChild = nil
			} else {
				parent.rightChild = nil
			}
		}
	} else if find.lightChild != nil && find.rightChild == nil { // Node to be removed has 1 left child node
		if parent == nil {
			bt.root = find.lightChild
		} else {
			if parent.lightChild == find {
				parent.lightChild = find.lightChild
			} else {
				parent.rightChild = find.lightChild
			}
		}
	} else if find.lightChild == nil && find.rightChild != nil { // Node to be removed has 1 right child node
		if parent == nil {
			bt.root = find.rightChild
		} else {
			if parent.lightChild == find {
				parent.lightChild = find.rightChild
			} else {
				parent.rightChild = find.rightChild
			}
		}
	} else if find.lightChild != nil && find.rightChild != nil { // Node to be removed has 2 child node
		leftMost, leftMostParent := findLeftMostNode(find.rightChild, find)
		find.data = leftMost.data
		if leftMostParent == find {
			if leftMost.rightChild != nil {
				leftMostParent.rightChild = leftMost.rightChild
			} else {
				leftMostParent.rightChild = nil
			}
		} else {
			if leftMost.rightChild != nil {
				leftMostParent.lightChild = leftMost.rightChild
			} else {
				leftMostParent.lightChild = nil
			}
		}
	}
	bt.size--
	return nil
}

func findLeftMostNode(ret *Tnode, parent *Tnode) (*Tnode, *Tnode) {
	// if using ret == nil to end recursion calls, it can't get the
	// parent of the leftmost node.
	if ret.lightChild == nil {
		return ret, parent
	}

	return findLeftMostNode(ret.lightChild, ret)
}

func findRightMostNode(ret *Tnode, parent *Tnode) (*Tnode, *Tnode) {
	if ret.rightChild == nil {
		return ret, parent
	}

	return findRightMostNode(ret.rightChild, ret)
}

// Update updates the value associated with key to val. If the tree is empty, returns ErrEmptyTree. If not found, returns ErrNotExist.
func (bt *BSTree) Update(key interface{}, val interface{}) error {
	bt.mux.Lock()
	defer bt.mux.Unlock()
	if bt.root == nil && bt.size == 0 {
		return container.ErrEmptyTree
	}

	find, _ := lookup(bt.root, nil, key)
	if find == nil {
		return container.ErrNotExist
	}
	v := find.data.(container.Interface)
	v.Set(val)
	return nil
}

func inorderTraversal(tn *Tnode, ch chan<- interface{}) {
	if tn == nil {
		return
	}

	inorderTraversal(tn.lightChild, ch)
	ch <- tn.data
	inorderTraversal(tn.rightChild, ch)
}

func preorderTraversal(tn *Tnode, ch chan<- interface{}) {
	if tn == nil {
		return
	}

	ch <- tn.data
	preorderTraversal(tn.lightChild, ch)
	preorderTraversal(tn.rightChild, ch)
}

func postorderTraversal(tn *Tnode, ch chan<- interface{}) {
	if tn == nil {
		return
	}

	postorderTraversal(tn.lightChild, ch)
	postorderTraversal(tn.rightChild, ch)
	ch <- tn.data
}

func levelTraversal(tn *Tnode, ch chan<- interface{}) {
	if tn == nil {
		return
	}

	lq := queue.NewLQueue()
	lq.EnQueue(tn)

	for !lq.Empty() {
		ret := lq.LeQueue().(*Tnode)
		ch <- ret.data
		if ret.lightChild != nil {
			lq.EnQueue(ret.lightChild)
		}
		if ret.rightChild != nil {
			lq.EnQueue(ret.rightChild)
		}
	}
}

// Traversal traverses the tree using predefined traverse method, which specified by the predefined constants,
// such as InorderTrav, PreorderTrav, PostorderTrav, LevelTrav.
func (bt *BSTree) Traversal(TravType int) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		bt.mux.Lock()
		defer bt.mux.Unlock()
		defer close(ch)
		if bt.root == nil && bt.size == 0 {
			return
		}

		switch TravType {
		case InorderTrav:
			inorderTraversal(bt.root, ch)
		case PreorderTrav:
			preorderTraversal(bt.root, ch)
		case PostorderTrav:
			postorderTraversal(bt.root, ch)
		case LevelTrav:
			levelTraversal(bt.root, ch)
		default:
			close(ch)
			return
		}
	}()

	return ch
}

// TravWith traverses the tree using user-defined function. Note that if you use recursion inside anonymouse function you
// must declare it first.
func (bt *BSTree) TravWith(trave TravFunc) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		bt.mux.Lock()
		defer bt.mux.Unlock()
		defer close(ch)
		if bt.root == nil && bt.size == 0 {
			return
		}

		trave(bt.root, ch)
	}()
	return ch
}

// Size returns the size of the tree, which refers to the total number of nodes of tree.
func (bt *BSTree) Size() int {
	bt.mux.Lock()
	defer bt.mux.Unlock()

	return bt.size
}

// Clear empty the tree, it will drop all of data then becomes an empty tree.
func (bt *BSTree) Clear() {
	bt.mux.Lock()
	defer bt.mux.Unlock()

	bt.root = nil
	bt.size = 0
}

// Empty returns true if the tree is empty tree, otherwise false.
func (bt *BSTree) Empty() bool {
	bt.mux.Lock()
	defer bt.mux.Unlock()

	return bt.root == nil && bt.size == 0
}

func getHeight(tn *Tnode) int {
	if tn == nil {
		return -1
	}

	lHeight := getHeight(tn.lightChild)
	rHeight := getHeight(tn.rightChild)

	if lHeight < rHeight {
		return rHeight + 1
	}

	return lHeight + 1
}

// Height returns the height of tree, it also refers to the height of root, which is the number of edges
// on the longest downward path between the root and a leaf. If return -1, it means that the tree is empty.
func (bt *BSTree) Height() int {
	bt.mux.Lock()
	defer bt.mux.Unlock()
	if bt.root == nil && bt.size == 0 {
		return -1
	}

	return getHeight(bt.root)
}

// HeightOf returns the height of specified node of tree, which is the largest number of edges in the
// path from that node to leaf. If data can't be found using key, return -1 and ErrNotExist error.
func (bt *BSTree) HeightOf(key interface{}) (int, error) {
	bt.mux.Lock()
	defer bt.mux.Unlock()
	if bt.root == nil && bt.size == 0 {
		return -1, container.ErrEmptyTree
	}

	find, _ := lookup(bt.root, nil, key)
	if find == nil {
		return -1, container.ErrNotExist
	}

	return getHeight(find), nil
}

func getDepth(from *Tnode, key interface{}) (int, error) {
	if from == nil {
		return -1, container.ErrNotExist
	}
	v := from.data.(container.Interface)
	if v.Find(key) {
		return 0, nil
	}

	var d int
	var err error
	toRight := v.Less(key)
	if toRight {
		d, err = getDepth(from.rightChild, key)
	} else {
		d, err = getDepth(from.lightChild, key)
	}

	if err != nil {
		return -1, container.ErrNotExist
	}
	return d + 1, nil
}

// Depth returns the depth of the tree, which is equal to the height of the tree.
func (bt *BSTree) Depth() int {
	bt.mux.Lock()
	defer bt.mux.Unlock()
	if bt.root == nil && bt.size == 0 {
		return -1
	}

	return getHeight(bt.root)
}

// DepthOf returns the depth of data found by key in the tree, which is the number of
// edges in the path from the root to that node.
func (bt *BSTree) DepthOf(key interface{}) (int, error) {
	bt.mux.Lock()
	defer bt.mux.Unlock()
	if bt.root == nil && bt.size == 0 {
		return -1, container.ErrEmptyTree
	}

	return getDepth(bt.root, key)
}

// FullTree returns true if the tree is full tree, note that beacuse of property of
// binary tree, a full tree is also a complete tree.
func (bt *BSTree) FullTree() bool {
	bt.mux.Lock()
	defer bt.mux.Unlock()

	h := getHeight(bt.root)
	if bt.size == int(math.Pow(2, float64(h+1)))-1 {
		return true
	}
	return false
}

// Compare compares whether two trees is the same, if be the same, return true, otherwise false.
func Compare(bt1, bt2 *BSTree) bool {
	// As long as the data to be inserted are the same, the result of inorder traversal of binary tree is irrelevant with
	// order of insertion.
	ch1, ch2 := bt1.Traversal(PostorderTrav), bt2.Traversal(PostorderTrav)

	for {
		itf1, ok1 := <-ch1
		itf2, ok2 := <-ch2

		if !ok1 || !ok2 {
			return ok1 == ok2
		}

		if !reflect.DeepEqual(itf1, itf2) {
			break
		}

	}
	return false
}

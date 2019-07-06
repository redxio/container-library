package tree_test

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/NzKSO/container"
	"github.com/NzKSO/container/queue"
	"github.com/NzKSO/container/tree"
)

type corp struct {
	id      int
	company string
}

var testCase = []corp{
	{5, "Netflix"},   // 0
	{3, "Google"},    // 1
	{7, "Intel"},     // 2
	{0, "Amazon"},    // 3
	{9, "Dell"},      // 4
	{2, "Facebook"},  // 5
	{10, "Qualcomm"}, // 6
	{4, "Apple"},     // 7
	{6, "IBM"},       // 8
	{1, "Microsoft"}, // 9
	{8, "Samsung"},   // 10
}

func (c *corp) Set(i interface{}) {
	v, ok := i.(string)
	if !ok {
		return
	}
	c.company = v
}

func (c *corp) Less(kv interface{}) bool {
	switch v := kv.(type) {
	case corp:
		if c.id <= v.id {
			return true
		}
		return false
	case *corp:
		if c.id <= (*v).id {
			return true
		}
		return false
	case int:
		if c.id <= v {
			return true
		}
		return false
	case *int:
		if c.id <= *v {
			return true
		}
		return false
	}
	return false
}

func (c *corp) Find(key interface{}) bool {
	switch v := key.(type) {
	// make no sense, used for testing purpose only
	case corp:
		if c.id == v.id {
			return true
		}
		return false
	// comment same as above
	case *corp:
		if c.id == v.id {
			return true
		}
		return false
	case int:
		if c.id == v {
			return true
		}
		return false
	case *int:
		if c.id == *v {
			return true
		}
		return false
	default:
		return false
	}
}

var (
	r = rand.New(rand.NewSource(time.Now().UnixNano()))

	// Predefined index order of various traversal in testCase
	index = [][]int{
		{3, 9, 5, 1, 7, 0, 8, 2, 10, 4, 6}, // Inorder
		{0, 1, 3, 5, 9, 7, 2, 8, 4, 10, 6}, // Preorder
		{9, 5, 3, 7, 1, 8, 10, 6, 4, 2, 0}, // Postorder
		{0, 1, 2, 3, 7, 8, 4, 5, 10, 6, 9}, // Levelorder
	}
)

func createTree(ri []int) (*tree.BSTree, error) {
	bt := tree.NewBSTree()

	for _, iv := range ri {
		err := bt.Insert(&testCase[iv])
		if err != nil {
			return bt, err
		}
	}

	return bt, nil
}

func TestBSTreeInsert(t *testing.T) {
	ri := r.Perm(len(testCase))

	bt, err := createTree(ri)
	if err != nil || bt.Size() != len(testCase) {
		t.Errorf("(%v != nil) or (%v != %v)", err, bt.Size(),
			len(testCase))
	}

	ri = r.Perm(len(testCase))
	for _, iv := range ri {
		err = bt.Insert(&testCase[iv])
		if err != container.ErrDataExists {
			t.Errorf("%v != %v", err, container.ErrDataExists)
		}
	}
}

func TestBSTreeSearch(t *testing.T) {
	ri := r.Perm(len(testCase))
	bt, err := createTree(ri)

	for _, iv := range ri {
		itf, err := bt.Search(testCase[iv].id)
		pval, ok := itf.(*corp)
		t.Log(pval, ok)
		if *pval != testCase[iv] || err == container.ErrNotExist {
			t.Errorf("(%v != %v) or (%v == %v)", *pval, testCase[iv],
				err, container.ErrNotExist)
		}
	}

	for _, iv := range ri {
		itf, err := bt.Search(&testCase[iv].id)
		pval, ok := itf.(*corp)
		t.Log(pval, ok)
		if *pval != testCase[iv] || err == container.ErrNotExist {
			t.Errorf("(%v != %v) or (%v == %v)", itf, testCase[iv],
				err, container.ErrNotExist)
		}
	}

	i := r.Intn(len(testCase))
	// make no sense, used for testing purpose only
	inf, err := bt.Search(testCase[i])
	pval := inf.(*corp)
	if *pval != testCase[i] || err == container.ErrNotExist {
		t.Errorf("(%v != %v) or (%v == %v)", inf, testCase[i],
			err, container.ErrNotExist)
	}

	// comment same as above
	i = r.Intn(len(testCase))
	inf, err = bt.Search(&testCase[i])
	pval = inf.(*corp)
	if *pval != testCase[i] || err == container.ErrNotExist {
		t.Errorf("(%v != %v) or (%v == %v)", inf, testCase[i],
			err, container.ErrNotExist)
	}

	overflow := len(testCase) + r.Intn(len(testCase))
	inf, err = bt.Search(overflow)
	if inf != nil || err != container.ErrNotExist {
		t.Errorf("(%v != nil) or (%v != %v)", inf, err,
			container.ErrNotExist)
	}
}

func TestBSTreeDelete(t *testing.T) {

	ri := r.Perm(len(testCase))
	t.Log("ri:", ri)

	bt, err := createTree(ri)
	if err != nil {
		t.Errorf("%v != nil", err)
	}

	ri = r.Perm(len(testCase))
	t.Log("ri:", ri)
	for _, iv := range ri {
		err := bt.Delete(testCase[iv].id)
		if err != nil {
			t.Errorf("%v != nil", testCase[iv].id)
		}

		ret, err := bt.Search(testCase[iv].id)
		if ret != nil || err == nil {
			t.Errorf("(%v != nil) or (%v != %v)", ret, err, nil)
		}
	}
}

func TestBSTreeUpdate(t *testing.T) {
	ri := r.Perm(len(testCase))
	bt, _ := createTree(ri)

	for _, iv := range ri {
		before, _ := bt.Search(iv)
		v, ok := before.(*corp)
		t.Log(v.company, ok)
		newStr := strings.ToUpper(v.company)

		bt.Update(iv, newStr)

		after, _ := bt.Search(iv)
		v, ok = after.(*corp)
		t.Log(v.company, ok)

		if v.company != newStr {
			t.Errorf("%v != %v", v.company, newStr)
		}
	}
}

func TestBSTreeTraversal(t *testing.T) {
	defidx := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	bt, _ := createTree(defidx)

	for i := 0; i <= 3; i++ {
		ch := bt.Traversal(tree.TraversalType(i))
		go func(n int, ch <-chan interface{}) {
			var j int
			for itf := range ch {
				pn, ok := itf.(*corp)
				t.Log(pn, ok)
				if *pn != testCase[index[n][j]] {
					t.Errorf("%v != %v", *pn, testCase[index[n][j]])
				}
				j++
			}
		}(i, ch)
	}
}

func TestBSTreeTravWith(t *testing.T) {
	defidx := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	bt, _ := createTree(defidx)

	levelTraversal := func(tn *tree.Tnode, ch chan<- interface{}) {
		if tn == nil {
			return
		}

		lq := queue.NewLQueue()
		lq.EnQueue(tn)

		for !lq.Empty() {
			ret := lq.LeQueue().(*tree.Tnode)
			ch <- ret.GetData()
			lchild := ret.GetLchild()
			if lchild != nil {
				lq.EnQueue(lchild)
			}

			rchild := ret.GetRchild()
			if rchild != nil {
				lq.EnQueue(rchild)
			}
		}
	}

	ch := bt.TravWith(tree.TravFunc(levelTraversal))
	var next int
	for itf := range ch {
		pn, ok := itf.(*corp)
		t.Log(pn, ok)
		if *pn != testCase[index[3][next]] {
			t.Errorf("%v != %v", *pn, index[3][next])
		}
		next++
	}

	var inorderTraversal tree.TravFunc
	inorderTraversal = func(tn *tree.Tnode, ch chan<- interface{}) {
		if tn == nil {
			return
		}

		inorderTraversal(tn.GetLchild(), ch)
		ch <- tn.GetData()
		inorderTraversal(tn.GetRchild(), ch)
	}

	next = 0
	ch = bt.TravWith(inorderTraversal)
	for itf := range ch {
		pn, ok := itf.(*corp)
		t.Log(pn, ok)
		if *pn != testCase[index[0][next]] {
			t.Errorf("%v != %v", *pn, index[0][next])
		}
		next++
	}
}

func TestBSTreeSize(t *testing.T) {
	ri := r.Perm(len(testCase))
	bt, _ := createTree(ri)

	if bt.Size() != len(testCase) {
		t.Errorf("%v != %v", bt.Size(), len(testCase))
	}

	var c int = 1
	for _, iv := range ri {
		err := bt.Delete(testCase[iv].id)
		if err != nil || bt.Size() != len(testCase)-c {
			t.Errorf("(err != %v) or (%v != %v)", err, bt.Size(), len(testCase)-c)
		}
		c++
	}
}

func TestBSTreeClear(t *testing.T) {
	ri := r.Perm(len(testCase))
	bt, _ := createTree(ri)

	bt.Clear()
	if bt.Size() != 0 {
		t.Errorf("%v != 0", bt.Size())
	}
}

func TestBSTreeEmpty(t *testing.T) {
	ri := r.Perm(len(testCase))
	bt, _ := createTree(ri)

	bt.Clear()
	if bt.Empty() != true {
		t.Errorf("%v != true", bt.Empty())
	}
}

func TestBSTreeHeight(t *testing.T) {
	idx := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	height := 4
	bt, _ := createTree(idx)

	h := bt.Height()
	if h != height {
		t.Errorf("(%v != %v)", h, height)
	}
}

func TestBSTreeHeightOf(t *testing.T) {
	idx := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	bt, _ := createTree(idx)

	type heightOfNode struct {
		key    int
		height int
	}

	sliOfHeight := []heightOfNode{
		{5, 4},
		{3, 3},
		{7, 2},
		{0, 2},
		{4, 0},
		{6, 0},
		{9, 1},
		{2, 1},
		{8, 0},
		{10, 0},
		{1, 0},
	}
	for i, _ := range sliOfHeight {
		h, _ := bt.HeightOf(sliOfHeight[i].key)
		if sliOfHeight[i].height != h {
			t.Errorf("%v != %v\n", sliOfHeight[i].height, h)
		}
	}
}

func TestBSTreeDepth(t *testing.T) {
	idx := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	bt, _ := createTree(idx)
	depth := 4

	if bt.Depth() != depth {
		t.Errorf("%v != %v\n", bt.Depth(), depth)
	}
}

func TestBSTreeDepthOf(t *testing.T) {
	idx := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	bt, _ := createTree(idx)

	type depthOfNode struct {
		key   int
		depth int
	}

	sliOfDepth := []depthOfNode{
		{5, 0},
		{3, 1},
		{7, 1},
		{0, 2},
		{4, 2},
		{6, 2},
		{9, 2},
		{2, 3},
		{8, 3},
		{10, 3},
		{1, 4},
	}

	for i, _ := range sliOfDepth {
		d, _ := bt.DepthOf(sliOfDepth[i].key)
		if sliOfDepth[i].depth != d {
			t.Errorf("%v != %v\n", sliOfDepth[i].depth, d)
		}
	}
}

func TestBSTreeFullTree(t *testing.T) {
	idx := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	bt, _ := createTree(idx)

	if bt.FullTree() != false {
		t.Errorf("%v != false\n", bt.FullTree())
	}
}

func TestBSTreeCompare(t *testing.T) {
	predefIdx1 := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	predefIdx2 := predefIdx1

	bt1, _ := createTree(predefIdx1)
	bt2, _ := createTree(predefIdx2)

	same := tree.Compare(bt1, bt2)
	if !same {
		t.Errorf("Two trees is the same but return false")
	}

	ri1 := r.Perm(len(testCase))
	ri2 := r.Perm(len(testCase))
	t.Logf("ri1: %v\n", ri1)
	t.Logf("ri2: %v\n", ri2)

	id := make([]int, len(testCase))
	//ri2 = ri2[:len(ri2)-3]
	for i, v := range ri1 {
		id[i] = testCase[v].id
	}
	t.Logf("id: %v\n", id)

	for i, v := range ri2 {
		id[i] = testCase[v].id
	}
	t.Logf("id: %v\n", id)

	t.Log(compareIntSlice(ri1, ri2))

	bt1, _ = createTree(ri1)
	bt2, _ = createTree(ri2)

	same = tree.Compare(bt1, bt2)
	t.Log(same)

	if compareIntSlice(ri1, ri2) && !same {
		t.Errorf("Two trees is the same but return false")
	}
	if !compareIntSlice(ri1, ri2) && same {
		t.Errorf("Two trees is not the same but return true")
	}
}

func compareIntSlice(s1, s2 []int) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i, _ := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

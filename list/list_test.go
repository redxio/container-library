package list_test

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/NzKSO/container"
	"github.com/NzKSO/container/list"
)

type corp struct {
	id      int
	company string
}

var testCase = []corp{
	{22, "General Electric"},
	{19, "Siemens"},
	{2, "Alphabet"},
	{6, "Softbank"},
	{30, "Wells Fargo"},
	{7, "Intel"},
	{33, "Morgan Stanley"},
	{23, "Amazon"},
	{17, "Goldman Sachs Group"},
	{29, "Dell"},
	{0, "Facebook"},
	{24, "Coca-Cola"},
	{12, "Qualcomm"},
	{32, "Verizon Communications"},
	{28, "Apple"},
	{3, "IBM"},
	{11, "Microsoft"},
	{18, "Samsung Electronics"},
	{8, "Booking"},
	{1, "Twitter"},
	{13, "Berkshire Hathaway"},
	{35, "JPMorgan Chase"},
	{31, "Bank of America"},
	{14, "Netflix"},
	{5, "Toyota Motor"},
	{21, "AT&T"},
	{27, "Citigroup"},
	{9, "Wal-Mart Stores"},
	{16, "Uber"},
	{20, "Royal Dutch Shell"},
	{34, "Comcast"},
	{4, "Cisco Systems"},
	{26, "Walt Disney"},
	{15, "Oracle"},
	{10, "Boeing"},
	{25, "HP"},
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

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
	// make no sense, used for test purpose only
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

func createAndFillList(ri []int) *list.LinkedList {
	ll := list.NewLinkedList()
	for _, iv := range ri {
		ll.Insert(&testCase[iv])
	}

	return ll
}

func TestListInsert(t *testing.T) {
	ri := r.Perm(len(testCase))
	ll := createAndFillList(ri)
	if ll.Empty() {
		t.Errorf("Tree is empty? %v", ll.Empty())
	}

	if ll.Size() != len(testCase) {
		t.Errorf("%v != %v", ll.Size(), len(testCase))
	}
}

func TestListDelete(t *testing.T) {
	ri := r.Perm(len(testCase))
	ll := list.NewLinkedList()
	if err := ll.Delete(testCase[r.Intn(len(testCase))].id); err != container.ErrEmptyList {
		t.Errorf("%v != %v", err, container.ErrEmptyList)
	}
	for _, iv := range ri {
		ll.Insert(&testCase[iv])
	}

	ri = r.Perm(len(testCase))
	for _, iv := range ri {
		if err := ll.Delete(testCase[iv].id); err != nil {
			t.Errorf("%v != nil", err)
		}
	}
}

func TestListSearch(t *testing.T) {
	ri := r.Perm(len(testCase))
	ll := list.NewLinkedList()
	itf, err := ll.Search(testCase[r.Intn(len(testCase))].id)
	if itf != nil || err != container.ErrEmptyList {
		t.Errorf("%v != nil or %v != %v", itf, err, container.ErrEmptyList)
	}
	for _, iv := range ri {
		ll.Insert(&testCase[iv])
	}

	ri = r.Perm(len(testCase))
	for _, iv := range ri {
		itf, err = ll.Search(testCase[iv].id)
		if itf.(*corp) != &testCase[iv] || err != nil {
			t.Errorf("%v != %p or %v != nil", itf.(*corp), &testCase[iv], err)
		}
	}
}

func TestListUpdate(t *testing.T) {
	ri := r.Perm(len(testCase))
	ll := createAndFillList(ri)
	for _, iv := range ri {
		testStr := strings.ToUpper(testCase[iv].company)
		key := testCase[iv].id
		ll.Update(key, testStr)
		itf, _ := ll.Search(key)
		if itf.(*corp).company != testStr {
			t.Errorf("%v != %v", itf.(*corp).company, testStr)
		}
	}
}

func TestListTraversal(t *testing.T) {
	ri := r.Perm(len(testCase))
	ll := createAndFillList(ri)
	ch := ll.Traversal()
	i := len(ri) - 1
	for v := range ch {
		if v.(*corp) != &testCase[ri[i]] {
			t.Errorf("%v != %v", v.(*corp), &testCase[ri[i]])
		}
		i--
	}
}

func TestListReverse(t *testing.T) {
	ri := r.Perm(len(testCase))
	ll := createAndFillList(ri)
	ll.Reverse()
	ch := ll.Traversal()
	i := 0
	for v := range ch {
		if v.(*corp) != &testCase[ri[i]] {
			t.Errorf("%v != %v", v.(*corp), &testCase[ri[i]])
		}
		i++
	}
}

func TestListEmpty(t *testing.T) {
	ri := r.Perm(len(testCase))
	ll := list.NewLinkedList()
	if !ll.Empty() {
		t.Errorf("List is empty? %v", ll.Empty())
	}

	ll.Insert(&testCase[ri[0]])
	if ll.Empty() {
		t.Errorf("List is empty? %v", ll.Empty())
	}
	ll.Delete(testCase[ri[0]].id)
	if !ll.Empty() {
		t.Errorf("List is empty? %v", ll.Empty())
	}

	for _, iv := range ri {
		ll.Delete(testCase[iv].id)
	}
	if !ll.Empty() {
		t.Errorf("List is empty? %v", ll.Empty())
	}
}

func TestListSize(t *testing.T) {
	ll := list.NewLinkedList()
	if ll.Size() != 0 {
		t.Errorf("%v != 0", ll.Size())
	}

	ri := r.Perm(len(testCase))
	size := 0
	for _, iv := range ri {
		ll.Insert(&testCase[iv])
		size++
		if ll.Size() != size {
			t.Errorf("%v != %v", ll.Size(), size)
		}
	}

	for _, iv := range ri {
		ll.Delete(testCase[iv].id)
		size--
		if ll.Size() != size {
			t.Errorf("%v != %v", ll.Size(), size)
		}
	}
}

func TestListClear(t *testing.T) {
	ri := r.Perm(len(testCase))
	ll := createAndFillList(ri)
	ll.Clear()
	if !ll.Empty() || ll.Size() != 0 {
		t.Errorf("%v or %v != 0", !ll.Empty(), ll.Size())
	}
}

func TestListSort(t *testing.T) {
	ri := r.Perm(len(testCase))
	ll := createAndFillList(ri)

	ll.Sort()
	ch := ll.Traversal()
	id := 0
	for itf := range ch {
		data := itf.(*corp)
		if data.id != id {
			t.Errorf("%v != %v", data.id, id)
		}
		id++
	}
}

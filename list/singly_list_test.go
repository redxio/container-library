package list_test

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/NzKSO/container"
	"github.com/NzKSO/container/list"
	"github.com/NzKSO/container/testdata"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func createAndFillList(ri []int) *list.SinglyList {
	ll := list.NewSinglyList()
	for _, iv := range ri {
		ll.Insert(&testdata.TestCases[iv])
	}

	return ll
}

func TestListInsert(t *testing.T) {
	ri := r.Perm(len(testdata.TestCases))
	ll := createAndFillList(ri)
	if ll.Empty() {
		t.Errorf("Tree is empty? %v", ll.Empty())
	}

	if ll.Size() != len(testdata.TestCases) {
		t.Errorf("%v != %v", ll.Size(), len(testdata.TestCases))
	}
}

func TestListDelete(t *testing.T) {
	ri := r.Perm(len(testdata.TestCases))
	ll := list.NewSinglyList()
	ll.NumPerGoroutine = 10

	if err := ll.Delete(testdata.TestCases[r.Intn(len(testdata.TestCases))].ID); err != container.ErrEmptyList {
		t.Errorf("%v != %v", err, container.ErrEmptyList)
	}
	for _, iv := range ri {
		ll.Insert(&testdata.TestCases[iv])
	}

	ri = r.Perm(len(testdata.TestCases))
	for _, iv := range ri {
		if err := ll.Delete(testdata.TestCases[iv].ID); err != nil {
			t.Errorf("%v != nil", err)
		}
	}
}

func TestListSearch(t *testing.T) {
	ri := r.Perm(len(testdata.TestCases))
	ll := list.NewSinglyList()
	ll.NumPerGoroutine = 9

	itf, err := ll.Search(testdata.TestCases[r.Intn(len(testdata.TestCases))].ID)
	if itf != nil || err != container.ErrEmptyList {
		t.Errorf("%v != nil or %v != %v", itf, err, container.ErrEmptyList)
	}
	for _, iv := range ri {
		ll.Insert(&testdata.TestCases[iv])
	}

	ri = r.Perm(len(testdata.TestCases))
	for _, iv := range ri {
		itf, err = ll.Search(testdata.TestCases[iv].ID)
		if itf.(*testdata.Corp) != &testdata.TestCases[iv] || err != nil {
			t.Errorf("%v != %p or %v != nil", itf.(*testdata.Corp), &testdata.TestCases[iv], err)
		}
	}
}

func TestListUpdate(t *testing.T) {
	ri := r.Perm(len(testdata.TestCases))
	ll := createAndFillList(ri)
	ll.NumPerGoroutine = 11

	for _, iv := range ri {
		testStr := strings.ToUpper(testdata.TestCases[iv].Name)
		key := testdata.TestCases[iv].ID
		ll.Update(key, testStr)
		itf, _ := ll.Search(key)
		if itf.(*testdata.Corp).Name != testStr {
			t.Errorf("%v != %v", itf.(*testdata.Corp).Name, testStr)
		}
	}
}

func TestListTraversal(t *testing.T) {
	ri := r.Perm(len(testdata.TestCases))
	ll := createAndFillList(ri)
	ch := ll.Traversal()
	i := len(ri) - 1
	for v := range ch {
		if v.(*testdata.Corp) != &testdata.TestCases[ri[i]] {
			t.Errorf("%v != %v", v.(*testdata.Corp), &testdata.TestCases[ri[i]])
		}
		i--
	}
}

func TestListReverse(t *testing.T) {
	ri := r.Perm(len(testdata.TestCases))
	ll := createAndFillList(ri)
	ll.NumPerGoroutine = 12

	ll.Reverse()
	ch := ll.Traversal()
	i := 0
	for v := range ch {
		if v.(*testdata.Corp) != &testdata.TestCases[ri[i]] {
			t.Errorf("%v != %v", v.(*testdata.Corp), &testdata.TestCases[ri[i]])
		}
		i++
	}
}

func TestListEmpty(t *testing.T) {
	ri := r.Perm(len(testdata.TestCases))
	ll := list.NewSinglyList()
	if !ll.Empty() {
		t.Errorf("List is empty? %v", ll.Empty())
	}

	ll.Insert(&testdata.TestCases[ri[0]])
	if ll.Empty() {
		t.Errorf("List is empty? %v", ll.Empty())
	}
	ll.Delete(testdata.TestCases[ri[0]].ID)
	if !ll.Empty() {
		t.Errorf("List is empty? %v", ll.Empty())
	}

	for _, iv := range ri {
		ll.Delete(testdata.TestCases[iv].ID)
	}
	if !ll.Empty() {
		t.Errorf("List is empty? %v", ll.Empty())
	}
}

func TestListSize(t *testing.T) {
	ll := list.NewSinglyList()
	if ll.Size() != 0 {
		t.Errorf("%v != 0", ll.Size())
	}

	ri := r.Perm(len(testdata.TestCases))
	size := 0
	for _, iv := range ri {
		ll.Insert(&testdata.TestCases[iv])
		size++
		if ll.Size() != size {
			t.Errorf("%v != %v", ll.Size(), size)
		}
	}

	for _, iv := range ri {
		ll.Delete(testdata.TestCases[iv].ID)
		size--
		if ll.Size() != size {
			t.Errorf("%v != %v", ll.Size(), size)
		}
	}
}

func TestListReset(t *testing.T) {
	ri := r.Perm(len(testdata.TestCases))
	ll := createAndFillList(ri)
	ll.Reset()
	if !ll.Empty() || ll.Size() != 0 {
		t.Errorf("%v or %v != 0", !ll.Empty(), ll.Size())
	}
}

func TestListSort(t *testing.T) {
	ri := r.Perm(len(testdata.TestCases))
	ll := createAndFillList(ri)

	ll.Sort()
	ch := ll.Traversal()
	ID := 0
	for itf := range ch {
		data := itf.(*testdata.Corp)
		if data.ID != ID {
			t.Errorf("%v != %v", data.ID, ID)
		}
		ID++
	}
}

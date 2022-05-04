package skipgo

import (
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	m := NewMapOrdered[int, int]()

	if m.Len() != 0 {
		t.Fatalf("Len got %#v want %0v", m.Len(), 0)
	}

	v, ok := m.Contains(1)
	if ok != false {
		t.Errorf("ok got %#v want %#v", ok, false)
	}

	if v != 0 {
		t.Errorf("v got %#v want %#v", v, 0)
	}

	for i := 1; i <= 10; i++ {
		old, overwrite := m.Store(i, i*100)
		if overwrite != false {
			t.Errorf("Store %d overwrite got %#v want %#v", i, overwrite, false)
		}
		if old != 0 {
			t.Errorf("Store %d old got %#v want %#v", i, old, 0)
		}

		if m.Len() != i {
			t.Errorf("Store %d Len got %#v want %#v", i, m.Len(), i)
		}

	}

	for i := 1; i <= 10; i++ {
		v, ok = m.Contains(i)
		if ok != true {
			t.Errorf("Contains %d ok got %#v want %#v", i, ok, true)
		}
		if v != i*100 {
			t.Errorf("Contains %d v got %#v want %#v", i, v, i*100)
		}
	}

	for i := 1; i <= 10; i++ {
		if i%2 == 1 {
			v, ok = m.Delete(i)
			if ok != true {
				t.Errorf("Delete %d ok got %#v want %#v", i, ok, true)
			}
			if v != i*100 {
				t.Errorf("Delete %d v got %#v want %#v", i, v, i*100)
			}
		}
	}

	if m.Len() != 5 {
		t.Errorf("After deletes Len got %#v want %#v", m.Len(), 5)
	}

	// Delete already deleted key, should fail
	v, ok = m.Delete(1)
	if ok != false {
		t.Errorf("Delete already deleted 1 ok got %#v want %#v", ok, false)
	}
	if v != 0 {
		t.Errorf("Delete already deleted 1 v got %#v want %#v", v, 0)
	}

	if m.Len() != 5 {
		t.Errorf("After repeated delete Len got %#v want %#v", m.Len(), 5)
	}

	// Delete key that never existed, should fail.
	v, ok = m.Delete(20)
	if ok != false {
		t.Errorf("Delete non existant got %#v want %#v", ok, false)
	}
	if v != 0 {
		t.Errorf("Delete non existant got %#v want %#v", v, 0)
	}

	keys := []int{}
	vals := []int{}

	it := m.Iterator()
	for {
		p, ok := it.Next()
		if !ok {
			break
		}
		keys = append(keys, p.Key)
		vals = append(vals, p.Val)
	}

	expect := []int{2, 4, 6, 8, 10}
	if !reflect.DeepEqual(keys, expect) {
		t.Errorf("keys got %#v want %#v", keys, expect)
	}

	expect = []int{200, 400, 600, 800, 1000}
	if !reflect.DeepEqual(vals, expect) {
		t.Errorf("vals got %#v want %#v", vals, expect)
	}

	for i := 1; i <= 10; i++ {
		old, overwrite := m.Store(i, i*100+1)

		shouldOverwrite := i%2 == 0
		if overwrite != shouldOverwrite {
			t.Errorf("%d overwrite got %#v want %#v", i, overwrite, shouldOverwrite)
		}

		if shouldOverwrite {
			if old != i*100 {
				t.Errorf("%d old got %#v want %#v", i, old, i*100)
			}
		}
	}

	if m.Len() != 10 {
		t.Errorf("Len got %#v want %#v", m.Len(), 10)
	}

	keys = []int{}
	vals = []int{}

	kit := Keys(m.Iterator())
	for {
		k, ok := kit.Next()
		if !ok {
			break
		}
		keys = append(keys, k)
	}

	expect = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if !reflect.DeepEqual(keys, expect) {
		t.Errorf("keys got %#v want %#v", keys, expect)
	}

	vit := Vals(m.Iterator())
	for {
		v, ok := vit.Next()
		if !ok {
			break
		}
		vals = append(vals, v)
	}

	expect = []int{101, 201, 301, 401, 501, 601, 701, 801, 901, 1001}
	if !reflect.DeepEqual(vals, expect) {
		t.Errorf(" got %#v want %#v", vals, expect)
	}
}

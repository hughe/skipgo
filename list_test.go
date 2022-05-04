package skipgo

import (
	"reflect"
	"testing"
)

func TestList(t *testing.T) {
	l := NewList[int](func(x, y int) int { return x - y })

	if l.Len() != 0 {
		t.Errorf("Len got %#v want %#v", l, 0)
	}

	_, found := l.Find(42)
	if found != false {
		t.Errorf("Find before insert got %#v want %#v", found, false)
	}

	l.Insert(42)

	if l.Len() != 1 {
		t.Errorf("Len got %#v want %#v", l.Len(), 1)
	}

	if l.h.root.nexts[0].item != 42 {
		t.Errorf("First item got %#v want %#v", l.h.root.nexts[0].item, 42)
	}

	if l.h.root.nexts[0].nexts[0] != &l.h.root {
		t.Errorf("List ends after first item, got %#v want %#v", l.h.root.nexts[0].nexts[0], &l.h.root)
	}

	x, found := l.Find(42)
	if found != true {
		t.Errorf("Find after first insert got %#v want %#v", found, true)
	}
	if x != 42 {
		t.Errorf("Found value got %#v want %#v", x, 42)
	}

	_, found = l.Find(41)
	if found != false {
		t.Errorf("Fiind 41 got %#v want %#v", found, false)
	}

}

type op int8

const (
	insert op = iota
	find
	delete
	list
)

func TestListTable(t *testing.T) {
	tab := []struct {
		o     op
		item  int
		found bool
		items []int
	}{
		{
			o:     find,
			item:  42,
			found: false,
		},
		{
			o:    insert,
			item: 42,
		},
		{
			o:     find,
			item:  42,
			found: true,
		},
		{
			o:     find,
			item:  41,
			found: false,
		},
		{
			o:     list,
			items: []int{42},
		},
		{
			o:     find,
			item:  78,
			found: false,
		},
		{
			o:    insert,
			item: 78,
		},
		{
			o:     find,
			item:  78,
			found: true,
		},
		{
			o:     list,
			items: []int{42, 78},
		},
	}

	l := NewList[int](func(x, y int) int { return x - y })

	for i, x := range tab {
		switch x.o {
		case insert:
			l.Insert(x.item)

		case find:
			item, found := l.Find(x.item)
			if found != x.found {
				t.Errorf("%d %d found got %#v want %#v", i, x.o, found, x.found)
			}
			if found {
				if item != x.item {
					t.Errorf("%d %d item got %#v want %#v", i, x.o, item, x.item)
				}
			}

		case list:
			items := []int{}
			iter := l.Iterator()
			for {
				item, ok := iter.Next()
				if !ok {
					break
				}
				items = append(items, item)
			}
			if !reflect.DeepEqual(items, x.items) {
				t.Errorf("%d list got %#v want %#v", i, items, x.items)
			}
		}
	}
}

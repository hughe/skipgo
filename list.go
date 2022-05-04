package skipgo

import (
	"golang.org/x/exp/constraints"
)

type List[T any] struct {
	h head[T]
}

func NewListOrdered[T constraints.Ordered](opts ...Option) *List[T] {
	return NewList[T](compareOrdered[T], opts...)
}

func NewList[T any](compare func(T, T) int, opts ...Option) *List[T] {
	ret := &List[T]{}
	ret.h.init(compare, opts)
	return ret
}

func (l *List[T]) Len() int {
	return l.h.len
}

// Insert inserts item into the list.  If there are other items in the
// list that compare equal to item, then item will be inserted before
// those items.
func (l *List[T]) Insert(item T) {
	l.doInsert(item, findFirst)
}

// Insert inserts item into the list.  If there are other items in the
// list that compare equal to item, then item will be inserted after
// those items.
func (l *List[T]) InsertAfter(item T) {
	l.doInsert(item, findLast)
}

func (l *List[T]) doInsert(item T, mode findMode) {
	var scratch [maxHeight]**node[T]
	_, ptrs := l.h.find(item, mode, scratch[:])

	n := l.h.allocateNode()
	n.item = item

	splice(n, ptrs)

	l.h.len++
}

// Delete tries to delte the first item in the list that compares
// equal to T.
//
// If an item is deleted then it returns that item and true.  If no
// matching item is found then it returns the zero value for T and
// false.
func (l *List[T]) Delete(item T) (T, bool) {
	return l.doDelete(item, findFirst)
}

func (l *List[T]) doDelete(item T, mode findMode) (T, bool) {
	var scratch [maxHeight]**node[T]
	n, ptrs := l.h.find(item, mode, scratch[:])

	if n == nil {
		var zero T
		return zero, false
	}

	unsplice(n, ptrs)
	return n.item, true
}

// Find tries to find the first item in the list that compares equal to T.
//
// If an item is found then it returns that item and true.  If no
// matching item is found then it returns the zero value for T and
// false.
func (l *List[T]) Find(item T) (T, bool) {
	n, _ := l.h.find(item, findFirst, nil)
	if n == nil {
		var zero T
		return zero, false
	}
	return n.item, true
}

// Iterator returns an iterator over all the items in the list,
// starting with one that compares smallest.
func (l *List[T]) Iterator() Iter[T] {
	return l.h.iterator()
}

// IteratorStartingAt an iterator over all the items in the list
// greater than or equal to T.
//
// If item is found then it returns an iterator whose first value will
// be item and true.  If itemis not found then it returns an iterator
// whose first value will be the next greater value and false.
func (l *List[T]) IteratorStartingAt(item T) (Iter[T], bool) {
	var scratch [maxHeight]**node[T]

	n, ptrs := l.h.find(item, findFirst, scratch[:])
	if n != nil {
		// n is the first node that compares equal
		return l.h.iteratorAt(n), true
	}

	// No exact match was found, return an iterator starting with the
	// first node whose item is greater than T.
	return l.h.iteratorAt(*ptrs[0]), false
}

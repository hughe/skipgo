package skipgo

import "golang.org/x/exp/constraints"

type Map[K any, V any] struct {
	h head[Pair[K, V]]
}

func makePairCompare[K any, V any](compare func(K, K) int) func(Pair[K, V], Pair[K, V]) int {
	return func(a Pair[K, V], b Pair[K, V]) int {
		return compare(a.Key, b.Key)
	}
}

func compareOrdered[K constraints.Ordered](a K, b K) int {
	switch {
	case a < b:
		return -1
	case a == b:
		return 0
	default:
		return 1
	}
}

func NewMapOrdered[K constraints.Ordered, V any](opts ...Option) *Map[K, V] {
	return NewMap[K, V](compareOrdered[K], opts...)
}

func NewMap[K any, V any](compare func(K, K) int, opts ...Option) *Map[K, V] {
	ret := &Map[K, V]{}

	ret.h.init(makePairCompare[K, V](compare), opts)

	return ret
}

func (m *Map[K, V]) Contains(key K) (V, bool) {
	var zeroV V
	n, _ := m.h.find(Pair[K, V]{Key: key, Val: zeroV}, findAny, nil)
	if n == nil {
		return zeroV, false
	}
	return n.item.Val, true
}

// Store the pair (key, val) in the map.
//
// Return the old value and true if the there was a matching key in
// the map, or the zero value and false if a new entry was added to
// the map.
func (m *Map[K, V]) Store(key K, val V) (V, bool) {
	var scratch [maxHeight]**node[Pair[K, V]]
	n, ptrs := m.h.find(Pair[K, V]{Key: key, Val: val}, findAny, scratch[:])

	if n != nil {
		// We found a node with matching key
		oldVal := n.item.Val
		n.item.Val = val
		return oldVal, true
	}

	n = m.h.allocateNode()
	n.item = Pair[K, V]{Key: key, Val: val}

	splice(n, ptrs)

	m.h.len++

	var zeroV V
	return zeroV, false
}

// Delete key from the map.
//
// If an entry with key key existed, then return the value associated
// with key and true.  If the entry did not exist, then return the
// zero value and false.
func (m *Map[K, V]) Delete(key K) (V, bool) {
	var scratch [maxHeight]**node[Pair[K, V]]
	var zeroV V

	n, ptrs := m.h.find(Pair[K, V]{Key: key, Val: zeroV}, findAny, scratch[:])

	if n == nil {
		return zeroV, false
	}

	unsplice(n, ptrs)

	m.h.len--

	return n.item.Val, true
}

func (m *Map[K, V]) Len() int {
	return m.h.len
}

func (m *Map[K, V]) Iterator() Iter[Pair[K, V]] {
	return m.h.iterator()
}

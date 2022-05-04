package skipgo

type Pair[K any, V any] struct {
	Key K
	Val V
}

type keyIter[K any, V any] struct {
	it Iter[Pair[K, V]]
}

func (ki *keyIter[K, V]) Next() (K, bool) {
	p, ok := ki.it.Next()
	return p.Key, ok
}

// Return an iterator that will iterate over the the Keys in interator
// it.
func Keys[K any, V any](it Iter[Pair[K, V]]) Iter[K] {
	return &keyIter[K, V]{it: it}
}

type valIter[K any, V any] struct {
	it Iter[Pair[K, V]]
}

func (vi *valIter[K, V]) Next() (V, bool) {
	p, ok := vi.it.Next()
	return p.Val, ok
}

// Return an iterator that will iterate over the the Vals in interator
// it.
func Vals[K any, V any](it Iter[Pair[K, V]]) Iter[V] {
	return &valIter[K, V]{it: it}
}

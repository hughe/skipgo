package skipgo

type node[T any] struct {
	item  T
	nexts []*node[T]
}

func makeNode[T any](height int) *node[T] {
	n := &node[T]{}
	n.init(height)
	return n
}

func (n *node[T]) init(height int) {
	n.nexts = make([]*node[T], height)
}

func splice[T any](n *node[T], ptrs []**node[T]) {
	for i := range n.nexts {
		n.nexts[i] = *ptrs[i]
		*ptrs[i] = n
	}
}

func unsplice[T any](n *node[T], ptrs []**node[T]) {
	for i := range n.nexts {
		*ptrs[i] = n.nexts[i]
	}
}

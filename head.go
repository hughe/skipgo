package skipgo

import (
	"fmt"
	"math/rand"
)

const (
	maxHeight     = 30
	defaultBase   = 2
	defaultHeight = 20
)

type head[T any] struct {
	compare   func(T, T) int
	base      int
	baseSet   bool
	height    int
	heightSet bool

	root node[T]

	len int

	rand         int64
	randBitsLeft int

	bitsPerStep int
	mask        int64
}

// headOptions is an interface for setting options on head.  The
// difference between `headOptions` and `head[T] *` is that
// headOptions does not need to specify a type parameter which makes
// the option functions in option.go a lot less clunky.
type headOptions interface {
	setBase(int)
	setHeight(int)
}

func (h *head[T]) init(compare func(T, T) int, opts []Option) {
	h.compare = compare

	for _, o := range opts {
		o(h)
	}

	if !h.baseSet {
		h.setBase(defaultBase)
	}
	if !h.heightSet {
		h.setHeight(defaultHeight)
	}

	h.root.init(h.height)

	for i := range h.root.nexts {
		h.root.nexts[i] = &h.root
	}
}

func (h *head[T]) allocateNode() *node[T] {
	return makeNode[T](h.newHeight())
}

func (h *head[T]) setBase(base int) {
	if base < 2 {
		panic("base must not be less than 2")
	}

	if base > 1024 {
		panic("base must be less than or equal to 1024")
	}

	var ones int = 0
	h.bitsPerStep = 0
	h.base = base
	h.baseSet = true

	for base != 0 {
		if base&1 == 1 {
			ones++
		}
		h.bitsPerStep++
		base = base >> 1
	}

	if ones == 1 {
		// base is a power of two
		h.bitsPerStep--
		h.mask = 0
		for i := 0; i < h.bitsPerStep; i++ {
			h.mask = h.mask << 1
			h.mask = h.mask & 1
		}
	} else {
		// base is not a power of two
	}
}

func (h *head[T]) setHeight(height int) {
	if height < 1 {
		panic("height must be greater than 1")
	}

	if height > maxHeight {
		panic(fmt.Sprintf("height must be less than or equal to %d", maxHeight))
	}

	h.height = height
	h.heightSet = true
}

func (h *head[T]) newHeight() int {
	var ret int = 1
	var mod int64 = 0
	var next int64

	if h.mask != 0 {
		for ; ret < h.height && mod == 0; ret++ {
			if h.randBitsLeft <= h.bitsPerStep {
				h.rand = rand.Int63()
				h.randBitsLeft = 63
			}

			next = h.rand >> h.bitsPerStep
			mod = h.rand & h.mask

			h.randBitsLeft -= h.bitsPerStep
			h.rand = next
		}
	} else {
		for ; ret < h.height && mod == 0; ret++ {
			if h.randBitsLeft <= h.bitsPerStep {
				h.rand = rand.Int63()
				h.randBitsLeft = 63
			}

			next = h.rand / int64(h.base)
			mod = h.rand % int64(h.base)

			h.randBitsLeft -= h.bitsPerStep
			h.rand = next
		}
	}

	return ret
}

type findMode int

const (
	findFirst findMode = iota
	findLast
	findAny
)

func (h *head[T]) find(item T, mode findMode, ptrs []**node[T]) (*node[T], []**node[T]) {
	prev := &h.root

	// TODO: optimize start position to the height of the list, truncate ptrs accordingly
	i := len(h.root.nexts) - 1

	recordPtr := func() {
		if ptrs != nil {
			ptrs[i] = &prev.nexts[i]
		}
	}

	var cmp int = 1

	for i >= 0 {
		curr := prev.nexts[i]

		if curr == &h.root {
			recordPtr()
			i-- // Go down
			continue
		}

		cmp = h.compare(item, curr.item)

		switch {
		case mode == findFirst && cmp <= 0:
			recordPtr()
			i-- // Go back & down
			continue

		case mode == findLast && cmp < 0:
			recordPtr()
			i-- // Go back & down
			continue

		case mode == findAny && cmp == 0 && ptrs == nil:
			return curr, nil

		case mode == findAny && cmp <= 0:
			recordPtr()
			i-- // Go back & down
			continue
		}

		prev = curr // Go right
	}

	if i != -1 {
		panic(fmt.Sprintf("i [%d] != -1", i))
	}

	if cmp == 0 {
		return prev.nexts[0], ptrs
	}

	return nil, ptrs
}

func (h *head[T]) iterator() *nodeIter[T] {
	return &nodeIter[T]{
		curr:     h.root.nexts[0],
		sentinel: &h.root,
	}
}

func (h *head[T]) iteratorAt(start *node[T]) *nodeIter[T] {
	return &nodeIter[T]{
		curr:     start,
		sentinel: &h.root,
	}
}

type nodeIter[T any] struct {
	curr     *node[T]
	sentinel *node[T]
}

func (i *nodeIter[T]) Next() (T, bool) {
	var zero T
	if i.curr == i.sentinel {
		return zero, false
	}
	ret := i.curr.item
	i.curr = i.curr.nexts[0]

	return ret, true
}

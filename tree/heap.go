package tree

import "github.com/iv-menshenin/fusion/stack"

type Heap[IDX Ordered, D any] struct {
	heap Storage[Node[IDX, D]]
}

type Storage[T any] interface {
	Len() int
	Push(T) *T
	Get(int) *T
	Pop() T
}

func NewHeap[IDX Ordered, D any](init Storage[Node[IDX, D]]) *Heap[IDX, D] {
	if init == nil {
		init = &stack.Stack[Node[IDX, D]]{}
	}
	var tree = Heap[IDX, D]{heap: init}
	if tree.Len() > 0 {
		tree.makeBalance()
	}
	return &tree
}

func (t *Heap[IDX, D]) makeBalance() {
	var n = 0
	for idxOfRightChild(n) < t.Len() {
		n = idxOfRightChild(n)
	}
	for ; n >= 0; n-- {
		t.bDown(n)
	}
}

func (t *Heap[IDX, D]) Len() int {
	return t.heap.Len()
}

func (t *Heap[IDX, D]) Put(nodes ...Node[IDX, D]) {
	for _, node := range nodes {
		heapLen := t.Len()
		t.heap.Push(node)
		t.bUp(heapLen)
	}
}

func (t *Heap[IDX, D]) PopMax() (val Node[IDX, D], ok bool) {
	heapLen := t.Len()
	if heapLen == 0 {
		return
	}

	first := t.heap.Get(0)
	val = *first
	*first = t.heap.Pop()
	t.bDown(0)
	return val, true
}

func (t *Heap[IDX, D]) Search(idx IDX) (val *Node[IDX, D], ok bool) {
	if t.heap.Len() == 0 {
		return nil, false
	}
	return t.searchFrom(idx, 0)
}

func (t *Heap[IDX, D]) searchFrom(idx IDX, from int) (*Node[IDX, D], bool) {
	if t.heap.Get(from).idx == idx {
		return t.heap.Get(from), true
	}
	if t.heap.Get(from).idx < idx {
		return nil, false
	}
	if left := idxOfLeftChild(from); left < t.Len() {
		if found, ok := t.searchFrom(idx, left); ok {
			return found, ok
		}
	}
	if right := idxOfRightChild(from); right < t.Len() {
		if found, ok := t.searchFrom(idx, right); ok {
			return found, ok
		}
	}
	return nil, false
}

func (t *Heap[IDX, D]) bUp(currentIdx int) {
	for currentIdx > 0 {
		parentIdx := idxOfParent(currentIdx)
		parRef, curRef := t.heap.Get(parentIdx), t.heap.Get(currentIdx)
		if parRef.idx < curRef.idx {
			*parRef, *curRef = *curRef, *parRef // swap
			currentIdx = parentIdx
			continue
		}
		return
	}
}

func (t *Heap[IDX, D]) bDown(currentIdx int) {
	for heapLen := t.Len(); currentIdx < heapLen; {
		var (
			leftChildIdx       = idxOfLeftChild(currentIdx)
			rightChildIdx      = idxOfRightChild(currentIdx)
			currElement        = t.heap.Get(currentIdx)
			isLeftChildExists  = leftChildIdx < heapLen
			isRightChildExists = rightChildIdx < heapLen
			leftChild          *Node[IDX, D]
			rightChild         *Node[IDX, D]
		)
		if isLeftChildExists {
			leftChild = t.heap.Get(leftChildIdx)
		}
		if isRightChildExists {
			rightChild = t.heap.Get(rightChildIdx)
		}
		var (
			isLeftGreaterThanCurrent  = isLeftChildExists && currElement.idx < leftChild.idx
			isRightGreaterThanCurrent = isRightChildExists && currElement.idx < rightChild.idx
			isRightGreaterThanLeft    = isRightChildExists && leftChild.idx < rightChild.idx
		)
		switch {

		case isLeftGreaterThanCurrent && !isRightGreaterThanLeft:
			*currElement, *leftChild = *leftChild, *currElement // swap
			currentIdx = leftChildIdx
			continue

		case isRightGreaterThanCurrent:
			*currElement, *rightChild = *rightChild, *currElement // swap
			currentIdx = rightChildIdx
			continue
		}
		return
	}
}

func idxOfLeftChild(idx int) int {
	return idx*2 + 1
}

func idxOfRightChild(idx int) int {
	return idx*2 + 2
}

func idxOfParent(idx int) int {
	return (idx - 1) / 2
}

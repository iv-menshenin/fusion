package sparseset

import (
	"slices"

	"github.com/iv-menshenin/fusion/collection"
)

type Key interface {
	int | int64 | uint64
}

type SparseSet[K Key, T any] struct {
	sparse []int
	dense  *collection.Collection[backRef[K, T]]
	size   int
}

type backRef[K Key, T any] struct {
	ref  K
	data T
}

func New[K Key, T any](p int) *SparseSet[K, T] {
	return &SparseSet[K, T]{
		sparse: slices.Repeat([]int{NULL}, p),
		dense:  collection.New[backRef[K, T]](p),
		size:   p,
	}
}

func (s *SparseSet[K, T]) Len() int {
	return len(s.sparse)
}

const NULL = -1

func (s *SparseSet[K, T]) Set(key K, val T) (ref *T) {
	id := int(key)
	if s.size <= id {
		s.sparse = append(s.sparse, slices.Repeat([]int{NULL}, s.size)...)
		s.size = len(s.sparse)
	}

	if densePos := s.sparse[id]; densePos == NULL {
		s.sparse[id] = s.dense.Len()
		br := s.dense.Push(backRef[K, T]{ref: key, data: val})
		ref = &br.data
	} else {
		br := s.dense.Get(densePos)
		br.data = val
	}
	return ref
}

func (s *SparseSet[K, T]) Get(key K) *T {
	id := int(key)
	if s.size <= id {
		return nil
	}
	if s.sparse[id] == NULL {
		return nil
	}
	br := s.dense.Get(s.sparse[id])
	return &br.data
}

func (s *SparseSet[K, T]) Pop() T {
	br := s.dense.Pop()
	s.sparse[int(br.ref)] = NULL
	return br.data
}

func (s *SparseSet[K, T]) Delete(key K) {
	id := int(key)
	if s.size <= id {
		panic("out of bounds")
	}

	deleted := s.sparse[id]
	dd := s.dense.Get(deleted)
	ld := s.dense.Pop() // removed last from dense
	*dd = ld
	s.sparse[id] = NULL
	s.sparse[int(ld.ref)] = deleted // ld.ref referenced to entity was poped from dense
	// sparse:     [1|2|3|4|5|6]
	// dense:      [a|b|c|d|e|f]
	// remove      -----^
	// pop it      -----------^  // ld := s.dense.Pop()
	// place here  -----^        // *dd = ld
	// swap sparse [1|2|-|4|5|3] // s.sparse[id] = NULL
	// dense:      [a|b|f|d|e]   // s.sparse[int(ld.ref)] = deleted
}

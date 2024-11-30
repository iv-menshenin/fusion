package sparseset

import (
	"github.com/iv-menshenin/fusion/collection"
)

type Key interface {
	int | int64 | uint64
}

// SparseSet is designed to save memory and improve performance when dealing with large datasets
// that are mostly empty. Instead of allocating space for every possible element, it only stores
// the elements that are present.
type SparseSet[K Key, T any] struct {
	sparse []int
	dense  *collection.Collection[backRef[K, T]]
	size   int
}

type backRef[K Key, T any] struct {
	ref  K
	data T
}

func New[K Key, T any](p, bsz int) *SparseSet[K, T] {
	s := SparseSet[K, T]{
		sparse: make([]int, p),
		dense:  collection.New[backRef[K, T]](bsz),
	}
	for n := 0; n < len(s.sparse); n++ {
		s.sparse[n] = NULL
	}
	return &s
}

func (s *SparseSet[K, T]) Len() int {
	return s.size
}

const NULL = -1

// Set stores the object under a specific identifier, returning a reference to it.
func (s *SparseSet[K, T]) Set(key K, val T) (ref *T) {
	id := int(key)
	if currSz := len(s.sparse); currSz <= id {
		newSize := len(s.sparse) * 2
		if newSize <= id {
			newSize = id + 1
		}
		old := s.sparse
		s.sparse = make([]int, 0, newSize)
		s.sparse = append(s.sparse, old...)
		s.sparse = s.sparse[:newSize]
		for n := currSz; n < newSize; n++ {
			s.sparse[n] = NULL
		}
	}

	if densePos := s.sparse[id]; densePos == NULL {
		s.sparse[id] = s.dense.Len()
		br := s.dense.Push(backRef[K, T]{ref: key, data: val})
		ref = &br.data
		s.size++
	} else {
		br := s.dense.Get(densePos)
		br.data = val
	}
	return ref
}

// Get returns a reference to the object associated with the identifier `key`.
//
// The reference data can be modified, but avoid saving the reference, it may become invalid after calling methods that modify it, such as Delete.
func (s *SparseSet[K, T]) Get(key K) *T {
	id := int(key)
	if len(s.sparse) <= id {
		return nil
	}
	if s.sparse[id] == NULL {
		return nil
	}
	br := s.dense.Get(s.sparse[id])
	return &br.data
}

// Delete deletes an object by its Key from the SparseSet.
//
// Note that to improve performance, there is a side effect: the deleted object is replaced by the last object,
// not the next in line. This avoids large data movement when deleting values.
//
// Please note that after calling this method, the links that were obtained earlier by the Get method may be invalid.
func (s *SparseSet[K, T]) Delete(key K) {
	id := int(key)
	if len(s.sparse) <= id {
		panic("out of bounds")
	}
	s.size--
	deleted := s.sparse[id]
	s.sparse[id] = NULL
	dd := s.dense.Get(deleted)
	ld := s.dense.Pop() // removed last from dense
	if ld.ref == key {
		return
	}
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

// Each iterates through all the elements in the Collection and calls the provided callback function for each of
// the elements. If the callback function returns false, the iteration will be stopped.
func (s *SparseSet[K, T]) Each(callback func(key K, val *T) bool) {
	for k, v := range s.sparse {
		if v == NULL {
			continue
		}
		gh := callback(K(k), &s.dense.Get(v).data)
		if !gh {
			return
		}
	}
}

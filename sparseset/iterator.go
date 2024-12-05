package sparseset

import "context"

// Iterator returns a chan-iterator for iterating over all elements without using a callback function.
// This allows you to iterate through all elements using a standard `for range` loop.
// If you need to prematurely terminate the iteration, call the cancel function of the context.
func (s *SparseSet[K, T]) Iterator(ctx context.Context, buf int) <-chan Pair[K, T] {
	ch := make(chan Pair[K, T], buf)
	go s.sendElements(ctx, ch)
	return ch
}

type Pair[K Key, T any] struct {
	Key K
	Val *T
}

func (s *SparseSet[K, T]) sendElements(ctx context.Context, ch chan<- Pair[K, T]) {
	defer close(ch)
	for _, v := range s.sparse {
		if v == NULL {
			continue
		}
		val := s.dense.Get(v)
		select {
		case <-ctx.Done():
			return
		case ch <- Pair[K, T]{Key: val.ref, Val: &val.data}:
			// go ahead
		}
	}
}

// Fetcher allows for a sequential traversal of all elements in the SparseSet.
//
// Call the Next method to move the cursor and access the data. The returned value serves as an indicator that the value exists and can be accessed.
//
// The basic usage is as follows:
//
//	f := s.Fetcher()
//	for f.Next() {
//	  _ = f.Fetch()
//	}
func (s *SparseSet[K, T]) Fetcher() *Fetcher[K, T] {
	return &Fetcher[K, T]{s: s}
}

// Fetcher allows for a sequential traversal of all elements in the SparseSet.
type Fetcher[K Key, T any] struct {
	s *SparseSet[K, T]
	i int
	x int
}

// Next advances the cursor forward, returning true if the end has not yet been reached; otherwise, it returns false.
func (f *Fetcher[K, T]) Next() bool {
	for {
		f.i++
		if f.s.Len() < f.i {
			return false
		}
		idx := f.s.sparse[f.i-1]
		if idx != NULL {
			f.x = idx
			return true
		}
	}
}

// Fetch allows access to the current element and it's key of the SparseSet.
//
// Be careful, this function should not be called before the Next function has been called for the first time;
// otherwise, you will attempt to fetch an element at position -1, which will result in an error.
func (f *Fetcher[K, T]) Fetch() (K, *T) {
	el := f.s.dense.Get(f.x)
	return el.ref, &el.data
}

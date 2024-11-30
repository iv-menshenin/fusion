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

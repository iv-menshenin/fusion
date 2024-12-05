package collection

import "context"

// Iterator returns a chan-iterator for iterating over all elements without using a callback function.
// This allows you to iterate through all elements using a standard `for range` loop.
// If you need to prematurely terminate the iteration, call the cancel function of the context.
func (c *Collection[T]) Iterator(ctx context.Context, buf int) <-chan *T {
	ch := make(chan *T, buf)
	go c.sendElements(ctx, ch)
	return ch
}

func (c *Collection[T]) sendElements(ctx context.Context, ch chan<- *T) {
	defer close(ch)
	c.Each(func(val *T) bool {
		select {
		case <-ctx.Done():
			return false
		case ch <- val:
			return true
		}
	})
}

// Fetcher allows for a sequential traversal of all elements in the Collection.
//
// Call the Next method to move the cursor and access the data. The returned value serves as an indicator that the value exists and can be accessed.
//
// The basic usage is as follows:
//
//	f := c.Fetcher()
//	for f.Next() {
//	  _ = f.Fetch()
//	}
func (c *Collection[T]) Fetcher() *Fetcher[T] {
	return &Fetcher[T]{c: c}
}

// Fetcher allows for a sequential traversal of all elements in the Collection.
type Fetcher[T any] struct {
	c *Collection[T]
	i int
}

// Next advances the cursor forward, returning true if the end has not yet been reached; otherwise, it returns false.
func (f *Fetcher[T]) Next() bool {
	f.i++
	return f.i <= f.c.Len()
}

// Fetch allows access to the current element of the Collection.
//
// Be careful, this function should not be called before the Next function has been called for the first time;
// otherwise, you will attempt to fetch an element at position -1, which will result in an error.
func (f *Fetcher[T]) Fetch() *T {
	return f.c.Get(f.i - 1)
}

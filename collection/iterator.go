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

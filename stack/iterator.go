package stack

import "context"

// Iterator returns a chan-iterator for iterating over all elements without using a callback function.
// This allows you to iterate through all elements using a standard `for range` loop.
// If you need to prematurely terminate the iteration, call the cancel function of the context.
func (c *Stack[T]) Iterator(ctx context.Context, backward bool, buf int) <-chan *T {
	ch := make(chan *T, buf)
	if backward {
		go c.sendBackward(ctx, ch)
	} else {
		go c.sendForward(ctx, ch)
	}
	return ch
}

func (c *Stack[T]) sendBackward(ctx context.Context, ch chan<- *T) {
	defer close(ch)
	var (
		cur = c.last
		idx = cur.count
	)
	for {
		if idx--; idx < 0 {
			cur = cur.prev
			if cur == nil {
				break
			}
			idx = cur.count - 1
		}
		select {
		case <-ctx.Done():
			return
		case ch <- &cur.cont[idx]:
			// next
		}
	}
}

func (c *Stack[T]) sendForward(ctx context.Context, ch chan<- *T) {
	defer close(ch)
	var q = make([]*bucket[T], 0)
	for cur := c.last; cur != nil; cur = cur.prev {
		q = append(q, cur)
	}
	if len(q) == 0 {
		return
	}
	var idx int
	for {
		select {
		case <-ctx.Done():
			return
		case ch <- &q[len(q)-1].cont[idx]:
			idx++
			if idx >= q[len(q)-1].count {
				if q = q[:len(q)-1]; len(q) == 0 {
					return
				}
				idx = 0
			}
		}
	}
}

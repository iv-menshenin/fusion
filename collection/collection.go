package collection

type (
	Collection[T any] struct {
		len     int
		buckets []*bucket[T]
	}
	bucket[T any] struct {
		data []T
	}
)

func Init[T any](val []T) *Collection[T] {
	var buckets []*bucket[T]
	for len(val) > bucketSz {
		buckets = []*bucket[T]{
			{data: append(make([]T, 0, bucketSz), val[:bucketSz]...)},
		}
		val = val[bucketSz:]
	}
	buckets = append(buckets, &bucket[T]{data: append(make([]T, 0, bucketSz), val...)})
	return &Collection[T]{
		len:     len(val),
		buckets: buckets,
	}
}

func New[T any]() *Collection[T] {
	return &Collection[T]{
		buckets: []*bucket[T]{
			newBucket[T](),
		},
	}
}

func (c *Collection[T]) Len() int {
	return c.len
}

const bucketSz = 1000

func newBucket[T any]() *bucket[T] {
	return &bucket[T]{
		data: make([]T, 0, bucketSz),
	}
}

func (c *Collection[T]) Push(val T) *T {
	if len(c.buckets) == 0 {
		c.buckets = append(make([]*bucket[T], 0, 1000), newBucket[T]())
	}
	b := c.buckets[len(c.buckets)-1]
	l := len(b.data)
	if l == cap(b.data) {
		b = newBucket[T]()
		c.buckets = append(c.buckets, b)
		l = 0
	}
	b.data = append(b.data, val)
	c.len++
	return &b.data[l]
}

func (c *Collection[T]) Get(idx int) *T {
	if idx >= c.len {
		return nil
	}
	return &c.buckets[idx/bucketSz].data[idx%bucketSz]
}

func (c *Collection[T]) Delete(idx int) {
	if idx >= c.len {
		panic("out of bounds")
	}
	c.len--
	if c.len == idx {
		// removed last element
		return
	}
	ref := &c.buckets[idx/bucketSz].data[idx%bucketSz]
	last := c.Pop()
	*ref = last
}

func (c *Collection[T]) Pop() T {
	if c.len < 1 {
		panic("called Pop on empty Collection")
	}
	c.len--
	b := c.buckets[len(c.buckets)-1]
	l := len(b.data)
	val := b.data[l-1]
	b.data = b.data[:l-1]
	if len(b.data) == 0 {
		c.buckets[len(c.buckets)-1] = nil
		c.buckets = c.buckets[:len(c.buckets)-1]
	}
	return val
}

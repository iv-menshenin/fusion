package collection

import "github.com/iv-menshenin/fusion/errors"

type (
	Collection[T any] struct {
		len     int
		buckets []*bucket[T]
	}
	bucket[T any] struct {
		data []T
	}
)

const defaultBucketSz = 1000

func Init[T any](val []T) *Collection[T] {
	var (
		l       = len(val)
		buckets = make([]*bucket[T], 0, 1+(l/defaultBucketSz))
	)
	for len(val) > 0 {
		if len(val) < defaultBucketSz {
			buckets = append(buckets, &bucket[T]{data: append(make([]T, 0, defaultBucketSz), val...)})
			break
		}
		buckets = append(buckets, &bucket[T]{
			data: val[:defaultBucketSz], // no copy data
		})
		val = val[defaultBucketSz:]
	}
	return &Collection[T]{
		len:     l,
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

func newBucket[T any]() *bucket[T] {
	return &bucket[T]{
		data: make([]T, 0, defaultBucketSz),
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

func (c *Collection[T]) Get(i int) *T {
	if i >= c.len {
		return nil
	}
	return &c.buckets[i/defaultBucketSz].data[i%defaultBucketSz]
}

func (c *Collection[T]) Delete(i int) {
	if i >= c.len {
		panic(errors.OutOfBounds(c.len, i))
	}
	if (c.len - 1) == i {
		// removed last element
		c.len--
		b := c.buckets[len(c.buckets)-1]
		b.data = b.data[:len(b.data)-1]
		return
	}
	ref := &c.buckets[i/defaultBucketSz].data[i%defaultBucketSz]
	last := c.Pop()
	*ref = last
}

func (c *Collection[T]) Pop() T {
	if c.len < 1 {
		panic(errors.OutOfBounds(c.len, 0))
	}
	c.len--
	b := c.buckets[len(c.buckets)-1]
	l := len(b.data)
	val := b.data[l-1]
	// clean cell
	var empty T
	b.data[l-1] = empty
	// reduce
	if len(b.data) == 1 {
		c.buckets[len(c.buckets)-1] = nil
		c.buckets = c.buckets[:len(c.buckets)-1]
	} else {
		b.data = b.data[:l-1]
	}
	return val
}

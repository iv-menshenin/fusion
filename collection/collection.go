package collection

import "github.com/iv-menshenin/fusion/errors"

type (
	Collection[T any] struct {
		len     int
		bsz     int
		buckets []*bucket[T]
	}
	bucket[T any] struct {
		data []T
	}
)

const defaultBucketSz = 1000

func Init[T any](val []T, bucketSz int) *Collection[T] {
	if bucketSz == 0 {
		bucketSz = defaultBucketSz
	}
	var (
		l       = len(val)
		buckets = make([]*bucket[T], 0, 1+(l/bucketSz))
	)
	for len(val) > 0 {
		if len(val) < bucketSz {
			buckets = append(buckets, &bucket[T]{data: append(make([]T, 0, bucketSz), val...)})
			break
		}
		buckets = append(buckets, &bucket[T]{
			data: val[:bucketSz], // no copy data
		})
		val = val[bucketSz:]
	}
	return &Collection[T]{
		len:     l,
		bsz:     bucketSz,
		buckets: buckets,
	}
}

func New[T any](bucketSz int) *Collection[T] {
	if bucketSz == 0 {
		bucketSz = defaultBucketSz
	}
	return &Collection[T]{
		bsz: bucketSz,
	}
}

func (c *Collection[T]) Len() int {
	return c.len
}

func (c *Collection[T]) Push(val T) *T {
	if len(c.buckets) == 0 {
		c.buckets = append(make([]*bucket[T], 0, 1000), c.newBucket())
	}
	b := c.buckets[len(c.buckets)-1]
	l := len(b.data)
	if l == cap(b.data) {
		b = c.newBucket()
		c.buckets = append(c.buckets, b)
		l = 0
	}
	b.data = append(b.data, val)
	c.len++
	return &b.data[l]
}

func (c *Collection[T]) newBucket() *bucket[T] {
	if c.bsz == 0 {
		c.bsz = defaultBucketSz
	}
	return &bucket[T]{
		data: make([]T, 0, c.bsz),
	}
}

func (c *Collection[T]) Get(i int) *T {
	if i >= c.len {
		return nil
	}
	return &c.buckets[i/c.bsz].data[i%c.bsz]
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
	ref := &c.buckets[i/c.bsz].data[i%c.bsz]
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

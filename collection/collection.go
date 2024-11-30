package collection

import "github.com/iv-menshenin/fusion/errors"

type (
	// Collection is a special data structure designed to avoid excessive memory allocation when adding
	// a large number of values to a slice in situations where you do not know in advance the amount of data
	// that needs to be stored.
	//
	// The collection guarantees correct sorting as long as you do not remove a value from the middle.
	// When a value is removed, all subsequent values are not shifted; instead, the current value is replaced
	// with the last one, and the length is reduced by one.
	//
	// The Push and Get methods return a reference to an object that can be modified. Note that the reference
	// is guaranteed to be valid only until the first call to methods that delete values, such as Delete or even Pop.
	// Avoid storing the reference for a long time.
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

// Init allows you to create a Collection with pre-fulfilled data.
// The slice passed as an argument will be fully reused if its length is divisible by the bucket size.
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
			data := make([]T, bucketSz)
			copy(data[:len(val)], val)
			buckets = append(buckets, &bucket[T]{data: data})
			break
		}
		// no copy data
		buckets = append(buckets, &bucket[T]{
			data: val[:bucketSz],
		})
		val = val[bucketSz:]
	}
	return &Collection[T]{
		len:     l,
		bsz:     bucketSz,
		buckets: buckets,
	}
}

// New creates a new Collection with the specified bucket size. If the size is zero, the default value will be used.
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

// Push adds a new value to the end of the Collection and returns a reference to it.
func (c *Collection[T]) Push(val T) *T {
	if c.bsz == 0 {
		c.bsz = defaultBucketSz
	}
	id := c.len
	bId := id / c.bsz
	xId := id % c.bsz
	c.len++
	if len(c.buckets) <= bId {
		c.extendBuckets()
	}
	c.buckets[bId].data[xId] = val
	return &c.buckets[bId].data[xId]
}

func (c *Collection[T]) extendBuckets() {
	c.buckets = append(c.buckets, &bucket[T]{
		data: make([]T, c.bsz),
	})
}

// Get allows you to get a reference to an object located in a Collection.
//
// Avoid storing the link outside of the Collection for long periods of time.
func (c *Collection[T]) Get(id int) *T {
	if id >= c.len {
		return nil
	}
	bId := id / c.bsz
	xId := id % c.bsz
	return &c.buckets[bId].data[xId]
}

// Delete deletes an object by its index from the collection.
//
// Note that to improve performance, there is a side effect: the deleted object is replaced by the last object,
// not the next in line. This avoids large data movement when deleting values from the beginning.
//
// So if you need to delete several values from n to m, it is safe to do it only in the index decreasing direction,
// i.e. from m to n.
func (c *Collection[T]) Delete(id int) {
	if id >= c.len {
		panic(errors.OutOfBounds(c.len, id))
	}
	lId := c.len - 1
	xbId := lId / c.bsz
	xxId := lId % c.bsz
	bId := id / c.bsz
	xId := id % c.bsz
	if bId != xbId || xId != xxId {
		// swap
		c.buckets[bId].data[xId] = c.buckets[xbId].data[xxId]
	}
	c.len--

	// clear cell
	var empty T
	c.buckets[xbId].data[xxId] = empty
}

// Pop selects the last item in the collection and returns a copy of it. The original item is deleted.
func (c *Collection[T]) Pop() T {
	if c.len < 1 {
		panic(errors.OutOfBounds(c.len, 0))
	}
	id := c.len - 1
	bId := id / c.bsz
	xId := id % c.bsz
	c.len--
	b := c.buckets[bId]
	val := b.data[xId]
	// clean cell
	var empty T
	b.data[xId] = empty
	return val
}

// Prune clears unoccupied space. It can be used after a large number of calls to Delete or Pop method.
func (c *Collection[T]) Prune() {
	bId := c.len / c.bsz
	for n := bId + 1; n < len(c.buckets); n++ {
		c.buckets[n] = nil
	}
	c.buckets = c.buckets[:bId]
}

// Each iterates through all the elements in the Collection and calls the provided callback function for each of
// the elements. If the callback function returns false, the iteration will be stopped.
func (c *Collection[T]) Each(callback func(*T) bool) {
	bId := c.len / c.bsz
	xId := c.len % c.bsz
	for cbId := range c.buckets {
		for cxId := range c.buckets[cbId].data {
			if cbId == bId && xId <= cxId {
				return
			}
			if callback(&c.buckets[cbId].data[cxId]) {
				continue
			}
			return
		}
	}
}

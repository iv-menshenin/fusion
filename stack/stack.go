package stack

type (
	Stack[T any] struct {
		count int
		last  *bucket[T]
		cache []bucket[T]

		// search cache
		sclx  int
		sccur *bucket[T]
	}
	bucket[T any] struct {
		count int
		cont  []T
		prev  *bucket[T]
	}
)

func Init[T any](val []T) *Stack[T] {
	s := Stack[T]{
		count: len(val),
		last: &bucket[T]{
			count: len(val),
			cont:  val,
		},
	}
	return &s
}

// Get returns the top element of the stack without removing it. It's copy of Peek method for sorting support.
func (c *Stack[T]) Get(i int) *T {
	return c.Peek(i)
}

// Peek returns the top element of the stack without removing it.
func (c *Stack[T]) Peek(i int) *T {
	if i > c.count-1 {
		panic(ErrOutOfBounds{i: i, l: c.count})
	}
	var (
		lx  = c.count
		cur = c.last
	)
	if c.sccur != nil && c.sclx-c.sccur.count >= i {
		lx = c.sclx
		cur = c.sccur
	}
	for lx-cur.count > i {
		lx -= cur.count
		cur = cur.prev
	}
	c.sclx = lx
	c.sccur = cur
	return &cur.cont[cur.count-(lx-i)]
}

// Pop removes and returns the top element of the stack.
func (c *Stack[T]) Pop() T {
	if c.count == 0 {
		panic(ErrOutOfBounds{l: c.count})
	}
	if c.sccur != nil {
		c.sccur = nil
	}
	l := len(c.last.cont)
	for l == 0 {
		c.removeLast()
		l = len(c.last.cont)
	}
	v := c.last.cont[l-1]
	c.last.cont = c.last.cont[:l-1]
	c.count--
	c.last.count--
	return v
}

// Push adds an element to the top of the stack.
func (c *Stack[T]) Push(elem T) *T {
	if !c.capable() {
		c.extend()
	}
	if c.sccur == c.last {
		c.sccur = nil
	}
	l := len(c.last.cont)
	c.last.cont = c.last.cont[:l+1]
	c.last.cont[l] = elem
	c.last.count++
	c.count++
	return &c.last.cont[l]
}

func (c *Stack[T]) Len() int {
	return c.count
}

func (c *Stack[T]) capable() bool {
	if c.last == nil {
		return false
	}
	return len(c.last.cont) < cap(c.last.cont)
}

func (c *Stack[T]) extend() {
	if c.last == nil {
		c.last = c.newBucket()
		return
	}
	n := c.newBucket()
	n.prev = c.last
	c.last = n
}

const (
	bucketsCache  = 32
	firstBucketSz = 1_000
	maxBucketSz   = 1_000_000
)

func (c *Stack[T]) newBucket() *bucket[T] {
	if len(c.cache) == 0 {
		c.cache = make([]bucket[T], bucketsCache)
	}
	var (
		l = len(c.cache) - 1
		b = &c.cache[l]
	)
	if cap(b.cont) == 0 {
		b.cont = make([]T, 0, c.sz())
	}
	c.cache = c.cache[:l]
	return b
}

func (c *Stack[T]) sz() int {
	if c.last == nil {
		return firstBucketSz
	}
	var sz = cap(c.last.cont) * 2
	if sz > maxBucketSz {
		sz = maxBucketSz
	}
	return sz
}

func (c *Stack[T]) removeLast() {
	if c.last == nil {
		return
	}
	if c.last.count > 0 {
		panic("remove nonempty bucket")
	}
	removed := c.last
	c.last = c.last.prev
	// keep already allocated cont in cache
	for i := len(c.cache); i > 0; {
		i--
		if c.cache[i].cont == nil {
			c.cache[i] = *removed
			break
		}
	}
}

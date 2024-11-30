package stack

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStack(t *testing.T) {
	t.Parallel()
	t.Run("push_get", func(t *testing.T) {
		t.Parallel()

		var c Stack[int64]
		c.Push(101)
		c.Push(102)
		c.Push(103)
		require.Equal(t, int64(103), *c.Get(2))
		require.Equal(t, int64(102), *c.Get(1))
		require.Equal(t, int64(101), *c.Get(0))
	})
	t.Run("push_peek", func(t *testing.T) {
		t.Parallel()

		var c Stack[int64]
		c.Push(45)
		require.Equal(t, int64(45), *c.Peek())
		c.Push(46)
		require.Equal(t, int64(46), *c.Peek())
		c.Push(47)
		require.Equal(t, int64(47), *c.Peek())
	})
	t.Run("change_external", func(t *testing.T) {
		t.Parallel()

		var c Stack[int]
		c.Push(220)
		require.Equal(t, 1, c.Len())

		val := c.Get(0)
		require.Equal(t, 220, *val)
		*val = 400
		require.Equal(t, 400, *c.Get(0))

		val = c.Peek()
		*val = -1
		require.Equal(t, -1, *c.Get(0))
	})
	t.Run("len", func(t *testing.T) {
		t.Parallel()

		var c Stack[int64]
		require.Equal(t, 0, c.Len())
		c.Push(101)
		require.Equal(t, 1, c.Len())
		c.Push(102)
		require.Equal(t, 2, c.Len())
		c.Push(103)
		require.Equal(t, 3, c.Len())

		c.Pop()
		require.Equal(t, 2, c.Len())
		c.Pop()
		require.Equal(t, 1, c.Len())
		c.Pop()
		require.Equal(t, 0, c.Len())

		require.Equal(t, 0, c.Len())
		c.Push(88)
		require.Equal(t, 1, c.Len())
		c.Push(99)
		require.Equal(t, 2, c.Len())
	})
	t.Run("search-cache-bug", func(t *testing.T) {
		t.Parallel()
		c := Init([]int{67, 13, 54, 2, 1, 42})

		c.Push(43)

		require.Equal(t, *c.Get(4), 1)
		require.Equal(t, *c.Get(2), 54)
		require.Equal(t, *c.Get(5), 42)
		require.Equal(t, *c.Get(0), 67)
		require.Equal(t, *c.Get(6), 43)

		c.Push(2)
		require.Equal(t, *c.Get(7), 2)

		c.Push(67)
		require.Equal(t, *c.Get(2), 54)
	})
}

func TestStressStack(t *testing.T) {
	t.Parallel()
	type Elem struct {
		i int
		s string
	}
	var c Stack[Elem]

	const elemCount = 1_000_000

	// fill
	for n := 0; n < elemCount; n++ {
		require.Equal(t, n, c.Len())
		c.Push(Elem{i: n, s: strconv.Itoa(n)})
	}

	// forward
	for n := 0; n < elemCount; n++ {
		e := c.Get(n)
		require.Equal(t, n, e.i)
		require.Equal(t, strconv.Itoa(n), e.s)
	}

	// backward
	for n := elemCount; n > 0; n-- {
		e := c.Get(n - 1)
		require.Equal(t, n-1, e.i)
		require.Equal(t, strconv.Itoa(n-1), e.s)
	}

	// pop
	for n := elemCount; n > 0; n-- {
		e := c.Pop()
		require.Equal(t, n-1, c.Len())
		require.Equal(t, n-1, e.i)
		require.Equal(t, strconv.Itoa(n-1), e.s)
	}
}

func TestStackPushPop(t *testing.T) {
	t.Parallel()
	t.Run("push", func(t *testing.T) {
		t.Parallel()

		var c Stack[int]
		c.Push(100)
		require.Equal(t, 1, c.Len())
		require.Equal(t, 100, *c.Get(0))
	})
	t.Run("pop", func(t *testing.T) {
		t.Parallel()

		var c Stack[int]
		c.Push(200)
		require.Equal(t, 1, c.Len())
		require.Equal(t, 200, *c.Get(0))

		require.Equal(t, 200, c.Pop())
		require.Equal(t, 0, c.Len())
	})
	t.Run("push_pop", func(t *testing.T) {
		t.Parallel()
		var c Stack[string]

		c.Push("foo")
		require.Equal(t, 1, c.Len())
		c.Push("bar")
		require.Equal(t, 2, c.Len())
		require.Equal(t, "bar", c.Pop())
		require.Equal(t, 1, c.Len())
		require.Equal(t, "foo", c.Pop())
		require.Equal(t, 0, c.Len())

		c.Push("1")
		c.Push("3")
		require.Equal(t, "3", c.Pop())
		c.Push("2")
		c.Push("3")
		c.Push("4")
		require.Equal(t, 4, c.Len())
		require.Equal(t, "4", c.Pop())
		require.Equal(t, "3", c.Pop())
		require.Equal(t, "2", c.Pop())
		require.Equal(t, "1", c.Pop())
		require.Equal(t, 0, c.Len())
	})
}

func BenchmarkStackTravel(b *testing.B) {
	type Elem struct {
		s          string
		a, b, c, d int64
		n          int
	}
	var c Stack[Elem]
	const max = 1000000
	for n := 0; n < max; n++ {
		c.Push(Elem{n: n})
	}
	b.Run("Last", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_ = c.Get(max - 1)
		}
		b.ReportMetric(100, "items")
	})
	b.Run("First", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_ = c.Get(1)
		}
	})
	b.Run("Mid", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_ = c.Get(max / 2)
		}
	})
	b.Run("Forward", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_ = c.Get(n % max)
		}
	})
	b.Run("Backward", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_ = c.Get(max - (n % max) - 1)
		}
	})
	b.Run("Random", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			switch n % 5 {
			case 0:
				_ = c.Get(max / 2)
			case 1:
				_ = c.Get(n % max)
			case 2:
				_ = c.Get(max - (max / 4))
			case 3:
				_ = c.Get(max - (n % max) - 1)
			case 4:
				_ = c.Get(max / 4)
			}
		}
	})
}

func BenchmarkStackPushPop(b *testing.B) {
	type Elem struct {
		s          string
		a, b, c, d int64
		n          int
	}
	b.Run("PushAll", func(b *testing.B) {
		b.ReportAllocs()
		var c Stack[Elem]
		for n := 0; n < b.N; n++ {
			c.Push(Elem{n: n})
		}
	})
	b.Run("2Push1Pop", func(b *testing.B) {
		b.ReportAllocs()
		var c Stack[Elem]
		for n := 0; n < b.N; n++ {
			c.Push(Elem{n: n})
			if n%2 == 0 {
				_ = c.Pop()
			}
		}
	})
	b.Run("PopAll", func(b *testing.B) {
		b.ReportAllocs()
		var c Stack[Elem]
		for n := 0; n < b.N; n++ {
			c.Push(Elem{n: n})
		}
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			_ = c.Pop()
		}
	})
	b.Run("PushThenPop", func(b *testing.B) {
		b.ReportAllocs()
		var c Stack[Elem]
		for n := 0; n < b.N; n++ {
			c.Push(Elem{n: n})
			_ = c.Pop()
		}
	})
}

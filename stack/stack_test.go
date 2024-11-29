package stack

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFusionCollection(t *testing.T) {
	t.Parallel()
	t.Run("append_get", func(t *testing.T) {
		t.Parallel()

		var c Stack[int64]
		c.Push(101)
		c.Push(102)
		c.Push(103)
		require.Equal(t, int64(103), *c.Peek(2))
		require.Equal(t, int64(102), *c.Peek(1))
		require.Equal(t, int64(101), *c.Peek(0))
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
	})
	t.Run("force", func(t *testing.T) {
		t.Parallel()
		type Elem struct {
			i int
			s string
		}
		var c Stack[Elem]

		const elemCount = 100000

		for n := 0; n < elemCount; n++ {
			require.Equal(t, n, c.Len())
			c.Push(Elem{i: n, s: strconv.Itoa(n)})
		}

		// forward
		for n := 0; n < elemCount; n++ {
			e := c.Peek(n)
			require.Equal(t, n, e.i)
			require.Equal(t, strconv.Itoa(n), e.s)
		}

		// backward
		for n := elemCount; n > 0; n-- {
			e := c.Peek(n - 1)
			require.Equal(t, n-1, e.i)
			require.Equal(t, strconv.Itoa(n-1), e.s)
		}
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

func TestFusionCollectionPushPop(t *testing.T) {
	t.Parallel()
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
	t.Run("force", func(t *testing.T) {
		t.Parallel()

		var c Stack[int]
		const elemCount = 1000000

		for n := 0; n < elemCount; n++ {
			c.Push(n)
		}
		for n := elemCount - 1; n >= 0; n-- {
			require.Equal(t, n, c.Pop())
		}
	})
}

func BenchmarkFusionCollectionGet(b *testing.B) {
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
			_ = c.Peek(max - 1)
		}
		b.ReportMetric(100, "items")
	})
	b.Run("First", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_ = c.Peek(1)
		}
	})
	b.Run("Mid", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_ = c.Peek(max / 2)
		}
	})
	b.Run("Forward", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_ = c.Peek(n % max)
		}
	})
	b.Run("Backward", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_ = c.Peek(max - (n % max) - 1)
		}
	})
	b.Run("Random", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			switch n % 5 {
			case 0:
				_ = c.Peek(max / 2)
			case 1:
				_ = c.Peek(n % max)
			case 2:
				_ = c.Peek(max - (max / 4))
			case 3:
				_ = c.Peek(max - (n % max) - 1)
			case 4:
				_ = c.Peek(max / 4)
			}
		}
	})
}

func BenchmarkFusionCollectionPushPop(b *testing.B) {
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

func BenchmarkFusionCollectionInsert(b *testing.B) {
	b.Run("Push", func(b *testing.B) {
		b.ReportAllocs()
		type Elem struct {
			s          string
			a, b, c, d int64
			n          int
		}
		var c Stack[Elem]
		for n := 0; n < b.N; n++ {
			c.Push(Elem{n: n})
		}
	})
}

package collection

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFusionCollection(t *testing.T) {
	t.Parallel()
	t.Run("append_get", func(t *testing.T) {
		t.Parallel()

		var c Collection[int64]
		c.Push(101)
		c.Push(102)
		c.Push(103)
		require.Equal(t, int64(103), *c.Get(2))
		require.Equal(t, int64(102), *c.Get(1))
		require.Equal(t, int64(101), *c.Get(0))
	})
	t.Run("len", func(t *testing.T) {
		t.Parallel()

		var c Collection[int64]
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
		var c Collection[Elem]

		const elemCount = 100000

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
	})
	t.Run("random-search", func(t *testing.T) {
		t.Parallel()
		c := Init([]int{67, 13, 54, 2, 1, 42}, 500)

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
	t.Run("append_get_delete", func(t *testing.T) {
		t.Parallel()

		var c Collection[int64]
		c.Push(801)
		c.Push(802)
		c.Push(803)
		require.Equal(t, int64(803), *c.Get(2))
		require.Equal(t, int64(802), *c.Get(1))
		require.Equal(t, int64(801), *c.Get(0))
		c.Delete(1)
		require.Equal(t, int64(803), *c.Get(1))
		require.Equal(t, int64(801), *c.Get(0))
		c.Delete(0)
		c.Delete(0)
		require.Equal(t, 0, c.Len())

		c.Push(0)
		require.Equal(t, int64(0), *c.Get(0))
		require.Nil(t, c.Get(1))
		require.Nil(t, c.Get(2))
	})
}

func TestFusionCollectionPushPop(t *testing.T) {
	t.Parallel()
	t.Run("push_pop", func(t *testing.T) {
		t.Parallel()
		var c Collection[string]

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

		var c Collection[int]
		const elemCount = 1000000

		for n := 0; n < elemCount; n++ {
			c.Push(n)
		}
		for n := elemCount - 1; n >= 0; n-- {
			require.Equal(t, n, c.Pop())
		}
	})
}

func TestFusionCollectionInit(t *testing.T) {
	t.Parallel()
	t.Run("init_buckets", func(t *testing.T) {
		const count = (defaultBucketSz * 5) / 2
		var box = make([]uint64, count)
		for n := 0; n < count; n++ {
			box[n] = uint64(n)
		}
		c := Init(box, defaultBucketSz)
		require.Equal(t, count, c.Len())
		for n := 0; n < count; n++ {
			require.Equal(t, uint64(n), *c.Get(n))
		}
	})
	t.Run("half_bucket", func(t *testing.T) {
		const count = defaultBucketSz / 2
		var box = make([]uint64, count-1)
		for n := 0; n < count-1; n++ {
			box[n] = uint64(n)
		}
		c := Init(box, defaultBucketSz)
		require.Equal(t, count-1, c.Len())
		for n := 0; n < count-1; n++ {
			require.Equal(t, uint64(n), *c.Get(n))
		}

		c.Push(count - 1)
		require.Equal(t, count, c.Len())
		for n := 0; n < count; n++ {
			require.Equal(t, uint64(n), *c.Get(n))
		}
	})
}

func BenchmarkFusionCollectionGet(b *testing.B) {
	type Elem struct {
		s          string
		a, b, c, d int64
		n          int
	}
	var c Collection[Elem]
	const max = 1000000
	for n := 0; n < max; n++ {
		c.Push(Elem{n: n})
	}
	b.Run("Last", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_ = c.Get(max - 1)
		}
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

func BenchmarkFusionCollectionPushPop(b *testing.B) {
	type Elem struct {
		s          string
		a, b, c, d int64
		n          int
	}
	b.Run("PushAll", func(b *testing.B) {
		b.ReportAllocs()
		var c Collection[Elem]
		for n := 0; n < b.N; n++ {
			c.Push(Elem{n: n})
		}
	})
	b.Run("2Push1Pop", func(b *testing.B) {
		b.ReportAllocs()
		var c Collection[Elem]
		for n := 0; n < b.N; n++ {
			c.Push(Elem{n: n})
			if n%2 == 0 {
				_ = c.Pop()
			}
		}
	})
	b.Run("PopAll", func(b *testing.B) {
		b.ReportAllocs()
		var c Collection[Elem]
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
		var c Collection[Elem]
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
		var c Collection[Elem]
		for n := 0; n < b.N; n++ {
			c.Push(Elem{n: n})
		}
	})
}

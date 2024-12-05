package collection

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFusionCollection(t *testing.T) {
	t.Parallel()
	t.Run("one_append", func(t *testing.T) {
		t.Parallel()

		var c Collection[string]
		v := c.Push("test")
		require.Equal(t, 1, c.Len(), "length checking")
		require.Equal(t, "test", *v, "value checking")
	})
	t.Run("push_few", func(t *testing.T) {
		t.Parallel()

		var c Collection[string]
		c.Push("101")
		c.Push("102")
		c.Push("103")
		require.Equal(t, "103", *c.Get(2))
		require.Equal(t, "102", *c.Get(1))
		require.Equal(t, "101", *c.Get(0))
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
			require.NotNil(t, e)
			require.Equal(t, Elem{i: n, s: strconv.Itoa(n)}, *e, "value checking")
		}

		// backward
		for n := elemCount; n > 0; n-- {
			e := c.Get(n - 1)
			require.NotNil(t, e)
			require.Equal(t, Elem{i: n - 1, s: strconv.Itoa(n - 1)}, *e, "value checking")
		}
	})
	t.Run("random-get", func(t *testing.T) {
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

func TestCollectionPushPop(t *testing.T) {
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

func TestCollectionInit(t *testing.T) {
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

func TestFusionCollectionPrune(t *testing.T) {
	var c = New[int](10)
	for n := 0; n < 100; n++ {
		c.Push(n)
	}
	for n := 9; n >= 0; n-- {
		c.Delete(n)
	}
	require.Equal(t, 90, c.Len())

	require.Len(t, c.buckets, 10)

	c.Prune()

	require.Len(t, c.buckets, 9)

	// check values
	var x [100]bool
	for n := 0; n < c.Len(); n++ {
		x[*c.Get(n)] = true
	}
	for n := 0; n < len(x); n++ {
		require.Equal(t, n >= 10, x[n], "bad %d", n)
	}
}

func TestFusionCollectionEach(t *testing.T) {
	var c = New[int](33)
	for n := 0; n < 100; n++ {
		c.Push(n)
	}

	t.Run("each_value", func(t *testing.T) {
		var x [100]bool
		c.Each(func(i *int) bool {
			x[*i] = true
			return true
		})
		// check values
		for n := 0; n < len(x); n++ {
			require.True(t, x[n], "bad %d", n)
		}
	})
	t.Run("iter_cancellation", func(t *testing.T) {
		var cc = make(map[int]struct{})
		c.Each(func(i *int) bool {
			cc[*i] = struct{}{}
			return len(cc) < 3
		})
		require.Len(t, cc, 3)
	})
}

func BenchmarkCollectionGet_10mln(b *testing.B) {
	type Elem struct {
		s          string
		a, b, c, d int64
		n          int
	}
	var c Collection[Elem]
	const count = 10_000_000
	for n := 0; n < count; n++ {
		c.Push(Elem{n: n})
	}
	b.Run("Last", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_ = c.Get(count - 1)
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
			_ = c.Get(count / 2)
		}
	})
	b.Run("Forward", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_ = c.Get(n % count)
		}
	})
	b.Run("Backward", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_ = c.Get(count - (n % count) - 1)
		}
	})
	b.Run("Random", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			switch n % 5 {
			case 0:
				// static
				_ = c.Get(count / 2)
			case 1:
				// next
				_ = c.Get(n % count)
			case 2:
				// backward
				_ = c.Get(count - (n % count) - 1)
			case 3:
				// fast-backward
				_ = c.Get(count - ((n * 127) % count) - 1)
			case 4:
				// fast-forward
				_ = c.Get((n * 127) % count)
			}
		}
	})
}

func BenchmarkCollectionPushPop(b *testing.B) {
	type Elem struct {
		s          string
		a, b, c, d int64
		n          int
	}
	b.Run("Push_All", func(b *testing.B) {
		b.ReportAllocs()
		var c Collection[Elem]
		for n := 0; n < b.N; n++ {
			c.Push(Elem{n: n, s: "Push_All"})
		}
	})
	b.Run("Push_Len_All", func(b *testing.B) {
		b.ReportAllocs()
		var c Collection[Elem]
		for n := 0; n < b.N; n++ {
			c.Push(Elem{n: n, s: "Push_Len_All"})
			_ = c.Len()
		}
	})
	b.Run("Push_Then_Pop", func(b *testing.B) {
		b.ReportAllocs()
		var c Collection[Elem]
		for n := 0; n < b.N; n++ {
			c.Push(Elem{n: n, s: "Push_Then_Pop"})
			_ = c.Pop()
		}
	})
	b.Run("Push_Push_Then_Pop", func(b *testing.B) {
		b.ReportAllocs()
		var c Collection[Elem]
		for n := 0; n < b.N; n++ {
			c.Push(Elem{a: 0, n: n, s: "Push_Push_Then_Pop"})
			c.Push(Elem{a: 1, n: n, s: "Push_Push_Then_Pop"})
			_ = c.Pop()
		}
	})
	b.Run("Push_Push_Then_Pop_Len", func(b *testing.B) {
		b.ReportAllocs()
		var c Collection[Elem]
		for n := 0; n < b.N; n++ {
			c.Push(Elem{a: 0, n: n, s: "Push_Push_Then_Pop"})
			c.Push(Elem{a: 1, n: n, s: "Push_Push_Then_Pop"})
			_ = c.Pop()
			_ = c.Len()
		}
	})
	b.Run("Pop_All", func(b *testing.B) {
		b.ReportAllocs()
		var c Collection[Elem]
		for n := 0; n < b.N; n++ {
			c.Push(Elem{n: n, s: "Pop_All"})
		}
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			_ = c.Pop()
		}
	})
}

func BenchmarkCollectionInsert(b *testing.B) {
	type Elem struct {
		s          string
		a, b, c, d int64
		n          int
	}
	b.Run("Push_1K", func(b *testing.B) {
		b.ReportAllocs()
		var c Collection[Elem]
		for n := 0; n < b.N; n++ {
			for i := 0; i < 1_000; i++ {
				c.Push(Elem{n: n, a: int64(i), b: -1, s: "Push 1K"})
			}
		}
	})
	b.Run("Push_65K", func(b *testing.B) {
		b.ReportAllocs()
		var c Collection[Elem]
		for n := 0; n < b.N; n++ {
			for i := 0; i < 65_535; i++ {
				c.Push(Elem{n: n, a: int64(i), b: -1, s: "Push 65K"})
			}
		}
	})
	b.Run("Push_1M", func(b *testing.B) {
		b.ReportAllocs()
		var c Collection[Elem]
		for n := 0; n < b.N; n++ {
			for i := 0; i < 1_000_000; i++ {
				c.Push(Elem{n: n, a: int64(i), b: -1, s: "Push 1m"})
			}
		}
	})
	b.Run("Push_10M", func(b *testing.B) {
		b.ReportAllocs()
		var c Collection[Elem]
		for n := 0; n < b.N; n++ {
			for i := 0; i < 10_000_000; i++ {
				c.Push(Elem{n: n, a: int64(i), b: -1, s: "Push 10m"})
			}
		}
	})
}

func BenchmarkCollectionEach(b *testing.B) {
	type Elem struct {
		s          string
		a, b, c, d int64
		n          int
	}
	b.Run("Each_1K", func(b *testing.B) {
		b.ReportAllocs()
		var c Collection[Elem]
		for i := 0; i < 1_000; i++ {
			c.Push(Elem{n: i, a: int64(i), b: -1, s: "Each 1K"})
		}
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			c.Each(func(e *Elem) bool {
				e.a++
				return true
			})
		}
	})
	b.Run("Each_65K", func(b *testing.B) {
		b.ReportAllocs()
		var c Collection[Elem]
		for i := 0; i < 65_535; i++ {
			c.Push(Elem{n: i, a: int64(i), b: -1, s: "Each 65K"})
		}
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			c.Each(func(e *Elem) bool {
				e.a++
				return true
			})
		}
	})
	b.Run("Each_1M", func(b *testing.B) {
		b.ReportAllocs()
		var c Collection[Elem]
		for i := 0; i < 1_000_000; i++ {
			c.Push(Elem{n: i, a: int64(i), b: -1, s: "Each 1m"})
		}
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			c.Each(func(e *Elem) bool {
				e.a++
				return true
			})
		}
	})
	b.Run("Push_10M", func(b *testing.B) {
		b.ReportAllocs()
		var c Collection[Elem]
		for i := 0; i < 10_000_000; i++ {
			c.Push(Elem{n: i, a: int64(i), b: -1, s: "Push 10m"})
		}
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			c.Each(func(e *Elem) bool {
				e.a++
				return true
			})
		}
	})
}

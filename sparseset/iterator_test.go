package sparseset

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator(t *testing.T) {
	t.Parallel()
	t.Run("Empty", func(t *testing.T) {
		t.Parallel()
		var c = New[int, int](0, 1000)
		for x := range c.Iterator(context.Background(), 16) {
			t.Errorf("unexpected: %d", x)
		}
	})
	t.Run("100000", func(t *testing.T) {
		t.Parallel()
		var c = New[int, int](0, 1024)
		const count = 1_000_000
		for n := 0; n < count; n++ {
			c.Set(n, n)
		}
		var i int
		for x := range c.Iterator(context.Background(), 16) {
			require.Equal(t, i, *x.Val)
			require.Equal(t, i, x.Key)
			i++
		}
	})
}

func BenchmarkFetcher(b *testing.B) {
	type Elem struct {
		s          string
		a, b, c, d int64
		n          int
	}
	var c = New[int, Elem](0, 1024)
	const count = 1_000_000
	for n := 0; n < count; n++ {
		c.Set(n, Elem{
			s: "BenchmarkFetcher",
			n: n,
		})
	}
	b.Run("Forward", func(b *testing.B) {
		b.ReportAllocs()
		left := b.N
		for {
			f := c.Fetcher()
			for f.Next() && left > 0 {
				_, _ = f.Fetch()
				if left--; left < 1 {
					return
				}
			}
		}
	})
	b.Run("AllInOne", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			f := c.Fetcher()
			for f.Next() {
				_, _ = f.Fetch()
			}
		}
	})
}

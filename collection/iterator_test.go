package collection

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIterator(t *testing.T) {
	t.Parallel()
	t.Run("Empty", func(t *testing.T) {
		t.Parallel()
		var c = Collection[int]{
			bsz: 1,
		}
		for x := range c.Iterator(context.Background(), 16) {
			t.Errorf("unexpected: %d", x)
		}
	})
	t.Run("100000", func(t *testing.T) {
		t.Parallel()
		var c Collection[int]
		const count = 100000
		for n := 0; n < count; n++ {
			c.Push(n)
		}
		var i int
		for x := range c.Iterator(context.Background(), 16) {
			require.Equal(t, i, *x)
			i++
		}
	})
	t.Run("Cancel", func(t *testing.T) {
		t.Parallel()
		var c Collection[int]
		for n := 0; n < 10; n++ {
			c.Push(0)
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		ch := c.Iterator(ctx, 0)
		time.Sleep(time.Millisecond)
		_, ok := <-ch
		require.False(t, ok)
	})
}

func BenchmarkIterator(b *testing.B) {
	type Elem struct {
		s          string
		a, b, c, d int64
		n          int
	}
	var c Collection[Elem]
	const count = 1_000_000
	for n := 0; n < count; n++ {
		c.Push(Elem{n: n})
	}
	b.Run("Forward", func(b *testing.B) {
		b.ReportAllocs()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		left := b.N
		for {
			for range c.Iterator(ctx, 16) {
				if left--; left < 1 {
					return
				}
			}
		}
	})
	b.Run("AllInOne", func(b *testing.B) {
		b.ReportAllocs()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		for n := 0; n < b.N; n++ {
			for range c.Iterator(ctx, 16) {
				// iter
			}
		}
	})
}

func BenchmarkFetcher(b *testing.B) {
	type Elem struct {
		s          string
		a, b, c, d int64
		n          int
	}
	var c Collection[Elem]
	const count = 1_000_000
	for n := 0; n < count; n++ {
		c.Push(Elem{n: n})
	}
	b.Run("Forward", func(b *testing.B) {
		b.ReportAllocs()
		left := b.N
		for {
			f := c.Fetcher()
			for f.Next() && left > 0 {
				_ = f.Fetch()
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
				_ = f.Fetch()
			}
		}
	})
}

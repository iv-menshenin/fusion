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
		var c = New[int, int](0, 1000)
		const count = 100000
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

package fsort

import (
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iv-menshenin/fusion/stack"
)

func TestSortable(t *testing.T) {
	t.Parallel()
	t.Run("Sort10", func(t *testing.T) {
		t.Parallel()
		var c stack.Stack[int]
		s := Sortable[int](&c, func(i *int, j *int) bool {
			return *i < *j
		})
		c.Push(0)
		c.Push(2)
		c.Push(6)
		c.Push(3)
		c.Push(5)
		c.Push(8)
		c.Push(9)
		c.Push(4)
		c.Push(1)
		c.Push(7)
		sort.Sort(s)
		for n := 0; n < 10; n++ {
			require.Equal(t, n, *c.Peek(n))
		}
	})
	t.Run("Sort1Mln", func(t *testing.T) {
		t.Parallel()
		var c stack.Stack[int]
		s := Sortable[int](&c, func(i *int, j *int) bool {
			return *i < *j
		})
		const max = 1000000
		for n := max; n > 0; n-- {
			c.Push(n - 1)
		}
		sort.Sort(s)
		for n := 0; n < max; n++ {
			require.Equal(t, n, *c.Peek(n))
		}
	})
}

func BenchmarkSortable(b *testing.B) {
	var c stack.Stack[string]
	const max = 1000
	for n := 0; n < max; n++ {
		c.Push(strconv.Itoa(n))
	}
	s := Sortable[string](&c, func(i *string, j *string) bool {
		return *i < *j
	})
	b.Run("Sortable", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			if n%2 == 0 {
				sort.Sort(s)
			} else {
				sort.Sort(sort.Reverse(s))
			}
		}
	})
	b.Run("IsSorted", func(b *testing.B) {
		sort.Sort(s)
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			sort.IsSorted(s)
		}
	})
	b.Run("NotSorted", func(b *testing.B) {
		if b.N > 1 {
			s.Swap(0, 1)
		}
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			sort.IsSorted(s)
		}
	})
}

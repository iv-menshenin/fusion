package sparseset

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSparseSet(t *testing.T) {
	sp := New[int64, string](0)

	sp.Set(10, "foo")
	sp.Set(15, "bar")
	sp.Set(25, "baz")
	sp.Set(23, "qux")
	sp.Set(34, "quux")
	sp.Set(56, "corge")
	sp.Set(78, "grault")

	sp.Delete(23)

	sp.Set(1, "garply")
	sp.Set(44, "waldo")
	sp.Set(51, "fred")

	sp.Delete(10)

	sp.Set(0, "plugh")
	sp.Set(88, "xyzzy")
	sp.Set(91, "thud")

	require.Equal(t, "thud", sp.Pop())
	require.Equal(t, "xyzzy", sp.Pop())

	require.Equal(t, "fred", *sp.Get(51))
	require.Equal(t, "baz", *sp.Get(25))
	require.Equal(t, "corge", *sp.Get(56))

	require.Nil(t, sp.Get(10))
	require.Nil(t, sp.Get(23))
	require.Nil(t, sp.Get(88))
	require.Nil(t, sp.Get(91))
}

func TestSparseSetMass(t *testing.T) {
	sp := New[int64, string](1000)

	const count = 1000000
	for n := 0; n < count; n += 2 {
		sp.Set(int64(n), strconv.Itoa(n))
		sp.Set(int64(n+1), strconv.Itoa(n+1))
	}

	// test
	for n := 0; n < count; n++ {
		val := sp.Get(int64(n))
		if expected := strconv.Itoa(n); *val != expected {
			t.Errorf("expected %q, got %q", expected, *val)
		}
	}

	// pop last
	last := sp.Pop()
	if last != strconv.Itoa(count-1) {
		t.Fatal("wrong last value")
	}
	sp.Set(count-1, "last")

	// delete all thirds
	for n := 3; n < count; n += 3 {
		if n == count-1 {
			continue
		}
		sp.Delete(int64(n))
	}

	// replace all fifths
	for n := 5; n < count; n += 5 {
		if n == count-1 {
			continue
		}
		sp.Set(int64(n), "foo"+strconv.Itoa(n))
	}

	// test again
	for n := 0; n < count; n++ {
		val := sp.Get(int64(n))
		switch {
		case n == count-1:
			if expected := "last"; *val != expected {
				t.Errorf("expected %q, got %q", expected, *val)
			}
		case n > 0 && n%5 == 0:
			if expected := "foo" + strconv.Itoa(n); *val != expected {
				t.Errorf("expected %q, got %q", expected, *val)
			}
		case n > 0 && n%3 == 0:
			if val != nil {
				t.Errorf("expected nil, got %q", *val)
			}
		default:
			if expected := strconv.Itoa(n); *val != expected {
				t.Errorf("expected %q, got %q", expected, *val)
			}
		}
	}
}

func BenchmarkSparseSet(b *testing.B) {
	sp := New[int, string](1000)

	// init fill
	for i := 0; i < 1000000; i++ {
		sp.Set(i, fmt.Sprintf("inited-%d", i))
	}

	b.Run("insert", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			sp.Set(i, fmt.Sprintf("inserted-%d", i))
		}
	})

	b.Run("pop", func(b *testing.B) {
		var idx int
		for sp.Len() < b.N {
			sp.Set(idx, strconv.Itoa(idx))
			idx++
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			sp.Pop()
		}
	})
}

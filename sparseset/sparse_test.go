package sparseset

import (
	"strconv"
	"testing"
)

func TestSparseSet(t *testing.T) {
	sp := New[int64, string](1000)

	const count = 1000000
	for n := 0; n < count; n += 2 {
		sp.Set(int64(n+1), strconv.Itoa(n+1))
		sp.Set(int64(n), strconv.Itoa(n))
	}

	// test
	for n := 0; n < count; n++ {
		val := sp.Get(int64(n))
		if expected := strconv.Itoa(n); *val != expected {
			t.Errorf("expected %q, got %q", expected, *val)
		}
	}

	// delete all thirds
	for n := 3; n < count; n += 3 {
		sp.Delete(int64(n))
	}

	// replace all fifths
	for n := 5; n < count; n += 5 {
		sp.Set(int64(n), "foo"+strconv.Itoa(n))
	}

	// test again
	for n := 0; n < count; n++ {
		val := sp.Get(int64(n))
		switch {
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

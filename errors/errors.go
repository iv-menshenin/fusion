package errors

import "fmt"

type ErrOutOfBounds struct {
	l, i int
}

func (e ErrOutOfBounds) Error() string {
	if e.l == 0 {
		return fmt.Sprintf("index out of bounds: is empty")
	}
	return fmt.Sprintf("index out of bounds: requested %d element with %d length", e.i, e.l)
}

func OutOfBounds(l, i int) error {
	return ErrOutOfBounds{l: l, i: i}
}

package stack

type Cmp[T any] struct {
	sortable   *Stack[T]
	comparator func(*T, *T) bool
}

func Sortable[T any](s *Stack[T], less func(*T, *T) bool) Cmp[T] {
	return Cmp[T]{
		sortable:   s,
		comparator: less,
	}
}

func (c Cmp[T]) Len() int {
	return c.sortable.Len()
}

func (c Cmp[T]) Less(i, j int) bool {
	return c.comparator(c.sortable.Get(i), c.sortable.Get(j))
}

func (c Cmp[T]) Swap(i, j int) {
	a, b := c.sortable.Get(i), c.sortable.Get(j)
	*a, *b = *b, *a
}

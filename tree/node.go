package tree

type (
	Node[IDX Ordered, D any] struct {
		idx  IDX
		data *D
	}
	Ordinary interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string
	}
	Ordered interface {
		Ordinary
	}
)

func NewNode[IDX Ordered, D any](idx IDX, data *D) Node[IDX, D] {
	return Node[IDX, D]{
		idx:  idx,
		data: data,
	}
}

func (n Node[IDX, D]) ID() IDX {
	return n.idx
}

func (n Node[IDX, D]) Data() *D {
	return n.data
}

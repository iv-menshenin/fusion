package tree

type (
	Chains[T any] struct {
		chains []Chain[T]
	}
	Chain[T any] struct {
		data []T
	}
)

const (
	largeChainSize = 65535
	preAllocChains = 128
)

func (c *Chains[T]) Put(v T) {
	if len(c.chains) == 0 {
		c.chains = make([]Chain[T], 1, preAllocChains)
	}
	last := len(c.chains) - 1
	if c.chains[last].PutIfFits(v) {
		return
	}
	c.chains = append(c.chains, Chain[T]{})
	c.chains[last+1].PutIfFits(v)
}

func (c *Chains[T]) Get(v T) {

}

func (h *Chain[T]) PutIfFits(v T) bool {
	if len(h.data) > 0 {
		l, c := len(h.data), cap(h.data)
		large := c > largeChainSize
		filled := l >= c
		if filled && large {
			return false
		}
	}
	h.data = append(h.data, v)
	return true
}

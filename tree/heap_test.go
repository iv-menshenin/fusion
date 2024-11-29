package tree

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/iv-menshenin/fusion/collection"
	"github.com/iv-menshenin/fusion/stack"
)

func Test_NewHeap(t *testing.T) {
	var s stack.Stack[Node[int, int]]
	tree := NewHeap[int, int](&s)
	tree.Put(
		NewNode[int, int](42, nil), NewNode[int, int](13, nil), NewNode[int, int](54, nil),
		NewNode[int, int](2, nil), NewNode[int, int](1, nil), NewNode[int, int](67, nil),
		NewNode[int, int](43, nil), NewNode[int, int](23, nil),
	)
	if err := checkBalance(tree); err != nil {
		t.Error(err)
	}

	tree.Put(NewNode[int, int](16, nil))
	if err := checkBalance(tree); err != nil {
		t.Error(err)
	}

	tree.Put(NewNode[int, int](43, nil))
	if err := checkBalance(tree); err != nil {
		t.Error(err)
	}

	tree.Put(NewNode[int, int](67, nil))
	if err := checkBalance(tree); err != nil {
		t.Error(err)
	}

	var last = -1
	for {
		if err := checkBalance(tree); err != nil {
			t.Error(err)
		}
		x, ok := tree.PopMax()
		if !ok {
			break
		}
		if last > 0 && x.idx > last {
			t.Error("got greater than prev")
		}
		last = x.idx
	}
}

func Test_ExistingHeap(t *testing.T) {
	data := []Node[int, int]{
		NewNode[int, int](42, nil), NewNode[int, int](13, nil), NewNode[int, int](54, nil),
		NewNode[int, int](2, nil), NewNode[int, int](1, nil), NewNode[int, int](67, nil),
		NewNode[int, int](43, nil), NewNode[int, int](23, nil),
	}
	var s = stack.Init[Node[int, int]](data)
	tree := NewHeap[int, int](s)
	if err := checkBalance(tree); err != nil {
		t.Error(err)
	}

	tree.Put(NewNode[int, int](16, nil))
	if err := checkBalance(tree); err != nil {
		t.Error(err)
	}

	tree.Put(NewNode[int, int](43, nil))
	if err := checkBalance(tree); err != nil {
		t.Error(err)
	}

	tree.Put(NewNode[int, int](67, nil))
	if err := checkBalance(tree); err != nil {
		t.Error(err)
	}

	var last = -1
	for {
		if err := checkBalance(tree); err != nil {
			t.Error(err)
		}
		x, ok := tree.PopMax()
		if !ok {
			break
		}
		if last > 0 && x.idx > last {
			t.Error("got greater than prev")
		}
		last = x.idx
	}
}

func checkBalance(tree *Heap[int, int]) error {
	for i := 0; i < tree.heap.Len(); i++ {
		x := tree.heap.Get(i)
		if l := idxOfLeftChild(i); l < tree.heap.Len() {
			if tree.heap.Get(l).idx > x.idx {
				return fmt.Errorf("node #%d has %d which less than it`s child #%d eqal %d\nfull: %v", i, x, l, tree.heap.Get(l).idx, tree.heap)
			}
		}
		if r := idxOfRightChild(i); r < tree.heap.Len() {
			if tree.heap.Get(r).idx > x.idx {
				return fmt.Errorf("node #%d has %d which less than it`s child #%d eqal %d\nfull: %v", i, x, r, tree.heap.Get(r).idx, tree.heap)
			}
		}
	}
	return nil
}

func Test_SortByHeap(t *testing.T) {
	heap := NewHeap[int, int](nil)
	heap.Put(NewNode[int, int](199, nil))
	heap.Put(NewNode[int, int](5, nil), NewNode[int, int](6, nil), NewNode[int, int](8, nil))
	heap.Put(
		NewNode[int, int](23, nil), NewNode[int, int](7, nil), NewNode[int, int](99, nil),
		NewNode[int, int](1, nil), NewNode[int, int](9322, nil),
	)
	if l := heap.Len(); l != 9 {
		t.Errorf("unexpected length = %d", l)
	}
	var sorted = make([]int, 0, 9)
	for {
		x, ok := heap.PopMax()
		if !ok {
			break
		}
		sorted = append(sorted, x.ID())
	}
	if out := fmt.Sprintf("%+v", sorted); out != "[9322 199 99 23 8 7 6 5 1]" {
		t.Errorf("unexpected result: %q", out)
	}
}

func Test_MemAllocHeap(t *testing.T) {
	testData := []int64{
		11515, 11060, 10149, 14106, 8455, 15011, 351, 18648, 14977, 12031, 17133, 3267, 474, 13121, 10317, 7667, 13085,
		17525, 13960, 22954, 19786, 17876, 19255, 13697, 13904, 8465, 22562, 6410, 6091, 20571, 19691, 562, 14364, 17474,
		16705, 6627, 7865, 12715, 1725, 8424, 301, 17197, 14884, 11808, 1878, 7631, 13099, 22605, 20784, 10772, 15238,
		9156, 11909, 14731, 3905, 15847, 5545, 19264, 14313, 947, 621, 5447, 22300, 5628, 4158, 8241, 8991, 20868, 10700,
		9567, 13939, 5342, 11574, 8504, 9921, 18960, 4437, 6495, 11022, 21598, 21065, 10701, 5205, 21224, 12140, 6364,
		2266, 5964, 16297, 1725, 11793, 11707, 16795, 20712, 789, 15228, 7533, 357, 11296, 18173, 12, 34, 5, 100, 9090,
	}
	idx := 0
	heap := NewHeap[int64, int](nil)
	cnt := testing.AllocsPerRun(1000, func() {
		heap.Put(NewNode[int64, int](testData[idx%len(testData)], nil))
		idx++
	})
	if cnt > 0 {
		t.Errorf("allocation count: %0.4f", cnt)
	}
	var d int
	cnt = testing.AllocsPerRun(1000, func() {
		_, _ = heap.Search(testData[d%len(testData)])
		d++
	})
	if cnt > 1 {
		t.Errorf("allocation count: %0.4f", cnt)
	}
	cnt = testing.AllocsPerRun(1000, func() {
		_, _ = heap.PopMax()
	})
	if cnt > 0 {
		t.Errorf("allocation count: %0.4f", cnt)
	}
}

func Benchmark_SearchTree(b *testing.B) {
	const testCnt = 12_000_000
	var st = collection.Collection[Node[int, int]]{}
	var tree = NewHeap[int, int](&st)
	var dataUniq = make(map[int]struct{}, testCnt)
	for len(dataUniq) < testCnt {
		i := rand.Int()
		if _, ok := dataUniq[i]; ok {
			continue
		}
		dataUniq[i] = struct{}{}
		tree.Put(NewNode(i, &i))
	}

	b.Run("random_search", func(b *testing.B) {
		b.ReportAllocs()
		var cnt int
		for {
			for idx := range dataUniq {
				node, ok := tree.Search(idx)
				if !ok {
					b.Errorf("can't find #%d", idx)
					return
				}
				if node.ID() != idx {
					b.Errorf("fund bad %d != %d", idx, node.idx)
				}
				cnt++
				if cnt > b.N {
					return
				}
			}
		}
	})
}

func Benchmark_PutHeapTree(b *testing.B) {
	var testCnt = b.N
	var tree = NewHeap[int, int](nil)
	var dataUniq = make(map[int]struct{}, testCnt)
	for len(dataUniq) < testCnt {
		i := rand.Int()
		if _, ok := dataUniq[i]; ok {
			continue
		}
		dataUniq[i] = struct{}{}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for idx := range dataUniq {
		i := idx
		tree.Put(NewNode[int, int](i, nil))
	}
}

package skiplist

import (
	"bytes"
	"math/rand/v2"
	"sync/atomic"
	"time"
)

type SkipList struct {
	head *Node
	cmp  CompareFunc

	maxHeight       int
	branchingFactor int
}

type Node struct {
	next      []atomic.Pointer[Node]
	key       []byte
	value     []byte
	timestamp time.Time
}

// Compare returns
//
//	-1 if x is less than y,
//	 0 if x equals y,
//	+1 if x is greater than y.
type CompareFunc func(a, b []byte) int

func New(opts ...SkipListOpt) *SkipList {
	newSkipList := SkipList{
		maxHeight:       12,
		branchingFactor: 4,
		cmp:             bytes.Compare,
	}

	for _, opt := range opts {
		opt(&newSkipList)
	}

	newSkipList.head = &Node{
		next: make([]atomic.Pointer[Node], newSkipList.maxHeight),
	}

	return &newSkipList
}

func (sl *SkipList) NewIterator() *Iterator {
	return &Iterator{
		current: sl.head.next[0].Load(),
	}
}

func (sl *SkipList) Put(key, value []byte) {
	var predecessors, successors []*Node

	for {
		predecessors, successors = sl.findPredecessors(key)

		if successors[0] != nil && sl.cmp(successors[0].key, key) == 0 {
			successors[0].value = value
			return
		}
		height := sl.randomHeight()
		newNode := newNode(key, value, height)

		for i := 0; i < height; i++ {
			newNode.next[i].Store(successors[i])
		}

		success := true
		for i := 0; i < height; i++ {
			if !predecessors[i].next[i].CompareAndSwap(successors[i], newNode) {
				success = false
				break
			}
		}

		if success {
			break
		}
	}

}

func (sl *SkipList) Get(key []byte) ([]byte, bool) {
	current := sl.head
	for i := sl.maxHeight - 1; i >= 0; i-- {
		next := current.next[i].Load()
		for next != nil && sl.cmp(next.key, key) < 0 {
			current = next
			next = current.next[i].Load()
		}
	}

	current = current.next[0].Load()
	if current != nil && sl.cmp(current.key, key) == 0 {
		return current.value, true
	}
	return nil, false
}

func (sl *SkipList) findPredecessors(key []byte) ([]*Node, []*Node) {
	pre := make([]*Node, sl.maxHeight)
	suc := make([]*Node, sl.maxHeight)

	curr := sl.head
	for i := sl.maxHeight - 1; i >= 0; i-- {
		next := curr.next[i].Load()
		for next != nil && sl.cmp(key, next.key) > 0 {
			curr = next
			next = curr.next[i].Load()
		}
		suc[i] = next
		pre[i] = curr
	}
	return pre, suc
}

func (sl *SkipList) randomHeight() int {
	height := 1
	for height < int(sl.maxHeight) && rand.IntN(sl.branchingFactor) == 0 {
		height++
	}
	return height
}

func newNode(key, value []byte, height int) *Node {
	return &Node{
		key:       key,
		value:     value,
		next:      make([]atomic.Pointer[Node], height),
		timestamp: time.Now().UTC(),
	}
}

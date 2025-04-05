package skiplist

import (
	"sync"
	"sync/atomic"
)

// TODO: implement
type NodePool struct {
	pool sync.Pool
}

func newPool(maxHeight int) NodePool {
	return NodePool{
		sync.Pool{
			New: func() interface{} {
				return &Node{
					next: make([]atomic.Pointer[Node], maxHeight),
				}
			},
		},
	}
}

func (np *NodePool) get() *Node {
	return np.pool.Get().(*Node)
}

func (np *NodePool) put(n *Node) {
	np.pool.Put(n)
}

package cache

import (
	"sync"
	"weak"
)

type Weak[K comparable, V any] struct {
	mu    sync.Mutex
	items map[K]weak.Pointer[V]
}

func NewWeak[K comparable, V any]() *Weak[K, V] {
	return &Weak[K, V]{
		items: make(map[K]weak.Pointer[V], 64),
	}
}

func (c *Weak[K, V]) Get(key K) (*V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ptr, exists := c.items[key]
	if !exists {
		return nil, false
	}

	val := ptr.Value()
	if val == nil {
		delete(c.items, key)
		return nil, false
	}

	return val, true
}

func (c *Weak[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = weak.Make(&value)
}

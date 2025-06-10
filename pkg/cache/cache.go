package cache

import (
	"sync"
)

type Hash[K comparable, V any] struct {
	keyCleaner func(K) K
	updateFn   func(V, V) V

	mu     *sync.RWMutex
	values map[K]V
}

func New[K comparable, V any](keyCleaner func(K) K, updateFn func(V, V) V) *Hash[K, V] {
	return &Hash[K, V]{
		keyCleaner: keyCleaner,
		updateFn:   updateFn,
		mu:         &sync.RWMutex{},
		values:     make(map[K]V),
	}
}

func (c *Hash[K, V]) Has(key K) bool {
	if c.keyCleaner != nil {
		key = c.keyCleaner(key)
	}
	var has bool
	c.mu.RLock()
	_, has = c.values[key]
	c.mu.RUnlock()
	return has
}

func (c *Hash[K, V]) Add(key K, value V) {
	if c.keyCleaner != nil {
		key = c.keyCleaner(key)
	}
	c.mu.Lock()
	if existing, ok := c.values[key]; ok && c.updateFn != nil {
		c.values[key] = c.updateFn(existing, value)
	} else {
		c.values[key] = value
	}
	c.mu.Unlock()
}

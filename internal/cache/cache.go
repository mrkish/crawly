package cache

import (
	"sync"
)

type Cache[K comparable, V any] struct {
	updateFn func(V, V) V

	mu     sync.RWMutex
	values map[K]V
}

func New[K comparable, V any](updateFn func(V, V) V) *Cache[K, V] {
	return &Cache[K, V]{
		updateFn: updateFn,
		mu:       sync.RWMutex{},
		values:   make(map[K]V),
	}
}

func (c *Cache[K, V]) Has(key K) bool {
	var has bool
	c.mu.RLock()
	_, has = c.values[key]
	c.mu.RUnlock()
	return has
}

func (c *Cache[K, V]) Add(key K, value V) {
	c.mu.Lock()
	if existing, ok := c.values[key]; ok {
		c.values[key] = c.updateFn(existing, value)
	} else {
		c.values[key] = value
	}
	c.mu.Unlock()
}

package cache

import (
	"sync"
)

type Hashed[K comparable, V any] struct {
	keyCleaner func(K) K
	updateFn   func(V, V) V

	mu     *sync.RWMutex
	values map[K]V
}

func NewHashed[K comparable, V any](keyCleaner func(K) K, updateFn func(V, V) V) *Hashed[K, V] {
	return &Hashed[K, V]{
		keyCleaner: keyCleaner,
		updateFn:   updateFn,
		mu:         &sync.RWMutex{},
		values:     make(map[K]V),
	}
}

func (c *Hashed[K, V]) Has(key K) (has bool) {
	if c.keyCleaner != nil {
		key = c.keyCleaner(key)
	}
	c.mu.RLock()
	_, has = c.values[key]
	c.mu.RUnlock()
	return has
}

func (c *Hashed[K, V]) Get(key K) V {
	if c.keyCleaner != nil {
		key = c.keyCleaner(key)
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.values[key]
}

func (c *Hashed[K, V]) Add(key K, value V) {
	if c.keyCleaner != nil {
		key = c.keyCleaner(key)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if existing, ok := c.values[key]; ok && c.updateFn != nil {
		c.values[key] = c.updateFn(existing, value)
	} else {
		c.values[key] = value
	}
}

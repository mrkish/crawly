package cache

import (
	"crawly/internal/model"
	"sync"
)

type Cache struct {
	mu sync.RWMutex
	crawled map[string]*model.Fetched
}

func New() *Cache {
	return &Cache{
		mu: sync.RWMutex{},
		crawled: make(map[string]*model.Fetched),
	}
}

func (c *Cache) Has(url string) bool {
	var has bool
	c.mu.RLock()
	_, has = c.crawled[url]
	c.mu.RUnlock()
	return has
}

func (c *Cache) Add(fetched *model.Fetched) {
	c.mu.Lock()
	c.crawled[fetched.URL] = fetched
	c.mu.Unlock()
}

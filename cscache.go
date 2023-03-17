package gcache

import (
	"gcache/lru"
	"sync"
)

// cache 把lru.Cache封装成并发安全 Concurrent security
type csCache struct {
	mu     sync.Mutex
	lru    *lru.Cache
	maxCap int
}

func (c *csCache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		//延时加载lru.Cache
		c.lru = lru.New(c.maxCap)
	}
	c.lru.Add(key, value)
}

func (c *csCache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}

func (c *csCache) delete(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return false
	}
	ok := c.lru.Delete(key)
	return ok
}

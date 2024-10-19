package geecache

import (
	"Gee/GeeCache/lru"
	"sync"
)

// 为数据添加并发特性

// 毕竟只有增、改的时候会发生并发并发错误，查询、删除不会。但是注意：在数据库数据访问中是不行的，需要mvcc进行版本粒度控制。

type cache struct {
	m          sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) Add(key string, value ByteView) {
	c.m.Lock()
	defer c.m.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) Get(key string) (value ByteView, ok bool) {
	c.m.Lock()
	defer c.m.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}

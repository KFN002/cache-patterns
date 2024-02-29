package main

import (
	"sync"
)

type BCache struct {
	data  map[string]interface{}
	mutex sync.RWMutex
}

func NewBCache() *BCache {
	return &BCache{
		data: make(map[string]interface{}),
	}
}

func (c *BCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, ok := c.data[key]
	return value, ok
}

func (c *BCache) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = value
}

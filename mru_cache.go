package main

import (
	"container/list"
	"sync"
)

type MRUCache struct {
	capacity int
	cache    map[string]*list.Element
	order    *list.List
	mutex    sync.Mutex
}

type cacheItem struct {
	key   string
	value string
}

func NewMRUCache(capacity int) *MRUCache {
	return &MRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		order:    list.New(),
	}
}

func (c *MRUCache) Set(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if element, ok := c.cache[key]; ok {
		element.Value.(*cacheItem).value = value
		c.order.MoveToFront(element)
	} else {
		item := &cacheItem{key: key, value: value}
		element := c.order.PushFront(item)
		c.cache[key] = element
		if c.order.Len() > c.capacity {
			lastElement := c.order.Back()
			delete(c.cache, lastElement.Value.(*cacheItem).key)
			c.order.Remove(lastElement)
		}
	}
}

func (c *MRUCache) Get(key string) (string, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if element, ok := c.cache[key]; ok {
		c.order.MoveToFront(element)
		return element.Value.(*cacheItem).value, true
	}
	return "", false
}

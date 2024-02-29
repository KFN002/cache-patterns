package main

import (
	"container/list"
	"fmt"
	"sync"
)

type LRUCache struct {
	capacity int
	cache    map[string]*list.Element
	list     *list.List
	mutex    sync.Mutex
}

type CacheEntry struct {
	key   string
	value interface{}
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (lru *LRUCache) Get(key string) (interface{}, bool) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	if element, ok := lru.cache[key]; ok {
		lru.list.MoveToFront(element)
		return element.Value.(*CacheEntry).value, true
	}

	return nil, false
}

func (lru *LRUCache) Set(key string, value interface{}) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	if element, ok := lru.cache[key]; ok {
		element.Value.(*CacheEntry).value = value
		lru.list.MoveToFront(element)
	} else {
		entry := &CacheEntry{key: key, value: value}
		element := lru.list.PushFront(entry)
		lru.cache[key] = element

		if lru.list.Len() > lru.capacity {
			oldest := lru.list.Back()
			if oldest != nil {
				delete(lru.cache, oldest.Value.(*CacheEntry).key)
				lru.list.Remove(oldest)
			}
		}
	}
}

func (lru *LRUCache) PrintCache() {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	fmt.Printf("LRU Cache (Capacity: %d, Size: %d): [", lru.capacity, lru.list.Len())
	for element := lru.list.Front(); element != nil; element = element.Next() {
		entry := element.Value.(*CacheEntry)
		fmt.Printf("(%s: %v) ", entry.key, entry.value)
	}
	fmt.Println("]")
}

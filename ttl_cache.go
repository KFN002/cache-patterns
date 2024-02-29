package main

import (
	"fmt"
	"sync"
	"time"
)

type TCacheEntry struct {
	value    interface{}
	expireAt int64
}

func NewCacheEntry(value interface{}, expireAt int64) TCacheEntry {
	return TCacheEntry{
		value:    value,
		expireAt: expireAt,
	}
}

func (ce TCacheEntry) IsExpired() bool {
	return ce.expireAt < time.Now().UnixNano()
}

type TCache struct {
	kvstore  map[string]TCacheEntry
	locker   sync.RWMutex
	interval time.Duration
	stop     chan struct{}
}

func NewTCache(cleanUpInterval time.Duration) *TCache {
	cache := &TCache{
		kvstore:  make(map[string]TCacheEntry),
		interval: cleanUpInterval,
		stop:     make(chan struct{}),
	}

	if cleanUpInterval > 0 {
		go cache.cleaning()
	}
	return cache
}

func (c *TCache) cleaning() {
	fmt.Println("cleaner starting...")
	ticker := time.NewTicker(c.interval)
	fmt.Println("cleaner was started")
	for {
		select {
		case <-ticker.C:
			c.purge()
		case <-c.stop:
			ticker.Stop()
			fmt.Println("cleaner was stopped")
			return
		}
	}
}

func (c *TCache) purge() {
	c.locker.Lock()
	defer c.locker.Unlock()
	for key, data := range c.kvstore {
		if data.IsExpired() {
			delete(c.kvstore, key)
		}
	}
}

func (c *TCache) Set(key string, value interface{}, expiryDuration time.Duration) {
	expireAt := time.Now().Add(expiryDuration).UnixNano()
	c.locker.Lock()
	defer c.locker.Unlock()
	c.kvstore[key] = NewCacheEntry(value, expireAt)

}

func (c *TCache) Get(key string) (interface{}, bool) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	data, found := c.kvstore[key]
	if !found || data.IsExpired() {
		return nil, false
	}

	return data.value, true
}

func (c *TCache) Close() {
	close(c.stop)
}

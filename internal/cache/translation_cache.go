package cache

import (
	"log"
	"sync"
	"time"

	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/models"
)

type Cache interface {
	Get(key string) ([]models.Translation, bool)
	Set(key string, data []models.Translation)
	Invalidate(key string)
	Close()
}

type cacheItem struct {
	Data      []models.Translation
	ExpiresAt time.Time
}

type TranslationCache struct {
	mu sync.RWMutex

	items map[string]cacheItem

	ttl time.Duration

	cleanupTicker *time.Ticker
	stopCleanup   chan struct{}

	hits   uint64
	misses uint64
}

func NewTranslationCache(
	ttl time.Duration,
) *TranslationCache {

	cache := &TranslationCache{
		items:         make(map[string]cacheItem),
		ttl:           ttl,
		cleanupTicker: time.NewTicker(ttl),
		stopCleanup:   make(chan struct{}),
	}

	go cache.startCleanupLoop()

	return cache
}

func (c *TranslationCache) Get(
	key string,
) ([]models.Translation, bool) {

	c.mu.RLock()
	item, exists := c.items[key]
	c.mu.RUnlock()

	if !exists {
		c.recordMiss()
		return nil, false
	}

	if c.isExpired(item) {

		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		c.recordMiss()

		return nil, false
	}

	c.recordHit()

	return item.Data, true
}

func (c *TranslationCache) Set(
	key string,
	data []models.Translation,
) {

	c.mu.Lock()

	c.items[key] = cacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(c.ttl),
	}

	c.mu.Unlock()
}

func (c *TranslationCache) Invalidate(
	key string,
) {

	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

func (c *TranslationCache) Close() {
	close(c.stopCleanup)
	c.cleanupTicker.Stop()
}

func (c *TranslationCache) startCleanupLoop() {

	for {

		select {

		case <-c.cleanupTicker.C:
			c.cleanupExpired()

		case <-c.stopCleanup:
			return
		}
	}
}

func (c *TranslationCache) cleanupExpired() {

	now := time.Now()

	c.mu.Lock()
	defer c.mu.Unlock()

	for key, item := range c.items {

		if now.After(item.ExpiresAt) {
			delete(c.items, key)
		}
	}

	log.Printf(
		"cache cleanup completed items=%d hits=%d misses=%d",
		len(c.items),
		c.hits,
		c.misses,
	)
}

func (c *TranslationCache) isExpired(
	item cacheItem,
) bool {

	return time.Now().After(
		item.ExpiresAt,
	)
}

func (c *TranslationCache) recordHit() {
	c.hits++
}

func (c *TranslationCache) recordMiss() {
	c.misses++
}


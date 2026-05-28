package cache

import (
	"log"
	"sync"
	"time"

	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/models"
)

type cacheItem struct {
	Data      []models.Translation
	ExpiresAt time.Time
}

type TranslationCache struct {
	mu       sync.RWMutex
	items    map[string]cacheItem
	ttl      time.Duration
	interval time.Duration
}

func NewTranslationCache(
	ttl time.Duration,
) *TranslationCache {

	cache := &TranslationCache{
		items:    make(map[string]cacheItem),
		ttl:      ttl,
		interval: ttl,
	}

	go cache.startCleanup()
	return cache
}

func (c *TranslationCache) Get(
	key string,
) ([]models.Translation, bool) {

	c.mu.RLock()
	item, exists := c.items[key]
	c.mu.RUnlock()

	if !exists {
		log.Printf(
			"translation cache miss key=%s",
			key,
		)

		return nil, false
	}

	if time.Now().After(item.ExpiresAt) {
		c.Invalidate(key)
		log.Printf(
			"translation cache expired key=%s",
			key,
		)

		return nil, false
	}

	log.Printf(
		"translation cache hit key=%s",
		key,
	)

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

func (c *TranslationCache) startCleanup() {

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {

		now := time.Now()
		c.mu.Lock()
		for key, item := range c.items {

			if now.After(item.ExpiresAt) {
				delete(c.items, key)
			}
		}

		c.mu.Unlock()
	}
}
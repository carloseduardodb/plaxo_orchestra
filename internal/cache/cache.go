package cache

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

type CacheEntry struct {
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	TTL       time.Duration `json:"ttl"`
}

type DistributedCache struct {
	local  map[string]*CacheEntry
	mutex  sync.RWMutex
	stats  *CacheStats
}

type CacheStats struct {
	Hits   int64 `json:"hits"`
	Misses int64 `json:"misses"`
	mutex  sync.RWMutex
}

func NewDistributedCache() *DistributedCache {
	return &DistributedCache{
		local: make(map[string]*CacheEntry),
		stats: &CacheStats{},
	}
}

func (c *DistributedCache) Get(ctx context.Context, key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	entry, exists := c.local[key]
	if !exists || time.Since(entry.Timestamp) > entry.TTL {
		c.stats.recordMiss()
		return nil, false
	}
	
	c.stats.recordHit()
	return entry.Data, true
}

func (c *DistributedCache) Set(ctx context.Context, key string, data interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.local[key] = &CacheEntry{
		Data:      data,
		Timestamp: time.Now(),
		TTL:       ttl,
	}
}

func (c *DistributedCache) GenerateKey(input string) string {
	hash := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", hash)
}

func (s *CacheStats) recordHit() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Hits++
}

func (s *CacheStats) recordMiss() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Misses++
}

func (s *CacheStats) GetStats() (int64, int64) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Hits, s.Misses
}

func (c *DistributedCache) GetStats() (int64, int64) {
	return c.stats.GetStats()
}

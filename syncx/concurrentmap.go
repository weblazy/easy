package syncx

import (
	"sync"

	"github.com/weblazy/goutil"
)

type (
	ConcurrentMap struct {
		cMap       []goutil.Map // Add the len method to sync.Map to get the total number of elements
		shareCount int          // Number of shards
	}
)

const (
	defaultShareCount = 32 // Default number of shards
)

// NewConcurrentMap creates a new ConcurrentMap.
func NewConcurrentMap(shareCount int) *ConcurrentMap {
	if shareCount == 0 {
		shareCount = defaultShareCount
	}
	m := make([]goutil.Map, shareCount)
	for i := 0; i < shareCount; i++ {
		m[i] = goutil.AtomicMap()
	}
	concurrentMap := &ConcurrentMap{
		shareCount: shareCount,
		cMap:       m,
	}
	return concurrentMap
}

// GetShard return a goutil.Map for a key
func (c *ConcurrentMap) GetShard(key string) goutil.Map {
	return c.cMap[uint(fnv32(key))%uint(c.shareCount)]
}

// fnv32 FNV hash algorithm
func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (c *ConcurrentMap) Load(key string) (interface{}, bool) {
	shard := c.GetShard(key)
	value, ok := shard.Load(key)
	return value, ok
}

// Store sets the value for a key.
func (c *ConcurrentMap) Store(key string, value interface{}) {
	shard := c.GetShard(key)
	shard.Store(key, value)
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (c *ConcurrentMap) LoadOrStore(key string, value interface{}) (actual interface{}, loaded bool) {
	shard := c.GetShard(key)
	actual, loaded = shard.LoadOrStore(key, value)
	return actual, loaded
}

// Clear clears all current data in the map.
func (c *ConcurrentMap) Clear() {
	for _, shard := range c.cMap {
		shard.Clear()
	}
}

// Len get the total number of elements
func (c *ConcurrentMap) Len() int {
	count := 0
	for _, shard := range c.cMap {
		length := shard.Len()
		count += length
	}
	return count
}

// Keys Get all the keys
func (c *ConcurrentMap) Keys() []interface{} {
	count := c.Len()
	ch := make(chan interface{}, count)
	// Each shard initiates a coroutine traverse elements
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(c.shareCount)
		for _, shard := range c.cMap {
			go func(shard goutil.Map) {
				shard.Range(func(k, v interface{}) bool {
					defer wg.Done()
					ch <- k
					return true
				})
			}(shard)
		}
		wg.Wait()
		close(ch)
	}()

	keys := make([]interface{}, count)
	// Collects the key of the map shard from each shard
	for k := range ch {
		keys = append(keys, k)
	}
	return keys
}

// Values Get all the values
func (c *ConcurrentMap) Values() []interface{} {
	count := c.Len()
	ch := make(chan interface{}, count)
	// Each shard initiates a coroutine traverse elements
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(c.shareCount)
		for _, shard := range c.cMap {
			shard.Range(func(k, v interface{}) bool {
				defer wg.Done()
				ch <- v
				return true
			})
		}
		wg.Wait()
		close(ch)
	}()

	values := make([]interface{}, count)
	// Collects the key of the map shard from each shard
	for value := range ch {
		values = append(values, value)
	}
	return values
}

// Delete deletes the value for a key.
func (c *ConcurrentMap) Delete(key string) {
	shard := c.GetShard(key)
	shard.Delete(key)
}

// Range traverse elements
func (c *ConcurrentMap) Range(f func(key, value interface{}) bool) {
	for _, shard := range c.cMap {
		if !shard.Range(f) {
			break
		}
	}
}

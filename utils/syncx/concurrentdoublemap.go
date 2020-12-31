package syncx

import (
	"sync"
	"sync/atomic"
)

const (
	maxDeleteCount = 10000
)

type (
	ConcurrentDoubleMap struct {
		cMap       []*ConcurrentMapShared
		length     int32
		shareCount int
	}
	ConcurrentMapShared struct {
		items        map[string]map[string]interface{}
		deleteCount  uint32
		sync.RWMutex // 各个分片Map各自的锁
	}
)

func NewConcurrentDoubleMap(shareCount int) *ConcurrentDoubleMap {
	if shareCount == 0 {
		shareCount = defaultShareCount
	}
	m := make([]*ConcurrentMapShared, shareCount)
	for i := 0; i < shareCount; i++ {
		m[i] = &ConcurrentMapShared{
			items: make(map[string]map[string]interface{}),
		}
	}
	return &ConcurrentDoubleMap{
		shareCount: shareCount,
		cMap:       m,
	}
}

func (c *ConcurrentDoubleMap) GetShard(key string) *ConcurrentMapShared {
	return c.cMap[uint(fnv32(key))%uint(c.shareCount)]
}

func (c *ConcurrentDoubleMap) Load(key1, key2 string) (interface{}, bool) {
	shard := c.GetShard(key1)
	shard.RLock()
	defer shard.RUnlock()
	oldMap, ok := shard.items[key1]
	if !ok {
		return nil, false
	}
	value, ok := oldMap[key2]
	return value, ok
}

func (c *ConcurrentDoubleMap) LoadMap(key1 string) (map[string]interface{}, bool) {
	shard := c.GetShard(key1)
	shard.RLock()
	defer shard.RUnlock()
	oldMap, ok := shard.items[key1]
	return oldMap, ok
}

func (c *ConcurrentDoubleMap) Store(key1, key2 string, value interface{}) {
	shard := c.GetShard(key1) // 段定位找到分片
	shard.Lock()
	oldMap, ok := shard.items[key1]
	if !ok {
		oldMap = make(map[string]interface{})
	}
	_, ok = oldMap[key2]
	if !ok {
		atomic.AddInt32(&c.length, 1)
	}
	oldMap[key2] = value
	shard.items[key1] = oldMap
	shard.Unlock()
}

func (c *ConcurrentDoubleMap) StoreWithPlugin(key1, key2 string, value interface{}, plugin func()) {
	shard := c.GetShard(key1) // 段定位找到分片
	shard.Lock()
	oldMap, ok := shard.items[key1]
	if !ok {
		oldMap = make(map[string]interface{})
	}
	_, ok = oldMap[key2]
	if !ok {
		atomic.AddInt32(&c.length, 1)
	}
	oldMap[key2] = value
	shard.items[key1] = oldMap
	plugin()
	shard.Unlock()
}

func (c *ConcurrentDoubleMap) LoadOrStore(key1, key2 string, value interface{}) (interface{}, bool) {
	shard := c.GetShard(key1)
	shard.Lock()
	defer shard.Unlock()
	oldMap, loaded := shard.items[key1]
	if !loaded {
		atomic.AddInt32(&c.length, 1)
		oldMap = make(map[string]interface{})
		oldMap[key2] = value
		shard.items[key1] = oldMap
		return value, loaded
	}
	v, loaded := oldMap[key2]
	if !loaded {
		atomic.AddInt32(&c.length, 1)
		oldMap[key2] = value
		shard.items[key1] = oldMap
		v = value
	}
	return v, loaded
}

// func (c *ConcurrentDoubleMap) Clear() {
// 	for _, shard := range c.CMap {
// 		shard.Lock()
// 		for k := range shard.items {
// 			delete(shard.items, k)
// 		}
// 		shard.Unlock()
// 	}
// }

// 统计当前分段map中item的个数
func (c *ConcurrentDoubleMap) Len() int32 {
	return atomic.LoadInt32(&c.length)
}

func (c *ConcurrentDoubleMap) Delete(key1, key2 string) {
	shard := c.GetShard(key1)
	shard.Lock()
	defer shard.Unlock()
	if shard.delete(key1, key2) {
		atomic.AddInt32(&c.length, -1)
	}
}

func (c *ConcurrentDoubleMap) DeleteWithPlugin(key1, key2 string, plugin func()) {
	shard := c.GetShard(key1)
	shard.Lock()
	defer shard.Unlock()
	if shard.delete(key1, key2) {
		atomic.AddInt32(&c.length, -1)
	}
	plugin()
}

func (c *ConcurrentDoubleMap) DeleteWithoutLock(key1, key2 string) {
	shard := c.GetShard(key1)
	if shard.delete(key1, key2) {
		atomic.AddInt32(&c.length, -1)
	}
}

func (c *ConcurrentDoubleMap) Range(f func(key1, key2 string, value interface{}) bool) bool {
	for _, shard := range c.cMap {
		shard.RLock()
		defer shard.RUnlock()
		for k1, oldMap := range shard.items {
			for k2, v := range oldMap {
				if !f(k1, k2, v) {
					return false
				}
			}
		}
	}
	return true
}

// func (c *ConcurrentDoubleMap) RangeShard(key1 string, f func(key2 string, value interface{}) bool) bool {
// 	shard := c.GetShard(key1)
// 	shard.RLock()
// 	defer shard.RUnlock()
// 	for k, v := range shard.items {
// 		if !f(k, v) {
// 			return false
// 		}
// 	}
// 	return true
// }

func (c *ConcurrentDoubleMap) RangeNextMap(key1 string, f func(key1, key2 string, value interface{}) bool) bool {
	shard := c.GetShard(key1)
	shard.RLock()
	defer shard.RUnlock()
	oldMap, ok := shard.items[key1]
	if !ok {
		return false
	}
	for k, v := range oldMap {
		if !f(key1, k, v) {
			return false
		}
	}

	return true
}

func (shard *ConcurrentMapShared) delete(key1, key2 string) bool {
	m1, ok := shard.items[key1]
	if !ok {
		return false
	}
	delete(m1, key2)
	shard.deleteCount++
	if len(m1) == 0 {
		delete(shard.items, key1)
	}
	if shard.deleteCount > maxDeleteCount {
		items := make(map[string]map[string]interface{})
		for k1 := range shard.items {
			items[k1] = shard.items[k1]
		}
		shard.items = items
		shard.deleteCount = 0
	}
	return true
}

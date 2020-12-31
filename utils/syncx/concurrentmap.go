package syncx

import (
	"github.com/weblazy/goutil"
	"sync"
)

type (
	ConcurrentMap struct {
		cMap       []goutil.Map //分段sync map
		shareCount int
	}
)

const (
	defaultShareCount = 32
)

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

func (c *ConcurrentMap) GetShard(key string) goutil.Map {
	return c.cMap[uint(fnv32(key))%uint(c.shareCount)]
}

// FNV hash
func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
func (c *ConcurrentMap) Load(key string) (interface{}, bool) {
	shard := c.GetShard(key)
	value, ok := shard.Load(key)
	return value, ok
}

func (c *ConcurrentMap) Store(key string, value interface{}) {
	shard := c.GetShard(key) // 段定位找到分片
	shard.Store(key, value)
}

func (c *ConcurrentMap) LoadOrStore(key string, value interface{}) (actual interface{}, loaded bool) {
	shard := c.GetShard(key)
	actual, loaded = shard.LoadOrStore(key, value)
	return actual, loaded
}

func (c *ConcurrentMap) Clear() {
	for _, shard := range c.cMap {
		shard.Clear()
	}
}

// 统计当前分段map中item的个数
func (c *ConcurrentMap) Len() int {
	count := 0
	for _, shard := range c.cMap {
		length := shard.Len()
		count += length
	}
	return count
}

// 获取所有的key
func (c *ConcurrentMap) Keys() []interface{} {
	count := c.Len()
	ch := make(chan interface{}, count)
	// 每一个分片启动一个协程 遍历key
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
	// 统计各个协程并发读取Map分片的key
	for k := range ch {
		keys = append(keys, k)
	}
	return keys
}

// 获取所有的value
func (c *ConcurrentMap) Values() []interface{} {
	count := c.Len()
	ch := make(chan interface{}, count)
	// 每一个分片启动一个协程 遍历key
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
	// 统计各个协程并发读取Map分片的key
	for value := range ch {
		values = append(values, value)
	}
	return values
}

func (c *ConcurrentMap) Delete(key string) {
	shard := c.GetShard(key)
	shard.Delete(key)
}

func (c *ConcurrentMap) Range(f func(key, value interface{}) bool) {
	for _, shard := range c.cMap {
		if !shard.Range(f) {
			break
		}
	}
}

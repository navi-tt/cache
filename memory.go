package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	// DefaultEvery means the clock time of recycling the expired cache items in memory.
	DefaultEvery = 60 // 1 minute
)

// MemoryItem store memory cache item.
type MemoryItem struct {
	val         interface{}
	createdTime time.Time
	lifespan    time.Duration
}

func (mi *MemoryItem) isExpire() bool {
	// 0 means forever
	if mi.lifespan == 0 {
		return false
	}
	return time.Now().Sub(mi.createdTime) > mi.lifespan
}

// MemoryCache is Memory cache adapter.
// it contains a RW locker for safe map storage.
type MemoryCache struct {
	sync.RWMutex
	dur   time.Duration
	items map[string]*MemoryItem
	Every int // run an expiration check Every clock time
}

// NewMemoryCache returns a new MemoryCache.
func NewMemoryCache() Cache {
	memoryCache := MemoryCache{items: make(map[string]*MemoryItem)}
	return &memoryCache
}

// Get cache from memory.
// if non-existed or expired, return nil.
func (bc *MemoryCache) Get(name string) (interface{}, error) {
	bc.RLock()
	defer bc.RUnlock()
	if itm, ok := bc.items[name]; ok {
		if itm.isExpire() {
			return nil, fmt.Errorf("key not existed")
		}
		return itm.val, nil
	}
	return nil, fmt.Errorf("server internal error")
}

// Put cache to memory.
// if lifespan is 0, it will be forever till restart.
func (bc *MemoryCache) Set(name string, value interface{}, expireTime int) error {
	bc.Lock()
	defer bc.Unlock()
	bc.items[name] = &MemoryItem{
		val:         value,
		createdTime: time.Now(),
		lifespan:    time.Duration(expireTime) * time.Second,
	}
	return nil
}

// Delete cache in memory.
func (bc *MemoryCache) Delete(name string) error {
	bc.Lock()
	defer bc.Unlock()
	if _, ok := bc.items[name]; !ok {
		return errors.New("key not exist")
	}
	delete(bc.items, name)
	if _, ok := bc.items[name]; ok {
		return errors.New("delete key error")
	}
	return nil
}

// IsExist check cache exist in memory.
func (bc *MemoryCache) IsExist(name string) bool {
	bc.RLock()
	defer bc.RUnlock()
	if v, ok := bc.items[name]; ok {
		return !v.isExpire()
	}
	return false
}

// StartAndGC start memory cache. it will check expiration in every clock time.
func (bc *MemoryCache) Init(config string) error {
	var cf map[string]int
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["interval"]; !ok {
		cf = make(map[string]int)
		cf["interval"] = DefaultEvery
	}
	dur := time.Duration(cf["interval"]) * time.Second
	bc.Every = cf["interval"]
	bc.dur = dur

	return nil
}

func init() {
	Register("memory", NewMemoryCache)
}

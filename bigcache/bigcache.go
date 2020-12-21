package bigcache

import (
	"encoding/json"
	"errors"
	"github.com/allegro/bigcache"
	"github.com/navi-tt/cache"
	"strings"
	"time"
)

var (
	invalidValueType = errors.New("val only support []byte")
)

// Cache Memcache adapter.
type BigCache struct {
	cache *bigcache.BigCache
}

// NewMemCache create new memcache adapter.
func NewBigCache() cache.Cache {
	return &BigCache{}
}

// Get get value from bigcache.
func (b *BigCache) Get(key string) (interface{}, error) {
	item, err := b.cache.Get(key)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// Set set value to bigcache.
func (b *BigCache) Set(key string, val interface{}, expireTimeout int) error {

	v, ok := val.([]byte)
	if !ok {
		return invalidValueType
	}

	return b.cache.Set(key, v)
}

// Delete delete value in bigcache.
func (b *BigCache) Delete(key string) error {
	return b.cache.Delete(key)
}

// IsExist check value exists in bigcache.
func (b *BigCache) IsExist(key string) bool {
	_, err := b.cache.Get(key)
	if err != nil {
		if strings.EqualFold(err.Error(), "Entry not found") {
			return false
		}
		return false
	}

	return true
}

// start bigcache adapter.
func (b *BigCache) Init(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)

	cfg := bigcache.Config{
		Shards:             1024,
		LifeWindow:         10 * time.Minute, //超时时间
		CleanWindow:        5 * time.Minute,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		Verbose:            true,
		HardMaxCacheSize:   8192,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	}

	c, err := bigcache.NewBigCache(cfg)
	if err != nil {
		return err
	}

	b.cache = c

	return nil
}

func init() {
	cache.Register(cache.BIG_CACHE, NewBigCache)
}

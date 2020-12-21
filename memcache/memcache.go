package memcache

import (
	"encoding/json"
	"errors"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/navi-tt/cache"
	"strings"
	"time"
)

var (
	invalidAddress   = errors.New("memcache addr invalid")
	invalidValueType = errors.New("val only support string and []byte")
)

// Cache Memcache adapter.
type MemCache struct {
	conn     *memcache.Client
	conninfo []string
}

// NewMemCache create new memcache adapter.
func NewMemCache() cache.Cache {
	return &MemCache{}
}

// Get get value from memcache.
func (m *MemCache) Get(key string) (interface{}, error) {
	item, err := m.conn.Get(key)
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}

// Put put value to memcache.
func (m *MemCache) Set(key string, val interface{}, expireTimeout int) error {
	item := memcache.Item{Key: key, Expiration: int32(time.Duration(expireTimeout) / time.Second)}
	if v, ok := val.([]byte); ok {
		item.Value = v
	} else if str, ok := val.(string); ok {
		item.Value = []byte(str)
	} else {
		return invalidValueType
	}
	return m.conn.Set(&item)
}

// Delete delete value in memcache.
func (m *MemCache) Delete(key string) error {
	return m.conn.Delete(key)
}

// IsExist check value exists in memcache.
func (m *MemCache) IsExist(key string) bool {
	_, err := m.conn.Get(key)
	return !(err != nil)
}

// start memcache adapter.
func (m *MemCache) Init(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["addr"]; !ok {
		return invalidAddress
	}
	m.conninfo = strings.Split(cf["addr"], "-")
	if m.conn == nil {
		m.conn = memcache.New(m.conninfo...)
	}

	return nil
}

func init() {
	cache.Register(cache.MEM_CACHE, NewMemCache)
}

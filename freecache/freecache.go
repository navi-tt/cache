package freecache

import (
	"fmt"
	fc "github.com/coocood/freecache"
	"github.com/navi-tt/cache"
)

var defaultSize = 0

type FreeCache struct {
	size  int
	cache *fc.Cache
}

func NewFreeCache() cache.Cache {
	return &FreeCache{
		size: defaultSize,
	}
}

func (f *FreeCache) Init(conf interface{}) error {
	freeCacheCfg, ok := conf.(cache.FreeCacheConf)
	if !ok {
		return cache.InvalidConfig
	}

	f.size = freeCacheCfg.Size
	f.cache = fc.NewCache(freeCacheCfg.Size)

	return nil
}

func (f *FreeCache) Get(key string) (interface{}, error) {
	fmt.Printf("[freecache Get key:{%s}]\n", key)

	if len(key) == 0 {
		return nil, cache.InvalidKey
	}

	val, err := f.cache.Get([]byte(key))
	if err != nil {
		return val, err
	}

	return val, nil
}

func (f *FreeCache) Set(key string, val interface{}, expireTime int) error {
	fmt.Printf("[freecache set key:{%s} value:{%v} timeout:{%d}]\n", key, val, expireTime)

	if len(key) == 0 {
		return cache.InvalidKey
	}

	v := fmt.Sprintf("%v", val)

	//v, ok := val.([]byte)
	//if !ok {
	//	return InvalidValue
	//}

	err := f.cache.Set([]byte(key), []byte(v), expireTime)
	if err != nil {
		return err
	}

	return nil
}

func (f *FreeCache) Delete(key string) error {
	fmt.Printf("[freecache del key:{%s}]\n", key)

	if len(key) == 0 {
		return cache.InvalidKey
	}

	affected := f.cache.Del([]byte(key))
	if !affected {
		return cache.DelFail
	}

	return nil
}

func (f *FreeCache) IsExist(key string) bool {
	_, err := f.Get(key)
	if err != nil {
		if err == cache.EntryNotFound {
			return false
		} else {
			return false
		}
	}
	return true
}

func init() {
	fmt.Println("fk")
	cache.Register(cache.FREE_CACHE, NewFreeCache)
}

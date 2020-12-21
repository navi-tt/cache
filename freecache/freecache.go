package freecache

import (
	"encoding/json"
	"fmt"
	fc "github.com/coocood/freecache"
	"github.com/navi-tt/cache"
	"strings"
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

func (f *FreeCache) Init(conf string) error {
	var cf map[string]int
	json.Unmarshal([]byte(conf), &cf)
	if _, ok := cf["size"]; !ok {
		cf = make(map[string]int)
		cf["size"] = defaultSize
	}

	f.size = cf["size"]
	f.cache = fc.NewCache(cf["size"])

	return nil
}

func (f *FreeCache) Get(key string) (interface{}, error) {
	fmt.Printf("[freecache Get key:{%s}]\n", key)

	if strings.EqualFold(key, "") {
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

	if strings.EqualFold(key, "") {
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

	if strings.EqualFold(key, "") {
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
		if strings.EqualFold(err.Error(), "Entry not found") {
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

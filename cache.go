package cache

import (
	"fmt"
)

type Cache interface {

	// cache initialize
	Init(config interface{}) error

	// get cached value by key.
	Get(key string) (interface{}, error)

	// set cached value with key and expire time.
	Set(key string, val interface{}, expireTime int) error

	// delete cached value by key.
	Delete(key string) error

	// check key is existed
	IsExist(key string) bool
}

// Instance is a function create a new Cache Instance
type Instance func() Cache

var adapters = make(map[string]Instance)

// Register makes a cache adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("cache: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("cache: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// NewCache Create a new cache driver by adapter name and config string.
// config need to be correct JSON as string: {"interval":360}.
// it will start gc automatically.
func NewCache(adapterName string, config interface{}) (adapter Cache, err error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("cache: unknown adapter name %q (forgot to import?)", adapterName)
		return
	}
	adapter = instanceFunc()
	err = adapter.Init(config)
	if err != nil {
		adapter = nil
	}
	return
}

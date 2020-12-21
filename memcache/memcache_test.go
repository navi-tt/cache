package memcache

import (
	"fmt"
	_ "github.com/bradfitz/gomemcache/memcache"
	"github.com/navi-tt/cache"
	"testing"
)

var c cache.Cache

func TestMain(m *testing.M) {
	ca, err := cache.NewCache(cache.MEM_CACHE, `{"addr":"127.0.0.1:11211"}`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	c = ca

	m.Run()
}

func TestSet(t *testing.T) {
	err := c.Set("navi", []byte("test"), 10)
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Logf("set success")
}

func TestGet(t *testing.T) {
	val, err := c.Get("navi")
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Logf("get success : %s", string(val.([]byte)))
}

func TestDelete(t *testing.T) {
	err := c.Delete("navi")
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Logf("del success")
}

func TestIsExist(t *testing.T) {
	exist := c.IsExist("navi")

	t.Log(exist)
}

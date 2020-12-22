package freecache

import (
	"fmt"
	"github.com/navi-tt/cache"
	"testing"
)

var c cache.Cache

func TestMain(m *testing.M) {
	ca, err := cache.NewCache(cache.FREE_CACHE, cache.FreeCacheConf{Size: 100})
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

	t.Logf("get success : %s", val.([]byte))
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

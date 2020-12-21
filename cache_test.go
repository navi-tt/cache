package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	bm, err := NewCache(MEMORY_CACHE, `{"interval":20}`)
	if err != nil {
		t.Fatal("init err")
	}

	if err = bm.Set("navi", 1, 5); err != nil {
		t.Error("set Error", err)
	}

	v, err := bm.Get("navi")
	if err != nil {
		t.Fatal("get err")
	}

	t.Logf("%d", v.(int))

	time.Sleep(5 * time.Second)

	v, err = bm.Get("navi")
	if err != nil {
		t.Fatal("get err")
	}

	bm.Delete("navi")
	if bm.IsExist("navi") {
		t.Fatal("delete err")
	}

}

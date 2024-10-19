package geecache

import (
	"fmt"
	"log"
	"testing"
	"time"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGeeCacheGet(t *testing.T) {
	loadCount := make(map[string]int, len(db))
	gee := NewGroup("socre", 2<<10, GetterFunc(func(key string) ([]byte, error) {
		log.Println("[slow DB] search key ", key)
		if v, ok := db[key]; ok {
			// 到本地获取其余应用场景的函数中获取数据
			if _, ok := loadCount[key]; !ok {
				loadCount[key] = 0
			}
			loadCount[key] += 1
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))

	now := time.Now()
	for k, v := range db {
		if view, err := gee.Get(k); err != nil || v != view.String() {
			t.Fatal("failed to get value of Tom")
		}
		if _, err := gee.Get(k); err != nil || loadCount[k] > 1 {
			t.Fatalf("cache key miss %s", k)
		}
	}

	fmt.Println("do again", time.Since(now))

	now = time.Now()
	for k, v := range db {
		if view, err := gee.Get(k); err != nil || v != view.String() {
			t.Fatal("failed to get value of Tom")
		}
		if _, err := gee.Get(k); err != nil || loadCount[k] > 1 {
			t.Fatalf("cache key miss %s", k)
		}
	}
	fmt.Println(time.Since(now))

	if view, err := gee.Get("unkonwn"); err == nil {
		t.Fatalf("unkonwn should be empty, but %s got", view)
	}
}

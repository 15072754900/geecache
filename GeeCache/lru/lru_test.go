package lru

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	now := time.Now()
	lru := New(int64(0), nil)
	lru.Add("hufeng1", String("1234"))
	fmt.Println(time.Since(now))
	if v, ok := lru.Get("hufeng1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("error on get function")
	} else {
		fmt.Println(string(v.(String)))
	}
	if _, ok := lru.Get("hufeng2"); ok {
		t.Fatalf("cache missed key failed")
	}
}

func BenchmarkGet(b *testing.B) {
	lru := New(int64(0), nil)
	lru.Add("hufeng1", String("1234"))
	for i := 0; i < b.N; i++ {
		if v, ok := lru.Get("hufeng1"); !ok || string(v.(String)) != "1234" {
			b.Fatalf("error on get function")
		} else {
			fmt.Println(string(v.(String)))
		}
	}
}

// func TestRemoveOldList(t *testing.T) {
// 	// 进行删除操作的模拟，以及超出设定可用空间的测试
// 	key1 := "hufeng1"
// 	key2 := "hufeng2"
// 	key3 := "hufeng3"
// 	cap := len(key1 + key2 + key3)
// 	value1, value2, value3 := "value", "value2", "value3"
// 	lru := New(int64(cap), nil)
// 	lru.Add(key1, String(value1))
// 	lru.Add(key2, String(value2))
// 	lru.Add(key3, String(value3))

// 	if _, ok := lru.Get(key1); ok || lru.Len() != 2 {
// 		t.Fatalf("error on get")
// 	}
// }

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	// 建立hooks函数
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}

	lru := New(int64(30), callback)
	// lru.Add("key1", String("value1"))
	// lru.Add("key2", String("value2"))
	// lru.Add("key3", String("value3"))

	// expect := []string{"key1", "key2"}
	// if !reflect.DeepEqual(expect, keys) {
	// 	t.Fatalf("error on hooks")
	// } else {
	// 	fmt.Println("right")
	// }
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expect := []string{"key1", "k2", "k3", "k4"}
	fmt.Println(keys)

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}

func TestOnEvicted1(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expect := []string{"key1", "k2"}

	// if !reflect.DeepEqual(expect, keys) {
	// 	t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	// }
	fmt.Println(keys, expect)
}

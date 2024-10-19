package geecache

import (
	"reflect"
	"testing"
)

func TestGetterFunc(t *testing.T) {
	var f Getter = GetterFunc(func(s string) ([]byte, error) {
		return []byte(s), nil
	})
	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(expect, v) {
		t.Errorf("callback failed")
	}
}

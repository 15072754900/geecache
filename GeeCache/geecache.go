package geecache

import (
	"fmt"
	"log"
	"sync"
)

// 获取源数据的内容

type Getter interface {
	Get(key string) ([]byte, error)
}

// 接口型函数

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// 上面是获取不到数据的时候使用函数获取，下面是获取的流程

type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// 进行三件事：1.判断是否符合建立条件；2.加上读写锁；3.初始化并给全局map一个值

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(s string) *Group {
	mu.RLock()
	g := groups[s]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key if required")
	}
	if v, ok := g.mainCache.Get(key); ok {
		log.Println("cache hit")
		return v, nil
	}
	// 一个核心步骤：如果不能再缓存中查找到信息，就需要在本地或者其他地方调用回调函数获取信息
	return g.load(key)
}

// 不同场景这个中间层的实现不一样，当处于本地时使用该方式，当处于分布式情况调用getFromPeer调取其他节点中的数据

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}

	// 从其余地方获取之后还会建立一个新的缓存，这一次建立之后其余请求从缓存中获取数据而不会像之前一样重新获取并建立一个缓存，建立缓存的过程是并发安全的
	g.populateKey(key, value)

	return value, nil
}

func (g *Group) populateKey(key string, value ByteView) {
	g.mainCache.Add(key, value)
}

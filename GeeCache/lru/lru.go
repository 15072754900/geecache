package lru

import (
	"container/list"
)

// 定义缓存结构，包含字典和双向链表，用于实现后续的增删改查，字典用于键值对管理（数据管理），双向链表用于fifo、lfu、lru的实现
type Cache struct {
	maxBytes  int64                         // 可使用的最大内存
	nbytes    int64                         // 已使用容量
	ll        *list.List                    // 一个双向链表，内部包含一个计数器和一个元素类型，还有很多面向该对象的方法
	cache     map[string]*list.Element      // 缓存队列
	OnEvicted func(key string, value Value) // hook函数，作为一个可以选择进行处理的函数
}

type entry struct { // 双向链表中的数据结构
	key   string
	value Value // 是要实现了len()方法的数据结构就可以是这个的value
}

type Value interface {
	Len() int
}

// 使用懒汉式单例模式进行数据实例化

func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 一个数据需要进行增删改查

// 查

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		// 当存在的时候，先获取数据，然后根据lru算法，将获取的数据放到指定位置
		c.ll.MoveToFront(c.cache[key])
		// 然后进行断言，由于本身是一个接口，数据是实现了该接口的值，读取数据使用断言转换为其他的数据结构进行解耦。同时可以使用类型开关和反射进行类型判断和后续的数据处理。
		kv := ele.Value.(*entry)
		return kv.value, ok
	}
	return
}

// 删

func (c *Cache) RemoveLodest() {
	ele := c.ll.Back()
	if ele != nil {
		// 删除队列中的、删除map中的
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		// 还要修改使用的空间大小
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// 增、改。对数据进行查询、判断存在则修改，否则新建

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// 说明存在于缓存中，进行修改，将数据移动到队列的尾部
		c.ll.MoveToFront(ele)
		// 后面还是一贯的修改数据
		kv := ele.Value.(*entry)
		c.nbytes += int64(len(kv.key)) - int64(kv.value.Len())
		kv.value = value
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	} else {
		// 不存在，则添加
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) - int64(value.Len())
		kv := ele.Value.(*entry)
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
	// 还要判断是否具备可以修改的条件
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveLodest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}

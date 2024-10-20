package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int
	keys     []int
	hashMap  map[int]string
}

// 懒汉式单例模式

func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
		// 这里的map类型数据一定要使用make初始化了，不然无法进行数据存储和获取
	}
	// 建立一个默认的fn
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Map是一致性哈希算法的主数据结构，包含4个成员变量：hash函数，虚拟节点倍数，哈希环，映射表

// 添加真实节点

func (m *Map) Add(keys ...string) {
	// 向里边添加很多个节点进行处理
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// 定义存储的哈希环个体信息
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
			// 一个节点key对应多个虚拟节点，hash是计算出来的虚拟节点，将其对应到hashmap上并加入到哈希环上
		}
		sort.Ints(m.keys)
	}
}

func (m *Map) Get(key string) string {
	// 计算key的哈希值，获取对应虚拟节点的真实节点，返回真实节点
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}

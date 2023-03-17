package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash 将bytes映射到uint32
type Hash func(data []byte) uint32

// Map 包含所有 hashed keys
type Map struct {
	hash     Hash
	replicas int
	//keys 为哈希环
	keys []int // 已经排序
	//虚拟节点和真实节点的映射
	hashMap map[int]string
}

// New 新建一个Map实例
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// AddPeers 添加一些keys-->为IP+port到一致性哈希。
func (m *Map) AddPeers(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// GetPeer 获取哈希中与提供的key最近的项(节点)。
func (m *Map) GetPeer(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	//二分查找合适的副本(index为第一个下标)。
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}

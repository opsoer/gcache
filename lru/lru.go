package lru

import (
	"container/list"
	"time"
)

var TimeMap map[string]time.Time

// Cache Cache为LRU缓存,并发访问是不安全的。
type Cache struct {
	//缓存热点数据
	old lruList
	//缓存新数据
	young lruList
}

type lruList struct {
	//最大容量
	maxCap int
	//已经使用了的容量
	usedMem int
	ll      *list.List
	mp      map[string]*list.Element
}

type entry struct {
	key   string
	value Value
}

// Value 使用Len来计算需要多少字节,lru缓存的value是Value接口类型。
type Value interface {
	Len() int
}

// New 创造一个缓存。
func New(maxCap int) *Cache {
	return &Cache{
		old:   lruList{maxCap: maxCap / 8 * 5, ll: list.New(), mp: map[string]*list.Element{}},
		young: lruList{maxCap: maxCap / 8 * 3, ll: list.New(), mp: map[string]*list.Element{}},
	}
}

// Add 添加一个值lruList。
func (lru *lruList) add(key string, value Value) {
	if ele, ok := lru.mp[key]; ok {
		lru.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		lru.usedMem += value.Len() - kv.value.Len()
		kv.value = value
	} else {
		ele := lru.ll.PushFront(&entry{key, value})
		lru.mp[key] = ele
		lru.usedMem += len(key) + value.Len()
	}
	for lru.maxCap != 0 && lru.maxCap < lru.usedMem {
		lru.removeTail()
	}
}

// removeTail 根据lru的规则删除lruList一个k-v。
func (lru *lruList) removeTail() {
	ele := lru.ll.Back()
	if ele != nil {
		lru.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(lru.mp, kv.key)
		lru.usedMem -= len(kv.key) + kv.value.Len()
	}
}

// Len 返回lruList k-v的个数。
func (lru *lruList) Len() int {
	return lru.ll.Len()
}

// Delete 删除lruList 里面某key
func (lru *lruList) delete(key string) bool {
	var ele *list.Element
	if ele = lru.mp[key]; ele == nil {
		return false
	}
	val := ele.Value.(*entry)
	lru.ll.Remove(ele)
	delete(lru.mp, key)
	lru.usedMem -= len(val.key) + val.value.Len()
	return true
}

func (c *Cache) Add(key string, value Value) {
	if _, ok := c.old.mp[key]; ok {
		c.old.add(key, value)
		return
	}
	if ele, ok := c.young.mp[key]; ok {

		addTime := TimeMap[key]
		if ok := addTime.Add(time.Second).Before(time.Now()); ok {
			//如果加入youngList 一秒钟之后又被访问，就加入olsList
			//加入oldList
			val := ele.Value.(*entry)
			c.old.add(key, value)
			//删除youngList
			c.young.delete(val.key)
			delete(TimeMap, key)
			return
		} else {
			//如果加入youngList 一秒钟之内又被访问，只在youngList 变化位置，加入youngList 头部
			c.young.add(key, value)
			c.young.ll.MoveToFront(ele)
			return
		}
	}

	TimeMap[key] = time.Now()
	c.young.add(key, value)

}

// Get 查找的key的值。
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.old.mp[key]; ok {
		val := ele.Value.(*entry)
		c.old.ll.MoveToFront(ele)
		return val.value, true
	}

	if ele, ok := c.young.mp[key]; ok {
		addTime := TimeMap[key]
		val := ele.Value.(*entry)
		if ok := addTime.Add(time.Second).Before(time.Now()); ok {
			//如果加入youngList 一秒钟之后又被访问，就加入olsList
			//加入oldList
			c.old.add(val.key, val.value)
			//删除youngList
			c.young.delete(val.key)
			delete(TimeMap, key)
		}
		return val.value, true
	}

	return nil, false
}

func (c *Cache) Delete(key string) bool {
	if _, ok := c.old.mp[key]; ok {
		return c.old.delete(key)
	}

	if _, ok := c.young.mp[key]; ok {
		delete(TimeMap, key)
		return c.young.delete(key)
	}

	return false
}

// Len 返回缓存k-v的个数。
func (c *Cache) Len() int {
	return c.young.Len() + c.old.Len()
}

func (c *Cache) Clear() {
	c.old = lruList{
		ll: list.New(),
		mp: map[string]*list.Element{},
	}
	c.young = lruList{
		ll: list.New(),
		mp: map[string]*list.Element{},
	}
	TimeMap = make(map[string]time.Time)
}

func init() {
	TimeMap = make(map[string]time.Time)
}

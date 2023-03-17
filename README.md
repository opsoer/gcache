# gcache
简介：一款基于LRU缓存淘汰策略和一致性哈希算法进行负载均衡的分布式缓存框架，可用于单机版缓存，也可以用于分布式版，通过HTTP协议进行通信。

**特性：**

- 单机缓存和基于HTTP的分布式缓存。
- 参考MySQL Buffer Pool，独立实现old、young两个lru链表防止缓存污染。
- 使用Go的锁和一秒钟的map缓存防止缓存击穿。
- 使用一致性哈希算法选择节点，实现负载均衡。

### API

```go
func NewCache(maxCap int, getter Getter) *GCache
```

```go
func (c *GCache) Get(key string) (ByteView, error)
```

```go
func (c *GCache) Delete(key string) bool
```

```go
func (c *GCache) RegisterHTTPPool(peers PeerPicker)
```

```go
func NewHTTPPool(self string) *HTTPPool
```

```go
func (p *HTTPPool) AddPeers(peers ...string)
```



### simple demo

```go
package main

import (
   "fmt"
   "gcache"
   "log"
)

var db = map[string]string{
   "a": "aa",
   "b": "bb",
   "c": "cc",
   "d": "dd",
   "e": "ee",
   "f": "ff",
}

func simple() {
   gc := gcache.NewCache(1<<5, gcache.GetterFunc(
      func(key string) ([]byte, error) {
         log.Println("[SlowDB] search key", key)
         if v, ok := db[key]; ok {
            return []byte(v), nil
         }
         return nil, fmt.Errorf("%s not exist", key)
      }))
   selfUrl := gcache.NewHTTPPool("127.0.0.1:8081")
   selfUrl.AddPeers("127.0.0.1:8081")
   val, err := gc.Get("a")
   if err != nil {
      fmt.Println(err)
   }
   fmt.Printf("key %s get value %s\n", "a", val.String())
   val, err = gc.Get("a")
   if err != nil {
      fmt.Println(err)
   }
   fmt.Printf("key %s get value %s\n", "a", val.String())
   fmt.Println(gc.Delete("a"))
   gc.Get("a")
}
```

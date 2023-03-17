package main

import (
	"flag"
	"fmt"
	"gcache"
	"log"
	"net/http"
)

func creatCache(cacheCap int, getter gcache.Getter) *gcache.GCache {
	if getter == nil {
		return gcache.NewCache(cacheCap, gcache.GetterFunc(
			func(key string) ([]byte, error) {
				log.Println("[SlowDB] search key", key)
				if v, ok := db[key]; ok {
					return []byte(v), nil
				}
				return nil, fmt.Errorf("%s not exist", key)
			}))
	}
	return gcache.NewCache(cacheCap, getter)
}

func startCacheServer(addr string, addrs []string, c *gcache.GCache) {
	httpPool := gcache.NewHTTPPool(addr)
	httpPool.AddPeers(addrs...)
	c.RegisterHTTPPool(httpPool)
	log.Println("gcache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], httpPool))
}

func startAPIServer(apiAddr string, cache *gcache.GCache) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := cache.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("api server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

func distributed() {
	var port int
	flag.IntVar(&port, "port", 8081, "StoneCache server port")
	flag.Parse()

	apiAddr := "http://localhost:8084"
	addrMap := map[int]string{
		8081: "http://localhost:8081",
		8082: "http://localhost:8082",
		8083: "http://localhost:8083",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	cache := creatCache(1<<5, nil)
	go startAPIServer(apiAddr, cache)
	startCacheServer(addrMap[port], addrs, cache)
}

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
	gc.Peers = selfUrl
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

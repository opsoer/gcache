package lru

import (
	"fmt"
	"testing"
	"time"
)

type myValue string

func (v myValue) Len() int {
	return len(v)
}

func TestCase1(t *testing.T) {
	cache := New(80)
	if cache.young.maxCap != 30 || cache.old.maxCap != 50 {
		t.Error("old、young 内存分配错误")
	}
	cache.Add("hello", myValue("world"))
	val, ok := cache.Get("hello")
	v := val.(myValue)
	if !ok || v != "world" {
		t.Error("获得的 kv 不是预期的")
	}

	cache.Add("hello", myValue("WORLD"))
	val, ok = cache.Get("hello")
	v = val.(myValue)
	if !ok || v != "WORLD" {
		t.Error("获得的 kv 不是预期的")
	}

	if cache.young.Len() != 1 || cache.young.usedMem != 10 {
		t.Error("kv 没有分配到 young list")
	}
	if cache.old.Len() != 0 || cache.old.usedMem != 0 {
		t.Error("old list 不应该分配kv")
	}
	if cache.Len() != 1 {
		t.Error("cache.Len err")
	}
}

func TestCase2(t *testing.T) {
	cache := New(80)
	cache.Add("hello", myValue("world"))
	//sleep second
	time.Sleep(time.Second)
	val, ok := cache.Get("hello")
	v := val.(myValue)
	if !ok || v != "world" {
		t.Error("获得的 kv 不是预期的")
	}

	cache.Add("hello", myValue("WORLD"))
	val, ok = cache.Get("hello")
	v = val.(myValue)
	if !ok || v != "WORLD" {
		t.Error("获得的 kv 不是预期的")
	}

	//order case
	cache.Add("HELLO", myValue("world"))
	//sleep second
	time.Sleep(time.Second)

	cache.Add("HELLO", myValue("WORLD"))
	val, ok = cache.Get("HELLO")
	v = val.(myValue)
	if !ok || v != "WORLD" {
		t.Error("获得的 kv 不是预期的")
	}

	if cache.old.Len() != 2 || cache.old.usedMem != 20 {
		t.Error("kc 没有分配到 old list")
	}

	if cache.young.Len() != 0 || cache.young.usedMem != 0 {
		t.Error("young list 不应该分配kv")
	}
	if len(TimeMap) != 0 {
		t.Error("TimeMap 没有删除")
	}
}

func TestCase3(t *testing.T) {
	cache := New(80)

	for i := 0; i <= 3; i++ {
		cache.Add(fmt.Sprintf("hell%d", i), myValue(fmt.Sprintf("worl%d", i)))
	}

	for i := 1; i <= 3; i++ {
		val, ok := cache.Get(fmt.Sprintf("hell%d", i))
		v := val.(myValue)
		target := fmt.Sprintf("worl%d", i)
		if !ok || myValue(target) != v {
			t.Error("获得的 kv 不是预期的")
		}
	}
	if cache.young.Len() != 3 || cache.young.usedMem != 30 {
		t.Error("kv 没有全部分配到 young list")
	}
	if cache.old.Len() != 0 || cache.old.usedMem != 0 {
		t.Error("old list 不应该分配kv")
	}
}

func TestCase4(t *testing.T) {
	cache := New(80)

	for i := 0; i <= 3; i++ {
		cache.Add(fmt.Sprintf("hell%d", i), myValue(fmt.Sprintf("worl%d", i)))
	}

	time.Sleep(time.Second)

	for i := 1; i <= 3; i++ {
		val, ok := cache.Get(fmt.Sprintf("hell%d", i))
		v := val.(myValue)
		target := fmt.Sprintf("worl%d", i)
		if !ok || myValue(target) != v {
			t.Error("获得的 kv 不是预期的")
		}
	}
	if cache.young.Len() != 0 || cache.young.usedMem != 0 {
		t.Error("kv 应该迁移到 old list 了")
	}
	if cache.old.Len() != 3 || cache.old.usedMem != 30 {
		t.Error("kv 应该迁移到 old list 了")
	}
}

func TestCase5(t *testing.T) {
	cache := New(80)
	cache.Add("hell1", myValue("worl1"))
	cache.Add("hell2", myValue("worl2"))
	cache.Add("hell3", myValue("worl3"))
	cache.Delete("hell1")
	if cache.young.Len() != 2 || cache.young.usedMem != 20 {
		t.Error("kv 没有全部分配到 young list")
	}
	if cache.old.Len() != 0 || cache.old.usedMem != 0 {
		t.Error("old list 不应该分配kv")
	}

	cache.Clear()
	if cache.young.Len() != 0 || cache.young.usedMem != 0 {
		t.Error("clear err")
	}
	if cache.old.Len() != 0 || cache.old.usedMem != 0 {
		t.Error("clear err")
	}
	if len(TimeMap) != 0 {
		t.Error("clear err")
	}
}

package singleflight

import (
	"sync"
	"time"
)

var N = 0

// call is an in-flight(正在进行) or completed Do call
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Ones struct {
	//mu保证m不会被并发读写
	mu sync.Mutex       // protects m
	m  map[string]*call //延迟初始化
}

// Do 执行并返回给定函数的结果
// 确定一个键只执行一次
// 时间。如果有一个重复调用，则重复调用者等待
// 原始的完成和收相同的结果。接
func (g *Ones) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()         //有请求在进行，等待
		return c.val, c.err //请求结束，返回结果
	}
	c := new(call)
	c.wg.Add(1)  //请求前加锁
	g.m[key] = c //添加到g.m表示key已经有请求在处理
	g.mu.Unlock()

	c.val, c.err = fn() //调用fn，发起请求
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key) //更新g.m
	g.mu.Unlock()

	return c.val, c.err //返回结果
}

func (g *Ones) DoOK(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		return c.val, c.err //请求结束，返回结果
	}
	c := new(call)
	g.m[key] = c        //添加到g.m表示key已经有请求在处理
	c.val, c.err = fn() //调用fn，发起请求
	N++
	go func() {
		time.Sleep(time.Second)
		delete(g.m, key) //更新g.m
	}()

	return c.val, c.err //返回结果
}

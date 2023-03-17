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

func (g *Ones) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
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

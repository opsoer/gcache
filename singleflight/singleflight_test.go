package singleflight

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var ExecTimes = 0

func TestDoCase1(t *testing.T) {
	ExecTimes = 0
	var g Ones
	wg := sync.WaitGroup{}
	wg.Add(10000)
	fn := func() (interface{}, error) {
		ExecTimes++
		return "", nil
	}
	for i := 0; i < 10000; i++ {
		go func() {
			defer wg.Done()
			g.Do("key", fn)
		}()
	}
	wg.Wait()
	fmt.Println(ExecTimes)
	if ExecTimes != 1 {
		t.Error("singleFlight err")
	}
}

func TestDoCase2(t *testing.T) {
	ExecTimes = 0
	var g Ones
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			i, _ := g.Do("key", func() (interface{}, error) {
				ExecTimes++
				return "bar", nil
			})
			v := i.(string)
			if v != "bar" {
				t.Errorf("Expect: %v, Actual: %v", "bar", v)
			}

			i, _ = g.Do("key1", func() (interface{}, error) {
				ExecTimes++
				return "bar", nil
			})
			v = i.(string)
			if v != "bar" {
				t.Errorf("Expect: %v, Actual: %v", "bar", v)
			}
		}()
	}
	wg.Wait()
	fmt.Println(ExecTimes)
	if ExecTimes != 2 {
		t.Error("singleFlight err")
	}
}

func TestDoCase3(t *testing.T) {
	ExecTimes = 0
	var g Ones
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			i, _ := g.Do("key", func() (interface{}, error) {
				ExecTimes++
				return "bar", nil
			})
			v := i.(string)
			if v != "bar" {
				t.Errorf("Expect: %v, Actual: %v", "bar", v)
			}
		}()
	}
	time.Sleep(time.Second)
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			i, _ := g.Do("key", func() (interface{}, error) {
				ExecTimes++
				return "bar", nil
			})
			v := i.(string)
			if v != "bar" {
				t.Errorf("Expect: %v, Actual: %v", "bar", v)
			}
		}()
	}
	wg.Wait()
	fmt.Println(ExecTimes)
	if ExecTimes != 2 {
		t.Error("singleFlight err")
	}
}

package singleflight

import (
	"fmt"
	"sync"
	"testing"
)

var TN = 0

func TestDoCase1(t *testing.T) {
	TN = 0
	var g Ones
	wg := sync.WaitGroup{}
	wg.Add(10000)
	fn := func() (interface{}, error) {
		TN++
		return "", nil
	}
	for i := 0; i < 10000; i++ {
		go func() {
			defer wg.Done()
			g.DoOK("key", fn)
		}()
	}
	wg.Wait()
	fmt.Println(TN)
	if TN != 1 {
		t.Error("singleFlight err")
	}
}

//
//func TestDoCase2(t *testing.T) {
//	TN = 0
//	var g Ones
//	wg := sync.WaitGroup{}
//	wg.Add(100)
//	for i := 0; i < 100; i++ {
//		go func() {
//			defer wg.Done()
//			i, _ := g.DoOK("key", func() (interface{}, error) {
//				TN++
//				return "bar", nil
//			})
//			fmt.Println(i)
//			g.DoOK("key1", func() (interface{}, error) {
//				TN++
//				return "", nil
//			})
//		}()
//	}
//	wg.Wait()
//	fmt.Println(TN)
//	if TN != 2 {
//		t.Error("singleFlight err")
//	}
//}
//
//func TestDoCase3(t *testing.T) {
//	TN = 0
//	var g Ones
//	wg := sync.WaitGroup{}
//	wg.Add(100)
//	for i := 0; i < 100; i++ {
//		go func() {
//			defer wg.Done()
//			g.DoOK("key", func() (interface{}, error) {
//				TN++
//				return "", nil
//			})
//		}()
//	}
//	time.Sleep(time.Second)
//	wg.Add(100)
//	for i := 0; i < 100; i++ {
//		go func() {
//			defer wg.Done()
//			g.DoOK("key", func() (interface{}, error) {
//				TN++
//				return "", nil
//			})
//		}()
//	}
//	wg.Wait()
//	fmt.Println(TN)
//	if TN != 2 {
//		t.Error("singleFlight err")
//	}
//}

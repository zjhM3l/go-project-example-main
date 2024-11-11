package Goroutine

import (
	"fmt"
	"sync"
	"time"
)

func hello(i int) {
	println("Hello from goroutine ", fmt.Sprint(i))
}

func HelloGoRoutine() {
	for i := 0; i < 10; i++ {
		go func(j int) {
			hello(j)
		}(i)
	}
	// 暴力阻塞，等待子协程执行完毕
	time.Sleep(time.Second)
}

func HelloGoRoutineWithWaitGroup() {
	// 使用sync.WaitGroup实现协程的同步
	var wg sync.WaitGroup
	// 开启五个协程
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(j int) {
			// 子协程执行完毕，通过Done通知WaitGroup
			defer wg.Done()
			hello(j)
		}(i)
	}
	// 阻塞等待所有子协程执行完毕
	wg.Wait()
}

func test() {
	// HelloGoRoutine()
	HelloGoRoutineWithWaitGroup()
}

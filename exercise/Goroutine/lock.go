package Goroutine

// 通过共享内存实现协程间的通信，需要保证并发安全

// 对变量进行2000次+1操作，5个协程并发执行，预期结果10000
import (
	"sync"
	"time"
)

var (
	x int64
	// 互斥锁，通过对临界区权限的控制，保证并发安全的共享内存
	lock sync.Mutex
)

func addWithLock() {
	for i := 0; i < 2000; i++ {
		// 在每次对x进行操作时，通过lock获取临界区的资源
		lock.Lock()
		x += 1
		// 计算完成，将临界区的权限释放
		lock.Unlock()
	}
}

func addWithoutLock() {
	for i := 0; i < 2000; i++ {
		x += 1
	}
}

func add10000() {
	x = 0
	for i := 0; i < 5; i++ {
		go addWithoutLock()
	}
	// 使用sleep实现暴力阻塞，因为不知道子协程的执行时间
	// 更优解：使用sync.WaitGroup
	time.Sleep(time.Second)
	println("Without lock: ", x)
	x = 0
	for i := 0; i < 5; i++ {
		go addWithLock()
	}
	time.Sleep(time.Second)
	println("With lock: ", x)
}

package Goroutine

// 通过通信来共享内存，而不是通过共享内存来通信

// A子协程发送0~9数字
// B子协程计算输入数字的平方
// 住协程输出最后的平方数

// 并发安全的channel
func CalcuSquare() {
	src := make(chan int)
	// 接收有缓冲，考虑消费者消费速度慢于生产者生产速度，导致生产者阻塞
	dest := make(chan int, 3)
	go func() {
		// 延迟资源关闭
		defer close(src)
		for i := 0; i < 10; i++ {
			src <- i
		}
	}()
	go func() {
		defer close(dest)
		for i := range src {
			dest <- i * i
		}
	}()
	for i := range dest {
		println(i)
	}
}

package test

import "testing"

// 基准测试select函数，命名规则：Benchmark+函数名

func BenchmarkSelect(b *testing.B) {
	// 不属于测试函数的损耗
	InitServerIndex()
	// 所以要重置计时器
	b.ResetTimer()
	// 串行执行select函数，基准测试
	for i := 0; i < b.N; i++ {
		Select()
	}
}

func BencahmarkSelectParallel(b *testing.B) {
	InitServerIndex()
	b.ResetTimer()
	// 并行执行select函数，基准测试
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Select()
		}
	})
}

// 最后是串行性能更好，因为select用到了sdk的rand函数
// rand函数为了保证全局的随机性，并发安全，所以加了锁，降低了并发性能
// 可以用fastrand包替代rand包，提高并发性能

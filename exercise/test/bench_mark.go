package test

// 基准测试的例子，随机选择执行服务器

import (
	"math/rand"

	"github.com/bytedance/gopkg/lang/fastrand"
)

var ServerIndex [10]int

func InitServerIndex() {
	for i := 0; i < 10; i++ {
		ServerIndex[i] = i + 100
	}
}

func Select() int {
	return ServerIndex[rand.Intn(10)]
}

// fastrand包替代rand包，提高并发性能
func FastSelect() int {
	return ServerIndex[fastrand.Uint32n(10)]
}

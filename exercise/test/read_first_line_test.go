package test

import (
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestProcessFirstLine(t *testing.T) {
	firstLine := ProcessFirstLine()
	assert.Equal(t, "Hello, World!", firstLine)
}

// mock测试，不对file强依赖
func TestProcessFirstLineWithMock(t *testing.T) {
	// 打桩
	monkey.Patch(ReadFirstLine, func() string {
		return "line110"
	})
	// 卸载打桩
	defer monkey.Unpatch(ReadFirstLine)
	line := ProcessFirstLine()
	assert.Equal(t, "line000", line)
}

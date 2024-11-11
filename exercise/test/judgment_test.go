package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test judgment_test.go judgment.go --cover
// 66.7%因为测试只执行了10,11行，而没有执行第13行
func TestJudgePassLineTrue(t *testing.T) {
	isPass := JudgePassLine(70)
	assert.Equal(t, true, isPass)
}

// 100%因为测试执行了所有的代码分支
func TestJudgePassLineFalse(t *testing.T) {
	isPass := JudgePassLine(50)
	assert.Equal(t, false, isPass)
}

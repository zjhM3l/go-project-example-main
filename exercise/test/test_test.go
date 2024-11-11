package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func HelloTom() string {
	return "Jerry"
}

func TestHelloTom(t *testing.T) {
	output := HelloTom()
	expectOutput := "Tom"
	if output != expectOutput {
		t.Errorf("Output: %s, ExpectOutput: %s", output, expectOutput)
	}
}

// 单元测试-assert
func TestHelloTomWithAssert(t *testing.T) {
	output := HelloTom()
	expectOutput := "Tom"
	assert.Equal(t, expectOutput, output)
}

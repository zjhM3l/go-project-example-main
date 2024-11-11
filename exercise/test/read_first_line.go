package test

import (
	"bufio"
	"os"
	"strings"
)

// 问题：一旦文件被人为删除，测试就会失败
// 应该mock测试
func ReadFirstLine() string {
	open, err := os.Open("test.txt")
	defer open.Close()
	if err != nil {
		return ""
	}
	scanner := bufio.NewScanner(open)
	for scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

func ProcessFirstLine() string {
	line := ReadFirstLine()
	destLine := strings.ReplaceAll(line, "11", "00")
	return destLine
}

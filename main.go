//go:build !windows
// +build !windows

package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	fmt.Printf("错误: 此工具仅支持Windows平台\n")
	fmt.Printf("当前平台: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("请在Windows系统上运行此工具\n")
	os.Exit(1)
}
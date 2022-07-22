package test

import (
	"fmt"
	"testing"
)

func errorHandle() {
	defer func() {
		fmt.Println("defer panic")
		if err := recover(); err != nil {
			fmt.Println("recover success")
		}
	}()
	arr := make([]int, 3)
	//数组越界 直接panic
	fmt.Print(arr[3])
	fmt.Println("after recover")
}

func TestErrorHandle(t *testing.T) {
	errorHandle()
	fmt.Println("after defer")
}

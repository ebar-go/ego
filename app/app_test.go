package app

import (
	"fmt"
	"testing"
)

type Obj struct {
	Number int
}

func NewObj() (*Obj, error) {
	return &Obj{Number: 1}, fmt.Errorf("xxx")
}

func TestProvide(t *testing.T) {
	err := Container().Provide(NewObj)
	fmt.Println("err:", err)

	// 当构造函数NewObj返回err时，provide不会立刻报错，而是再invoke得时候报错
	err1 := Container().Invoke(func(obj *Obj) {
		fmt.Println(obj.Number)
	})
	fmt.Println("err1:", err1)
}

package library

import (
	"fmt"
	"testing"
)

type conversionTest struct {
	Name  string
	State int
	Node  []conversionTest2
}

type conversionTest2 struct {
	Name  string
	State int
}

func TestStructToMap(t *testing.T) {

	obj2 := []conversionTest2{
		{
			Name:  "22",
			State: 0,
		},
		{
			Name:  "33",
			State: 1,
		},
	}
	obj := conversionTest{
		Name:  "aa",
		State: 0,
		Node:  obj2,
	}
	res := Struct2Map(obj)
	fmt.Println(res)
}

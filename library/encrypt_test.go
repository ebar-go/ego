package library

import (
	"fmt"
	"testing"

	"github.com/ebar-go/ego/test"
)

func getMd5StringDataProvider() map[string]string {
	items := make(map[string]string)
	items["123456"] = "e10adc3949ba59abbe56e057f20f883e"
	return items
}

// 测试获取Md5
func TestGetMd5String(t *testing.T) {
	items := getMd5StringDataProvider()
	for key, value := range items {
		test.AssertEqual(t, GetMd5String(key), value)
	}
}

func TestUniqueId(t *testing.T) {
	fmt.Println(UniqueId())
}

func TestGetHash(t *testing.T) {
	fmt.Println(GetHash("123456"))
}

func TestPasswordHash(t *testing.T) {
	combin, err := PasswordHash("123456")
	if err != nil {
		t.Error("PasswordHash 错误：", err)
	}
	fmt.Println(string(combin))
	res := PasswordVerify(combin, "123456")
	if res == false {
		t.Error("PasswordHash 验证失败")
	}
}

func TestRandString(t *testing.T) {
	str, err := RandString(16)
	if err != nil {
		t.Error(err)
	}
	if len(str) != 16 {
		t.Error("长度错误:", len(str))
	}
	fmt.Println(str)
}

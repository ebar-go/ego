package library

import (
	"testing"
	"fmt"
	"github.com/magiconair/properties/assert"
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
		assert.Equal(t, GetMd5String(key), value, "加密错误")
	}
}

func TestUniqueId(t *testing.T) {
	fmt.Println(UniqueId())
}

func TestGetHash(t *testing.T) {
	fmt.Println(GetHash("123456"))
}

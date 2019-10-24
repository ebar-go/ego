package library

import (
	"testing"
	"git.epetbar.com/go-package/ego/test"
	"fmt"
)

func TestGetLocalIp(t *testing.T) {
	ip, err := GetLocalIp()
	test.AssertNil(t, err)
	fmt.Println(ip)
}

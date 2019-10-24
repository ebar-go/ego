package library

import (
	"testing"
	"github.com/ebar-go/ego/test"
	"fmt"
)

func TestGetLocalIp(t *testing.T) {
	ip, err := GetLocalIp()
	test.AssertNil(t, err)
	fmt.Println(ip)
}

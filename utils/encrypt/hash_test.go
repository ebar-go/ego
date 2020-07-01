package encrypt

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	fmt.Println(Md5("test"))
	fmt.Println(Sha1("test"))
}
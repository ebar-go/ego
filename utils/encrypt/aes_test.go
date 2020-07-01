package encrypt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAes(t *testing.T) {
	encipher := Aes("hgfedcba87654321")
	source := []byte("test")
	encode, err := encipher.Encrypt(source)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(encode)
	decode, err := encipher.Decrypt(encode)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(decode))
	assert.Equal(t, source, decode)
}

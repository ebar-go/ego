package component

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPost(t *testing.T) {
	res, err := NewCurl().Post("http://pf.mayiyahei.net/app/login/track_list", nil)
	assert.Nil(t, err)
	fmt.Println(string(res.Bytes()))
}

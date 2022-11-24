package curl

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPost(t *testing.T) {
	res, err := New().Post(context.TODO(), "http://pf.mayiyahei.net/app/login/track_list", nil)
	assert.Nil(t, err)
	fmt.Println(string(res.Bytes()))
}

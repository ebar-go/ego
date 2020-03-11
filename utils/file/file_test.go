package file

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFile(t *testing.T)  {
	fmt.Println(GetCurrentDir())
	fmt.Println(GetCurrentPath())
	fmt.Println(GetExecuteDir())
	assert.True(t, Exist(GetCurrentPath()))
	assert.Nil(t, Mkdir("/tmp/test", true))
}
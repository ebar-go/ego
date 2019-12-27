package validator

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidator_ValidateStruct(t *testing.T) {
	vd := new(Validator)

	type User struct {
		Name string `json:"phone" binding:"required,omitempty" comment:"名称"`
		Age  uint   `json:"age" binding:"required,min=10" comment:"年龄"`
	}

	var user = User{Name: "test", Age: 9}

	err := vd.ValidateStruct(user)
	fmt.Println(err)
	assert.NotNil(t, err)

}

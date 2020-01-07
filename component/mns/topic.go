package mns

import (
	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
)

// topic
type Topic struct {
	Name     string
	Instance ali_mns.AliMNSTopic
}

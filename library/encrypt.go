package library

import (
	"encoding/hex"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/ebar-go/ego/http/constant"
)

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// GetHash 获取hash加密内容
func GetHash(s string) string {
	Sha1Inst := sha1.New()
	Sha1Inst.Write([]byte(s))
	result := Sha1Inst.Sum([]byte(""))
	return fmt.Sprintf("%x", result)
}

//UniqueId 生成Guid字串
func UniqueId() string {
	return uuid.NewV4().String()
}

// GetTraceId 获取全局ID
func GetTraceId() string {
	return constant.TraceIdPrefix + UniqueId()
}

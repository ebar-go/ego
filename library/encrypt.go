package library

import (
	"encoding/hex"
	"io"
	"crypto/rand"
	"encoding/base64"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"github.com/satori/go.uuid"
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

//生成Guid字串
func UniqueId() string {
	uuidInstance, err := uuid.NewV4()
	if err == nil {
		return uuidInstance.String()
	}

	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}

	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}


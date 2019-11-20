package library

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/ebar-go/ego/http/constant"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"hash/crc32"
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

// PasswordHash 通过bcrypt算法生成hash
func PasswordHash(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashBytes), err
}

// PasswordVerify 通过bcrypt算法验证hash
func PasswordVerify(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// RandString 生成指定长度的随机字符串
func RandString(num uint) (string, error) {
	if num == 0 {
		return "", nil
	}
	b := make([]byte, num/2)
	n, err := rand.Read(b)
	if uint(n) != num/2 || err != nil {
		return "", err
	}
	str := hex.EncodeToString(b)
	return str[0:num], nil
}

// NewTraceId 生成全局ID
func NewTraceId() string {
	return constant.TraceIdPrefix + UniqueId()
}

func HashCode(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

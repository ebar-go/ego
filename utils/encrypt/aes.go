package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"github.com/ebar-go/ego/utils/conv"
)

type aesEncrypt struct {
	key []byte
}

func Aes(key string) *aesEncrypt {
	return &aesEncrypt{
		key: conv.Str2Byte(key),
	}
}

// padding 填充数据
func padding(src []byte, blockSize int) []byte {
	padNum := blockSize - len(src) % blockSize
	pad := bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append(src, pad...)
}

// unpadding 去掉填充数据
func unpadding(src []byte) []byte {
	n := len(src)
	unPadNum := int(src[n-1])
	return src[:n-unPadNum]
}

// Encrypt 加密
func (encrypt aesEncrypt) Encrypt(src []byte) ([]byte, error) {
	block, err := aes.NewCipher(encrypt.key)
	if err != nil {
		return nil, err
	}
	src = padding(src, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, encrypt.key)
	blockMode.CryptBlocks(src, src)
	return src, nil
}

// Decrypt 解密
func (encrypt aesEncrypt) Decrypt(src []byte) ([]byte, error) {
	block, err := aes.NewCipher(encrypt.key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, encrypt.key)
	blockMode.CryptBlocks(src, src)
	src = unpadding(src)
	return src, nil
}

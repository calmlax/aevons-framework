package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// DefaultSymmetricKey 默认的 AES 对称加密密钥。
// 在生产环境中，建议从环境变量或配置文件中注入提取此密钥。
var DefaultSymmetricKey = []byte("Aevons-32-byte-secret-key-123456")

// PKCS7Padding 为明文数据进行 PKCS7 填充，使其长度为 blockSize 的整数倍。
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding 去除数据尾部的 PKCS7 填充字符。
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length < unpadding {
		return origData
	}
	return origData[:(length - unpadding)]
}

// EncryptAES 使用 AES-CBC 模式加密明文，并返回 base64 编码的密文字符串。
func EncryptAES(plainText string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	origData := PKCS7Padding([]byte(plainText), blockSize)

	// 在 CBC 模式中，IV 向量的长度与块大小一致。为方便起见，默认使用密钥的前 blockSize 字节作为 IV。
	// 生产环境下更好的做法是生成随机 IV 并附加在密文开头。
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)

	return base64.StdEncoding.EncodeToString(crypted), nil
}

// DecryptAES 使用 AES-CBC 模式解密 base64 编码格式的密文，返回原明文字符串。
func DecryptAES(cryptoText string, key []byte) (string, error) {
	crypted, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	if len(crypted) < blockSize || len(crypted)%blockSize != 0 {
		return "", errors.New("ciphertext block size is not valid")
	}

	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)

	return string(origData), nil
}

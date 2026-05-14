package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

// GenerateRSAKey 生成 RSA 密钥对，默认推荐 2048 位
func GenerateRSAKey(bits int) (privateKeyStr, publicKeyStr string, err error) {
	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}

	// 编码私钥为 PKCS1 ASN.1 DER 格式并 PEM 编码
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	}
	privateKeyStr = string(pem.EncodeToMemory(privBlock))

	// 生成公钥
	publicKey := &privateKey.PublicKey
	pubDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}
	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubDER,
	}
	publicKeyStr = string(pem.EncodeToMemory(pubBlock))

	return privateKeyStr, publicKeyStr, nil
}

// RSADecrypt 使用预存的私钥解密前端公钥加密后经 Base64 编码的密文字符串
func RSADecrypt(cryptoText, privateKeyStr string) (string, error) {
	// 解析 base64 密文
	ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", errors.New("invalid base64 ciphertext")
	}

	// 解析 PEM 块
	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return "", errors.New("failed to decode PEM block containing private key")
	}

	// 解析 RSA 私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 使用 PKCS1v15 标准解密（与 JSEncrypt 默认行为兼容）
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

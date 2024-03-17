package tblschema

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

// hashKey 使用SHA-256散列函数将任意长度的密钥转换为256位密钥
func hashKey(key []byte) []byte {
	hash := sha256.Sum256(key)
	return hash[:]
}

// encryptString 使用AES-256对字符串进行加密
func EncryptString(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(hashKey(key))
	if err != nil {
		return "", err
	}

	plaintextBytes := []byte(plaintext)
	ciphertext := make([]byte, aes.BlockSize+len(plaintextBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintextBytes)

	return hex.EncodeToString(ciphertext), nil
}

// decryptString 使用AES-256对加密字符串进行解密
func DecryptString(encrypted string, key []byte) (string, error) {
	ciphertext, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(hashKey(key))
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		//密文太短
		return "", fmt.Errorf("解密失败")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

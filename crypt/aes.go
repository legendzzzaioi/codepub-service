package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Encrypt AES加密
func Encrypt(key []byte, text string) (string, error) {
	if len(key) != 32 {
		return "", errors.New("invalid key size, must be 32 bytes for AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintext := []byte(text)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt AES解密
func Decrypt(key []byte, cryptoText string) (string, error) {
	if len(key) != 32 {
		return "", errors.New("invalid key size, must be 32 bytes for AES-256")
	}

	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}

// GetEncryptionKey 密钥
func GetEncryptionKey() []byte {
	// Example key, should be exactly 32 bytes for AES-256
	// You can modify this to fetch from a secure source or environment variable
	key := "your-32-byte-long-encryption-key-"
	return []byte(key[:32])
}

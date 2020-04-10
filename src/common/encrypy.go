package common

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"sync"
)
import "crypto/cipher"

var (
	commonKey = []byte("xtkshishijieshangzuishuaideren@W")
	syncMutex sync.Mutex
)

func SetAesKey(key string) {
	syncMutex.Lock()
	defer syncMutex.Unlock()
	commonKey = []byte(key)
}

func AesEncrypt(plainText string) (string, error) {
	block, err := aes.NewCipher(commonKey)
	if err != nil {
		return "", err
	}
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cipher.NewCFBEncrypter(block, iv).XORKeyStream(cipherText[aes.BlockSize:],
		[]byte(plainText))

	return hex.EncodeToString(cipherText), nil

}
func AesDecrypt(d string) (string, error) {
	cipherText, err := hex.DecodeString(d)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(commonKey)
	if err != nil {
		return "", err
	}
	if len(cipherText) < aes.BlockSize {
		return "", errors.New("cipherText too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	cipher.NewCFBDecrypter(block, iv).XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

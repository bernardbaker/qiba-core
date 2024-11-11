package infrastructure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
)

type Encrypter struct {
	key []byte
}

func NewEncrypter(key []byte) *Encrypter {
	return &Encrypter{key: key}
}

func (e *Encrypter) EncryptGameData(data interface{}) (string, string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", "", err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, jsonData, nil)
	encryptedData := base64.StdEncoding.EncodeToString(ciphertext)

	// Generate HMAC
	mac := hmac.New(sha256.New, e.key)
	mac.Write(jsonData)
	hmacValue := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return encryptedData, hmacValue, nil
}

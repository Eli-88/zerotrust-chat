package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// interface compliance
var _ Key = &aesImpl{}

type aesImpl struct {
	secretKey []byte
}

func GenerateKey() (Key, error) {
	secretKey := make([]byte, 32)
	_, err := rand.Read(secretKey)
	if err != nil {
		return nil, err
	}

	return &aesImpl{
		secretKey: secretKey,
	}, nil
}

func GenerateKeyFromSecret(secretKey string) (Key, error) {
	decodedSecretKey, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		return nil, err
	}

	if len(decodedSecretKey) != 32 {
		return nil, errors.New("invalid secret key")
	}
	return &aesImpl{
		secretKey: decodedSecretKey,
	}, nil
}

func (a aesImpl) ToString() string {
	return base64.StdEncoding.EncodeToString(a.secretKey)
}

func (a aesImpl) Encrypt(plainText []byte) (string, error) {
	block, err := aes.NewCipher(a.secretKey)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]

	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (a aesImpl) Decrypt(cipherText string) ([]byte, error) {
	decodedCipherText, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(a.secretKey)
	if err != nil {
		return nil, err
	}

	if len(decodedCipherText) < aes.BlockSize {
		return nil, errors.New("ciphertext blocksize is too short")
	}

	iv := decodedCipherText[:aes.BlockSize]
	decodedCipherText = decodedCipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decodedCipherText, decodedCipherText)

	return decodedCipherText, nil
}

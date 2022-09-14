package rsa

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"
	"zerotrust_chat/crypto/rsa"
	test "zerotrust_chat/test/helper"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	key, err := rsa.GeneratePrivateKey()
	assert.NotNil(t, key)
	assert.NoError(t, err)
}

func TestEncryptDecrypt(t *testing.T) {
	key, _ := rsa.GeneratePrivateKey()
	plainText := test.GenerateRandPlainText(100)
	cipherText, err := key.GetPublicKey().Encrypt(plainText)
	println(cipherText)
	assert.NoError(t, err)

	decryptedPlainText, err := key.Decrypt(cipherText)
	assert.NoError(t, err)
	assert.True(t, bytes.Equal(decryptedPlainText, plainText))
}

func TestEncryptDecryptLargePlainText(t *testing.T) {
	key, _ := rsa.GeneratePrivateKey()
	plainText := test.GenerateRandPlainText(2048 / 8) // larger than the rsa key bit size

	cipherText, err := key.GetPublicKey().Encrypt(plainText)
	assert.Error(t, err)
	assert.True(t, cipherText == "")
}

func TestGeneratePublicKeyFromString(t *testing.T) {
	priKey, _ := rsa.GeneratePrivateKey()
	pubKeyStr := priKey.GetPublicKey().ToString()

	plainText := make([]byte, 100)
	io.ReadFull(rand.Reader, plainText)

	pubKey, err := rsa.GeneratePublicKey(pubKeyStr)
	assert.NoError(t, err)
	assert.NotNil(t, pubKey)

	cipherText, err := pubKey.Encrypt(plainText)
	assert.NoError(t, err)
	assert.NotNil(t, cipherText)

	decryptedPlainText, err := priKey.Decrypt(cipherText)
	assert.NoError(t, err)

	assert.True(t, bytes.Equal(decryptedPlainText, plainText))
}

package aes

import (
	"testing"
	"zerotrust_chat/crypto/aes"
	test "zerotrust_chat/test/helper"

	"github.com/stretchr/testify/assert"
)

func TestGenerateKey(t *testing.T) {
	key, err := aes.GenerateKey()
	assert.NoError(t, err)
	assert.NotNil(t, key)
}

func TestDecrypt(t *testing.T) {
	plainText := test.GenerateRandPlainText(100)
	key, _ := aes.GenerateKey()
	cipherText, err := key.Encrypt(plainText)
	assert.NoError(t, err)
	decryptedPlainText, err := key.Decrypt(cipherText)
	assert.NoError(t, err)
	assert.True(t, test.CompareByteSlices(decryptedPlainText, plainText))
}

func TestGenerateKeyFromString(t *testing.T) {
	key, _ := aes.GenerateKey()
	keyStr := key.ToString()
	newKey, err := aes.GenerateKeyFromSecret(keyStr)
	assert.NoError(t, err)

	plainText := test.GenerateRandPlainText(100)
	cipherText, _ := key.Encrypt(plainText)
	decryptedPlainText, _ := newKey.Decrypt(cipherText)
	assert.True(t, test.CompareByteSlices(decryptedPlainText, plainText))
}

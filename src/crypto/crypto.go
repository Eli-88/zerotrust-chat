package crypto

import (
	"zerotrust_chat/crypto/aes"
	"zerotrust_chat/crypto/rsa"
)

// interface compliance
var _ KeyFactory = keyFactory{}

type keyFactory struct{}

func NewKeyFactory() KeyFactory {
	return &keyFactory{}
}

func (k keyFactory) GenerateRsaPrivateKey() (rsa.PrivateKey, error) {
	return rsa.GeneratePrivateKey()
}

func (k keyFactory) GenerateAesSecretKey() (aes.Key, error) {
	return aes.GenerateKey()
}

func (k keyFactory) ConstructAesSecretKey(secretKey string) (aes.Key, error) {
	return aes.GenerateKeyFromSecret(secretKey)
}

func (k keyFactory) ConstructRsaPublicKey(pubKey string) (rsa.PublicKey, error) {
	return rsa.GeneratePublicKey(pubKey)
}

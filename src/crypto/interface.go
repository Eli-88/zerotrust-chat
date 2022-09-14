package crypto

import (
	"zerotrust_chat/crypto/aes"
	"zerotrust_chat/crypto/rsa"
)

//go:generate mockgen -destination=../test/mocks/mock_key_factory.go -package=mocks zerotrust_chat/crypto KeyFactory

type KeyFactory interface {
	GenerateRsaPrivateKey() (rsa.PrivateKey, error)
	GenerateAesSecretKey() (aes.Key, error)
	ConstructAesSecretKey(string) (aes.Key, error)
	ConstructRsaPublicKey(string) (rsa.PublicKey, error)
}

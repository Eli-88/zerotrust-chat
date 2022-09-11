package crypto

import (
	"zerotrust_chat/crypto/aes"
	"zerotrust_chat/crypto/rsa"
)

type KeyFactory interface {
	GenerateRsaPrivateKey() (rsa.PrivateKey, error)
	GenerateAesSecretKey() (aes.Key, error)
	ConstructAesSecretKey(string) (aes.Key, error)
	ConstructRsaPublicKey(string) (rsa.PublicKey, error)
}

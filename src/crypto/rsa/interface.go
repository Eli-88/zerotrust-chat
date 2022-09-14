package rsa

//go:generate mockgen -destination=../../test/mocks/mock_rsa_private.go -package=mocks zerotrust_chat/crypto/rsa PrivateKey
//go:generate mockgen -destination=../../test/mocks/mock_rsa_public.go -package=mocks zerotrust_chat/crypto/rsa PublicKey

type PrivateKey interface {
	Decrypt(string) ([]byte, error)
	GetPublicKey() PublicKey
	GetLabel() []byte
}

type PublicKey interface {
	Encrypt([]byte) (string, error)
	ToString() string
}

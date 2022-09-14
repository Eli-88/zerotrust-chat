package aes

//go:generate mockgen -destination=../../test/mocks/mock_aes_key.go -package=mocks zerotrust_chat/crypto/aes Key

type Key interface {
	ToString() string
	Encrypt([]byte) (string, error)
	Decrypt(string) ([]byte, error)
}

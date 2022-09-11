package aes

type Key interface {
	ToString() string
	Encrypt([]byte) (string, error)
	Decrypt(string) ([]byte, error)
}

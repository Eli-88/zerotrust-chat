package rsa

type PrivateKey interface {
	Decrypt(string) ([]byte, error)
	GetPublicKey() PublicKey
	GetLabel() []byte
}

type PublicKey interface {
	Encrypt([]byte) (string, error)
	ToString() string
}

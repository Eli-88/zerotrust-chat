package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
)

// interface compliance
var _ PrivateKey = &rsaPrivateKey{}
var _ PublicKey = &rsaPublicKey{}

const LABEL = "just a label"

type rsaPrivateKey struct {
	priKey *rsa.PrivateKey
	pubKey *rsaPublicKey
	label  []byte
}

type rsaPublicKey struct {
	pubKey rsa.PublicKey
	label  []byte
}

func GeneratePrivateKey() (PrivateKey, error) {
	priKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return nil, err
	}

	rsaPubKey := &rsaPublicKey{
		pubKey: priKey.PublicKey,
		label:  []byte(LABEL),
	}

	return &rsaPrivateKey{
		priKey: priKey,
		pubKey: rsaPubKey,
		label:  []byte(LABEL),
	}, nil
}

func GeneratePublicKey(pubKey string) (PublicKey, error) {
	decodedPubKey, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}
	pub, err := x509.ParsePKCS1PublicKey(decodedPubKey)
	if err != nil {
		return nil, err
	}
	return &rsaPublicKey{
		pubKey: *pub,
		label:  []byte(LABEL),
	}, nil
}

func (r *rsaPrivateKey) Decrypt(cipherText string) ([]byte, error) {
	decodedCipherText, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	priKey := r.priKey
	label := r.label

	return rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		priKey,
		decodedCipherText,
		label,
	)
}

func (r *rsaPublicKey) Encrypt(plainText []byte) (string, error) {

	pubKey := &r.pubKey
	label := r.label

	cipherText, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		pubKey,
		plainText,
		label,
	)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (r *rsaPublicKey) ToString() string {
	return base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&r.pubKey))
}

func (r rsaPrivateKey) Encrypt(plainText []byte) (string, error) {
	return r.pubKey.Encrypt(plainText)
}

func (r rsaPrivateKey) GetPublicKey() PublicKey {
	return r.pubKey
}

func (r rsaPrivateKey) GetLabel() []byte {
	return r.label
}

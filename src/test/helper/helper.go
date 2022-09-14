package test

import (
	"bytes"
	"crypto/rand"
	"io"
)

func GenerateRandPlainText(size int) []byte {
	result := make([]byte, size)
	io.ReadFull(rand.Reader, result)
	return result
}

func CompareByteSlices(a, b []byte) bool {
	return bytes.Equal(a, b)
}

type KeyExchangeRequest struct {
	PubKey string `json:"public_key"`
}

type KeyExchangeResponse struct {
	SecretKey string `json:"secret_key"`
}

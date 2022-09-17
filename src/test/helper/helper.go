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

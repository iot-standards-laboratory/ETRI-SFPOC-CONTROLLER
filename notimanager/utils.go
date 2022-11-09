package notimanager

import (
	"crypto/rand"
	"errors"
	"hash/crc64"
)

func GenerateSecureToken(length int) (uint64, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return 0, errors.New("error to generate randome value")
	}
	return getHashValue(b), nil
}

func getHashValue(b []byte) uint64 {
	return crc64.Checksum(b, crc64.MakeTable(crc64.ISO))
}

package jcrypt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func _make_uid(size int) string {
	randomBytes := make([]byte, size)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	h := sha256.Sum256(randomBytes)
	return hex.EncodeToString(h[:])
}

func MakeUID() string {
	return _make_uid(256)
}

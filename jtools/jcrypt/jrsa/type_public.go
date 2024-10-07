package jrsa

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
)

//cPublic :
type cPublic struct {
	key *rsa.PublicKey
}

//KeyToBytes :
func (my cPublic) KeyToBytes() []byte {
	b, err := PublicKeyToBytes(my.key)
	if err != nil {
		fmt.Println(err)
	}
	return b
}

//ToString :
func (my cPublic) ToString() string {
	b, err := PublicKeyToBytes(my.key)
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}

//Encrypt :
func (my cPublic) Encrypt(msg []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, my.key, msg)
}

//EncBase64 :
func (my cPublic) EncBase64(msg []byte) (string, error) {
	b, err := rsa.EncryptPKCS1v15(rand.Reader, my.key, msg)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

//EncBase64String :
func (my cPublic) EncBase64String(text string) (string, error) {
	return my.EncBase64([]byte(text))
}

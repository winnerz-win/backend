package jrsa

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
)

//cPrivate : IPrivate
type cPrivate struct {
	key *rsa.PrivateKey
}

//KeyToBytes :
func (my cPrivate) KeyToBytes() []byte {
	return PrivateKeyToBytes(my.key)
}

//ToString :
func (my cPrivate) ToString() string {
	return string(PrivateKeyToBytes(my.key))
}

//PublicKey :
func (my cPrivate) PublicKey() IPublic {
	return cPublic{
		key: &my.key.PublicKey,
	}
}

//Decrypt :
func (my cPrivate) Decrypt(ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, my.key, ciphertext)
}

//DecBase64 :
func (my cPrivate) DecBase64(b64 string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, my.key, ciphertext)
}

//DecBase64String :
func (my cPrivate) DecBase64String(b64 string) (string, error) {
	b, err := my.DecBase64(b64)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

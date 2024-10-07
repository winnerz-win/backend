package jrsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// GenerateKeyPair : generates a new key pair
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		fmt.Println(err)
	}
	return privkey, &privkey.PublicKey
}

// PrivateKeyToBytes : private key to bytes
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)
	return privBytes
}

// PublicKeyToBytes : public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) ([]byte, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes, nil
}

// BytesToPrivateKey : bytes to private key
func BytesToPrivateKey(priv []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		//fmt.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// BytesToPublicKey : bytes to public key
func BytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		//fmt.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		//fmt.Println("not ok")
	}
	return key, nil
}

// EncryptB64WithPublicKey : encrypts data with public key
func EncryptB64WithPublicKey(msg []byte, pub *rsa.PublicKey) (string, error) {
	chiperBuffer, err := rsa.EncryptPKCS1v15(rand.Reader, pub, msg)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(chiperBuffer), nil
}

// DecryptB64WithPrivateKey : jsencrypt가 OAEP가 아닌 PKCS1을 사용한다는 점
func DecryptB64WithPrivateKey(b64 string, priv *rsa.PrivateKey) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

// EncryptWithPublicKey : encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	//ciphertext, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, pub, msg, nil)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, msg)
}

// DecryptWithPrivateKey : decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	//plaintext, err := rsa.DecryptOAEP(sha512.New(), rand.Reader, priv, ciphertext, nil)
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

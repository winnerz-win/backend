package jrsa

//IPrivate :
type IPrivate interface {
	KeyToBytes() []byte
	ToString() string
	PublicKey() IPublic
	Decrypt(ciphertext []byte) ([]byte, error)
	DecBase64(b64 string) ([]byte, error)
	DecBase64String(b64 string) (string, error)
}

//IPublic :
type IPublic interface {
	KeyToBytes() []byte
	ToString() string
	Encrypt(msg []byte) ([]byte, error)
	EncBase64(msg []byte) (string, error)
	EncBase64String(text string) (string, error)
}

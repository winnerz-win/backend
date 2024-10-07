package jaes

//Openssl :
type Openssl interface {
	ToString() string
	EncryptBytesBytes(target []byte) ([]byte, error)
	EncryptBytesString(target []byte) (string, error)
	EncryptStringString(plain string) (string, error)

	EncryptBytesString1(target []byte) string
	EncryptStringString1(plain string) string

	DecryptBytesBytes(b64Bytes []byte) ([]byte, error)
	DecryptStringString(b64 string) (string, error)
	DecryptStringBytes(b64 string) ([]byte, error)
}

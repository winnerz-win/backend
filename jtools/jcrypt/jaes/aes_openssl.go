package jaes

//openssl "github.com/Luzifer/go-openssl"

//openAES : Openssl interface{}
type openAES string

var (
	digestMD5    = DigestMD5Sum //java
	digestSHA1   = DigestSHA1Sum
	digestSHA256 = DigestSHA256Sum

	currentDigest = digestSHA256
)

//New : Openssl
func New(key string) Openssl {
	return openAES(key)
}

//ToString :
func (my openAES) ToString() string {
	return string(my)
}

//EncryptBytesBytes :
func (my openAES) EncryptBytesBytes(target []byte) ([]byte, error) {
	op := NewOpenSSL()
	salt, _ := op.GenerateSalt()
	return op.EncryptBytesWithSaltAndDigestFunc(string(my), salt, target, currentDigest)
	//return op.EncryptBytes(string(my), target, currentDigest)
}

//DecryptBytesBytes :
func (my openAES) DecryptBytesBytes(b64Bytes []byte) ([]byte, error) {
	op := NewOpenSSL()
	return op.DecryptBytes(string(my), b64Bytes)
	// salt, _ := op.GenerateSalt()
	// return op.EncryptBytesWithSaltAndDigestFunc(string(my), salt, b64Bytes, currentDigest)
	//return openssl.New().DecryptBytes(string(my), b64Bytes, currentDigest)
}

//EncryptBytesString : butter -> base64
func (my openAES) EncryptBytesString(target []byte) (string, error) {
	if buf, err := my.EncryptBytesBytes(target); err != nil {
		return "<nil>", err
	} else {
		return string(buf), nil
	}
}

//EncryptStringString : string -> base64
func (my openAES) EncryptStringString(plain string) (string, error) {
	return my.EncryptBytesString([]byte(plain))
}

//EncryptBytesString1 :
func (my openAES) EncryptBytesString1(target []byte) string {
	s, e := my.EncryptBytesString(target)
	if e != nil {
		return ""
	}
	return s
}

//EncryptStringString1 :
func (my openAES) EncryptStringString1(plain string) string {
	s, e := my.EncryptStringString(plain)
	if e != nil {
		return ""
	}
	return s
}

//DecryptStringString :
func (my openAES) DecryptStringString(b64 string) (string, error) {
	if buf, err := my.DecryptBytesBytes([]byte(b64)); err != nil {
		return "<nil>", err
	} else {
		return string(buf), nil
	}
}

//DecryptStringByte :
func (my openAES) DecryptStringBytes(b64 string) ([]byte, error) {
	buf, err := my.DecryptBytesBytes([]byte(b64))
	if err != nil {
		return nil, err
	}
	return buf, err
}

package jrsa

///////////////////////////////////////////////////////////////////////////////////////
// Module Area
///////////////////////////////////////////////////////////////////////////////////////

//GenerateKeys :
func GenerateKeys(bits int) (IPrivate, IPublic) {
	prv, pub := GenerateKeyPair(bits)
	return cPrivate{prv}, cPublic{pub}
}

//ToPublicKey :
func ToPublicKey(pub []byte) (IPublic, error) {
	rKey := cPublic{}
	key, err := BytesToPublicKey(pub)
	if err != nil {
		return rKey, err
	}
	return cPublic{key: key}, nil
}

//ToPublicKeyString :
func ToPublicKeyString(pem string) (IPublic, error) {
	return ToPublicKey([]byte(pem))
}

//ToPrivateKey :
func ToPrivateKey(prv []byte) (IPrivate, error) {
	rKey := cPrivate{}
	key, err := BytesToPrivateKey(prv)
	if err != nil {
		return rKey, err
	}
	return cPrivate{key: key}, nil
}

//ToPrivateKeyString :
func ToPrivateKeyString(pem string) (IPrivate, error) {
	return ToPrivateKey([]byte(pem))
}

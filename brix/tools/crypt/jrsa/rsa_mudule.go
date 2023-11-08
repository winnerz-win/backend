package jrsa

import "txscheduler/brix/tools/crypt/jrsa/rsacore"

///////////////////////////////////////////////////////////////////////////////////////
// Module Area
///////////////////////////////////////////////////////////////////////////////////////

//GenerateKeys :
func GenerateKeys(bits int) (IPrivate, IPublic) {
	prv, pub := rsacore.GenerateKeyPair(bits)
	return cPrivate{prv}, cPublic{pub}
}

//ToPublicKey :
func ToPublicKey(pub []byte) (IPublic, error) {
	rKey := cPublic{}
	key, err := rsacore.BytesToPublicKey(pub)
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
	key, err := rsacore.BytesToPrivateKey(prv)
	if err != nil {
		return rKey, err
	}
	return cPrivate{key: key}, nil
}

//ToPrivateKeyString :
func ToPrivateKeyString(pem string) (IPrivate, error) {
	return ToPrivateKey([]byte(pem))
}

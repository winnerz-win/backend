package ebcm

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"txscheduler/brix/tools/dbg"

	"golang.org/x/crypto/sha3"
)

/*
https://gist.github.com/qbig/313b06424d96619144b9fd593ae64dac

MakeSign:
1.암호화할 원본 데이터
2.원본을 개인키로 사인
3.개인키의 공개키

------------------------------

Valid:
1.데이터와 사인 검증
2.데이터와사인으로 추출한 공개키 검증
3.데이터와 공개키 검증
*/
type GetSignTooler func(message_prefix MessagePrefix) SignTool
type SignTool struct {
	HexToECDSA             func(prvHex string) (*ecdsa.PrivateKey, error)
	FromECDSAPub           func(pub *ecdsa.PublicKey) []byte
	Keccak256HashBytes     func(data []byte) []byte
	GetEthereumMessageHash func(data []byte) []byte
	MessageV_addVal        int
	MessageV_subVal        int
	Sign                   func(keccak256Hash []byte, prv *ecdsa.PrivateKey) (sig []byte, err error)
	Ecrecover              func(keccak256Hash, sig []byte) (pub []byte, err error)
	SigToPub               func(keccak256Hash, sig []byte) (*ecdsa.PublicKey, error)
	VerifySignature        func(pubkey, digestHash, signature []byte) bool
}

func (my *SignTool) MakeSign(prvHex string, data []byte, callback func(SignDataBytes)) error {
	privateKey, err := my.HexToECDSA(prvHex)
	if err != nil {
		return err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("error casting public key to ECDSA")
	}
	publicKeyBytes := my.FromECDSAPub(publicKeyECDSA)

	dataHashBytes := my.Keccak256HashBytes(data)
	fmt.Println("dataHashBytes", dataHashBytes)
	signature, err := my.Sign(dataHashBytes, privateKey)

	//fmt.Println("원본데이타 사인 :", Hexutil_Encode(signature))
	if err != nil {
		return err
	}

	dataClone := make([]byte, len(data))
	copy(dataClone, data)
	callback(
		SignDataBytes{
			Data: dataClone,
			Sig:  signature,
			Pub:  publicKeyBytes,
		},
	)
	return nil
}
func (my *SignTool) Valid(signdata interface{}) (string, error) {

	box := SignDataBytes{}
	switch v := signdata.(type) {
	case SignDataBytes:
		box = v
	case SignDataHex:
		box = v.Bytes()
	default:
		return "", errors.New("signData Type is not SignDataBytes or SignDataHex")
	}
	if !box.Valid() {
		return "", errors.New("SignDataBytes Self Valid Fail")
	}

	dataHashBytes := my.Keccak256HashBytes(box.Data)
	sigPublicKey, err := my.Ecrecover(dataHashBytes, box.Sig)
	if err != nil {
		return "", err
	}
	if !bytes.Equal(sigPublicKey, box.Pub) {
		return "", errors.New("invalid signature")
	}
	sigPublicKeyECDSA, err := my.SigToPub(dataHashBytes, box.Sig)
	if err != nil {
		return "", err
	}
	sigPublicKeyBytes := my.FromECDSAPub(sigPublicKeyECDSA)
	if !bytes.Equal(sigPublicKeyBytes, box.Pub) {
		return "", errors.New("invalid publicKey")
	}

	signatureNoRecoverID := box.Sig[:len(box.Sig)-1] // remove recovery id
	if !my.VerifySignature(box.Pub, dataHashBytes, signatureNoRecoverID) {
		return "", errors.New("miss match both publicKeys")
	}

	publicKeyBytes := my.FromECDSAPub(sigPublicKeyECDSA)
	keccak := sha3.NewLegacyKeccak256()
	keccak.Write(publicKeyBytes[1:])
	a1 := fmt.Sprintf("%v", Hexutil_Encode(keccak.Sum(nil)[12:]))
	address := dbg.TrimToLower(a1)

	return address, nil
}

type SignDataBytes struct {
	Data []byte `bson:"data" json:"data"`
	Sig  []byte `bson:"sig" json:"sig"`
	Pub  []byte `bson:"pub,omitempty" json:"pub,omitempty"`
}

func (my SignDataBytes) String() string {
	// return dbg.Cat(
	// 	"SignDataBytes:", dbg.ENTER,
	// 	"Data:", my.Data, dbg.ENTER,
	// 	"Sig:", my.Sig, dbg.ENTER,
	// 	"Pub:", my.Pub, dbg.ENTER,
	// )
	return dbg.ToJsonString(my)
}

func (my SignDataBytes) Valid() bool {
	return len(my.Data) != 0 && len(my.Sig) != 0 && len(my.Pub) != 0
}

func (my SignDataBytes) Hex() SignDataHex {
	return SignDataHex{
		DataHex: "0x" + hex.EncodeToString(my.Data),
		SigHex:  "0x" + hex.EncodeToString(my.Sig),
		PubHex:  "0x" + hex.EncodeToString(my.Pub),
	}
}

type SignDataHex struct {
	DataHex string `bson:"data_hex" json:"data_hex"`
	SigHex  string `bson:"sig_hex" json:"sig_hex"`
	PubHex  string `bson:"pub_hex" json:"pub_hex"`
}

func (my SignDataHex) String() string { return dbg.ToJsonString(my) }

func (SignDataHex) TagString() []string {
	return []string{
		"data_hex", "0x.... (사이닝 할 데이터의 원본 hex값)",
		"sig_hex", "0x.... (사닝한 sign hex값 )",
		"pub_hex", "0x.... (퍼블릭 키 hex값)",
		"", "",
	}
}

func (my SignDataHex) Bytes() SignDataBytes {
	item := SignDataBytes{}

	//remove 0x
	if v, err := hex.DecodeString(my.DataHex[2:]); err == nil {
		item.Data = v
	}
	if v, err := hex.DecodeString(my.SigHex[2:]); err == nil {
		item.Sig = v
	}
	if v, err := hex.DecodeString(my.PubHex[2:]); err == nil {
		item.Pub = v
	}
	return item
}

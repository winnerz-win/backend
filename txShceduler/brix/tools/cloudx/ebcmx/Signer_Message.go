package ebcmx

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"

	"golang.org/x/crypto/sha3"
)

func signatureADDSUB(signature []byte, val int) []byte {
	clone := make([]byte, len(signature))
	copy(clone, signature)

	r := clone[:32]
	clone = clone[32:]
	s := clone[:32]
	clone = clone[32:]
	v := clone
	ss := jmath.ADD(v, val)

	vHexString := jmath.HEX(ss, true)
	if vHexString == "" {
		vHexString = "00"
	}
	msgSignatureBuf := dbg.Cat(
		//"0x",
		hex.EncodeToString(r),
		hex.EncodeToString(s),
		vHexString,
	)
	msgSignature, _ := hex.DecodeString(msgSignatureBuf)
	return msgSignature
}

func debugBuffer(b []byte) {
	//dbg.Purple("<debug>", "0x"+hex.EncodeToString(b))
}

func (my *SignTool) MessageSign(prvHex string, text string, callback func(SignMessageHex)) error {
	privateKey, err := my.HexToECDSA(prvHex)
	if err != nil {
		return err
	}

	dataHashBytes := my.GetEthereumMessageHash([]byte(text))
	debugBuffer(dataHashBytes)

	signature, err := my.Sign(dataHashBytes, privateKey)
	if err != nil {
		return err
	}
	debugBuffer(signature)

	msgSignature := signatureADDSUB(signature, my.MessageV_addVal)

	callback(
		SignMessageHex{
			Message: text,
			SigHex:  "0x" + hex.EncodeToString(msgSignature),
		},
	)

	return nil
}

type SignMessageHex struct {
	Message string `bson:"message" json:"message"`
	SigHex  string `bson:"sig_hex" json:"sig_hex"`
}

func (my SignMessageHex) String() string { return dbg.ToJSONString(my) }

func (SignMessageHex) TagString() []string {
	return []string{
		"message", "사이닝할 원본 데이타 텍스트",
		"sig_hex", "0x.... (사닝한 sign hex값 )",
		"", "",
	}
}
func (SignMessageHex) Help(prefix string, message string) string {
	if !strings.HasSuffix(prefix, ".") {
		prefix = prefix + "."
	}
	msg := dbg.Cat(
		prefix+"messge : 원본 데이타 텍스트(", message, "), ",
		prefix+"sig_hex : 0x....(sign hex)",
	)
	return msg
}

func (my *SignTool) MessageValid(signdata SignMessageHex) (string, error) {
	box := signdata

	if len(box.SigHex) < 2 {
		return "", errors.New("SignMessageHex Self Valid Fail")
	}

	sigBuf, err := hex.DecodeString(box.SigHex[2:])
	if err != nil {
		return "", err
	}

	boxSig := signatureADDSUB(sigBuf, my.MessageV_subVal)

	dataHashBytes := my.GetEthereumMessageHash([]byte(box.Message))
	debugBuffer(dataHashBytes)
	debugBuffer(boxSig)

	sigPublicKey, err := my.Ecrecover(dataHashBytes, boxSig)
	if err != nil {
		return "", err
	}
	sigPublicKeyECDSA, err := my.SigToPub(dataHashBytes, boxSig)
	if err != nil {
		return "", err
	}
	sigPublicKeyBytes := my.FromECDSAPub(sigPublicKeyECDSA)
	if !bytes.Equal(sigPublicKeyBytes, sigPublicKey) {
		return "", errors.New("invalid publicKey")
	}

	publicKeyBytes := my.FromECDSAPub(sigPublicKeyECDSA)
	keccak := sha3.NewLegacyKeccak256()
	keccak.Write(publicKeyBytes[1:])
	a1 := fmt.Sprintf("%v", Hexutil_Encode(keccak.Sum(nil)[12:]))
	address := dbg.TrimToLower(a1)

	return address, nil
}

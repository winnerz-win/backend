package ecs

import (
	"crypto/ecdsa"
	"txscheduler/brix/tools/cloud/ebcm"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"

	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

/*
	// <ebcm package : Sender.go
	const (
		MESSAGE_PREFIX_METAMASK = "\u0019Ethereum Signed Message:\n"
		MESSAGE_PREFIX_KLAYTN   = "\u0019Klaytn Signed Message:\n"
	)
*/

func getSignMessagePrefix(prefixMsg ebcm.MessagePrefix, messageLength int) []byte {
	prefix := prefixMsg.String() + dbg.Cat(messageLength)
	return []byte(prefix)
}

func getSignTool(message_prefix ebcm.MessagePrefix) ebcm.SignTool {

	keccak256HashBytes := func(data []byte) []byte {
		dataHash := crypto.Keccak256Hash(data)
		buf := dataHash.Bytes()
		return buf
	}

	signTool := ebcm.SignTool{
		HexToECDSA:         crypto.HexToECDSA,
		FromECDSAPub:       crypto.FromECDSAPub,
		Keccak256HashBytes: keccak256HashBytes,
		GetEthereumMessageHash: func(message []byte) []byte {
			prefix := getSignMessagePrefix(message_prefix, len(message))
			result := make([]byte, len(prefix)+len(message))
			size := copy(result, prefix)
			copy(result[size:], message)
			return keccak256HashBytes(result)
		},
		MessageV_addVal: 27,
		MessageV_subVal: -27,
		Sign:            crypto.Sign,
		Ecrecover: func(keccak256Hash, sig []byte) (pub []byte, err error) {
			return crypto.Ecrecover(keccak256Hash, sig)
		},
		SigToPub: func(keccak256Hash, sig []byte) (*ecdsa.PublicKey, error) {
			return crypto.SigToPub(keccak256Hash, sig)
		},
		VerifySignature: func(pubkey, digestHash, signature []byte) bool {
			return crypto.VerifySignature(pubkey, digestHash, signature)
		},
	}
	return signTool
}

type GasResult struct {
	val *big.Int
}

func (my GasResult) String() string       { return jmath.VALUE(my.val) }
func (my GasResult) GetFast() *big.Int    { return my.val }
func (my GasResult) GetFastest() *big.Int { return my.val }
func (my GasResult) GetSafeLow() *big.Int { return my.val }
func (my GasResult) GetAverage() *big.Int { return my.val }
func (my GasResult) GetBegger() *big.Int  { return my.val }

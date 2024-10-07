package EtherClient

import (
	"context"
	"crypto/ecdsa"
	"errors"

	"txscheduler/brix/tools/dbg"

	"github.com/ethereum/go-ethereum/crypto"
)

/*
	calculate nounce Value
*/

//NounceValue :
func NounceValue(client *EClient, privateKeyString string) (uint64, error) {
	privatekey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		return 0, dbg.MakeError("MakeSignTx@HexToECDSA", err)
	}
	publicKey := privatekey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return 0, errors.New("error casting public key to ECDSA")
	}
	fromHexAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client._PendingNonceAt(context.Background(), fromHexAddress)
	if err != nil {
		return 0, dbg.MakeError("MakeEtherTx@_PendingNonceAt", err)
	}
	return nonce, nil
}

package EtherClient

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"strings"

	"txscheduler/brix/tools/dbg"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// Transfer :
func Transfer() {

	//sender := "0x8CE5bb2013887eD586e6a87211aa126453368b7A"
}

// GetAddress :
func GetAddress(privateKeyString string) (string, error) {
	privatekey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		return "nil", err
	}
	publicKey := privatekey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "nil", errors.New("error casting public key to ECDSA")
	}
	fromHexAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return fromHexAddress.String(), nil
}

// GetHexToAddress :
func GetHexToAddress(hexString string) (string, error) {
	if strings.HasPrefix(hexString, "0x") {
		hexString = hexString[2:]
	}
	b, err := hex.DecodeString(hexString)
	if err != nil {
		return "", err
	}
	privatekey, err := crypto.ToECDSA(b)
	if err != nil {
		return "nil", err
	}

	publicKey := privatekey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "nil", errors.New("error casting public key to ECDSA")
	}

	// fmt.Println(publicKey)
	// fmt.Println(publicKeyECDSA)
	fromHexAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return fromHexAddress.String(), nil
}

// MakeSignTx : ERC-20
func MakeSignTx(client *EClient, privateKeyString string, tx *types.Transaction) (*types.Transaction, error) {
	privatekey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		return nil, dbg.MakeError("MakeSignTx@HexToECDSA", err)
	}

	chainID, err := client._NetworkID(context.Background())
	if err != nil {
		return nil, dbg.MakeError("MakeSignTx@_NetworkID", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privatekey)
	if err != nil {
		return nil, dbg.MakeError("MakeSignTx@SignTx", err)
	}

	return signedTx, nil
}

// SendTransaction :
func SendTransaction(client *EClient, signedTx *types.Transaction) error {
	err := client._SendTransaction(context.Background(), signedTx)
	if err != nil {
		return dbg.MakeError("SendTransaction@_SendTransaction", err)
	}

	//fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())
	return nil
}

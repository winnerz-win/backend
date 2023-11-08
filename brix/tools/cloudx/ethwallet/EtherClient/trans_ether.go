package EtherClient

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"

	"txscheduler/brix/tools/dbg"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

//MakeEtherTx : 유저끼리만일경우 ( 컨트랙에 송금시 로직 변경 필...)
func MakeEtherTx(client *EClient, privateKeyString, toAddress string, wei string, skipSigned ...bool) (*types.Transaction, error) {
	privatekey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		return nil, dbg.MakeError("MakeSignTx@HexToECDSA", err)
	}
	publicKey := privatekey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("error casting public key to ECDSA")
	}
	fromHexAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client._PendingNonceAt(context.Background(), fromHexAddress)
	if err != nil {
		return nil, dbg.MakeError("MakeEtherTx@_PendingNonceAt", err)
	}

	toHexAddress := common.HexToAddress(toAddress)

	value := new(big.Int)
	value.SetString(wei, 10)

	//gasLimit := uint64(21000) // in units
	gasLimit, err := client._EstimateGas(context.Background(), ethereum.CallMsg{
		From:  fromHexAddress,
		To:    &toHexAddress,
		Value: value,
		Data:  []byte{},
	})
	gasStation := NewGasStation()
	gasFastPrice := gasStation.GetFast()
	gasLowPrice := gasStation.GetSafeLow()

	var data []byte
	//tx := types.NewTransaction(nonce, toHexAddress, value, gasLimit, gasFastPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, dbg.MakeError("MakeSignTx@NetworkID", err)
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasLowPrice,
		GasFeeCap: gasFastPrice,
		Gas:       gasLimit,
		To:        &toHexAddress,
		Value:     value,
		Data:      data,
	})

	//사이닝을 스킵할경우 그냥 트랜잭션만 만들어서 반환 한다.
	if len(skipSigned) > 0 && skipSigned[0] == true {
		return tx, nil
	}

	var signer types.Signer
	if tx.Type() == 0 {
		signer = types.NewEIP155Signer(chainID)
	} else {
		signer = types.NewLondonSigner(chainID)
	}
	signedTx, err := types.SignTx(tx, signer, privatekey)
	if err != nil {
		return nil, dbg.MakeError("MakeSignTx@SignTx", err)
	}

	// if err := client.SendTransaction(context.Background(), signedTx); err != nil {
	// 	return nil, dbg.MakeError("MakeEtherTx@SendTransaction", err)
	// }
	return signedTx, nil
}

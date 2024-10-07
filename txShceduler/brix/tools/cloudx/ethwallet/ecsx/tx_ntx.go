package ecsx

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// NTXList :
type NTXList []*NTX

// Append :
func (my *NTXList) Append(ntx *NTX) {
	*my = append(*my, ntx)
}

// NTX  unsigned transaction
type NTX struct {
	client     *ethclient.Client
	tx         *types.Transaction
	privatekey *ecdsa.PrivateKey
	chainID    *big.Int

	from string
	to   string
	wei  string

	gasLimit uint64
	gasRight *big.Int

	nonceCount uint64

	gasFeeWei string //트랜젝션 수수료

	UserData TxUserData

	snap GasSnapShot
}

func (my NTX) SnapShot() GasSnapShot { return my.snap }

// From :
func (my NTX) From() string {
	return my.from
}

// To :
func (my NTX) To() string {
	return my.to
}

// Wei :
func (my NTX) Wei() string {
	return my.wei
}

// GasFeeWei : 가스 수수료
func (my NTX) GasFeeWei() string {
	return my.gasFeeWei
}
func (my NTX) GasFeeETH() string {
	return WeiToETH(my.gasFeeWei)
}

// Hash : no signed Tx Hash
func (my NTX) Hash() string {
	return my.tx.Hash().Hex()
}

// Transaction :
func (my NTX) Transaction() *types.Transaction {
	return my.tx
}

// String :
func (my NTX) String() string {
	view := txViewer{my.tx}.Map()
	b, _ := json.MarshalIndent(view, "", "  ")
	return string(b)
}

// Tx : signed transaction
func (my *NTX) Tx() (*STX, error) {

	var signer types.Signer
	if my.tx.Type() == 2 {
		signer = types.NewLondonSigner(my.chainID)
	} else {
		signer = types.NewEIP155Signer(my.chainID)
	}

	signedTx, err := types.SignTx(my.tx, signer, my.privatekey)
	if err != nil {
		return nil, fmt.Errorf("NTX.Tx : %v", err)
	}
	stx := &STX{
		client:      my.client,
		tx:          signedTx,
		isSend:      false,
		from:        my.from,
		blockNumber: "",
		gasFeeWei:   my.gasFeeWei,

		UserData: my.UserData.Clone(),
	}
	return stx, nil
}

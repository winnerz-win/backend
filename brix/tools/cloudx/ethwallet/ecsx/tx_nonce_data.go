package ecsx

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"txscheduler/brix/tools/cloudx/ethwallet/ecsaa"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx/jwalletx"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

// NonceData :
type NonceData struct {
	Sender *Sender
	client *ethclient.Client

	privateKeyString string
	privatekey       *ecdsa.PrivateKey
	fromHexAddress   common.Address
	nonceValue       uint64 //nonce value
	chainID          *big.Int

	etherTags []string
}

// CoreClient : ICore
func (my *NonceData) CoreClient() *ethclient.Client { return my.client }

func (my NonceData) String() string {
	j := map[string]interface{}{}
	j["from"] = my.fromHexAddress.Hex()
	j["chain_id"] = jmath.VALUE(my.chainID)
	j["nonce"] = my.nonceValue
	j["tags"] = my.etherTags
	b, _ := json.MarshalIndent(j, "", "  ")
	return string(b)
}

func (my NonceData) isEtherTag(tokenAddress string) bool {
	for _, tag := range my.etherTags {
		if tokenAddress == tag {
			return true
		}
	}
	return false
}

// NonceCount :
func (my NonceData) NonceCount() uint64 { return my.nonceValue }

// BoxTx :
func (my *NonceData) BoxTx(boxdata GasBoxData, nAdd ...uint64) (*NTX, error) {
	if boxdata.fromAddress != strings.ToLower(my.fromHexAddress.Hex()) {
		return nil, errors.New(fmt.Sprintf("NonceData.BoxTx -- Invalid fromAddress : [%v/%v]",
			boxdata.fromAddress,
			strings.ToLower(my.fromHexAddress.Hex())),
		)
	}

	ntx := &NTX{
		client: my.client,
		//tx:         tx,
		privatekey: my.privatekey,
		chainID:    my.chainID,
		from:       strings.ToLower(my.fromHexAddress.Hex()),
		to:         boxdata.toAddress,
		wei:        boxdata.wei,
		gasFeeWei:  boxdata.GasWei(),

		UserData: newTxUserData(),
	}

	nonce := my.nonceValue
	if len(nAdd) > 0 {
		nonce += nAdd[0]
	}

	if boxdata.isETH {
		toHexAddress := common.HexToAddress(strings.ToLower(boxdata.toAddress))
		value := new(big.Int)
		value.SetString(boxdata.wei, 10)

		var data []byte

		switch my.Sender.txnType {
		case TXN_EIP_1559:

			tip_cap := ecsaa.SUGGEST_TIP_PRICE(my.Sender)
			ntx.tx = types.NewTx(&types.DynamicFeeTx{
				ChainID:   my.Sender.chainID,
				Nonce:     nonce,
				GasTipCap: tip_cap,
				GasFeeCap: boxdata.gasPrice,
				Gas:       boxdata.gasLimit,
				To:        &toHexAddress,
				Value:     value,
				Data:      data,
			})

		default:
			ntx.tx = types.NewTransaction(nonce, toHexAddress, value, boxdata.gasLimit, boxdata.gasPrice, data)
		}

	} else {
		tokenHexAddress := common.HexToAddress(strings.ToLower(boxdata.tokenAddress))

		ethValue := big.NewInt(0) // in wei (0 eth)

		switch my.Sender.txnType {
		case TXN_EIP_1559:
			tip_cap := ecsaa.SUGGEST_TIP_PRICE(my.Sender)
			ntx.tx = types.NewTx(&types.DynamicFeeTx{
				ChainID:   my.Sender.chainID,
				Nonce:     nonce,
				GasTipCap: tip_cap,
				GasFeeCap: boxdata.gasPrice,
				Gas:       boxdata.gasLimit,
				To:        &tokenHexAddress,
				Value:     ethValue,
				Data:      boxdata.inputData,
			})

		default:
			ntx.tx = types.NewTransaction(nonce, tokenHexAddress, ethValue, boxdata.gasLimit, boxdata.gasPrice, boxdata.inputData)
		}

	}

	ntx.nonceCount = nonce
	return ntx, nil
}

// TxETH : no signed to ether
func (my *NonceData) TxETH(toAddress, wei string, speed GasSpeed, nAdd ...uint64) (*NTX, error) {
	toHexAddress := common.HexToAddress(toAddress)

	value := new(big.Int)
	value.SetString(wei, 10)

	//gasLimit := uint64(21000) // in units
	gasLimit, err := my.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  my.fromHexAddress,
		To:    &toHexAddress,
		Value: value,
		Data:  []byte{},
	})
	dbg.Yellow("gasLimit :", gasLimit)
	if err != nil {
		return nil, fmt.Errorf("NonceData.EstimateGas : %v", err)
	}

	//gasPrice := NewGasStation().Price(speed.Value())
	gasPrice := ecsaa.SUGGEST_GAS_PRICE(my.Sender)
	dbg.Yellow("gasPrice :", gasPrice.String())

	nonce := my.nonceValue
	if len(nAdd) > 0 {
		nonce += nAdd[0]
	}

	var data []byte

	var tx *types.Transaction
	switch my.Sender.txnType {
	case TXN_EIP_1559:
		tip_cap := ecsaa.SUGGEST_TIP_PRICE(my.Sender)

		tx = types.NewTx(&types.DynamicFeeTx{
			ChainID:   my.Sender.chainID,
			Nonce:     nonce,
			GasTipCap: tip_cap,
			GasFeeCap: gasPrice,
			Gas:       gasLimit,
			To:        &toHexAddress,
			Value:     value,
			Data:      data,
		})

	default:
		tx = types.NewTransaction(nonce, toHexAddress, value, gasLimit, gasPrice, data)
	}

	ntx := &NTX{
		client:     my.client,
		tx:         tx,
		privatekey: my.privatekey,
		chainID:    my.chainID,
		from:       strings.ToLower(my.fromHexAddress.Hex()),
		to:         toAddress,
		wei:        wei,
	}
	return ntx, nil
}

// TxToken : [ tokenAddress == "eth", "ether", "ethereum" -> TxETH ]
func (my *NonceData) TxToken(tokenAddress, toAddress, wei string, speed GasSpeed, nAdd ...uint64) (*NTX, error) {
	ethName := strings.ToLower(tokenAddress)

	if my.isEtherTag(ethName) {
		return my.TxETH(toAddress, wei, speed, nAdd...)
	}

	tokenHexAddress := common.HexToAddress(strings.ToLower(tokenAddress))
	toHexAddress := common.HexToAddress(strings.ToLower(toAddress))

	//gasPrice := NewGasStation().Price(speed.Value())
	gasPrice := ecsaa.SUGGEST_GAS_PRICE(my.Sender)

	paddedAddress := common.LeftPadBytes(toHexAddress.Bytes(), 32)

	amount := new(big.Int)
	amount.SetString(wei, 10)

	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit, err := my.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: my.fromHexAddress,
		To:   &tokenHexAddress,
		//Value: amount,	//이더value만 들어간다.
		Data: data,
	})
	if err != nil {
		return nil, fmt.Errorf("NonceData.TxToken : %v", err)
	}
	//gasLimit += EtherClient.TXLimitTokenAddValue

	nonce := my.nonceValue
	if len(nAdd) > 0 {
		nonce += nAdd[0]
	}

	ethValue := big.NewInt(0) // in wei (0 eth)

	var tx *types.Transaction
	switch my.Sender.txnType {
	case TXN_EIP_1559:
		tip_cap := ecsaa.SUGGEST_TIP_PRICE(my.Sender)
		tx = types.NewTx(&types.DynamicFeeTx{
			ChainID:   my.Sender.chainID,
			Nonce:     nonce,
			GasTipCap: tip_cap,
			GasFeeCap: gasPrice,
			Gas:       gasLimit,
			To:        &toHexAddress,
			Value:     ethValue,
			Data:      data,
		})

	default:
		tx = types.NewTransaction(nonce, tokenHexAddress, ethValue, gasLimit, gasPrice, data)
	}

	newSnap := GasSnapShot{
		Limit:  gasLimit,
		Price:  gasPrice.String(),
		FeeWei: gasPrice.Mul(gasPrice, big.NewInt(int64(gasLimit))).String(),
	}
	ntx := &NTX{
		client:     my.client,
		tx:         tx,
		privatekey: my.privatekey,
		chainID:    my.chainID,
		from:       strings.ToLower(my.fromHexAddress.Hex()),
		gasFeeWei:  newSnap.FeeWei,
		snap:       newSnap,
	}
	return ntx, nil
}

// NonceValuePrivate :
func NonceValuePrivate(ic ICore, hexPrivate string) uint64 {
	w, err := jwalletx.Get(hexPrivate)
	if err != nil {
		dbg.Red("ecsx.NonceValuePrivate :", err)
		return 0
	}
	return NonceValue(ic, w.Address())
}

// NonceValue :
func NonceValue(ic ICore, hexAddress string) uint64 {
	address := common.HexToAddress(hexAddress)
	client := ic.CoreClient()
	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		dbg.Red("ecsx.NonceValue :", err)
	}
	return nonce
}

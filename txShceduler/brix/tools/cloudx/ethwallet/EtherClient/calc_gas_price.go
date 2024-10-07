package EtherClient

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"txscheduler/brix/tools/cloudx/ethwallet/EtherScanAPI"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

const (
	//TXLimitTokenAddValue : Token 송금시 gas limit에 추가할 값.
	TXLimitTokenAddValue = 1000000000 //"1" //GWEI(1) //
)

// GetGasPriceEther :
func GetGasPriceEther(client *EClient, fromAddress, toAddress string, wei string, isLog ...bool) *big.Int {

	fromHexAddress := common.HexToAddress(strings.ToLower(fromAddress))
	toHexAddress := common.HexToAddress(strings.ToLower(toAddress))

	value := new(big.Int)
	value.SetString(wei, 10)
	gasLimit, err := client._EstimateGas(context.Background(), ethereum.CallMsg{
		From:  fromHexAddress,
		To:    &toHexAddress,
		Value: value,
		Data:  []byte{},
	})
	if err != nil {
		fmt.Println("GetGasPriceEther@_EstimateGas", err)
		return nil
	}
	//gasLimit := uint64(21000) // in units

	limit := big.NewInt(int64(gasLimit))
	gasPrice := NewGasStation().GetFast()

	if len(isLog) > 0 && isLog[0] {
		fmt.Println("limit    :", gasLimit, "(", EtherScanAPI.WeiToEtherString(fmt.Sprintf("%v", gasLimit), 18), ")")
		fmt.Println("gasPrice :", gasPrice, "(", EtherScanAPI.WeiToEtherString(gasPrice.String(), 18), ")")
	}

	fee := gasPrice.Mul(gasPrice, limit)
	return fee
}

// GetGasPriceToken :
func GetGasPriceToken(client *EClient, tokenAddress, fromAddress, toAddress string, value string, isLog ...bool) *big.Int {
	tokenHexAddress := common.HexToAddress(strings.ToLower(tokenAddress))
	fromHexAddress := common.HexToAddress(strings.ToLower(fromAddress))
	toHexAddress := common.HexToAddress(strings.ToLower(toAddress))

	gasPrice := NewGasStation().GetFast()

	paddedAddress := common.LeftPadBytes(toHexAddress.Bytes(), 32)

	amount := new(big.Int)
	amount.SetString(value, 10)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit, err := client._EstimateGas(context.Background(), ethereum.CallMsg{
		From: fromHexAddress,
		To:   &tokenHexAddress,
		//Value: amount,	//이더value만 들어간다.
		Data: data,
	})
	if err != nil {
		fmt.Println("GetGasPriceToken@_EstimateGas", err)
		return nil
	}
	gasLimit += TXLimitTokenAddValue

	if len(isLog) > 0 && isLog[0] {
		fmt.Println("limit    :", gasLimit, "(", EtherScanAPI.WeiToEtherString(fmt.Sprintf("%v", gasLimit), 18), ")")
		fmt.Println("gasPrice :", gasPrice, "(", EtherScanAPI.WeiToEtherString(gasPrice.String(), 18), ")")
	}

	limit := big.NewInt(int64(gasLimit))
	fee := gasPrice.Mul(gasPrice, limit)
	return fee
}

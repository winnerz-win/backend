package EtherClient

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"txscheduler/brix/tools/cloudx/ethwallet/EtherScanAPI"
	"txscheduler/brix/tools/dbg"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/crypto/sha3"
)

// MakeTokenTx : ERC-20 토큰 trx
func MakeTokenTx(client *EClient, tokenAddress, fromAddress, toAddress string, value string, decimalFullValue bool) (*types.Transaction, error) {
	tokenHexAddress := common.HexToAddress(strings.ToLower(tokenAddress))
	fromHexAddress := common.HexToAddress(strings.ToLower(fromAddress))
	toHexAddress := common.HexToAddress(strings.ToLower(toAddress))

	gasPrice := NewGasStation().GetFast()

	paddedAddress := common.LeftPadBytes(toHexAddress.Bytes(), 32)

	amount := new(big.Int)
	if decimalFullValue == false {
		etoken := client.GetToken(tokenAddress, false)
		decimal := etoken.Decimals()
		fVal := EtherScanAPI.EtherToWeiString(value, decimal)
		amount.SetString(fVal, 10)
	} else {
		amount.SetString(value, 10)
	}
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
		return nil, dbg.MakeError("MakeTokenTx@_EstimateGas", err)
	}
	gasLimit += TXLimitTokenAddValue

	nonce, err := client._PendingNonceAt(context.Background(), fromHexAddress)
	if err != nil {
		return nil, dbg.MakeError("MakeTokenTx@_PendingNonceAt", err)
	}

	ethValue := big.NewInt(0) // in wei (0 eth)
	tx := types.NewTransaction(nonce, tokenHexAddress, ethValue, gasLimit, gasPrice, data)

	return tx, nil
}

//////////////////////////////////////////////////////////////////////////////

// PadBytes :
type PadBytes []byte

// NewPadBytes :
func NewPadBytes(methodName string) PadBytes {
	transferFnSignature := []byte(methodName) // ----- "transfer(address,uint256)"
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	pad := hash.Sum(nil)[:4]
	return pad
}

// Bytes :
func (my PadBytes) Bytes() []byte {
	return []byte(my)
}

// Append :
func (my *PadBytes) Append(data PadBytes) {
	*my = append(*my, data...)
}

// AppendAddress :
func (my *PadBytes) AppendAddress(address string) {
	hexAddress := common.HexToAddress(strings.ToLower(address))
	my.Append(common.LeftPadBytes(hexAddress.Bytes(), 32))
}

// AppendAmount :
func (my *PadBytes) AppendAmount(value string) {
	amount := new(big.Int)
	amount.SetString(value, 10)
	my.Append(common.LeftPadBytes(amount.Bytes(), 32))
}

// AppendAddressArray :
func (my *PadBytes) AppendAddressArray(addresslist ...string) {
	count := len(addresslist)
	size := new(big.Int)
	size.SetString(fmt.Sprint(count), 10)
	my.Append(common.LeftPadBytes(size.Bytes(), 32))

	for _, address := range addresslist {
		my.AppendAddress(address)
	} //for
}

// AppendAmountArray :
func (my *PadBytes) AppendAmountArray(values ...string) {
	count := len(values)
	size := new(big.Int)
	size.SetString(fmt.Sprint(count), 10)
	my.Append(common.LeftPadBytes(size.Bytes(), 32))

	for _, value := range values {
		my.AppendAmount(value)
	} //for
}

// MakeTokenTxCallback :
func MakeTokenTxCallback(
	client *EClient,
	contractAddress,
	signerAddress string,
	makePaddingFn func() PadBytes) (*types.Transaction, error) {

	contractHexAddress := common.HexToAddress(strings.ToLower(contractAddress))
	signerHexAddress := common.HexToAddress(strings.ToLower(signerAddress))

	paddedData := makePaddingFn()

	gasLimit, err := client._EstimateGas(context.Background(), ethereum.CallMsg{
		From: signerHexAddress,
		To:   &contractHexAddress,
		Data: paddedData.Bytes(),
	})
	if err != nil {
		return nil, dbg.MakeError("MakeTokenTxCallback@_EstimateGas", err)
	}
	gasLimit += TXLimitTokenAddValue

	nonce, err := client._PendingNonceAt(context.Background(), signerHexAddress)
	if err != nil {
		return nil, dbg.MakeError("MakeTokenTx@_PendingNonceAt", err)
	}

	gasPrice := NewGasStation().GetFast()
	ethValue := big.NewInt(0) // in wei (0 eth)
	tx := types.NewTransaction(nonce, contractHexAddress, ethValue, gasLimit, gasPrice, paddedData)

	return tx, nil
}

// GetGasPriceTokenCallback :
func GetGasPriceTokenCallback(
	client *EClient,
	contractAddress,
	signerAddress string,
	makePaddingFn func() PadBytes,
	isLog ...bool) *big.Int {

	contractHexAddress := common.HexToAddress(strings.ToLower(contractAddress))
	signerHexAddress := common.HexToAddress(strings.ToLower(signerAddress))

	paddedData := makePaddingFn()

	gasLimit, err := client._EstimateGas(context.Background(), ethereum.CallMsg{
		From: signerHexAddress,
		To:   &contractHexAddress,
		Data: paddedData.Bytes(),
	})
	if err != nil {
		fmt.Println("GetGasPriceTokenCallback@_EstimateGas", err)
		return nil
	}
	gasLimit += TXLimitTokenAddValue

	gasPrice := NewGasStation().Price(GasFast)

	if len(isLog) > 0 && isLog[0] {
		fmt.Println("limit    :", gasLimit, "(", EtherScanAPI.WeiToEtherString(fmt.Sprint(gasLimit), 18), ")")
		fmt.Println("gasPrice :", gasPrice, "(", EtherScanAPI.WeiToEtherString(gasPrice.String(), 18), ")")
	}

	limit := big.NewInt(int64(gasLimit))
	fee := gasPrice.Mul(gasPrice, limit)
	return fee
}

// CallContract :
func CallContract(
	client *EClient,
	contractAddress,
	signerAddress string,
	makePaddingFn func() PadBytes) ([]byte, error) {

	contractHexAddress := common.HexToAddress(strings.ToLower(contractAddress))
	signerHexAddress := common.HexToAddress(strings.ToLower(signerAddress))

	paddedData := makePaddingFn()

	blockNumber := big.NewInt(0)
	return client._CallContract(context.Background(), ethereum.CallMsg{
		From: signerHexAddress,
		To:   &contractHexAddress,
		Data: paddedData.Bytes(),
	}, blockNumber)

}

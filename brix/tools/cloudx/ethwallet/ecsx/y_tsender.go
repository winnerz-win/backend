package ecsx

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"txscheduler/brix/tools/cloudx/ethwallet/ecsaa"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type TSender struct {
	*Sender
	contractAddress string
}

// ContractAddress : contract-address
func (my TSender) ContractAddress() string { return my.contractAddress }

// TNew : tSender (token contract sender)
func TNew(mainnet bool, infuraKey string, contractAddress string) *TSender {
	return New(mainnet, infuraKey).TSender(contractAddress)
}

// TNewMainnet : tSender (token contract sender)
func TNewMainnet(host, infuraKey string, tags []string, contractAddress string) *TSender {
	return NewMainnet(host, infuraKey, tags).TSender(contractAddress)
}

// TSender :
func (my *Sender) TSender(contractAddress string) *TSender {
	return &TSender{
		Sender:          my,
		contractAddress: dbg.TrimToLower(contractAddress),
	}
}

func PadBytesFromHex(hexString string) PadBytes {
	buf, err := hex.DecodeString(hexString)
	if err != nil {
		dbg.Red(err)
	}
	return PadBytes{
		cPadBytes: buf,
	}
}

// MakePadBytes : ecsx.MakePadBytes
func (my *TSender) MakePadBytes(method string, callback func(Appender)) PadBytes {
	return MakePadBytes(method, callback)
}

// GasPrice :
func (my *TSender) GasPrice(fromAddress string, paddedData PadBytes, wei string, speed GasSpeed, isLog ...bool) (string, error) {
	logView := dbg.IsTrue2(isLog...)
	contractHexAddress := common.HexToAddress(my.contractAddress)
	signerHexAddress := common.HexToAddress(strings.ToLower(fromAddress))
	// fmt.Println(contractHexAddress)
	// fmt.Println(signerHexAddress)
	// fmt.Println(paddedData.Bytes())

	var ethValue *big.Int
	if !jmath.IsUnderZero(wei) {
		ethValue = new(big.Int)
		ethValue.SetString(jmath.VALUE(wei), 10)
	}

	gasLimit, err := my.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  signerHexAddress,
		To:    &contractHexAddress,
		Value: ethValue,
		Data:  paddedData.Bytes(),
	})
	if err != nil {
		if logView {
			dbg.RedUL("tSender.GasPrice@EstimateGas", err)
		}
		return "0", err
	}
	//gasLimit += EtherClient.TXLimitTokenAddValue

	gasPrice := ecsaa.SUGGEST_GAS_PRICE(my)
	if logView {
		dbg.PurpleItalic("limit    :", gasLimit, "(", WeiToETH(fmt.Sprintf("%v", gasLimit)), ")")
		dbg.PurpleItalic("gasPrice :", gasPrice, "(", WeiToETH(fmt.Sprintf("%v", gasPrice)), ")")
	}
	limit := big.NewInt(int64(gasLimit))
	fee := gasPrice.Mul(gasPrice, limit)
	return fee.String(), nil
}

// type TGasNonceData struct {
// 	Error error
// }

// func (my *TSender) GasNonceData(
// 	privateKeyString string,
// 	paddedData PadBytes,
// 	wei string,
// 	speed GasSpeed,
// 	isLog ...bool) TGasNonceData {

// 	logView := dbg.IsTrue2(isLog...)

// 	rError := func(tag string, e error) TGasNonceData {
// 		if logView {
// 			dbg.Red("[TSender.GasNonceData]::", tag, ":", e)
// 		}
// 		return TGasNonceData{Error: e}
// 	}

// 	privatekey, err := crypto.HexToECDSA(privateKeyString)
// 	if err != nil {
// 		return rError("PrivateKey", err)
// 	}
// 	publicKey := privatekey.Public()
// 	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
// 	if !ok {
// 		return rError("PrivateKey2:", errors.New("publicKeyECDSA"))
// 	}
// 	fromHexAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
// 	nonce, err := my.client.PendingNonceAt(context.Background(), fromHexAddress)
// 	if err != nil {
// 		return rError("PendingNonceAt", err)
// 	}
// 	chainID, err := my.client.NetworkID(context.Background())
// 	if err != nil {
// 		return rError("NetworkID", err)
// 	}

// 	data := TGasNonceData{}
// 	return data
// }

// tNonce :
func (my *TSender) Nonce(privateKeyString string) (*tNonce, error) {
	privatekey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		return nil, fmt.Errorf("tSender.Nonce@HexToECDSA : %v", err)
	}
	publicKey := privatekey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("tSender.Nonce@publicKeyECDSA")
	}
	fromHexAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := my.client.PendingNonceAt(context.Background(), fromHexAddress)
	if err != nil {
		return nil, fmt.Errorf("tSender.Nonce@PendingNonceAt : %v", err)
	}

	chainID, err := my.client.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("tSender.Nonce@NetworkID : %v", err)
	}

	nd := &tNonce{
		contractAddress:  my.contractAddress,
		Sender:           my,
		privateKeyString: privateKeyString,
		privatekey:       privatekey,
		fromHexAddress:   fromHexAddress,
		nonceValue:       nonce,
		chainID:          chainID,
	}

	for _, tag := range my.etherTags {
		nd.etherTags = append(nd.etherTags, tag)
	}
	return nd, nil
}

type tNonce struct {
	contractAddress string
	Sender          *TSender

	privateKeyString string
	privatekey       *ecdsa.PrivateKey
	fromHexAddress   common.Address
	nonceValue       uint64 //nonce value
	chainID          *big.Int

	etherTags []string
}

// NonceCount :
func (my tNonce) NonceCount() uint64 { return my.nonceValue }

func (my tNonce) String() string {
	j := map[string]interface{}{}
	j["contract_address"] = my.contractAddress
	j["from"] = my.fromHexAddress.Hex()
	j["chain_id"] = jmath.VALUE(my.chainID)
	j["nonce"] = my.nonceValue
	j["tags"] = my.etherTags
	b, _ := json.MarshalIndent(j, "", "  ")
	return string(b)
}

// NTX :
func (my *tNonce) NTX(paddedData PADBYTES, wei string, speed GasSpeed, nAdd ...uint64) (*NTX, error) {
	contractHexAddress := common.HexToAddress(my.contractAddress)
	signerHexAddress := my.fromHexAddress
	// fmt.Println(contractHexAddress)
	// fmt.Println(signerHexAddress)
	// fmt.Println(paddedData.Bytes())

	var ethValue *big.Int
	if !jmath.IsUnderZero(wei) {
		ethValue = new(big.Int)
		ethValue.SetString(jmath.VALUE(wei), 10)
	}

	gasLimit, err := my.Sender.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  signerHexAddress,
		To:    &contractHexAddress,
		Value: ethValue,
		Data:  paddedData.Bytes(),
	})
	if err != nil {
		return nil, dbg.MakeError("tNonce.NTX@EstimateGas", err)
	}
	//gasLimit += EtherClient.TXLimitTokenAddValue

	nonce := my.nonceValue
	if len(nAdd) > 0 {
		nonce += nAdd[0]
	}

	gasPrice := ecsaa.SUGGEST_GAS_PRICE(my.Sender)

	if ethValue == nil {
		ethValue = big.NewInt(0) // in wei (0 eth)
	}

	var tx *types.Transaction
	switch my.Sender.txnType {
	case TXN_EIP_1559:
		lowPrice := ecsaa.SUGGEST_GAS_PRICE(my.Sender)
		if gasPrice.Cmp(lowPrice) < 0 {
			gasPrice = big.NewInt(lowPrice.Int64())
		}

		tx = types.NewTx(&types.DynamicFeeTx{
			ChainID:   my.Sender.chainID,
			Nonce:     nonce,
			GasTipCap: lowPrice,
			GasFeeCap: gasPrice,
			Gas:       gasLimit,
			To:        &contractHexAddress,
			Value:     ethValue,
			Data:      paddedData.Bytes(),
		})
	default:
		tx = types.NewTransaction(
			nonce,
			contractHexAddress,
			ethValue,
			gasLimit,
			gasPrice,
			paddedData.Bytes(),
		)

	} //switch

	newSnap := GasSnapShot{
		Limit:  gasLimit,
		Price:  gasPrice.String(),
		FeeWei: gasPrice.Mul(gasPrice, big.NewInt(int64(gasLimit))).String(),
	}
	ntx := &NTX{
		client:     my.Sender.client,
		tx:         tx,
		privatekey: my.privatekey,
		chainID:    my.chainID,

		nonceCount: nonce,

		from:      strings.ToLower(my.fromHexAddress.Hex()),
		to:        my.contractAddress, //contract-Address
		wei:       wei,
		gasFeeWei: newSnap.FeeWei,

		snap: newSnap,
	}
	//dbg.Yellow(ntx.GasFeeETH())
	return ntx, nil
}

// NTX :
func (my *tNonce) NTX_FixedGAS(paddedData PADBYTES, wei string, gasPair []*big.Int, nAdd ...uint64) (*NTX, error) {
	contractHexAddress := common.HexToAddress(my.contractAddress)
	signerHexAddress := my.fromHexAddress
	// fmt.Println(contractHexAddress)
	// fmt.Println(signerHexAddress)
	// fmt.Println(paddedData.Bytes())

	var ethValue *big.Int
	if !jmath.IsUnderZero(wei) {
		ethValue = new(big.Int)
		ethValue.SetString(jmath.VALUE(wei), 10)
	}

	gasLimit, err := my.Sender.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  signerHexAddress,
		To:    &contractHexAddress,
		Value: ethValue,
		Data:  paddedData.Bytes(),
	})
	if err != nil {
		return nil, dbg.MakeError("tNonce.NTX@EstimateGas", err)
	}
	//gasLimit += EtherClient.TXLimitTokenAddValue

	nonce := my.nonceValue
	if len(nAdd) > 0 {
		nonce += nAdd[0]
	}

	if ethValue == nil {
		ethValue = big.NewInt(0) // in wei (0 eth)
	}

	var left *big.Int
	var right *big.Int
	left = gasPair[0]
	if len(gasPair) > 1 {
		right = gasPair[1]
	} else {
		right = jmath.New(left).ToBigInteger()
	}

	var tx *types.Transaction
	switch my.Sender.txnType {
	case TXN_EIP_1559:

		tx = types.NewTx(&types.DynamicFeeTx{
			ChainID:   my.Sender.chainID,
			Nonce:     nonce,
			GasTipCap: left,
			GasFeeCap: right,
			Gas:       gasLimit,
			To:        &contractHexAddress,
			Value:     ethValue,
			Data:      paddedData.Bytes(),
		})
	default:
		tx = types.NewTransaction(
			nonce,
			contractHexAddress,
			ethValue,
			gasLimit,
			left,
			paddedData.Bytes(),
		)

	} //switch

	newSnap := GasSnapShot{
		Limit:  gasLimit,
		Price:  jmath.VALUE(right),
		FeeWei: jmath.MUL(right, gasLimit),
	}
	ntx := &NTX{
		client:     my.Sender.client,
		tx:         tx,
		privatekey: my.privatekey,
		chainID:    my.chainID,

		nonceCount: nonce,

		from:      strings.ToLower(my.fromHexAddress.Hex()),
		to:        my.contractAddress, //contract-Address
		wei:       wei,
		gasFeeWei: newSnap.FeeWei,

		gasLimit: gasLimit,
		gasRight: right,

		snap: newSnap,
	}
	//dbg.Yellow(ntx.GasFeeETH())
	return ntx, nil
}

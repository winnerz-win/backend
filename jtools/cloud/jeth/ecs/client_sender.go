package ecs

import (
	"context"
	"encoding/hex"
	"fmt"
	"jtools/cc"
	"jtools/cloud/ebcm"
	"jtools/cloud/ebcm/abi"
	"jtools/dbg"
	"jtools/jmath"
	"jtools/jnet/cnet"
	"jtools/unix"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	client "github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type CSender struct {
	client *client.Client
	info   ebcm.Info

	chainID  *big.Int
	txn_type ebcm.TXNTYPE

	/////////////////////////////////////////
	ClientUtil
}

func NewClient(host, key string, cacheId ...interface{}) *CSender {
	connURL := ebcm.GetHostURL(host, key)

	http_client := new(http.Client)
	http_client.Transport = cnet.HttpTransport()
	rpc_http_client, err := rpc.DialHTTPWithClient(connURL, http_client)
	if err != nil {
		cc.RedItalic("Sender.New::DialHTTPWithClient :", err)
		return nil
	}
	session := client.NewClient(rpc_http_client)

	// session, err := client.Dial(connURL)
	// if err != nil {
	// 	cc.RedItalic("Dial :", err)
	// 	return nil
	// }

	info := ebcm.Info{
		Host: host,
		Key:  key,
	}

	if len(cacheId) > 0 {
		info.NetworkID = jmath.NEW(cacheId[0]).BigInt()

	} else {
		network_id, err := session.NetworkID(context.Background())
		if err != nil {
			cc.RedItalic("NetworkID :", err)
			return nil
		}
		info.NetworkID = network_id

	}

	client := &CSender{
		client:  session,
		info:    info,
		chainID: info.NetworkID,
	}
	client.SetTXNTYPE(ebcm.TXN_EIP_1559)

	return client
}
func (my CSender) Client() interface{} { return my.client }
func (my CSender) Host() string        { return my.info.Host }
func (my CSender) Info() ebcm.Info     { return my.info }
func (my CSender) ChainID() *big.Int {
	return my.chainID
}
func (my CSender) TXNTYPE() ebcm.TXNTYPE { return my.txn_type }

func (my *CSender) SetTXNTYPE(v ebcm.TXNTYPE) {
	my.txn_type = v
}
func (my CSender) IsDebug() bool { return my.info.IsDebug }
func (my *CSender) SetDebug(f ...bool) {
	my.info.IsDebug = dbg.IsTrue(f)
}

///////////////////////////////////////////////////////////////////////////////

func (my CSender) NewTransaction(
	nonce uint64,
	to interface{},
	amount interface{},
	gasLimit uint64,
	gasPrice ebcm.GasPrice,
	data []byte,
) ebcm.WrappedTransaction {
	value := jmath.BigInt(amount)
	if my.TXNTYPE() == ebcm.TXN_EIP_1559 {
		toAddress := wrappedAddress(to)

		return WrappedTransaction(
			types.NewTx(
				&types.DynamicFeeTx{
					ChainID:   nil,
					Nonce:     nonce,
					GasTipCap: gasPrice.Tip,
					GasFeeCap: gasPrice.Gas,
					Gas:       gasLimit,
					To:        &toAddress,
					Value:     value,
					Data:      data,
				},
			),
		)
	}

	return WrappedTransaction(
		types.NewTransaction(
			nonce,
			wrappedAddress(to),
			value,
			gasLimit,
			gasPrice.Gas,
			data,
		),
	)
}

func (my CSender) SignTx(
	tx ebcm.WrappedTransaction,
	prv interface{},
) (stx ebcm.WrappedTransaction, err error) {
	ntx := typesTransaction(tx)
	var singer types.Signer
	if ntx.Type() == ebcm.TXN_EIP_1559.Uint8() {
		singer = types.NewLondonSigner(my.chainID)

	} else {
		singer = types.NewEIP155Signer(my.chainID)
	}

	defer func() {
		if e := recover(); e != nil {
			err = dbg.Error(e)
		}
	}()

	_tx, err := types.SignTx(
		typesTransaction(tx),
		singer,
		my.WrappedPrivateKey(prv),
	)
	return WrappedTransaction(_tx), err

	// return types.SignTx(
	// 	wrappedTransaction(tx),
	// 	singer,
	// 	my.WrappedPrivateKey(prv),
	// )
}

func (my CSender) UnmarshalBinary(buf []byte) ebcm.WrappedTransaction {
	tx := types.Transaction{}
	if err := tx.UnmarshalBinary(buf); err != nil {
		cc.RedItalic("[ECS]", err)
		return nil
	}
	return WrappedTransaction(&tx)
}

func (my CSender) SendTransaction(ctx context.Context, tx ebcm.WrappedTransaction) (string, error) {
	stx := typesTransaction(tx)
	err := my.client.SendTransaction(
		ctx,
		stx,
	)
	return stx.Hash().Hex(), err
}

func (my CSender) CheckSendTxHashReceipt(
	tx ebcm.WrappedTransaction,
	limitSec int,
	is_debug ...bool,
) ebcm.CheckSendTxHashReceiptResult {
	hash := my.GetHash(tx)
	ack := ebcm.CheckSendTxHashReceiptResult{
		Hash: hash,
	}

	isDebug := dbg.IsTrue(is_debug)
	log := func(a ...interface{}) {
		if isDebug {
			cc.PurpleItalic(a...)
		}
	}

	for {
		if limitSec <= 0 {
			ack.FailMessage = "time_over"
			ack.IsTimeOver = true
			break
		}
		time.Sleep(time.Second)

		r, _, err := my.TransactionByHash(hash)
		if err != nil {
			limitSec--
			log("receipt wait -", limitSec)
			continue
		}

		if !r.IsReceiptedByHash {
			limitSec--
			log("receipt wait -", limitSec)
			continue
		}

		log("gas_fee :", r.TxFeeETH, " eth")
		ack.GasFeeETH = r.TxFeeETH
		if !r.IsError {
			ack.IsSuccess = true
		} else {
			ack.FailMessage = "tx_fail"
		}
		break
	} //for

	return ack
}

func (my CSender) CheckSendTxHashToNonce(
	from interface{},
	hash string,
	limitSec int,
	is_debug ...bool,
) ebcm.CheckSendTxHashReceiptResult {
	ack := ebcm.CheckSendTxHashReceiptResult{
		Hash: hash,
	}

	ctx := context.Background()
	isDebug := dbg.IsTrue(is_debug)

	log := func(a ...interface{}) {
		if isDebug {
			cc.PurpleItalic(a...)
		}
	}

	pending, err := my.PendingNonceAt(ctx, from)
	if err != nil {
		ack.FailMessage = "pending fail"
		return ack
	}

	is_nonce_end := false
	for {
		if limitSec <= 0 {
			ack.FailMessage = "time_over"
			ack.IsTimeOver = true
			break
		}
		time.Sleep(time.Second)

		if !is_nonce_end {
			nonce, err := my.NonceAt(ctx, from)
			if err != nil {
				limitSec--
				log("nonce wait -", limitSec)
				continue
			}

			log("nonce :", nonce, " , pending :", pending)
			if nonce < pending {
				limitSec--
				log("nonce cmp wait -", limitSec)
				continue
			}

			is_nonce_end = true
		}

		r, _, _ := my.TransactionByHash(hash)
		if !r.IsReceiptedByHash {
			limitSec--
			log("receipt wait -", limitSec)
			continue
		}

		log("gas_fee :", r.TxFeeETH, " eth")
		ack.GasFeeETH = r.TxFeeETH
		if !r.IsError {
			ack.IsSuccess = true
		} else {
			ack.FailMessage = "tx_fail"
		}
		break
	} //for

	return ack
}

func (my CSender) NetworkID(ctx context.Context) (*big.Int, error) {
	return my.client.NetworkID(ctx)
}

func (my CSender) BlockNumber(ctx context.Context) (*big.Int, error) {
	n, err := my.client.BlockNumber(ctx)
	if err != nil {
		return jmath.NEW(0).BigInt(), err
	}
	return jmath.NEW(n).BigInt(), nil
}

func (my CSender) BalanceAt(ctx context.Context, account interface{}, blockNumber *big.Int) (*big.Int, error) {
	return my.client.BalanceAt(
		ctx,
		wrappedAddress(account),
		blockNumber,
	)
}

func (my CSender) SuggestGasPrice(ctx context.Context, is_skip_tip_cap ...bool) (ebcm.GasPrice, error) {
	gasPrice := ebcm.GasPrice{}
	value, err := my.client.SuggestGasPrice(ctx)
	if err != nil {
		return gasPrice, err
	}
	gasPrice.Gas = value

	if my.TXNTYPE() == ebcm.TXN_EIP_1559 {
		if dbg.IsTrue(is_skip_tip_cap) {
			gasPrice.Tip = jmath.BigInt(value)

		} else {
			/*
				MaxFee = ( BaseFee * 2 ) + maxPriorityFee(tip_cap)

				: gas 와 basefee 값이 거의 비슷하므로 나는 (gas * 1.3) + maxPriorityFee 로 할것임.ㅋ
			*/
			tip, err := my.client.SuggestGasTipCap(ctx)
			if err != nil {
				return gasPrice, err
			}

			gasPrice.Gas = jmath.BigInt(
				jmath.DOTCUT(
					jmath.ADD(jmath.MUL(gasPrice.Gas, 1.3), tip),
					0,
				),
			)

			gasPrice.Tip = tip
		}
	}

	return gasPrice, nil
}

func (my CSender) EstimateGas(ctx context.Context, msg ebcm.CallMsg) (uint64, error) {
	to := wrappedAddress(msg.To)
	return my.client.EstimateGas(
		ctx,
		ethereum.CallMsg{
			From:  wrappedAddress(msg.From),
			To:    &to,
			Value: msg.Value,
			Data:  msg.Data,
		},
	)
}

func (my CSender) PendingNonceAt(ctx context.Context, account interface{}) (uint64, error) {
	return my.client.PendingNonceAt(
		ctx,
		wrappedAddress(account),
	)
}

func (my CSender) NonceAt(ctx context.Context, account interface{}) (uint64, error) {
	return my.client.NonceAt(
		ctx,
		wrappedAddress(account),
		nil,
	)
}

func (my CSender) CallContract(from, to string, data []byte) ([]byte, error) {
	fromHexAddress := common.HexToAddress(strings.ToLower(from))
	toHexAddress := common.HexToAddress(strings.ToLower(to))
	param := ethereum.CallMsg{
		From: fromHexAddress,
		To:   &toHexAddress,
		Data: data,
	}
	return my.client.CallContract(
		context.Background(),
		param,
		nil,
	)
}

func (my CSender) _block_by_number(number interface{}) (*types.Block, error) {
	num := jmath.NEW(number).BigInt()

	block, err := my.client.BlockByNumber(context.Background(), num)
	if err != nil {
		cc.RedItalic("blockByNumber[", num, "] :", err)
		return nil, err
	}
	return block, nil
}

func (my CSender) BlockByNumberSimple(number interface{}) *ebcm.BlockByNumberData {

	block, err := my._block_by_number(number)
	if err != nil {
		return nil
	}

	data := &ebcm.BlockByNumberData{
		BlockData: ebcm.BlockData{
			Number:       block.NumberU64(),
			NumberString: fmt.Sprintf("%v", block.NumberU64()),
			Time:         unix.Time(int64(block.Time())),
			Hash:         strings.ToLower(block.Hash().Hex()),
			PreHash:      strings.ToLower(block.ParentHash().Hex()),
			CoinBase:     strings.ToLower(block.Coinbase().Hex()),
			Difficulty:   block.Difficulty().String(),
			GasLimit:     block.GasLimit(),

			GasUsed: block.GasUsed(),
			Nonce:   block.Nonce(),

			Extra:       "0x" + hex.EncodeToString(block.Extra()),
			ReceiptHash: block.ReceiptHash().Hex(),
			Root:        block.Root().Hex(),
			Size:        jmath.VALUE(block.Size()),
			TxHash:      block.TxHash().Hex(),
			BaseFee:     jmath.VALUE(block.BaseFee()),
		},
	}

	return data
}

func checkEIP1559(txType ebcm.TxType, callback func()) bool {
	if txType.Uint16() == ebcm.TXN_EIP_1559.Uint16() {
		callback()
		return true
	}
	return false
}

func (my CSender) BlockByNumber(number interface{}) *ebcm.BlockByNumberData {
	block, err := my._block_by_number(number)
	if err != nil {
		return nil
	}

	data := &ebcm.BlockByNumberData{
		BlockData: ebcm.BlockData{
			Number:       block.NumberU64(),
			NumberString: fmt.Sprintf("%v", block.NumberU64()),
			Time:         unix.Time(int64(block.Time())),
			Hash:         strings.ToLower(block.Hash().Hex()),
			PreHash:      strings.ToLower(block.ParentHash().Hex()),
			CoinBase:     strings.ToLower(block.Coinbase().Hex()),
			Difficulty:   block.Difficulty().String(),
			GasLimit:     block.GasLimit(),

			GasUsed: block.GasUsed(),
			Nonce:   block.Nonce(),

			Extra:       "0x" + hex.EncodeToString(block.Extra()),
			ReceiptHash: block.ReceiptHash().Hex(),
			Root:        block.Root().Hex(),
			Size:        jmath.VALUE(block.Size()),
			TxHash:      block.TxHash().Hex(),
			BaseFee:     jmath.VALUE(block.BaseFee()),
		},
	}

	for _, tx := range block.Transactions() {
		txitem := my.NewTxBlock(tx, data.NumberString)
		txitem.Timestamp = data.Time
		txitem.BaseFee = data.BaseFee
		checkEIP1559(txitem.Type, func() {
			if txitem.GasTipCap != txitem.GasFeeCap {
				/*
					gas = block.Base + tx.GasTipCap(Max Priority)
					gas 가 FeeCap (MAX) 보다 크면 FeeCap이 적용됨.
				*/
				gas := jmath.ADD(txitem.GasTipCap, txitem.BaseFee)
				txitem.Gas = gas
				if jmath.CMP(gas, txitem.GasFeeCap) >= 0 {
					txitem.Gas = txitem.GasFeeCap
				}
				txitem.TxFeeETH = jmath.MUL(ebcm.WeiToETH(txitem.Gas), txitem.GasUsed)
			}
		})

		data.TxList = append(data.TxList, txitem)
	}

	return data
}

func (my CSender) TransactionByHash(hashString string) (ebcm.TransactionBlock, bool, error) {

	hashString = strings.TrimSpace(hashString)

	//ETHKJS
	blockNumber := "0"
	tx, isPending, err := my.client.TransactionByHash(context.Background(), common.HexToHash(hashString))
	if err != nil {
		return ebcm.TransactionBlock{}, false, err
	}

	txitem := my.NewTxBlock(tx, blockNumber)

	if !isPending {
		receipt := my.ReceiptByHash(hashString)
		my.InjectReceipt(&txitem, receipt)
		if receipt.IsNotFound {
			isPending = true
		}

	}

	if !isPending {
		if !checkEIP1559(txitem.Type, func() {
			if txitem.GasTipCap != txitem.GasFeeCap {
				if block, err := my._block_by_number(txitem.BlockNumber); err == nil {
					txitem.BaseFee = jmath.VALUE(block.BaseFee())

					txitem.Timestamp = unix.Time(int64(block.Time()))
				}

				if txitem.BaseFee != "" {
					/*
						gas = block.Base + tx.GasTipCap(Max Priority)
						gas 가 FeeCap (MAX) 보다 크면 FeeCap이 적용됨.
					*/
					gas := jmath.ADD(txitem.GasTipCap, txitem.BaseFee)
					txitem.Gas = gas
					if jmath.CMP(gas, txitem.GasFeeCap) >= 0 {
						txitem.Gas = txitem.GasFeeCap
					}
					txitem.TxFeeETH = jmath.MUL(ebcm.WeiToETH(txitem.Gas), txitem.GasUsed)
				}
			} else {
				if block, err := my._block_by_number(txitem.BlockNumber); err == nil {
					txitem.BaseFee = jmath.VALUE(block.BaseFee())
					txitem.Timestamp = unix.Time(int64(block.Time()))
				}
			}
		}) {
			if block, err := my._block_by_number(txitem.BlockNumber); err == nil {
				txitem.Timestamp = unix.Time(int64(block.Time()))
			}
		}
	}

	return txitem, isPending, nil
}

func (my CSender) ReceiptByHash(hexHash string) ebcm.TxReceipt {
	hash := common.HexToHash(strings.TrimSpace(hexHash))

	re := ebcm.TxReceipt{
		Logs: ebcm.TxLogList{},
	}

	getAddress := func(a common.Address) string {
		return dbg.TrimToLower(a.Hex())
	}
	encodeBytes := func(b []byte) string {
		return dbg.TrimToLower("0x" + hex.EncodeToString(b))
	}

	v, err := my.client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			re.IsNotFound = true
		}
		return re
	}
	//v.Status
	re.TransactionHash = hexHash

	re.BlockHash = v.BlockHash.Hex()
	re.BlockNumber = v.BlockNumber.String()

	re.Bloom = encodeBytes(v.Bloom.Bytes())
	re.Status = v.Status
	re.CumulativeGasUsed = v.CumulativeGasUsed

	re.TransactionIndex = v.TransactionIndex
	re.GasUsed = v.GasUsed
	re.ContractAddress = getAddress(v.ContractAddress)

	for _, l := range v.Logs {

		log := ebcm.TxLog{
			Address:     getAddress(l.Address),
			Data:        ebcm.MakeDataItemList(encodeBytes(l.Data)),
			BlockNumber: l.BlockNumber,
			TxHash:      l.TxHash.Hex(),
			BlockHash:   l.BlockHash.Hex(),
			LogIndex:    l.Index,
			TxIndex:     l.TxIndex,
			Removed:     l.Removed,
		}

		for _, topic := range l.Topics {
			c := dbg.TrimToLower(topic.Hex())
			log.Topics = append(log.Topics, ebcm.Topic(c))
		}

		re.Logs = append(re.Logs, log)
	} //for

	return re
}

func (my CSender) InjectReceipt(tx *ebcm.TransactionBlock, r ebcm.TxReceipt) {

	tx.IsReceptFailByTxInject = !r.Valid()

	tx.BlockNumber = r.BlockNumber
	tx.Number = jmath.Uint64(r.BlockNumber)

	if !r.IsNotFound {
		tx.IsError = r.Status != 1
		tx.IsReceiptedByHash = true
	}
	tx.TxIndex = r.TransactionIndex

	tx.CumulativeGasUsed = r.CumulativeGasUsed

	tx.Logs = r.Logs
	tx.GasUsed = jmath.VALUE(r.GasUsed)

	tx.TxFeeETH = jmath.MUL(ebcm.WeiToETH(tx.Gas), tx.GasUsed)

}

func (my *CSender) NewTxBlock(tx *types.Transaction, no ...string) ebcm.TransactionBlock {

	blockNumber := ""
	if len(no) > 0 {
		blockNumber = no[0]
	}

	from := ""
	if tx.Type() == ebcm.TXN_EIP_1559.Uint8() {
		if msg, err := types.NewLondonSigner(tx.ChainId()).Sender(tx); err == nil {
			from = msg.Hex()
		}

	} else {
		if msg, err := types.NewEIP155Signer(tx.ChainId()).Sender(tx); err == nil {
			from = msg.Hex()
		}
	}

	txInput := strings.ToLower(hex.EncodeToString(tx.Data()))

	txType := uint16(tx.Type())
	txitem := ebcm.TransactionBlock{
		/*
			Legacy  : 0 (0x00)
			EIP1559 : 2 (0x02)
		*/
		Type: ebcm.TxType(txType),

		IsContract:     false,
		ContractMethod: "",
		Hash:           strings.ToLower(tx.Hash().Hex()),
		From:           strings.ToLower(from),
		Nonce:          tx.Nonce(),
		Amount:         tx.Value().String(),

		BlockNumber: blockNumber,
		Number:      jmath.Uint64(blockNumber),

		Logs: ebcm.TxLogList{},
	}

	if tx.GasPrice() != nil {
		if txType == ebcm.TXN_EIP_1559.Uint16() {
			txitem.GasTipCap = jmath.VALUE(tx.GasTipCap())
			txitem.GasFeeCap = jmath.VALUE(tx.GasFeeCap())
		}
		// cc.Cyan(tx.Type())
		// cc.Cyan("Cost :", jmath.VALUE(tx.Cost()))
		// cc.Cyan("GasPrice  :", jmath.VALUE(tx.GasPrice()))  //2500000016
		// cc.Cyan("GasTipCap :", jmath.VALUE(tx.GasTipCap())) //2500000000
		// cc.Cyan("GasFeeCap :", jmath.VALUE(tx.GasFeeCap())) //2500000016
		// v := tx.EffectiveGasTipValue(tx.GasTipCap())
		// cc.Cyan("ccc :", jmath.VALUE(v))
		// v2 := tx.EffectiveGasTipValue(tx.GasFeeCap())
		// cc.Cyan("cc2 :", jmath.VALUE(v2))
		txitem.Limit = tx.Gas()
		txitem.Gas = jmath.VALUE(tx.GasPrice())
		txitem.GasUsed = txitem.Gas //read-only

	}

	toAddress := tx.To()
	if toAddress != nil { //거래의 수신자 주소를 반환
		txitem.To = strings.ToLower(toAddress.Hex())
	} else { //트랜잭션이 계약 생성 인 경우 nil을 반환합니다.
		txitem.ContractMethod = "deploy"
		txitem.To = ""

		txitem.ContractAddress = my.ContractAddressNonce(txitem.From, txitem.Nonce)
	}

	ebcm.CheckMethodERC20(my, txInput, &txitem)

	return txitem
}

func (my CSender) HeaderByNumber(number any) *ebcm.BlockHeader {
	data, err := my.client.HeaderByNumber(context.Background(), jmath.BigInt(number))
	if err != nil {
		return nil
	}

	block_number := jmath.NEW(data.Number)

	block_header := &ebcm.BlockHeader{
		Number:      block_number.Uint64(),
		BlockNumber: block_number.String(),

		ParentHash:       strings.ToLower(data.ParentHash.Hex()),
		RewardCoinBase:   strings.ToLower(data.Coinbase.Hex()),
		StateRoot:        strings.ToLower(data.Root.String()),
		TransactionsRoot: strings.ToLower(data.TxHash.Hex()),
		ReceiptsRoot:     strings.ToLower(data.ReceiptHash.Hex()),

		GasUsed:   jmath.VALUE(data.GasUsed),
		Timestamp: unix.Time(jmath.Int64(data.Time)),

		Hash: strings.ToLower(data.Hash().Hex()),
	}

	return block_header
}

func (my CSender) Call(
	contract string,
	method abi.Method,
	caller string,
	f func(rs abi.RESULT),
	isLogs ...bool,
) error {
	return abi.Call2(
		my,
		contract,
		method,
		caller,
		f,
		isLogs...,
	)
}

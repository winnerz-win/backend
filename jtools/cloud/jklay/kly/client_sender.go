package kly

import (
	"context"
	"encoding/hex"
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

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/networks/rpc"
)

type CSender struct {
	client *client.Client
	info   ebcm.Info

	chainID *big.Int

	/////////////////////////////////////////
	ClientUtil
}

func NewClient(host string, cacheId ...interface{}) *CSender {
	key := ""
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

	return client
}

func (my CSender) Client() interface{} { return my.client }
func (my CSender) Host() string        { return my.info.Host }
func (my CSender) Info() ebcm.Info     { return my.info }
func (my CSender) ChainID() *big.Int {
	return my.chainID
}
func (my CSender) TXNTYPE() ebcm.TXNTYPE { return ebcm.TXN_LEGACY }

func (my *CSender) SetTXNTYPE(v ebcm.TXNTYPE) {
}
func (my CSender) IsDebug() bool { return my.info.IsDebug }
func (my *CSender) SetDebug(f ...bool) {
	my.info.IsDebug = dbg.IsTrue(f)
}

////////////////////////////////////////////////////////////////////////

func (my CSender) NewTransaction(
	nonce uint64,
	to interface{},
	amount interface{},
	gasLimit uint64,
	gasPrice ebcm.GasPrice,
	data []byte,
) ebcm.WrappedTransaction {
	return WrappedTransaction(
		types.NewTransaction(
			nonce,
			wrappedAddress(to),
			jmath.BigInt(amount),
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

	defer func() {
		if e := recover(); e != nil {
			err = dbg.Error(e)
		}
	}()

	_tx, err := types.SignTx(
		typesTransaction(tx),
		types.NewEIP155Signer(my.chainID),
		my.WrappedPrivateKey(prv),
	)
	return WrappedTransaction(_tx), err
}
func (my CSender) UnmarshalBinary(buf []byte) ebcm.WrappedTransaction {
	tx := types.Transaction{}
	if err := tx.UnmarshalBinary(buf); err != nil {
		cc.RedItalic("[KLY]", err)
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
	return my.client.ChainID(ctx)
	//return my.client.NetworkID(ctx)
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
	gasPrice.Tip = jmath.BigInt(value)

	return gasPrice, nil
}

func (my CSender) EstimateGas(ctx context.Context, msg ebcm.CallMsg) (uint64, error) {
	to := wrappedAddress(msg.To)
	return my.client.EstimateGas(
		ctx,
		klaytn.CallMsg{
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
	param := klaytn.CallMsg{
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
			NumberString: dbg.Cat(block.NumberU64()),
			Time:         unix.Time(block.Time().Int64()),
			Hash:         block.Hash().Hex(),
			PreHash:      block.ParentHash().Hex(),
			RewardBase:   block.Rewardbase().Hex(),

			GasUsed: block.GasUsed(),

			Extra:       "0x" + hex.EncodeToString(block.Extra()),
			ReceiptHash: block.ReceiptHash().Hex(),
			Root:        block.Root().Hex(),
			Size:        block.Size().String(),
			TxHash:      block.TxHash().Hex(),

			BlockScore: block.BlockScore().Int64(),
		},
	}

	for _, tx := range block.Transactions() {
		txitem := my.NewTxBlock(tx)
		//checkCustomMethod(&txitem)

		//txitem.finder = my
		txitem.Number = data.BlockData.Number
		txitem.BlockNumber = data.BlockData.NumberString
		txitem.Timestamp = data.Time

		data.TxList = append(data.TxList, txitem)

	}

	return data
}

func (my CSender) BlockByNumber(number interface{}) *ebcm.BlockByNumberData {
	block, err := my._block_by_number(number)
	if err != nil {
		return nil
	}

	data := &ebcm.BlockByNumberData{
		BlockData: ebcm.BlockData{
			Number:       block.NumberU64(),
			NumberString: dbg.Cat(block.NumberU64()),
			Time:         unix.Time(block.Time().Int64()),
			Hash:         block.Hash().Hex(),
			PreHash:      block.ParentHash().Hex(),
			RewardBase:   block.Rewardbase().Hex(),

			GasUsed: block.GasUsed(),

			Extra:       "0x" + hex.EncodeToString(block.Extra()),
			ReceiptHash: block.ReceiptHash().Hex(),
			Root:        block.Root().Hex(),
			Size:        block.Size().String(),
			TxHash:      block.TxHash().Hex(),

			BlockScore: block.BlockScore().Int64(),
		},
	}

	for _, tx := range block.Transactions() {
		txitem := my.NewTxBlock(tx)
		//checkCustomMethod(&txitem)

		//txitem.finder = my
		txitem.Number = data.BlockData.Number
		txitem.BlockNumber = data.BlockData.NumberString
		txitem.Timestamp = data.Time

		data.TxList = append(data.TxList, txitem)

	}

	return data
}

func (my CSender) TransactionByHash(hashString string) (ebcm.TransactionBlock, bool, error) {

	hashString = strings.TrimSpace(hashString)

	tx, isPending, err := my.client.TransactionByHash(context.Background(), common.HexToHash(hashString))
	if err != nil {
		return ebcm.TransactionBlock{}, false, err
	}

	txitem := my.NewTxBlock(tx)

	if !isPending {
		receipt := my.ReceiptByHash(hashString)
		my.InjectReceipt(&txitem, receipt)
		if receipt.IsNotFound {
			isPending = true
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

	v, err := my.client.TransactionReceiptRpcOutput(context.Background(), hash)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			re.IsNotFound = true
		}
		return re
	}
	_ = v
	//v.Status
	re.TransactionHash = hexHash
	re.BlockNumber = jmath.VALUE(v["blockNumber"])
	re.TransactionIndex = uint(jmath.Uint64(v["transactionIndex"]))

	re.BlockHash = dbg.TrimToLower(v["blockHash"].(string))
	re.Gas = jmath.VALUE(v["gas"])
	re.GasPrice = jmath.VALUE(v["gasPrice"])
	re.GasUsed = jmath.Uint64(v["gasUsed"])

	re.From = dbg.TrimToLower(v["from"].(string))
	re.Bloom = dbg.TrimToLower(v["logsBloom"].(string))
	re.Nonce = jmath.Uint64(v["nonce"])

	re.SenderTxHash = dbg.TrimToLower(v["senderTxHash"].(string))
	re.Status = jmath.Uint64(v["status"])
	toAddress := ""
	if to, do := v["to"]; do {
		if to != nil {
			switch v := to.(type) {
			case string:
				toAddress = v
			default:
				toAddress = dbg.Cat(v)
			}
		}
	}

	re.To = dbg.TrimToLower(toAddress)
	re.TransactionIndex = uint(jmath.Uint64(v["transactionIndex"]))

	tx_type := ebcm.TxType(uint16(jmath.Int(v["typeInt"])))
	re.Type = &tx_type

	re.Amount = jmath.VALUE(v["value"])

	if v["contractAddress"] != nil {
		if ca, do := v["contractAddress"].(string); do {
			re.ContractAddress = dbg.TrimToLower(ca)
		} else {
			re.ContractAddress = dbg.TrimToLower(dbg.Cat(v["contractAddress"]))
		}
	}

	logs := []*types.Log{}
	dbg.ParseStruct(v["logs"], &logs)

	re.Logs = ebcm.TxLogList{}
	for _, l := range logs {
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
	tx.Logs = r.Logs
	tx.GasUsed = jmath.VALUE(r.GasUsed)

	tx.TxFeeKLAY = jmath.MUL(ebcm.WeiToETH(tx.GasPrice), tx.GasUsed)

	if tx.FeeRatio != nil {
		if !tx.FeeRatio.IsValid() {
			if strings.Contains(types.TxType(tx.Type).String(), "Delegated") {
				tx.TxFeeBySenderKLAY = "0"
				tx.TxFeeByFeePayerKLAY = tx.TxFeeKLAY
			} else {
				tx.TxFeeBySenderKLAY = tx.TxFeeKLAY
				tx.TxFeeByFeePayerKLAY = "0"
			}

		} else {
			fee := jmath.NEW(ebcm.ETHToWei(tx.TxFeeKLAY)).BigInt()
			if fee != nil {
				fee_payer_val, fee_sender_val := types.CalcFeeWithRatio(
					types.FeeRatio(*tx.FeeRatio),
					fee,
				)
				tx.TxFeeByFeePayerKLAY = ebcm.WeiToETH(jmath.VALUE(fee_payer_val))
				tx.TxFeeBySenderKLAY = ebcm.WeiToETH(jmath.VALUE(fee_sender_val))
			}
		}
	}

}

func makeFromTo(tx *types.Transaction, from, to func(string), contractCreation func()) {
	str := dbg.Cat(tx)
	sl := strings.Split(str, "\n")

	isContractCreation := false
	for _, v := range sl {
		//cc.PrintRed(v)
		if strings.Contains(v, "From:") {
			v = strings.ReplaceAll(v, "From:", "")
			name := dbg.TrimToLower(v)
			if !strings.HasPrefix(name, "0x") {
				name = "0x" + name
			}

			if from != nil {
				from(name)
			} else {
				cc.PurpleItalic("from :", name)
			}
		} else if strings.Contains(v, "To:") {
			v = strings.ReplaceAll(v, "To:", "")
			if strings.Contains(v, "[contract creation]") {
				isContractCreation = true
			}
			name := dbg.TrimToLower(v)
			if !strings.HasPrefix(name, "0x") {
				name = "0x" + name
			}

			if to != nil {
				to(name)
			} else {
				cc.PurpleItalic("to :", name)
			}

			break
		}
	} //for

	if isContractCreation {
		if contractCreation != nil {
			contractCreation()
		}
	}
}

func (my *CSender) NewTxBlock(tx *types.Transaction) ebcm.TransactionBlock {

	feePayer, _ := tx.FeePayer()
	feeRatio, _ := tx.FeeRatio()
	fromAddress := ""
	from, err := tx.From()
	if err != nil {
		//cc.PrintPurple("from :", err)
	} else {
		fromAddress = from.Hex()
	}

	txInput := strings.ToLower(hex.EncodeToString(tx.Data()))

	toAddress := ""
	if to := tx.To(); to != nil {
		toAddress = dbg.TrimToLower(to.Hex())
	}

	fee_ratio := ebcm.Klay_FeeRatio(uint8(feeRatio))
	role_type := ebcm.Klay_RoleType(int(tx.GetRoleTypeForValidation()))
	txitem := ebcm.TransactionBlock{
		ContractMethod: "",

		Hash:     dbg.TrimToLower(tx.Hash().Hex()),
		Type:     ebcm.TxType(uint16(tx.Type())),
		Fee:      tx.Fee().String(),
		FeePayer: dbg.TrimToLower(feePayer.Hex()),
		FeeRatio: &fee_ratio,
		Amount:   jmath.VALUE(tx.Value()),
		Cost:     jmath.VALUE(tx.Cost()),
		Gas:      jmath.VALUE(tx.Gas()),
		GasPrice: jmath.VALUE(tx.GasPrice()),
		Nonce:    tx.Nonce(),
		From:     dbg.TrimToLower(fromAddress),
		To:       toAddress,
		Limit:    jmath.Uint64(tx.GetTxInternalData().GetGasLimit()),

		RoleTypeForValidation: &role_type,
		ValidatedIntrinsicGas: tx.ValidatedIntrinsicGas(),
		ValidatedSender:       dbg.TrimToLower(tx.ValidatedSender().Hex()),

		Logs: ebcm.TxLogList{},

		//internalData: tx.GetTxInternalData(),
	}

	isContractCreate := false
	lazyContractCreateFunc := func() {
		isContractCreate = true
	}
	makeFromTo(tx,
		func(s string) { txitem.From = s },
		func(s string) { txitem.To = s },
		lazyContractCreateFunc, //contractCreation
	)

	ebcm.CheckMethodERC20(my, txInput, &txitem)
	if isContractCreate {
		if txitem.IsContract {
			txitem.IsContract = true
			if txitem.ContractMethod != "deploy" {
				txitem.ContractMethod = "deploy"
				txitem.ContractAddress = my.ContractAddressNonce(txitem.From, txitem.Nonce)
			}
		}
	}

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
		RewardCoinBase:   strings.ToLower(data.Rewardbase.Hex()),
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

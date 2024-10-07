package ecsx

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"txscheduler/brix/tools/cloudx/ebcmx"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
	"txscheduler/brix/tools/mms"

	"github.com/ethereum/go-ethereum/core/types"
)

type cBlockByNumberData struct {
	BlockData `bson:",inline" json:",inline"`
	TxList    TransactionBlockList `bson:"txlist" json:"txlist"`
}

func (my cBlockByNumberData) GetTxCount() int { return len(my.TxList) }
func (my cBlockByNumberData) String() string  { return dbg.ToJSONString(my) }

type BlockData struct {
	Number       uint64 `bson:"number" json:"number"`
	NumberString string `bson:"numberString" json:"numberString"`
	Time         int64  `bson:"time" json:"time"`
	Hash         string `bson:"hash" json:"hash"`
	PreHash      string `bson:"pre_hash" json:"pre_hash"`
	CoinBase     string `bson:"coin_base" json:"coin_base"`
	Difficulty   string `bson:"difficulty" json:"difficulty"`
	GasLimit     uint64 `bson:"gas_limit" json:"gas_limit"`
	GasUsed      uint64 `bson:"gas_used" json:"gas_used"`
	Nonce        uint64 `bson:"nonce" json:"nonce"`
	Extra        string `bson:"extra" json:"extra"`
	ReceiptHash  string `bson:"receipt_hash" json:"receipt_hash"`
	Root         string `bson:"root" json:"root"`
	Size         string `bson:"size" json:"size"`
	TxHash       string `bson:"tx_hash" json:"tx_hash"`
	BaseFee      string `bson:"baseFee" json:"baseFee"`
}

func (my BlockData) TimeMMS() mms.MMS { return mms.MMS(my.Time * 1000) }
func (my BlockData) String() string   { return dbg.ToJSONString(my) }

func (my Sender) BlockByNumberOrigin(number string, isLog ...bool) types.Transactions {

	num := big.NewInt(0)
	num, _ = num.SetString(number, 10)

	block, err := my.client.BlockByNumber(context.Background(), num)
	if err != nil {
		dbg.RedItalic("blockByNumber-fail :", err)
		return nil
	}
	return block.Transactions()
}

func (my *Sender) getBlockDataByNumber(number string) *BlockData {
	defer func() {
		if e := recover(); e != nil {
			dbg.PrintStack(e)
		}
	}()
	num := big.NewInt(0)
	num, _ = num.SetString(number, 10)

	block, err := my.client.BlockByNumber(context.Background(), num)
	if err != nil {
		return nil
	}
	return &BlockData{
		Number:       block.NumberU64(),
		NumberString: fmt.Sprintf("%v", block.NumberU64()),
		Time:         int64(block.Time()),
		Hash:         strings.ToLower(block.Hash().Hex()),
		PreHash:      strings.ToLower(block.ParentHash().Hex()),
		CoinBase:     strings.ToLower(block.Coinbase().Hex()),
		Difficulty:   block.Difficulty().String(),
		GasLimit:     block.GasLimit(),
		GasUsed:      block.GasUsed(),
		Nonce:        block.Nonce(),
		Extra:        "0x" + hex.EncodeToString(block.Extra()),
		ReceiptHash:  block.ReceiptHash().Hex(),
		Root:         block.Root().Hex(),
		Size:         block.Size().String(),
		TxHash:       block.TxHash().Hex(),
		BaseFee:      jmath.VALUE(block.BaseFee()),
	}
}

// BlockByNumber :  Tx- Fail 난것도 검색 된다. ( fail 여부 알수 없음)
func (my *Sender) BlockByNumber(number string, isLog ...bool) *cBlockByNumberData {
	defer func() {
		if e := recover(); e != nil {
			dbg.PrintStack(e)
		}
	}()
	num := big.NewInt(0)
	num, _ = num.SetString(number, 10)

	block, err := my.client.BlockByNumber(context.Background(), num)
	if err != nil {
		if dbg.IsTrue2(isLog...) {
			dbg.RedItalic("blockByNumber-fail :", err)
		}
		return nil
	}

	data := cBlockByNumberData{
		BlockData: BlockData{
			Number:       block.NumberU64(),
			NumberString: fmt.Sprintf("%v", block.NumberU64()),
			Time:         int64(block.Time()),
			Hash:         strings.ToLower(block.Hash().Hex()),
			PreHash:      strings.ToLower(block.ParentHash().Hex()),
			CoinBase:     strings.ToLower(block.Coinbase().Hex()),
			Difficulty:   block.Difficulty().String(),
			GasLimit:     block.GasLimit(),
			GasUsed:      block.GasUsed(),
			Nonce:        block.Nonce(),
			Extra:        "0x" + hex.EncodeToString(block.Extra()),
			ReceiptHash:  block.ReceiptHash().Hex(),
			Root:         block.Root().Hex(),
			Size:         block.Size().String(),
			TxHash:       block.TxHash().Hex(),
			BaseFee:      jmath.VALUE(block.BaseFee()),
		},
	}
	for _, tx := range block.Transactions() {
		txitem := NewTxBlock(my, tx, data.NumberString)
		my.checkCustomMethod(&txitem)

		txitem.finder = my
		txitem.BlockNumber = data.NumberString
		txitem.Timestamp = mms.MMS(data.Time)

		txitem.BaseFee = data.BaseFee

		checkEIP1559(txitem.Type, func() {
			if txitem.GasTipCap != txitem.GasFeeCap {
				txitem.Gas = jmath.ADD(txitem.GasTipCap, txitem.BaseFee)
				txitem.GasPriceETH = jmath.MUL(WeiToETH(txitem.Gas), txitem.GasUsed)
			}
		})

		// receipt := my.Receipt(txitem.Hash)
		// if receipt.Result() == TxResultFail {
		// 	txitem.IsError = true
		// }

		data.TxList = append(data.TxList, txitem)
	} //for

	if dbg.IsTrue2(isLog...) {
		fmt.Println(dbg.ToJSONString(data))
	}

	return &data
}

func (my *Sender) ebcm_BlockByNumber(number string) *ebcmx.BlockByNumberData {
	data := my.BlockByNumber(number, false)
	if data == nil {
		return nil
	}

	edata := ebcmx.BlockByNumberData{}
	dbg.ChangeStruct(data, &edata)
	// for i := range edata.TxList {
	// 	edata.TxList[i].Finder = my
	// }

	return &edata
}

// ToContractList :
func (my cBlockByNumberData) ToContractList(contract string, method ...string) TransactionBlockList {
	contract = strings.TrimSpace(strings.ToLower(contract))
	txlist := TransactionBlockList{}

	cname := ""
	if len(method) > 0 {
		cname = method[0]
	}

	for _, tx := range my.TxList {
		tx.BlockNumber = my.NumberString
		tx.Timestamp = mms.MMS(my.Time)

		if tx.To == contract {
			if cname != "" {
				if tx.ContractMethod == cname {
					txlist = append(txlist, tx)
				}
			} else {
				txlist = append(txlist, tx)
			}

		}
	} //for
	return txlist
}

// ToAddressList :
func (my cBlockByNumberData) ToAddressList(address string, method ...string) TransactionBlockList {
	txlist := TransactionBlockList{}
	address = strings.TrimSpace(strings.ToLower(address))

	cname := ""
	if len(method) > 0 {
		cname = method[0]
	}

	for _, tx := range my.TxList {
		if cname != "" {
			if tx.IsContract == true {
				if tx.ContractMethod != cname {
					continue
				}
			}
		}
		tx.BlockNumber = my.NumberString
		tx.Timestamp = mms.MMS(my.Time)

		if tx.To == address {
			txlist = append(txlist, tx)
		} else if tx.From == address {
			txlist = append(txlist, tx)
		} else if tx.ContractAddress == address {
			txlist = append(txlist, tx)
		}
	} //for
	return txlist
}

func (my cBlockByNumberData) ToAddressListByMethod(address string, method cMethodIDData) TransactionBlockList {
	return my.ToAddressList(address, method.FuncName)
}

// GetTransferList : transfer - to addresslist
func (my cBlockByNumberData) GetTransferList() TransactionBlockList {
	transfer := MethodTransfer.FuncName
	txlist := TransactionBlockList{}
	for _, tx := range my.TxList {
		if tx.IsContract && tx.ContractMethod != transfer {
			continue
		}
		tx.BlockNumber = my.NumberString
		tx.Timestamp = mms.MMS(my.Time)
		txlist = append(txlist, tx)
	} //for
	return txlist
}

///////////////////////////////////////////////////////////////////////////////////////////////

// BlockByNumberReceipt :  FAIL 된것들은 검색 안됨..;;
func (my Sender) BlockByNumberReceipt(number string) *cBlockByNumberData {
	blockdata := my.BlockByNumber(number)
	if blockdata != nil && blockdata.GetTxCount() > 0 {
		for i, tx := range blockdata.TxList {
			dbg.D(i, tx)
			if receipt := my.Receipt(tx.Hash); receipt == nil {
				blockdata.TxList[i].IsError = true
			} else {

				switch receipt.Result(tx.IsContract) {
				case TxResultOK:
					blockdata.TxList[i].IsError = false
				case TxResultPending:
					blockdata.TxList[i].IsPending = true

				case TxResultNoData:
					blockdata.TxList[i].IsError = false

				case TxResultFail:
					blockdata.TxList[i].IsError = true
				default:
					blockdata.TxList[i].IsError = true

				} //switch

			}
		} //for
	}

	return blockdata
}

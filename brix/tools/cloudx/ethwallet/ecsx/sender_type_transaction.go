package ecsx

import (
	"encoding/hex"
	"errors"
	"strings"

	"txscheduler/brix/tools/cloudx/ethwallet/EtherScanAPI"

	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/mms"

	"txscheduler/brix/tools/jmath"

	"github.com/ethereum/go-ethereum/core/types"
)

// TransactionBlock :
type TransactionBlock struct {
	finder          *Sender `bson:"-" json:"-"`
	IsContract      bool    `bson:"is_contract" json:"is_contract"`
	ContractAddress string  `bson:"contract_address,omitempty" json:"contract_address,omitempty"`
	ContractMethod  string  `bson:"contract_method,omitempty" json:"contract_method,omitempty"`
	Hash            string  `bson:"hash" json:"hash"`
	From            string  `bson:"from" json:"from"`
	Nonce           uint64  `bson:"nonce" json:"nonce"`
	To              string  `bson:"to" json:"to"`
	Amount          string  `bson:"amount" json:"amount"`
	TxIndex         uint    `bson:"tx_index" json:"tx_index"` //TransactionByHash , ErrorCheck

	Number            uint64  `bson:"number" json:"number"`
	BlockNumber       string  `bson:"blockNumber,omitempty" json:"blockNumber,omitempty"`
	Confirmations     string  `bson:"confirmations,omitempty" json:"confirmations,omitempty"`
	Timestamp         mms.MMS `bson:"timestamp,omitempty" json:"timestamp,omitempty"` //unix-time
	IsError           bool    `bson:"is_error,omitempty" json:"is_error,omitempty"`
	IsReceiptedByHash bool    `bson:"is_receipted_by_hash,omitempty" json:"is_receipted_by_hash,omitempty"`
	IsPending         bool    `bson:"is_pending,omitempty" json:"is_pending,omitempty"`
	Symbol            string  `bson:"symbol,omitempty" json:"symbol,omitempty"`
	Decimals          string  `bson:"decimals,omitempty" json:"decimals,omitempty"`
	//Input string `json:"input,omitempty"`
	CustomInput      string              `bson:"custom_input,omitempty" json:"custom_input,omitempty"`
	CustomInputParse CustomInputParseMap `bson:"custom_input_data,omitempty" json:"custom_input_data,omitempty"`

	IsInternal bool   `bson:"is_internal,omitempty" json:"is_internal,omitempty"`
	Type       string `bson:"type,omitempty" json:"type,omitempty"`
	Gas        string `bson:"gas,omitempty" json:"gas,omitempty"`             // userGasPrice
	GasTipCap  string `bson:"gasTipCap,omitempty" json:"gasTipCap,omitempty"` //min
	GasFeeCap  string `bson:"gasFeeCap,omitempty" json:"gasFeeCap,omitempty"` //max
	BaseFee    string `bson:"baseFee,omitempty" json:"baseFee,omitempty"`

	GasUsed           string `bson:"gasUsed,omitempty" json:"gasUsed,omitempty"`                     //
	Limit             uint64 `bson:"gas_limit,omitempty" json:"gas_limit,omitempty"`                 // GasLimit
	GasPriceETH       string `bson:"gas_price_eth,omitempty" json:"gas_price_eth,omitempty"`         // gasUsed * limit
	CumulativeGasUsed uint64 `bson:"cumulativeGasUsed,omitempty" json:"cumulativeGasUsed,omitempty"` //cumulativeGasUsed(apply of limit)
	TradeID           string `bson:"tradeId,omitempty" json:"tradeId,omitempty"`
	ErrCode           string `bson:"errCode,omitempty" json:"errCode,omitempty"`

	Logs TxLogList `bson:"logs" json:"logs"`
}

type CustomInputParseMap map[string]interface{}

func (my CustomInputParseMap) IsKey(key string) bool {
	if my == nil {
		return false
	}
	_, do := my[key]
	return do
}

func (my CustomInputParseMap) Get(key string) interface{} {
	return my[key]
}

func (my CustomInputParseMap) GetString(key string) string {
	if my == nil {
		return ""
	}
	return dbg.Cat(my[key])
}
func (my CustomInputParseMap) Address(key string) string {
	return dbg.TrimToLower(my.GetString(key))
}
func (my CustomInputParseMap) Number(key string) string {
	if my == nil {
		return "0"
	}
	return my.GetString(key)
}

func (my CustomInputParseMap) Int64(key string) int64 {
	if my == nil {
		return 0
	}
	return jmath.Int64(my[key])
}
func (my CustomInputParseMap) Int(key string) int {
	return int(my.Int64(key))
}

func (my CustomInputParseMap) UInt64(key string) uint64 {
	if my == nil {
		return 0
	}
	return jmath.Uint64(my[key])
}
func (my CustomInputParseMap) UInt(key string) uint {
	return uint(my.Int64(key))
}
func (my CustomInputParseMap) Bool(key string) bool {
	if v, do := my[key].(bool); do {
		return v
	}
	return false
}

func (my CustomInputParseMap) Parse(key string, p interface{}) error {
	if my == nil {
		return errors.New("CustomInputParseMap is nil")
	}
	return dbg.ChangeStruct(my[key], p)
}
func (my CustomInputParseMap) Value(key string) interface{} {
	if my == nil {
		return nil
	}
	return my[key]
}

// func (my CustomInputParseMap) Slice(key string, sp interface{}) error {
// 	r := reflect.TypeOf(sp)
// 	if r.Kind() != reflect.Ptr {
// 		return errors.New("Must sp is slice pointer type.")
// 	}
// 	if !r.AssignableTo(reflect.TypeOf(sp)) {
// 		return errors.New("MissMatch Type sp")
// 	}

// 	return nil
// }

func (my *TransactionBlock) NewCustomInputParse() CustomInputParseMap {
	my.CustomInputParse = CustomInputParseMap{}
	return my.CustomInputParse
}
func (my *TransactionBlock) SetCustomInput(key string, val interface{}) {
	if my.CustomInputParse == nil {
		my.NewCustomInputParse()
	}
	my.CustomInputParse[key] = val
}

// NewTxBlockInternal :
func NewTxBlockInternal(tx EtherScanAPI.InternalTransaction) TransactionBlock {
	item := TransactionBlock{
		IsInternal:      true,
		IsContract:      tx.IsContract,
		ContractAddress: tx.ContractAddress,
		Hash:            tx.Hash,
		From:            tx.From,
		To:              tx.To,
		Amount:          tx.Value,
		CustomInput:     tx.Input,
		IsError:         tx.IsError,
		Type:            tx.Type,
		Gas:             tx.Gas,
		GasUsed:         tx.GasUsed,
		TradeID:         tx.TradeID,
		ErrCode:         tx.ErrCode,
	}
	return item
}

// CustomInputData :
func (my TransactionBlock) CustomInputData() CustomInput {
	return newCustomInput(my.CustomInput)
}

// CustimInputParseTag : key check
func (my TransactionBlock) CustomInputParseKey(tag string) bool {
	if my.CustomInputParse == nil {
		return false
	}
	_, do := my.CustomInputParse[tag]
	return do
}

// CustomInputParseKeyVal :
func (my TransactionBlock) CustomInputParseKeyVal(key string, val interface{}) bool {
	if my.CustomInputParse == nil {
		return false
	}
	if v, do := my.CustomInputParse[key]; do {
		if v == val {
			return true
		}
	}
	return false
}

// String :
func (my TransactionBlock) String() string {
	return dbg.ToJSONString(my)
}

// GetConfirmCount :
func (my TransactionBlock) GetConfirmCount(no interface{}) string {
	return jmath.SUB(no, my.BlockNumber)
}

// ErrorCheck :
func (my *TransactionBlock) ErrorCheck(sender *Sender) {
	receipt := sender.Receipt(my.Hash)
	my.TxIndex = receipt.TransactionIndex
	my.IsError = false
	my.IsReceiptedByHash = false
	switch receipt.Result(my.IsContract) {
	case TxResultFail:
		my.IsReceiptedByHash = true
		if my.IsContract {
			my.IsError = true
		}
	case TxResultOK:
		my.IsReceiptedByHash = true
	}
}

// TransactionBlockList :
type TransactionBlockList []TransactionBlock

// NewTxBlock :
func NewTxBlock(my *Sender, tx *types.Transaction, no ...string) TransactionBlock {

	// dbg.Red("gasLimit :", tx.Gas())
	// dbg.Red("gasPrice :", tx.GasPrice())
	// dbg.Red("cost :", tx.Cost())

	blockNumber := ""
	if len(no) > 0 {
		blockNumber = no[0]
	}

	// if msg, err := types.NewEIP2930Signer(tx.ChainId()).Sender(tx); err == nil {
	// 	dbg.Yellow("from(EIP2930)", msg.Hex())
	// } else {
	// 	dbg.Red("eip2930 :", err)
	// }
	// if msg, err := types.NewLondonSigner(tx.ChainId()).Sender(tx); err == nil {
	// 	dbg.Yellow("from(London)", msg.Hex())
	// } else {
	// 	dbg.Red("London :", err)
	// }

	from := ""

	if tx.Type() == 2 {
		if msg, err := types.NewLondonSigner(tx.ChainId()).Sender(tx); err == nil {
			from = msg.Hex()
		}

	} else {
		if msg, err := types.NewEIP155Signer(tx.ChainId()).Sender(tx); err == nil {
			from = msg.Hex()
		}
	}
	// if msg, err := types.NewEIP155Signer(tx.ChainId()).Sender(tx); err == nil {
	// 	from = msg.Hex()
	// } else {
	// 	//dbg.RedItalic("[ecsx.NewTxBlock.go] fromAddress :", err)
	// 	if msg, err := types.NewLondonSigner(tx.ChainId()).Sender(tx); err == nil {
	// 		from = msg.Hex()
	// 	}

	// }
	//dbg.Yellow("TxType :", tx.Type())

	txInput := strings.ToLower(hex.EncodeToString(tx.Data()))
	txitem := TransactionBlock{
		/*
			Legacy  : 0 (0x00)
			EIP1559 : 2 (0x02)
		*/
		Type: jmath.VALUE(tx.Type()),

		IsContract:     false,
		ContractMethod: "",
		Hash:           strings.ToLower(tx.Hash().Hex()),
		From:           strings.ToLower(from),
		Nonce:          tx.Nonce(),
		Amount:         tx.Value().String(),
		//Gas:            WeiToETH(fmt.Sprint(tx.GasPrice())),

		BlockNumber: blockNumber,
		//Number:      jmath.Uint64(blockNumber),
		Logs: TxLogList{},
	}
	txitem.Number = jmath.Uint64(txitem.BlockNumber)

	if tx.GasPrice() != nil {
		if tx.Type() == 2 {
			txitem.GasTipCap = jmath.VALUE(tx.GasTipCap())
			txitem.GasFeeCap = jmath.VALUE(tx.GasFeeCap())
		}
		// dbg.Cyan(tx.Type())
		// dbg.Cyan("Cost :", jmath.VALUE(tx.Cost()))
		// dbg.Cyan("GasPrice  :", jmath.VALUE(tx.GasPrice()))  //2500000016
		// dbg.Cyan("GasTipCap :", jmath.VALUE(tx.GasTipCap())) //2500000000
		// dbg.Cyan("GasFeeCap :", jmath.VALUE(tx.GasFeeCap())) //2500000016
		// v := tx.EffectiveGasTipValue(tx.GasTipCap())
		// dbg.Cyan("ccc :", jmath.VALUE(v))
		// v2 := tx.EffectiveGasTipValue(tx.GasFeeCap())
		// dbg.Cyan("cc2 :", jmath.VALUE(v2))
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
	}

	CheckMethodERC20(txInput, &txitem)
	return txitem
}

// FindToAddress :
func (my TransactionBlockList) FindToAddress(address string) TransactionBlockList {
	txlist := TransactionBlockList{}
	address = dbg.TrimToLower(address)
	for _, tx := range my {
		if tx.To == address {
			txlist = append(txlist, tx)
		}
	} //for
	return txlist
}

// FindContract :
func (my TransactionBlockList) FindContract(contract string) TransactionBlockList {
	txlist := TransactionBlockList{}
	contract = dbg.TrimToLower(contract)
	for _, tx := range my {
		if tx.ContractAddress == contract {
			txlist = append(txlist, tx)
		}
	} //for
	return txlist
}

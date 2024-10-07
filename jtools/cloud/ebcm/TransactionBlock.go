package ebcm

import (
	"errors"
	"jtools/cc"
	"jtools/database/mongo"
	"jtools/dbg"
	"jtools/jmath"
	"jtools/unix"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionBlock struct {
	//Finder *Sender `bson:"-" json:"-"`

	IsContract      bool   `bson:"is_contract" json:"is_contract"`
	ContractAddress string `bson:"contract_address,omitempty" json:"contract_address,omitempty"`
	ContractMethod  string `bson:"contract_method,omitempty" json:"contract_method,omitempty"`
	Hash            string `bson:"hash" json:"hash"`
	From            string `bson:"from" json:"from"`
	Nonce           uint64 `bson:"nonce" json:"nonce"`
	To              string `bson:"to" json:"to"`
	Amount          string `bson:"amount" json:"amount"`
	TxIndex         uint   `bson:"tx_index" json:"tx_index"` //TransactionByHash , ErrorCheck

	Number            uint64    `bson:"number" json:"number"`
	BlockNumber       string    `bson:"blockNumber,omitempty" json:"blockNumber,omitempty"`
	Confirmations     string    `bson:"confirmations,omitempty" json:"confirmations,omitempty"`
	Timestamp         unix.Time `bson:"timestamp,omitempty" json:"timestamp,omitempty"` //unix-time
	IsError           bool      `bson:"is_error,omitempty" json:"is_error,omitempty"`
	IsReceiptedByHash bool      `bson:"is_receipted_by_hash,omitempty" json:"is_receipted_by_hash,omitempty"`
	IsPending         bool      `bson:"is_pending,omitempty" json:"is_pending,omitempty"`
	Symbol            string    `bson:"symbol,omitempty" json:"symbol,omitempty"`
	Decimals          string    `bson:"decimals,omitempty" json:"decimals,omitempty"`
	//Input string `json:"input,omitempty"`
	CustomInput      string              `bson:"custom_input,omitempty" json:"custom_input,omitempty"`
	CustomInputParse CustomInputParseMap `bson:"custom_input_data,omitempty" json:"custom_input_data,omitempty"`

	IsInternal bool   `bson:"is_internal,omitempty" json:"is_internal,omitempty"`
	Type       TxType `bson:"type" json:"type"`
	Gas        string `bson:"gas,omitempty" json:"gas,omitempty"`
	GasTipCap  string `bson:"gasTipCap,omitempty" json:"gasTipCap,omitempty"` //eth min (Max Priority)
	GasFeeCap  string `bson:"gasFeeCap,omitempty" json:"gasFeeCap,omitempty"` //eth max (Max)                           // userGasPrice
	BaseFee    string `bson:"baseFee,omitempty" json:"baseFee,omitempty"`     //eth eip1559
	GasUsed    string `bson:"gasUsed,omitempty" json:"gasUsed,omitempty"`     // limit (GasUsed == Limit)
	Limit      uint64 `bson:"gas_limit,omitempty" json:"gas_limit,omitempty"` // GasLimit
	//TxFeeETH       string `bson:"gas_price_eth,omitempty" json:"gas_price_eth,omitempty"`         //tx_fee        // gasUsed * limit
	TxFeeETH          string `bson:"tx_fee_eth,omitempty" json:"tx_fee_eth,omitempty"`               //tx_fee        // gasUsed * limit
	CumulativeGasUsed uint64 `bson:"cumulativeGasUsed,omitempty" json:"cumulativeGasUsed,omitempty"` //cumulativeGasUsed(apply of limit)
	TradeID           string `bson:"tradeId,omitempty" json:"tradeId,omitempty"`
	ErrCode           string `bson:"errCode,omitempty" json:"errCode,omitempty"`

	//////////// KLAY

	FeeRatio *Klay_FeeRatio `bson:"feeRatio,omitempty" json:"feeRatio,omitempty"` //klay
	Fee      string         `bson:"fee,omitempty" json:"fee,omitempty"`           //klay
	FeePayer string         `bson:"feePayer,omitempty" json:"feePayer,omitempty"` //klay
	GasPrice string         `bson:"gasPrice,omitempty" json:"gasPrice,omitempty"` //klay --- Gas Price (PEB)
	Cost     string         `bson:"cost,omitempty" json:"cost,omitempty"`         //klay

	EffectiveGasPrice string `bson:"effective_gas_price,omitempty" json:"effective_gas_price,omitempty"` //klay

	RoleTypeForValidation *Klay_RoleType `bson:"roleTypeForValidation,omitempty" json:"roleTypeForValidation,omitempty"` //klay
	ValidatedIntrinsicGas uint64         `bson:"validatedIntrinsicGas,omitempty" json:"validatedIntrinsicGas,omitempty"` //klay
	ValidatedSender       string         `bson:"validatedSender,omitempty" json:"validatedSender,omitempty"`             //klay

	TxFeeKLAY           string `bson:"tx_fee_klay,omitempty" json:"tx_fee_klay,omitempty"` //klay
	TxFeeBySenderKLAY   string `bson:"tx_fee_by_sender,omitempty" json:"tx_fee_by_sender,omitempty"`
	TxFeeByFeePayerKLAY string `bson:"tx_fee_by_fee_payer,omitempty" json:"tx_fee_by_fee_payer,omitempty"`
	TypeString          string `bson:"type_string,omitempty" json:"type_string,omitempty"`

	Logs TxLogList `bson:"logs" json:"logs"`
	//Logs interface{} `bson:"logs" json:"logs"`

	IsReceptFailByTxInject bool `bson:"-" json:"-"` //IsNotFound == true 일때도 true일수 있기때문에 block_data의 tx_list를 검증시에만 쓰도록 함.
}
type TransactionBlockList []TransactionBlock

func (my TransactionBlock) String() string     { return dbg.ToJsonString(my) }
func (my TransactionBlockList) String() string { return dbg.ToJsonString(my) }

func (my *TransactionBlock) TxBlockReceipt(sender *Sender, notFoundWait ...bool) TransactionBlock {
	r := sender.ReceiptByHash(my.Hash)
	//if len(notFou)
	//r.IsNotFound()

	if r.IsNotFound {
		if len(notFoundWait) > 0 && notFoundWait[0] {
			cnt := 1
			for {
				cc.YellowItalic("ebcm.TxBlockReceipt[not found wait]...", cnt)
				time.Sleep(time.Second)

				r = sender.ReceiptByHash(my.Hash)
				if !r.IsNotFound {
					cc.YellowItalic("ebcm.TxBlockReceipt[find it]:", my.Hash)
					break
				}
			} //for
		}
	}

	sender.InjectReceipt(my, r)
	return *my
}

// func (my *TransactionBlock) InjectReceipt(r TxReceipt) {
// 	my.Finder.InjectReceipt(my, r)
// }

// GetTransactionFee : eth / klay 에 따른 트랜젝션 수수료(코인) 반환
func (my TransactionBlock) GetTransactionFee() string {
	feePrice := "0"
	if my.TxFeeETH != "" {
		feePrice = my.TxFeeETH
	} else if my.TxFeeKLAY != "" {
		feePrice = my.TxFeeKLAY
	}
	return feePrice
}

// klay :
type Klay_FeeRatio uint8

func (my Klay_FeeRatio) IsValid() bool {
	return 1 <= my && my <= 99
}

type TxType uint16

func (my TxType) Uint16() uint16 { return uint16(my) }

type Klay_RoleType int

//////////////////////

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

func (my CustomInputParseMap) UnixTime(key string) unix.Time {
	if v, do := my[key]; do {
		switch val := v.(type) {
		case string:
			return unix.FromString(val)
		case unix.Time:
			return val
		}

	}
	return unix.ZERO
}

func (my CustomInputParseMap) GetStringArray(key string) []string {
	if val := my.Get(key); val != nil {
		switch array := val.(type) {
		case []string:
			return array

		case []bool:
			rs := []string{}
			for _, v := range array {
				rs = append(rs, dbg.Cat(v))
			}
			return rs
		case []int8:
			rs := []string{}
			for _, v := range array {
				rs = append(rs, dbg.Cat(v))
			}
			return rs
		case []int16:
			rs := []string{}
			for _, v := range array {
				rs = append(rs, dbg.Cat(v))
			}
			return rs
		case []int32:
			rs := []string{}
			for _, v := range array {
				rs = append(rs, dbg.Cat(v))
			}
			return rs
		case []int64:
			rs := []string{}
			for _, v := range array {
				rs = append(rs, dbg.Cat(v))
			}
			return rs
		case []int:
			rs := []string{}
			for _, v := range array {
				rs = append(rs, dbg.Cat(v))
			}
			return rs
		case []uint8:
			rs := []string{}
			for _, v := range array {
				rs = append(rs, dbg.Cat(v))
			}
			return rs
		case []uint16:
			rs := []string{}
			for _, v := range array {
				rs = append(rs, dbg.Cat(v))
			}
			return rs
		case []uint32:
			rs := []string{}
			for _, v := range array {
				rs = append(rs, dbg.Cat(v))
			}
			return rs
		case []uint64:
			rs := []string{}
			for _, v := range array {
				rs = append(rs, dbg.Cat(v))
			}
			return rs
		case []float32:
			rs := []string{}
			for _, v := range array {
				rs = append(rs, dbg.Cat(v))
			}
			return rs
		case []float64:
			rs := []string{}
			for _, v := range array {
				rs = append(rs, dbg.Cat(v))
			}
			return rs

		default:
			rs := []string{}
			rs = append(rs, dbg.Cat(array))
			return rs

		} //switch
	}
	return []string{}
}

func (my CustomInputParseMap) Address(key string) string {
	return dbg.TrimToLower(my.GetString(key))
}
func (my CustomInputParseMap) AddressArray(key string) []string {
	list := my.GetStringArray(key)
	rs := []string{}
	for _, v := range list {
		rs = append(rs, dbg.TrimToLower(v))
	} //for
	return rs
}

func (my CustomInputParseMap) Decimal128(key string) primitive.Decimal128 {
	if v, do := my[key].(primitive.Decimal128); do {
		return v
	}
	return mongo.Decimal128(my.Number(key))
}

func (my CustomInputParseMap) Decimal128Array(key string) []primitive.Decimal128 {
	if array, do := my[key].([]primitive.Decimal128); do {
		return array
	}
	sl := []primitive.Decimal128{}
	for _, v := range my.NumberArray(key) {
		sl = append(sl, mongo.Decimal128(v))
	}
	return sl

}

func (my CustomInputParseMap) Number(key string) string {
	if my == nil {
		return "0"
	}
	return my.GetString(key)
}
func (my CustomInputParseMap) NumberArray(key string) []string {
	return my.GetStringArray(key)
}

func (my CustomInputParseMap) Int64(key string) int64 {
	if my == nil {
		return 0
	}
	return jmath.Int64(my[key])
}
func (my CustomInputParseMap) Int64Array(key string) []int64 {
	list := my.GetStringArray(key)
	rs := []int64{}
	for _, v := range list {
		rs = append(rs, jmath.Int64(v))
	}
	return rs
}

func (my CustomInputParseMap) Int(key string) int {
	return int(my.Int64(key))
}
func (my CustomInputParseMap) IntArray(key string) []int {
	list := my.GetStringArray(key)
	rs := []int{}
	for _, v := range list {
		rs = append(rs, int(jmath.Int64(v)))
	}
	return rs
}

func (my CustomInputParseMap) UInt64(key string) uint64 {
	if my == nil {
		return 0
	}
	return jmath.Uint64(my[key])
}
func (my CustomInputParseMap) UInt64Array(key string) []uint64 {
	list := my.GetStringArray(key)
	rs := []uint64{}
	for _, v := range list {
		rs = append(rs, jmath.Uint64(v))
	}
	return rs
}

func (my CustomInputParseMap) UInt(key string) uint {
	return uint(my.Int64(key))
}
func (my CustomInputParseMap) UIntArray(key string) []uint {
	list := my.GetStringArray(key)
	rs := []uint{}
	for _, v := range list {
		rs = append(rs, uint(jmath.Int64(v)))
	}
	return rs
}

func (my CustomInputParseMap) Bool(key string) bool {
	if v, do := my[key].(bool); do {
		return v
	}
	return false
}
func (my CustomInputParseMap) BoolArray(key string) []bool {
	list := my.GetStringArray(key)
	rs := []bool{}
	for _, v := range list {
		rs = append(rs, dbg.IsTrue(v))
	} //for
	return rs
}

func (my CustomInputParseMap) Parse(key string, p interface{}) error {
	if my == nil {
		return errors.New("CustomInputParseMap is nil")
	}
	return dbg.ParseStruct(my[key], p)
}

func InputParse[T any](m CustomInputParseMap, key string, p *T) bool {
	if v, do := m[key].(T); do {
		*p = v
	}
	return false
}

func (my CustomInputParseMap) Value(key string) interface{} {
	if my == nil {
		return nil
	}
	return my[key]
}

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

// CustimInputParseTag : key check
func (my TransactionBlock) CustomInputParseKey(tag string) bool {
	if my.CustomInputParse == nil {
		return false
	}
	_, do := my.CustomInputParse[tag]
	return do
}

package ebcmx

import (
	"encoding/hex"
	"strings"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"

	"golang.org/x/crypto/sha3"
)

func normalizeEventFunc(event string) string {
	event = strings.TrimPrefix(
		strings.ReplaceAll(
			strings.TrimSpace(event),
			";",
			"",
		),
		"event ",
	)
	ss := strings.Split(event, "(")
	name := strings.TrimSpace(ss[0])
	body := strings.ReplaceAll(ss[1], ")", "")
	bodyc := strings.Split(body, " ")
	for i := 0; i < len(bodyc); i++ {
		if bodyc[i] == "uint" {
			bodyc[i] = "uint256"
		}
	}
	body = strings.Join(bodyc, " ")

	args := strings.Split(body, ",")
	params := []string{}
	for _, arg := range args {
		arg = strings.Split(
			strings.TrimSpace(arg),
			" ",
		)[0]
		params = append(params, arg)
	} //for
	return dbg.Cat(name, "(", strings.Join(params, ","), ")")
}
func EventFuncHash(normal_event string) string {
	fnSignature := []byte(normal_event)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(fnSignature)
	methodID := hash.Sum(nil)

	name := Hexutil_Encode(methodID)
	return name
}

// MakeTopicName :
func MakeTopicName(abiFuncString string) string {
	return EventFuncHash(normalizeEventFunc(abiFuncString))
}

// MakeTopicName2 : abifunc , topic-name-hex
func MakeTopicName2(abiFuncString string, isView ...bool) (string, string) {
	abiFuncString = normalizeEventFunc(abiFuncString)
	name := EventFuncHash(abiFuncString)

	//0xddf252ad1be2c89b69c2b068 fc378daa952ba7f163c4a11628f55a4df523b3ef
	//name = "0x" + name[2+24:]
	if dbg.IsTrue2(isView...) {
		dbg.Purple(name, "<--->", abiFuncString)
	}
	return abiFuncString, name
}

type Topic string //0x... //64
func (my Topic) Address() string {
	//0x0000000000000000000000009e99b3fd1e5558b304a2cecbe988c0404b3bd7e2
	return "0x" + string(my[26:])
}
func (my Topic) Value() string  { return string(my) }
func (my Topic) Number() string { return jmath.VALUE(string(my)) }
func (my Topic) Uint64() uint64 { return jmath.Uint64(string(my)) }
func (my Topic) Uint32() uint32 { return uint32(my.Uint64()) }
func (my Topic) Uint16() uint16 { return uint16(my.Uint64()) }
func (my Topic) Uint8() uint8   { return uint8(my.Uint64()) }
func (my Topic) Int() int       { return int(my.Uint64()) } //uint16 , uint8
func (my Topic) Bool() bool     { return jmath.CMP(my.Number(), 0) > 0 }

type Topics []Topic

type TopicParam interface {
	Address() string
	Value() string
	Number() string
}
type TopicParams []TopicParam

func (my TopicParams) Help() string {
	//return dbg.ToJSONString(my)
	msg := "< TopicParams >" + dbg.ENTER
	for i, v := range my {
		msg += dbg.Cat("[", i, "] ", v, dbg.ENTER)
	}
	return msg
}

func (my Topics) GetName() string {
	if len(my) == 0 {
		return ""
	}
	return string(my[0])
}
func (my Topics) IsMethod(topicName string) bool {
	return my.GetName() == topicName
}
func (my Topics) GetParams() TopicParams {
	ps := TopicParams{}
	tp := my[1:]
	for _, v := range tp {
		ps = append(ps, v)
	}
	return ps
}

type DataItem string

type DataItemList []DataItem

func (my DataItemList) Bytes() []byte {
	buf := []byte{}
	for _, v := range my {
		b, _ := hex.DecodeString(string(v)[2:])
		buf = append(buf, b...)
	}
	return buf
}

func (my DataItemList) Help() string {
	//return dbg.ToJSONString(my)
	msg := "< DataItemList >" + dbg.ENTER
	for i, v := range my {
		msg += dbg.Cat("[", i, "] ", v, dbg.ENTER)
	}
	return msg
}

func (my DataItem) Address() string { return "0x" + string(my[26:]) }

// Value : 0x00000xxxxx
func (my DataItem) Value() string  { return string(my) }
func (my DataItem) Number() string { return jmath.VALUE(string(my)) }
func (my DataItem) Uint64() uint64 { return jmath.Uint64(string(my)) }
func (my DataItem) Uint32() uint32 { return uint32(my.Uint64()) }
func (my DataItem) Uint16() uint16 { return uint16(my.Uint64()) }
func (my DataItem) Uint8() uint8   { return uint8(my.Uint64()) }
func (my DataItem) Int() int       { return int(my.Uint64()) } //uint16 , uint8
func (my DataItem) Bool() bool     { return my.Uint64() == 1 }

type TxLog struct {
	Address string `bson:"address" json:"address"`
	Topics  Topics `bson:"topics" json:"topics"`
	//Data        string       `bson:"data" json:"data"`
	Data        DataItemList `bson:"data" json:"data"`
	BlockNumber uint64       `bson:"block_number" json:"block_number"`
	TxHash      string       `bson:"tx_hash" json:"tx_hash"`
	TxIndex     uint         `bson:"tx_index" json:"tx_index"`
	BlockHash   string       `bson:"block_hash" json:"block_hash"`
	LogIndex    uint         `bson:"log_index" json:"log_index"`
	Removed     bool         `bson:"removed" json:"removed"`
}

func (my TxLog) Help() string {
	msg := "< Log >" + dbg.ENTER
	msg += dbg.Cat("address : ", my.Address, dbg.ENTER)
	msg += dbg.Cat("ABI     : ", my.Topics.GetName(), dbg.ENTER)
	msg += dbg.Cat(my.Topics.GetParams().Help(), dbg.ENTER)
	msg += dbg.Cat(my.Data.Help(), dbg.ENTER)
	return msg
}

type TxLogList []TxLog

func (my TxLog) String() string     { return dbg.ToJSONString(my) }
func (my TxLogList) String() string { return dbg.ToJSONString(my) }
func (my *TxLogList) Remove()       { (*my) = TxLogList{} }

type TxReceipt struct {
	BlockHash   string `bson:"block_hash" json:"block_hash"`
	BlockNumber string `bson:"block_number" json:"block_number"`
	From        string `bson:"from,omitempty" json:"from,omitempty"`   //klay
	To          string `bson:"to,omitempty" json:"to,omitempty"`       //klay
	Nonce       uint64 `bson:"nonce,omitempty" json:"nonce,omitempty"` //klay
	Bloom       string `bson:"bloom" json:"bloom"`
	Status      uint64 `bson:"status" json:"status"`
	//CumulativeGasUsed uint64    `bson:"cumulative_gas_used" json:"cumulative_gas_used"`
	Logs             TxLogList `bson:"logs" json:"logs"`
	TransactionHash  string    `bson:"transaction_hash" json:"transaction_hash"`
	TransactionIndex uint      `bson:"transaction_index" json:"transaction_index"`
	Gas              string    `bson:"gas,omitempty" json:"gas,omitempty"`           //klay
	GasPrice         string    `bson:"gasPrice,omitempty" json:"gasPrice,omitempty"` //klay
	GasUsed          uint64    `bson:"gasUsed" json:"gasUsed"`
	ContractAddress  string    `bson:"contract_address" json:"contract_address"`

	Type   interface{} `bson:"type,emitempty" json:"type,emitempty"`     //klay
	Amount string      `bson:"amount,omitempty" json:"amount,omitempty"` //klay

	SenderTxHash string `bson:"senderTxHash,omitempty" json:"senderTxHash,omitempty"` //klay

	//////ETH
	CumulativeGasUsed uint64 `bson:"cumulative_gas_used,omitempty" json:"cumulative_gas_used,omitempty"` //eth

	IsNotFound bool `bson:"-" json:"is_not_found"`
}

func (my TxReceipt) Valid() bool    { return my.BlockHash != "" }
func (my TxReceipt) IsError() bool  { return my.Status != 1 }
func (my TxReceipt) String() string { return dbg.ToJSONString(my) }

package ecsx

import (
	"context"
	"encoding/hex"
	"strings"
	"txscheduler/brix/tools/cloudx/ebcmx"
	ebcmABI "txscheduler/brix/tools/cloudx/ebcmx/abix"
	"txscheduler/brix/tools/cloudx/ethwallet/abmx"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/crypto/sha3"
)

// MakeTopicName :
func MakeTopicName(abiFuncString string, isView ...bool) string {
	_, n := MakeTopicName2(abiFuncString, isView...)
	return n
}

// MakeTopicName2 : abifunc , topic-name
func MakeTopicName2(abiFuncString string, isView ...bool) (string, string) {
	abiFuncString = strings.ReplaceAll(abiFuncString, " ", "")
	fnSignature := []byte(abiFuncString)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(fnSignature)
	methodID := hash.Sum(nil)
	name := hexutil.Encode(methodID)
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
func (my DataItemList) ParseABI(abiReturns interface{}, callback func(rs abmx.RESULT)) bool {
	rs := abmx.ReceiptDiv(
		my.Bytes(),
		abiReturns,
	)
	if !rs.IsError {
		callback(rs)
	} else {
		dbg.Red("DataItemList.ParseABI FAIL.")
		return false
	}
	return true
}

func EBCM_DataItemList_ParseABI(
	list ebcmx.DataItemList,
	typelist ebcmABI.TypeList,
	callback func(rs ebcmABI.RESULT),
) bool {
	rs := abmx.ReceiptDiv(
		list.Bytes(),
		abmx.EBCM_ABI_NewReturns(typelist...),
	)
	if !rs.IsError {
		callback(rs.RESULT)
	} else {
		dbg.Red("[ebcm]DataItemList.ParseABI FAIL.")
		return false
	}
	return true
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

func MakeDataItemList(data string) DataItemList {
	sl := DataItemList{}
	if len(data) <= 2+64 { //0x + 64
		sl = append(sl, DataItem(data))
		return sl
	}
	v := data[2:] // remove 0x
	for len(v) >= 64 {
		c := v[:64]
		sl = append(sl, DataItem("0x"+c))
		v = v[64:]
	}
	return sl
}

/*
	TxLogList : type TxLog struct {
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
*/
type TxLogList []TxLog

func (my TxLog) String() string     { return dbg.ToJSONString(my) }
func (my TxLogList) String() string { return dbg.ToJSONString(my) }
func (my *TxLogList) Remove()       { (*my) = TxLogList{} }

type TxReceipt struct {
	BlockHash         string    `bson:"block_hash" json:"block_hash"`
	BlockNumber       string    `bson:"block_number" json:"block_number"`
	Bloom             string    `bson:"bloom" json:"bloom"`
	Status            uint64    `bson:"status" json:"status"`
	CumulativeGasUsed uint64    `bson:"cumulative_gas_used" json:"cumulative_gas_used"`
	Logs              TxLogList `bson:"logs" json:"logs"`
	TransactionHash   string    `bson:"transaction_hash" json:"transaction_hash"`
	TransactionIndex  uint      `bson:"transaction_index" json:"transaction_index"`
	GasUsed           uint64    `bson:"gasUsed" json:"gasUsed"`
	ContractAddress   string    `bson:"contract_address" json:"contract_address"`

	Err        error `bson:"-" json:"err"`
	IsNotFound bool  `bson:"-" json:"is_not_found"`
}

func (my TxReceipt) Valid() bool     { return my.BlockHash != "" }
func (my TxReceipt) IsError() bool   { return my.Status != 1 }
func (my TxReceipt) IsSuccess() bool { return my.Status != 0 }
func (my TxReceipt) String() string  { return dbg.ToJSONString(my) }

func (my TxReceipt) IsAck() bool {
	if my.Err != nil {
		return false
	}
	return my.IsNotFound == false
}

// ReceiptByHash :
func (my *Sender) ReceiptByHash(hexHash string) TxReceipt {
	hexHash = strings.TrimSpace(hexHash)
	hash := common.HexToHash(hexHash)

	re := TxReceipt{
		Logs: TxLogList{},
	}

	getAddress := func(a common.Address) string {
		return dbg.TrimToLower(a.Hex())
	}
	encodeBytes := func(b []byte) string {
		return dbg.TrimToLower("0x" + hex.EncodeToString(b))
	}

	v, err := my.client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		re.Err = err
		dbg.Red("ecsx.ReceiptByHash : ", err, "[", hexHash, "]")
		if strings.Contains(err.Error(), "not found") {
			re.IsNotFound = true
		}
		return re
	}
	dbg.White(dbg.ToJSONString(v), err)

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

		log := TxLog{
			Address:     getAddress(l.Address),
			Data:        MakeDataItemList(encodeBytes(l.Data)),
			BlockNumber: l.BlockNumber,
			TxHash:      l.TxHash.Hex(),
			BlockHash:   l.BlockHash.Hex(),
			LogIndex:    l.Index,
			TxIndex:     l.TxIndex,
			Removed:     l.Removed,
		}

		for _, topic := range l.Topics {
			c := dbg.TrimToLower(topic.Hex())
			log.Topics = append(log.Topics, Topic(c))
		}

		re.Logs = append(re.Logs, log)
	} //for

	return re
}

func (my *Sender) ebcm_ReceiptByHash(hexHash string) ebcmx.TxReceipt {

	data := my.ReceiptByHash(hexHash)

	item := ebcmx.TxReceipt{}
	dbg.ChangeStruct(data, &item)
	return item
}

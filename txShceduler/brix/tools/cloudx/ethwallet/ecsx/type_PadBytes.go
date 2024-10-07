package ecsx

import (
	"encoding/hex"
	"math/big"
	"strings"
	"txscheduler/brix/tools/cloudx/ebcmx"
	ebcmABI "txscheduler/brix/tools/cloudx/ebcmx/abix"
	"txscheduler/brix/tools/cloudx/ethwallet/abmx"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

// Appender : PadBytes.Appender
type Appender interface {
	/////////////////////////////////////////////////

	IndexSize() int
	SetParamHeader(count int) int
	//SetIndex(header int, index int)

	SetAddress(header int, address string) int
	SetAddressArray(header, i int, address ...string) int

	SetAmount(header int, value string) int
	SetAmountArray(header, i int, value ...string) int

	SetBool(header int, b bool) int
	SetBoolArray(header int, i int, bs ...bool) int

	SetText(header, i int, text string) int
	SetTextArray(header, i int, text ...string) int

	SetBytes(header, i int, Bytes []byte) int
	SetBytesHex(header, i int, hexBytes string) int

	SetBytesArray(header, i int, Bytes ...[]byte) int
	SetBytesHexArray(header, i int, hexBytes ...string) int

	SetBytes32(header int, Bytes []byte) int
	SetBytes32Hex(header int, hexBytes string) int

	SetBytes32Array(header, i int, Bytes ...[]byte) int
	SetBytes32HexArray(header, i int, hexBytes ...string) int

	HexDebug() string
}

// PadBytes :
type cPadBytes []byte

// PadBytes : wrapped type
type PadBytes struct {
	cPadBytes
}

type PADBYTES interface {
	Bytes() []byte
	HexDebug() string
	Hex() string
}

func GetPadbytesFromEBCM(v ebcmx.PADBYTES) PadBytes {
	r := PadBytes{}
	r.cPadBytes = v.Bytes()
	return r
}

func MakePadBytesABI(pureName string, ebcmTypes ebcmABI.TypeList) PadBytes {
	ps := abmx.EBCM_ABI_NewParams(ebcmTypes...)
	ebcm_bytes := NewPadBytesABI(pureName, ps)
	return GetPadbytesFromEBCM(ebcm_bytes)
}

func NewPadBytesABI(pureName string, ps abmx.AbiParams) PadBytes {
	buf := abmx.PadBytes(pureName, ps)
	return PadBytes{cPadBytes(buf)}
}

// NewPadBytes :
func NewPadBytes(methodName string) cPadBytes {
	if !strings.Contains(methodName, "(") || !strings.Contains(methodName, ")") {
		dbg.Red("NewPadBytes Fail : ()")
	}
	methodName = strings.ReplaceAll(methodName, " ", "")
	transferFnSignature := []byte(methodName) // ----- "transfer(address,uint256)"
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	pad := hash.Sum(nil)[:4]
	return pad
}

// MakePadBytes :
func MakePadBytes(method string, callback func(Appender)) PadBytes {
	pad := NewPadBytes(method)
	data := PadBytes{cPadBytes: pad}
	callback(&data)
	return data
}
func MakePadBytes2(funcName string, types abmx.TypeList, callback func(Appender)) PadBytes {
	params := []string{}
	for _, v := range types {
		params = append(params, dbg.Cat(v))
	}
	method := strings.ReplaceAll(funcName+"("+strings.Join(params, ",")+")", " ", "")
	pad := NewPadBytes(method)
	data := PadBytes{cPadBytes: pad}
	callback(&data)
	return data
}

func ebcm_MakePadBytesABI(pureName string, ebcmTypes ebcmABI.TypeList) ebcmx.PADBYTES {
	ps := abmx.EBCM_ABI_NewParams(ebcmTypes...)
	return NewPadBytesABI(pureName, ps)
}

// ebcm_MakePadBytes :
func ebcm_MakePadBytes(method string, callback func(ebcmx.Appender)) ebcmx.PADBYTES {
	r := MakePadBytes(method, func(a Appender) {
		callback(a)
	})
	return r
}

func ebcm_MakePadBytes2(funcName string, ebcmTypes ebcmABI.TypeList, callback func(ebcmx.Appender)) ebcmx.PADBYTES {
	params := []string{}
	for _, v := range ebcmTypes {
		params = append(params, v.String())
	}
	method := strings.ReplaceAll(funcName+"("+strings.Join(params, ",")+")", " ", "")
	r := MakePadBytes(method, func(a Appender) {
		callback(a)
	})
	return r
}

// Bytes :
func (my cPadBytes) Bytes() []byte {
	return []byte(my)
}

func (my cPadBytes) Hex() string {
	return hex.EncodeToString(my.Bytes())
}
func (my cPadBytes) HexString() string {
	msg := "0x" + my.Hex()
	result := msg[:10] + dbg.ENTER
	msg = msg[10:]
	for len(msg) >= 64 {
		result += msg[:64] + dbg.ENTER
		msg = msg[64:]
	} //for

	return result
}
func (my cPadBytes) IndexSize() int {
	return (len(my) - 4) / 32
}

// Append :
func (my *cPadBytes) Append(data cPadBytes) {
	*my = append(*my, data...)
}

//AppendAddress :
// func (my *cPadBytes) AppendAddress(address string) {
// 	hexAddress := common.HexToAddress(strings.ToLower(address))
// 	my.Append(common.LeftPadBytes(hexAddress.Bytes(), 32))
// }

func (my *cPadBytes) setArea(i int) int {
	offset := 4 + (i * 32)
	size := offset + 32
	if size >= len(*my) {
		*my = append(*my, make([]byte, size-len(*my))...)
	}
	return offset
}
func (my *cPadBytes) SetAddress(header int, address string) int {
	offset := my.setArea(header)
	hexAddress := common.HexToAddress(strings.ToLower(address))
	buf := common.LeftPadBytes(hexAddress.Bytes(), 32)
	copy((*my)[offset:], buf)
	return my.IndexSize()
}

func (my *cPadBytes) SetAddressArray(header, i int, address ...string) int {
	my.__setIndex(header, i)

	cnt := len(address)
	my.SetAmount(i, jmath.VALUE(cnt))
	if cnt > 0 {
		j := 0
		for {
			i++
			my.SetAddress(i, address[j])
			j++
			if j >= cnt {
				break
			}
		} //for
	} else {
		dbg.Red("cPadBytes.SetAddressArray Size is Zero")
	}

	return my.IndexSize()
}

func (my *cPadBytes) SetAmountArray(header, i int, value ...string) int {
	my.__setIndex(header, i)

	cnt := len(value)
	my.SetAmount(i, jmath.VALUE(cnt))
	if cnt > 0 {
		j := 0
		for {
			i++
			my.SetAmount(i, value[j])
			j++
			if j >= cnt {
				break
			}
		} //for
	} else {
		dbg.Red("cPadBytes.SetAmountArray Size is Zero")
	}
	return my.IndexSize()
}

func (my *cPadBytes) SetAmount(header int, value string) int {
	offset := my.setArea(header)

	amount := new(big.Int)
	amount.SetString(value, 10)
	buf := common.LeftPadBytes(amount.Bytes(), 32)
	copy((*my)[offset:], buf)
	return my.IndexSize()
}

func (my *cPadBytes) SetBool(header int, b bool) int {
	v := "0"
	if b {
		v = "1"
	}
	return my.SetAmount(header, v)
}
func (my *cPadBytes) SetBoolArray(header, i int, bs ...bool) int {
	vals := []string{}
	for _, v := range bs {
		if v {
			vals = append(vals, "1")
		} else {
			vals = append(vals, "0")
		}
	}
	return my.SetAmountArray(header, i, vals...)
}

func (my *cPadBytes) SetBytes32(header int, Bytes []byte) int {
	if len(Bytes) > 32 {
		dbg.Red("cPadBytes.SetBytes32 Over Length :", len(Bytes))
		Bytes = Bytes[:32]
	}
	offset := my.setArea(header)
	copy((*my)[offset:], Bytes)
	return my.IndexSize()
}

func (my *cPadBytes) SetBytes32Hex(header int, hexBytes string) int {
	bigVal := jmath.New(hexBytes)
	buf := bigVal.ToBigInteger().Bytes()
	return my.SetBytes32(header, buf)
}

func (my *cPadBytes) SetBytes32Array(header, i int, Bytes ...[]byte) int {
	my.__setIndex(header, i)

	cnt := len(Bytes)
	my.SetAmount(i, jmath.VALUE(cnt))
	j := 0
	for {
		i++
		my.SetBytes32(i, Bytes[j])
		j++
		if j >= cnt {
			break
		}
	} //for

	return my.IndexSize()
}
func (my *cPadBytes) SetBytes32HexArray(header, i int, hexBytes ...string) int {
	bufs := [][]byte{}
	for _, str := range hexBytes {
		bigVal := jmath.New(str)
		buf := bigVal.ToBigInteger().Bytes()
		bufs = append(bufs, buf)
	}

	return my.SetBytes32Array(header, i, bufs...)
}

func (my *cPadBytes) __setString(i int, text string) {
	textBytes := []byte(text)
	sl := [][]byte{}
	for len(textBytes) > 0 {
		cut := 32
		if cut > len(textBytes) {
			cut = len(textBytes)
		}
		bb := textBytes[:cut]
		sl = append(sl, bb)
		textBytes = textBytes[cut:]
	}

	for _, v := range sl {
		offset := my.setArea(i)
		buf := common.RightPadBytes(v, 32)
		copy((*my)[offset:], buf)
		i++
	} //for
}

func (my *cPadBytes) SetText(header, i int, text string) int {
	my.__setIndex(header, i)

	my.SetAmount(i, jmath.VALUE(len(text)))
	my.__setString(i+1, text)
	return my.IndexSize()
}
func (my *cPadBytes) SetTextArray(header, i int, texts ...string) int {

	my.__setIndex(header, i)

	cnt := len(texts)
	my.SetAmount(i, jmath.VALUE(cnt))

	offIndex := []int{}
	offset := my.IndexSize()
	offText := offset
	offIndex = append(offIndex, offset)

	for j := 0; j < cnt; j++ {
		my.__setIndex(offset, 0)
		offset = my.IndexSize()
		offIndex = append(offIndex, offset)
	}

	pos := my.IndexSize()
	subIndex := pos - offText
	dbg.YellowBG(subIndex)
	for j := 0; j < cnt; j++ {
		my.__setIndex(offIndex[j], subIndex)

		my.SetAmount(pos, jmath.VALUE(len(texts[j])))
		my.__setString(pos+1, texts[j])

		pos = my.IndexSize()

		subIndex = pos - offText
		dbg.Yellow(subIndex)
		//dbg.Cyan(my.HexDebug())
	}

	return my.IndexSize()
}

func (my *cPadBytes) __setBytes(i int, buf []byte) {
	sl := [][]byte{}
	for len(buf) > 0 {
		cut := 32
		if cut > len(buf) {
			cut = len(buf)
		}
		bb := buf[:cut]
		sl = append(sl, bb)
		buf = buf[cut:]
	}

	for _, v := range sl {
		offset := my.setArea(i)
		buf := common.RightPadBytes(v, 32)
		copy((*my)[offset:], buf)
		i++
	} //for
}

func (my *cPadBytes) SetBytesHex(header, i int, hexBytes string) int {
	bigVal := jmath.New(hexBytes)
	buf := bigVal.ToBigInteger().Bytes()

	return my.SetBytes(header, i, buf)
}
func (my *cPadBytes) SetBytes(header, i int, Bytes []byte) int {
	my.__setIndex(header, i)

	my.SetAmount(i, jmath.VALUE(len(Bytes)))
	my.__setBytes(i+1, Bytes)
	return my.IndexSize()
}

func (my *cPadBytes) SetBytesArray(header, i int, Bytes ...[]byte) int {

	my.__setIndex(header, i)

	cnt := len(Bytes)
	my.SetAmount(i, jmath.VALUE(cnt))

	offIndex := []int{}
	offset := my.IndexSize()
	offText := offset
	offIndex = append(offIndex, offset)

	for j := 0; j < cnt; j++ {
		my.__setIndex(offset, 0)
		offset = my.IndexSize()
		offIndex = append(offIndex, offset)
	}

	pos := my.IndexSize()
	subIndex := pos - offText
	dbg.YellowBG(subIndex)
	for j := 0; j < cnt; j++ {
		my.__setIndex(offIndex[j], subIndex)

		my.SetAmount(pos, jmath.VALUE(len(Bytes[j])))
		my.__setBytes(pos+1, Bytes[j])

		pos = my.IndexSize()

		subIndex = pos - offText
		dbg.Yellow(subIndex)
		//dbg.Cyan(my.HexDebug())
	}

	return my.IndexSize()
}

func (my *cPadBytes) SetBytesHexArray(header, i int, hexBytes ...string) int {

	bufs := [][]byte{}
	for _, str := range hexBytes {
		bigVal := jmath.New(str)
		buf := bigVal.ToBigInteger().Bytes()
		bufs = append(bufs, buf)
	}

	return my.SetBytesArray(header, i, bufs...)
}

func (my *cPadBytes) __setIndex(header int, index int) {
	my.SetAmount(header, jmath.VALUE(32*index))
}
func (my *cPadBytes) SetParamHeader(count int) int {
	i := 0
	for {
		my.__setIndex(i, 0)
		i++
		if i >= count {
			break
		}
	} //for
	return my.IndexSize()
}
func (my cPadBytes) HexDebug() string {
	result := dbg.ENTER
	msg := hex.EncodeToString([]byte(my))

	result += dbg.Cat("[F] 0x", msg[:8], dbg.ENTER)

	msg = msg[8:] //funcName skip

	loop := 0
	for len(msg) >= 64 {
		result += dbg.Cat("[", loop, "] ", msg[:64], dbg.ENTER)
		msg = msg[64:]
		loop++
	}
	return result
}

//AppendAddressArray :
// func (my *cPadBytes) AppendAddressArray(addresslist ...string) {
// 	count := len(addresslist)
// 	size := new(big.Int)
// 	size.SetString(fmt.Sprint(count), 10)
// 	my.Append(common.LeftPadBytes(size.Bytes(), 32))

// 	for _, address := range addresslist {
// 		my.AppendAddress(address)
// 	} //for
// }

//AppendAmountArray :
// func (my *cPadBytes) AppendAmountArray(values ...string) {
// 	count := len(values)
// 	size := new(big.Int)
// 	size.SetString(fmt.Sprint(count), 10)
// 	my.Append(common.LeftPadBytes(size.Bytes(), 32))

// 	for _, value := range values {
// 		my.AppendAmount(value)
// 	} //for
// }

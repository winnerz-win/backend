package abmx

import (
	"encoding/hex"
	"math/big"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"

	"github.com/ethereum/go-ethereum/common"

	ebcmABI "txscheduler/brix/tools/cloudx/ebcmx/abix"
)

func ReceiptDiv(receipt []byte, Ireturns interface{}) (rs RESULT) {
	defer func() {
		if e := recover(); e != nil {
			dbg.Purple("abm.ReceiptDiv : ", e)
			rs = RESULT{}
			rs.IsError = true
		}
	}()
	if len(receipt)%32 != 0 {
		//dbg.Red("ReceiptDiv__Size :", len(receipt))
		rs = RESULT{}
		rs.IsError = true
		return rs
	}
	rts := Ireturns.(abiReturns)
	rs = receiptDiv(receipt, rts)
	return rs
}

func receiptDiv(receipt []byte, rts abiReturns) RESULT {
	if len(receipt) == 0 {
		rs := RESULT{}
		return rs
	}
	toAddress := func(data32 []byte) string {
		return dbg.TrimToLower(common.BytesToAddress(data32).Hex())
	}
	amountRe := func(data32 []byte) string {
		v := big.NewInt(0)
		v = v.SetBytes(data32)
		return v.String()
	}
	amountValue := func(data32 []byte) string {
		return jmath.VALUE(amountRe(data32))
	}
	boolValue := func(data32 []byte) bool {
		if data32[len(data32)-1] == 0 {
			return false
		} else {
			return true
		}
	}
	bytesHex := func(data32 []byte) string {
		return dbg.TrimToLower("0x" + hex.EncodeToString(data32))
	}
	getInt := func(data32 []byte) int {
		return jmath.Int("0x" + hex.EncodeToString(data32))
	}
	getIndex := func(data32 []byte) int {
		return getInt(data32) / 32
	}

	//dbg.Purple(receipt)
	sl := [][]byte{}
	for len(receipt) > 0 {
		v := receipt[:32]
		sl = append(sl, v)
		receipt = receipt[32:]
	} //for

	/*
		hex - index
		20 - 1
		40 - 2
		60 - 3
		a0 - 5
		c0 - 6
		e0 - 7
		100 - 8
		120 - 9
		140 - 10
	*/

	result := newResultArgs(rts)
	//var result RESULT

	// loopCount := 0
	// loopOffset := [][][]byte{}
	getValue := func(t Type, sl [][]byte, offset int) interface{} {
		var val interface{}
		data32 := sl[offset]
		switch t {
		case Address:
			val = toAddress(data32)
			offset++

		case Uint256, Uint128, Uint112:
			val = amountValue(data32)

		case Uint64:
			val = jmath.Uint64(amountRe(data32))

		case Uint32:
			val = uint32(jmath.Uint64(amountRe(data32)))

		case Uint16:
			val = uint16(jmath.Uint64(amountRe(data32)))

		case Uint8:
			val = uint8(jmath.Uint64(amountRe(data32)))

		case Bool:
			val = boolValue(data32)

		case Bytes32:
			val = bytesHex(data32)

		case String:
			buf := []byte{}
			pos := getIndex(data32)
			count := getInt(sl[pos])
			pos++

			loop := count / 32
			for j := 0; j < loop; j++ {
				data32 = sl[pos]
				buf = append(buf, data32...)
				pos++
			}
			if v := count % 32; v > 0 {
				data32 = sl[pos]
				buf = append(buf, data32[:v]...)
			}
			val = string(buf)

		case Bytes:
			buf := []byte{}
			pos := getIndex(data32)
			count := getInt(sl[pos])
			pos++

			loop := count / 32
			for j := 0; j < loop; j++ {
				data32 = sl[pos]
				buf = append(buf, data32...)
				pos++
			}
			if v := count % 32; v > 0 {
				data32 = sl[pos]
				buf = append(buf, data32[:v]...)
			}
			val = bytesHex(buf)

		case AddressArray:
			obj := []string{}
			pos := getIndex(data32)
			count := getInt(sl[pos])
			for j := 0; j < count; j++ {
				pos++
				data32 = sl[pos]
				val := toAddress(data32)
				obj = append(obj, val)
			} //for
			val = obj

		case Uint256Array, Uint128Array, Uint112Array:
			obj := []string{}
			pos := getIndex(data32)
			count := getInt(sl[pos])
			for j := 0; j < count; j++ {
				pos++
				data32 = sl[pos]
				val := amountValue(data32)
				obj = append(obj, val)
			} //for
			val = obj

		case Uint64Array:
			obj := []uint64{}
			pos := getIndex(data32)
			count := getInt(sl[pos])
			for j := 0; j < count; j++ {
				pos++
				data32 = sl[pos]
				val := jmath.Uint64(amountRe(data32))
				obj = append(obj, val)
			} //for
			val = obj

		case Uint32Array:
			obj := []uint32{}
			pos := getIndex(data32)
			count := getInt(sl[pos])
			for j := 0; j < count; j++ {
				pos++
				data32 = sl[pos]
				val := uint32(jmath.Uint64(amountRe(data32)))
				obj = append(obj, val)
			} //for
			val = obj

		case Uint16Array:
			obj := []uint16{}
			pos := getIndex(data32)
			count := getInt(sl[pos])
			for j := 0; j < count; j++ {
				pos++
				data32 = sl[pos]
				val := uint16(jmath.Uint64(amountRe(data32)))
				obj = append(obj, val)
			} //for
			val = obj

		case Uint8Array:
			obj := []uint8{}
			pos := getIndex(data32)
			count := getInt(sl[pos])
			for j := 0; j < count; j++ {
				pos++
				data32 = sl[pos]
				val := uint8(jmath.Uint64(amountRe(data32)))
				obj = append(obj, val)
			} //for
			val = obj

		case BoolArray:
			obj := []bool{}
			pos := getIndex(data32)
			count := getInt(sl[pos])
			for j := 0; j < count; j++ {
				pos++
				data32 = sl[pos]
				val := boolValue(data32)
				obj = append(obj, val)
			} //for
			val = obj

		case Bytes32Array:
			obj := []string{}
			pos := getIndex(data32)
			count := getInt(sl[pos])
			for j := 0; j < count; j++ {
				pos++
				data32 = sl[pos]
				val := bytesHex(data32)
				obj = append(obj, val)
			} //for
			val = obj

		case StringArray:
			ss := []string{}
			pos := getIndex(data32)
			size := getInt(sl[pos]) //2

			ttl := sl[pos+1:]
			//debugBufs(ttl)

			for i := 0; i < size; i++ {
				buf := []byte{}
				pos := getIndex(ttl[i])
				count := getInt(ttl[pos])
				pos++

				loop := count / 32
				for j := 0; j < loop; j++ {
					data32 = ttl[pos]
					buf = append(buf, data32...)
					pos++
				}
				if v := count % 32; v > 0 {
					data32 = ttl[pos]
					buf = append(buf, data32[:v]...)
				}
				ss = append(ss, string(buf))

			} //for

			val = ss

		case BytesArray:
			ss := []string{}
			pos := getIndex(data32)
			size := getInt(sl[pos]) //2

			ttl := sl[pos+1:]
			//debugBufs(ttl)

			for i := 0; i < size; i++ {
				buf := []byte{}
				pos := getIndex(ttl[i])
				count := getInt(ttl[pos])
				pos++

				loop := count / 32
				for j := 0; j < loop; j++ {
					data32 = ttl[pos]
					buf = append(buf, data32...)
					pos++
				}
				if v := count % 32; v > 0 {
					data32 = ttl[pos]
					buf = append(buf, data32[:v]...)
				}
				ss = append(ss, bytesHex(buf))

			} //for

			val = ss

		} //t
		return val
	}
	_ = getValue

	for headerIndex, v := range rts.args {
		vp := v.p
		//data32 := sl[headerIndex]
		switch vp.Name() {
		default:
			val := getValue(vp, sl, headerIndex)
			result.EBCMaddItem(vp.String(), val)
			//headerIndex = offset

		case iTuple, iTupleFlex:
			pos := getIndex(sl[headerIndex])

			flex := vp.(*flexType)
			tupleList := ebcmABI.TupleList{}
			list := ebcmABI.ResultItemList{}

			ttl := sl[pos:]
			//debugBufs(ttl)
			for i, data := range flex.datas {
				val := getValue(data, ttl, i)
				list.EBCMadd(data.String(), val)
			} //for

			tupleList = append(tupleList, list)
			result.EBCMaddItem(vp.String(), tupleList)

		case iTupleFixedArray:
			pos := getIndex(sl[headerIndex])
			count := getInt(sl[pos])
			pos++

			flex := vp.(*flexType)
			tupleList := ebcmABI.TupleList{}
			for i := 0; i < count; i++ {
				list := ebcmABI.ResultItemList{}
				for _, data := range flex.datas {
					val := getValue(data, sl, pos)
					list.EBCMadd(data.String(), val)
					pos++
				} //for
				tupleList = append(tupleList, list)
			} //for
			result.EBCMaddItem(vp.String(), tupleList)

		case iTupleFlexArray:
			pos := getIndex(sl[headerIndex])
			count := getInt(sl[pos])
			//pos++

			aal := [][][]byte{}

			{
				ttl := sl[pos+1:]
				//debugBufs(ttl)
				for i := 0; i < count; i++ {
					cut := getIndex(ttl[i])

					if i+1 < count {
						end := getIndex(ttl[i+1])
						aal = append(aal, ttl[cut:end])
					} else {
						aal = append(aal, ttl[cut:])
					}

				}
			}

			flex := vp.(*flexType)
			tupleList := ebcmABI.TupleList{}
			for i := 0; i < count; i++ {
				ssl := aal[i]
				//debugBufs(ssl)

				list := ebcmABI.ResultItemList{}
				for i, data := range flex.datas {
					val := getValue(data, ssl, i)
					list.EBCMadd(data.String(), val)
					pos++
				} //for
				tupleList = append(tupleList, list)

			} //for
			result.EBCMaddItem(vp.String(), tupleList)

		} //switch

	} //for

	return result

}

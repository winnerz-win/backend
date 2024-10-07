package abi

import (
	"encoding/hex"
	"jtools/cc"
	"jtools/jmath"
	"strings"

	"golang.org/x/crypto/sha3"
)

const (
	DEBUG_MODE = false
	ABI_CUSTOM = true
)

func MakeInputMethodParam(sl *[]string, args TypeList, isWrite bool) {
	for _, arg := range args {
		switch arg.HasTuple() {
		case 0:
			*sl = append(*sl, arg.String())
		case 1: //tuple
			tl := []string{}
			MakeInputMethodParam(&tl, arg.TupleTypes(), isWrite)
			*sl = append(*sl, "("+strings.Join(tl, ",")+")")

		case 2: //tuple-array
			tl := []string{}
			if isWrite {
				type_list := arg.TupleTypes()
				if len(type_list) > 0 {
					MakeInputMethodParam(&tl, type_list[0].TupleTypes(), isWrite)
				}

			} else {
				MakeInputMethodParam(&tl, arg.TupleTypes(), isWrite)
			}

			*sl = append(*sl, "("+strings.Join(tl, ",")+")[]")
		}
	} //for
}

func (my TypeList) GetBytes(funcName string, isWrite ...bool) []byte {

	hash := sha3.NewLegacyKeccak256()
	if len(my) == 0 {
		hash.Write([]byte(funcName + "()"))
		method_id_bytes := hash.Sum(nil)[:4]
		return method_id_bytes

	}
	is_write := false
	if len(isWrite) > 0 && isWrite[0] {
		is_write = true
	}
	sl := []string{}
	MakeInputMethodParam(&sl, my, is_write)

	if DEBUG_MODE {
		cc.YellowItalic(funcName + "(" + strings.Join(sl, ",") + ")")
	}

	hash.Write([]byte(funcName + "(" + strings.Join(sl, ",") + ")"))
	method_id_bytes := hash.Sum(nil)[:4]

	//make_padded_bytes
	buf := make([]byte, 32*atp_make_pos_size(my))
	index := 0
	for _, v := range my {
		index = atp_make_index_bytes(&buf, index, v)
	}
	//

	method_id_bytes = append(method_id_bytes, buf...)
	return method_id_bytes
}

/////////////////////////////////////////////////////////////////////////

func (my TypeList) ViewBytes() {

	buf := make([]byte, 32*atp_make_pos_size(my))
	index := 0
	for _, v := range my {
		index = atp_make_index_bytes(&buf, index, v)
	}
	if DEBUG_MODE {
		debug_receipt(buf)
	}
}

func atp_make_index_bytes(parent *[]byte, index int, _type Type) int {

	if DEBUG_MODE {
		cc.GreenItalic(_type.name)
	}
	if !_type.writeable {
		cc.Green("[", _type.name, "] writeable is false")
		return index
	}

	offset := index * 32
	switch _type.name {
	case "uint8",
		"uint16",
		"uint32",
		"uint64",
		"uint112",
		"uint128",
		"uint256":
		if _type.Data == nil {
			cc.Gray("[", _type.name, "] skip data")
			return index
		}

		if DEBUG_MODE {
			cc.CyanItalic(jmath.VALUE(_type.Data))
		}

		copy((*parent)[offset:], atp_left_padbytes(jmath.BYTES(_type.Data), 32))
		index++

	case "uint8[]":
		if _type.Data == nil {
			cc.Gray("[", _type.name, "] skip data")
			return index
		}

		item_list := _type.Data.([]uint8)

		buf := atp_set_fixed_array_item_value(
			len(item_list),
			func(buf *[]byte, pos_idx, i int) {
				copy((*buf)[pos_idx:], atp_left_padbytes(jmath.BYTES(item_list[i]), 32))
			},
		)

		copy((*parent)[offset:], atp_offset_bytes(*parent))
		*parent = append(*parent, buf...)

		index++

	case "uint16[]":
		if _type.Data == nil {
			cc.Gray("[", _type.name, "] skip data")
			return index
		}

		item_list := _type.Data.([]uint16)
		buf := atp_set_fixed_array_item_value(
			len(item_list),
			func(buf *[]byte, pos_idx, i int) {
				copy((*buf)[pos_idx:], atp_left_padbytes(jmath.BYTES(item_list[i]), 32))
			},
		)

		copy((*parent)[offset:], atp_offset_bytes(*parent))
		*parent = append(*parent, buf...)

		index++

	case "uint32[]":
		if _type.Data == nil {
			cc.Gray("[", _type.name, "] skip data")
			return index
		}

		item_list := _type.Data.([]uint32)
		buf := atp_set_fixed_array_item_value(
			len(item_list),
			func(buf *[]byte, pos_idx, i int) {
				copy((*buf)[pos_idx:], atp_left_padbytes(jmath.BYTES(item_list[i]), 32))
			},
		)

		copy((*parent)[offset:], atp_offset_bytes(*parent))
		*parent = append(*parent, buf...)

		index++

	case "uint64[]":
		if _type.Data == nil {
			cc.Gray("[", _type.name, "] skip data")
			return index
		}

		item_list := _type.Data.([]uint64)
		buf := atp_set_fixed_array_item_value(
			len(item_list),
			func(buf *[]byte, pos_idx, i int) {
				copy((*buf)[pos_idx:], atp_left_padbytes(jmath.BYTES(item_list[i]), 32))
			},
		)

		copy((*parent)[offset:], atp_offset_bytes(*parent))
		*parent = append(*parent, buf...)

		index++

	case "uint112[]",
		"uint128[]",
		"uint256[]":
		if _type.Data == nil {
			cc.Gray("[", _type.name, "] skip data")
			return index
		}

		item_list := _type.Data.([]string)
		buf := atp_set_fixed_array_item_value(
			len(item_list),
			func(buf *[]byte, pos_idx, i int) {
				copy((*buf)[pos_idx:], atp_left_padbytes(jmath.BYTES(item_list[i]), 32))
			},
		)

		copy((*parent)[offset:], atp_offset_bytes(*parent))
		*parent = append(*parent, buf...)

		index++

	case "bool":
		if _type.Data == nil {
			cc.Gray("[", _type.name, "] skip data")
			return index
		}

		if _type.Data.(bool) {
			(*parent)[offset+31] = 1
		}
		index++

	case "bool[]":
		if _type.Data == nil {
			cc.Gray("[", _type.name, "] skip data")
			return index
		}

		item_list := _type.Data.([]bool)
		buf := atp_set_fixed_array_item_value(
			len(item_list),
			func(buf *[]byte, pos_idx, i int) {
				if item_list[i] {
					(*buf)[pos_idx+31] = 1
				}
			},
		)

		copy((*parent)[offset:], atp_offset_bytes(*parent))
		*parent = append(*parent, buf...)

		index++

	case "address":
		if _type.Data == nil {
			cc.Gray("[", _type.name, "] skip data")
			return index
		}

		if DEBUG_MODE {
			cc.CyanItalic(_type.Data)
		}
		copy((*parent)[offset:], atp_left_padbytes(jmath.BYTES(_type.Data), 32))
		index++

	case "address[]":
		if _type.Data == nil {
			cc.Gray("[", _type.name, "] skip data")
			return index
		}

		item_list := _type.Data.([]string)

		buf := atp_set_fixed_array_item_value(
			len(item_list),
			func(buf *[]byte, pos_idx, i int) {
				copy((*buf)[pos_idx:], atp_left_padbytes(jmath.BYTES(item_list[i]), 32))
			},
		)

		copy((*parent)[offset:], atp_offset_bytes(*parent))
		*parent = append(*parent, buf...)

		index++

	case "byte32":
		if _type.Data == nil {
			cc.Gray("[", _type.name, "] skip data")
			return index
		}

		copy((*parent)[offset:], atp_right_padbytes(jmath.BYTES(_type.Data), 32))
		index++

	case "byte32[]":
		if _type.Data == nil {
			cc.Gray("[", _type.name, "] skip data")
			return index
		}

		item_list := _type.Data.([]string)

		buf := atp_set_fixed_array_item_value(
			len(item_list),
			func(buf *[]byte, pos_idx, i int) {
				copy((*buf)[pos_idx:], atp_right_padbytes(jmath.BYTES(item_list[i]), 32))
			},
		)

		copy((*parent)[offset:], atp_offset_bytes(*parent))
		*parent = append(*parent, buf...)

		index++

	case "string":
		if _type.Data == nil {
			cc.Gray("[", _type.name, "] skip data")
			return index
		}

		copy((*parent)[offset:], atp_offset_bytes(*parent))

		text := []byte(_type.Data.(string))
		*parent = append(*parent,
			atp_left_padbytes(jmath.BYTES(len(text)), 32)...,
		)
		loop := len(text) / 32
		for i := 0; i < loop; i++ {
			pos := i * 32
			*parent = append(*parent,
				text[pos:pos+32]...,
			)
		} //for
		if dot := len(text) % 32; dot > 0 {
			*parent = append(*parent,
				atp_right_padbytes(text[loop*32:], 32)...,
			)
		}
		index++

	case "string[]":
		if _type.Data == nil {
			cc.Gray("[", _type.name, "] skip data")
			return index
		}

		copy((*parent)[offset:], atp_offset_bytes(*parent))

		sl := _type.Data.([]string)

		array_count := atp_left_padbytes(jmath.BYTES(len(sl)), 32)
		*parent = append(*parent, array_count...)

		buf := make([]byte, 32*(len(sl)))

		i := 0
		for _, v := range sl {
			_type := NewString(v)
			i = atp_make_index_bytes(&buf, i, _type)
		}
		*parent = append(*parent, buf...)

		index++

	case "bytes":
		if _type.Data == nil {
			cc.Gray("[", _type.name, "] skip data")
			return index
		}

		copy((*parent)[offset:], atp_offset_bytes(*parent))

		text := _type.Data.([]byte)
		*parent = append(*parent,
			atp_left_padbytes(jmath.BYTES(len(text)), 32)...,
		)
		loop := len(text) / 32
		for i := 0; i < loop; i++ {
			pos := i * 32
			*parent = append(*parent,
				text[pos:pos+32]...,
			)
		} //for
		if dot := len(text) % 32; dot > 0 {
			*parent = append(*parent,
				atp_right_padbytes(text[loop*32:], 32)...,
			)
		}
		index++

	case "bytes[]":
		if _type.Data == nil {
			cc.Gray("[", _type.name, "] skip data")
			return index
		}

		copy((*parent)[offset:], atp_offset_bytes(*parent))

		sl := _type.Data.([][]byte)

		array_count := atp_left_padbytes(jmath.BYTES(len(sl)), 32)
		*parent = append(*parent, array_count...)

		buf := make([]byte, 32*(len(sl)))

		i := 0
		for _, v := range sl {
			_type := NewBytes(v)
			i = atp_make_index_bytes(&buf, i, _type)
		}
		*parent = append(*parent, buf...)

		index++

	case "tuple":
		for _, c := range _type.Datas {
			index = atp_make_index_bytes(parent, index, c.(Type))
		}

	case "tuple_flex":
		buf := make([]byte, 32*atp_make_pos_size(_type.Datas))
		i := 0
		for _, c := range _type.Datas {
			v := c.(Type)
			i = atp_make_index_bytes(&buf, i, v)
		}

		copy((*parent)[offset:], atp_offset_bytes(*parent))
		*parent = append(*parent, buf...)

		index++

	case "tuple_array":

		if len(_type.Datas) > 0 {
			if !_type.Datas[0].(Type).writeable {
				_type.Datas = []interface{}{} //empty data
			}
		}

		buf := make([]byte, 32*(atp_make_pos_size(_type.Datas)+1))

		array_count := atp_left_padbytes(jmath.BYTES(len(_type.Datas)), 32)
		copy(buf, array_count)

		i := 1
		for _, c := range _type.Datas {
			i = atp_make_index_bytes(&buf, i, c.(Type))
		}

		copy((*parent)[offset:], atp_offset_bytes(*parent))
		*parent = append(*parent, buf...)

		index++

	case "tuple_array_flex":
		copy((*parent)[offset:], atp_offset_bytes(*parent))

		pos_size := atp_make_pos_size(_type.Datas)
		*parent = append(
			*parent,
			atp_left_padbytes(jmath.BYTES(pos_size), 32)...,
		)

		buf := make([]byte, 32*pos_size)

		i := 0
		for _, v := range _type.Datas {
			__type := v.(Type)
			i = atp_make_index_bytes(&buf, i, __type)
		}
		*parent = append(*parent, buf...)

		index++

	} //switch

	return index
}

func atp_make_pos_size(void interface{}, isSelf ...bool) int {
	count := 0
	switch list := void.(type) {
	case TypeList:
		for _, v := range list {
			if v.name == "tuple" {
				count += atp_make_pos_size(v.Datas, true)
			} else {
				count++
			}
		}

	case []interface{}:
		for _, v := range list {
			item := v.(Type)
			if item.name == "tuple" {
				count += atp_make_pos_size(item.Datas, true)
			} else {
				count++
			}
		}
	} //switch

	if DEBUG_MODE {
		if len(isSelf) == 0 {
			cc.PurpleItalic("spot :", count)
		}
	}

	return count
}

func atp_offset_bytes(buf []byte) []byte {
	v := jmath.BYTES(len(buf))
	if DEBUG_MODE {
		cc.PurpleItalic("offsetBytes :", jmath.VALUE(v))
	}
	re := make([]byte, 32)
	copy(re[32-len(v):], v)
	return re
}

func atp_left_padbytes(sl []byte, cnt int) []byte {
	if cnt <= len(sl) {
		return sl
	}
	padded := make([]byte, cnt)
	copy(padded[cnt-len(sl):], sl)
	return padded
}
func atp_right_padbytes(sl []byte, cnt int) []byte {
	if cnt <= len(sl) {
		return sl
	}
	padded := make([]byte, cnt)
	copy(padded, sl)

	return padded
}

func atp_set_fixed_array_item_value(count int, f func(buf *[]byte, pos_idx, i int)) []byte {
	buf := make([]byte, 32*(count+1))
	array_count := atp_left_padbytes(jmath.BYTES(count), 32)
	copy(buf, array_count)

	pos_idx := 32
	for i := 0; i < count; i++ {
		f(&buf, pos_idx, i)
		pos_idx += 32
	}

	return buf
}

func debug_receipt(receipt []byte, isNameCut ...bool) {
	v := hex.EncodeToString(receipt)

	isCut := false
	if len(isNameCut) > 0 && isNameCut[0] {
		isCut = true
	}

	if isCut {
		name_area := v[:8]
		cc.Gray("ABI_NAME [", name_area, "]")
		v = v[8:]
	}

	estimateValue := func(v string) string {
		if string(v[0]) != "0" {
			b, _ := hex.DecodeString(v)
			return strings.TrimSpace(string(b))
		} else if strings.HasPrefix(v, "000000000000000000000000") && string(v[24]) != "0" {
			return "0x" + v[24:]
		}
		return jmath.VALUE("0x" + v)
	}

	loop := 0
	for len(v) > 0 {
		//cc.PurpleItalic("[", loop, "]", v[:64])
		cc.PurpleItalic("[", loop, "]", v[:64], "(", estimateValue(v[:64]), ")")
		v = v[64:]
		loop++
	} //for
}

package ebcm

import (
	"encoding/hex"
	"jtools/cc"
	"jtools/cloud/ebcm/abi"
	"jtools/dbg"
)

// type PADBYTES interface {
// 	Bytes() []byte
// 	Hex() string
// 	HexDebug() string
// }

type GasSpeed string

func (my GasSpeed) Value() GasSpeed { return my }
func (my GasSpeed) String() string  { return string(my) }

const (
	GasFastest = GasSpeed("fastest")
	GasFast    = GasSpeed("fast")
	GasAverage = GasSpeed("average")
	GasSafeLow = GasSpeed("safeLow")
	GasBegger  = GasSpeed("begger")
)

type PADBYTES []byte

func (my PADBYTES) Bytes() []byte { return []byte(my) }
func (my PADBYTES) Hex() string {
	return hex.EncodeToString(my.Bytes())
}
func (my PADBYTES) HexDebug() string {
	result := dbg.ENTER
	msg := hex.EncodeToString([]byte(my))
	if msg == "" {
		return "0x"
	}

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

func PADBYTESFromHex(hexString string) PADBYTES {
	buf, err := hex.DecodeString(hexString)
	if err != nil {
		cc.RedItalic(err)
	}
	return PADBYTES(buf)
}

///////////////////////////////////////////////////////////////////////////////////////////

func PadByteETH() PADBYTES {
	return PADBYTES([]byte{})
}

// PadByteTransfer : transfer(address,uint256)
func PadByteTransfer(to string, value string) PADBYTES {
	return MakePadBytesABI(
		"transfer",
		abi.TypeList{
			abi.NewAddress(to),
			abi.NewUint256(value),
		},
	)
}

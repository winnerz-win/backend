package ebcm

import (
	"math/big"
	"txscheduler/brix/tools/dbg"
)

type GasStation struct {
	caller func() GasResult
}

func NewGasStation(c func() GasResult) GasStation {
	return GasStation{
		caller: c,
	}
}

func (my GasStation) Call() GasResult {
	if my.caller != nil {
		return my.caller()
	}
	dbg.RedItalic("ebcm.GasStation.caller is nil")
	return nil
}

type GasResult interface {
	String() string
	GetFast() *big.Int
	GetFastest() *big.Int
	GetSafeLow() *big.Int
	GetAverage() *big.Int
	GetBegger() *big.Int
}

// func (my Sender) GasStation() GasStation {
// 	return NewGasStation(nil)
// }

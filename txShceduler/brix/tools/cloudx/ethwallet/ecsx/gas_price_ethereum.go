package ecsx

import (
	"math/big"
)

func (my *Sender) linkGasPriceFunc() {
	// doTag := func(tag string) bool {
	// 	for _, v := range my.etherTags {
	// 		if v == tag {
	// 			return true
	// 		}
	// 	}
	// 	return false
	// }

	// if doTag("tomo") {
	// 	my.gasFunc = tomoGasPrice
	// 	return
	// }
	// my.gasFunc = SuggestGasPrice
}
func (my *Sender) SetMinGasPrice(wei string) { my.minGasWei = wei }
func (my *Sender) SetGasURL(url string)      { my.gasURL = url }
func (my *Sender) SetLinkGasPriceFunc(f IGasPrice) {
	//	my.gasFunc = f
}

// IGasPrice :
type IGasPrice func(Sender, GasSpeed) *big.Int

func ethGasPrice(sender Sender, speed GasSpeed) *big.Int {
	return NewGasStation().Price(speed)
}

func tomoGasPrice(sender Sender, speed GasSpeed) *big.Int {
	return sender.tomoGasPrice()
}

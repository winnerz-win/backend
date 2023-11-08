package ecsx

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net/http"
	"sync"
	"time"
	"txscheduler/brix/tools/cloudx/ebcmx"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsaa"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
)

/*
	// https://api.etherscan.io/api?module=proxy&action=eth_gasPrice&apikey=YourApiKeyToken
	gasPrice, err := client._SuggestGasPrice(context.Background())
	if err != nil {
		return nil, dbg.MakeError("MakeTokenTx@_SuggestGasPrice", err)
	} else {
		fmt.Println("gasPrice :", gasPrice.String())
	}
*/

type GasSpeed string

func (my GasSpeed) Value() GasSpeed { return my }
func (my GasSpeed) String() string  { return string(my) }

const (
	gasAddFee          = 0
	gasToGweiMuliValue = 1000000000 //1GWEI
	limitDuration      = time.Duration(time.Minute * 1)

	GasFastest = GasSpeed("fastest")
	GasFast    = GasSpeed("fast")
	GasAverage = GasSpeed("average")
	GasSafeLow = GasSpeed("safeLow")
	GasBegger  = GasSpeed("begger")

	averageAdventage = 50 //(5 GWEI ++)
)

var (
	cachingFast = big.NewInt(125)                       //big.NewInt(6 * gasToGweiMuliValue)
	overTime    = time.Now().UTC().Add(time.Hour * -24) //과거 시간으로 초기화
	gasMu       = &sync.RWMutex{}

	pGas = &nGasStation{
		Fast:    float64(125),
		Fastest: float64(125),
		SafeLow: float64(125),
		Average: float64(125),
		Begger:  float64(1),
	}

	isOverLimitCheck = false
	overLimitWei     string
)

// SetOverLimitWEI : called by ecsx.init()
func SetOverLimitWEI(wei string) {
	isOverLimitCheck = true
	overLimitWei = wei
	dbg.PrintForce("overLimitWei :", wei)
}

// nGasStation :
type nGasStation struct {
	Fast    float64 `json:"fast"`
	Fastest float64 `json:"fastest"`
	SafeLow float64 `json:"safeLow"`
	Average float64 `json:"average"`
	Begger  float64 `json:"begger"`
}

func (my nGasStation) String() string {
	mm := map[string]interface{}{
		"fast":    my.GetFast(),
		"fastest": my.GetFastest(),
		"safeLow": my.GetSafeLow(),
		"average": my.GetAverage(),
		"begger":  my.GetBegger(),
	}
	return dbg.ToJSONString(mm)
}

func (my nGasStation) clone() nGasStation {
	return nGasStation{
		Fast:    my.Fast,
		Fastest: my.Fastest,
		SafeLow: my.SafeLow,
		Average: my.Average,
		Begger:  my.Begger,
	}
}

type GasResult struct {
	val *big.Int
}

func (my GasResult) String() string       { return jmath.VALUE(my.val) }
func (my GasResult) GetFast() *big.Int    { return my.val }
func (my GasResult) GetFastest() *big.Int { return my.val }
func (my GasResult) GetSafeLow() *big.Int { return my.val }
func (my GasResult) GetAverage() *big.Int { return my.val }
func (my GasResult) GetBegger() *big.Int  { return my.val }

// GasStation :
type GasStation interface {
	String() string
	ToString() string
	CalcString() string
	Price(speed GasSpeed) *big.Int
	GetFast() *big.Int
	GetFastest() *big.Int
	GetSafeLow() *big.Int
	GetAverage() *big.Int
	GetBegger() *big.Int
}

// ToString :
func (my nGasStation) ToString() string {
	return dbg.ToJSONString(my)
}

// CalcString :
func (my nGasStation) CalcString() string {
	msg := map[GasSpeed]string{}
	msg[GasFastest] = my.GetFastest().String()
	msg[GasFast] = my.GetFast().String()
	msg[GasAverage] = my.GetAverage().String()
	msg[GasSafeLow] = my.GetSafeLow().String()
	msg[GasBegger] = my.GetBegger().String()

	return dbg.ToJSONString(msg)
}

// isOver : refresh
func isOver() bool {
	d := time.Now().UTC().Sub(overTime)
	return d >= limitDuration
}
func refreshTime() {
	overTime = time.Now().UTC()
}

func ebcm_NewGasStation() ebcmx.GasResult {
	return NewGasStation()
}

func (my *Sender) SuggestGasPrice() *big.Int {
	gp, _ := ecsaa.SUGGEST_GAS_PRICE_2(my)
	return gp
}

// NewGasStation :
func NewGasStation() GasStation {
	defer gasMu.Unlock()
	gasMu.Lock()

	if isOver() == true {
		retryCnt := 5

		var err error
		var resp *http.Response
		for {
			if retryCnt == 0 {
				break
			}
			resp, err = http.Get("https://ethgasstation.info/json/ethgasAPI.json")
			if err != nil {
				time.Sleep(time.Millisecond * 500)
				retryCnt--
			} else {
				break
			}
		} //for

		if retryCnt > 0 {
			defer resp.Body.Close()

			buf, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				json.Unmarshal(buf, pGas)
				refreshTime()
				cachingFast = big.NewInt(pGas.GetFast().Int64())

				pGas.Begger = float64(10)
			} else {
				fval := cachingFast.Int64()
				pGas.Fast = float64(fval)
				pGas.Fastest = float64(fval)
				pGas.SafeLow = float64(fval)
				pGas.Average = float64(fval)
				pGas.Begger = float64(1)
			}
		} else {
			fval := cachingFast.Int64()
			pGas.Fast = float64(fval)
			pGas.Fastest = float64(fval)
			pGas.SafeLow = float64(fval)
			pGas.Average = float64(fval)
			pGas.Begger = float64(1)
		}
	} //if
	return pGas.clone()
}

// Price :
func (my nGasStation) Price(speed GasSpeed) *big.Int {
	var price *big.Int
	switch speed {
	case GasFastest:
		price = my.GetFastest()
	case GasFast:
		price = my.GetFast()
	case GasAverage:
		price = my.GetAverage()
		_ = price
		//dbg.Purple(price.String())

		price = getGasGwei(my.Average + averageAdventage) //5 gwei ++
		//dbg.Purple(price.String())

		if price.Cmp(my.GetFast()) > 0 {
			price = my.GetFast()
		}
	case GasSafeLow:
		price = my.GetSafeLow()

	case GasBegger:
		price = my.GetBegger()
		dbg.Yellow("nGasStation :: GasBeger :", price)

	default:
		price = my.GetFast()
	} //switch

	if isOverLimitCheck {
		if jmath.CompareTo(price.String(), overLimitWei) > 0 {
			dbg.Red("[EtherClient.nGasStation] Price Limit :", price.String(), "/", overLimitWei)
			v := big.NewInt(0)
			v.SetString(overLimitWei, 10)
			price = v
		}
	}

	return price
}

// GetFast :
func (my nGasStation) GetFast() *big.Int {
	return getGasGwei(my.Fast)
}

// GetFastest :
func (my nGasStation) GetFastest() *big.Int {
	return getGasGwei(my.Fastest)
}

// GetSafeLow :
func (my nGasStation) GetSafeLow() *big.Int {
	return getGasGwei(my.SafeLow)
}

// GetBegger :
func (my nGasStation) GetBegger() *big.Int {
	//return getGasGwei(my.Begger) // 3GWEI

	v := int64(gasToGweiMuliValue) * 1 // 1GWEI
	return big.NewInt(v)
}

// GetAverage :
func (my nGasStation) GetAverage() *big.Int {
	return getGasGwei(my.Average)
}

func getGasGwei(price float64) *big.Int {
	f := big.NewFloat(price)
	i64, _ := f.Int64()
	v := ((i64 / 10) + gasAddFee) * gasToGweiMuliValue
	return big.NewInt(v)
}

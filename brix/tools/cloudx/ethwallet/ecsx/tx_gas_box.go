package ecsx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"txscheduler/brix/tools/cloudx/ethwallet/ecsaa"
	"txscheduler/brix/tools/jmath"

	"txscheduler/brix/tools/dbg"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

// GasBoxData :
type GasBoxData struct {
	isETH        bool
	tokenAddress string
	fromAddress  string
	toAddress    string
	wei          string

	inputData []byte

	speed    GasSpeed
	gasLimit uint64
	gasPrice *big.Int
	Error    error

	etherTags []string
}

func (my GasBoxData) String() string {
	msg := map[string]interface{}{
		"is_eth":       my.isETH,
		"tokenAddress": my.tokenAddress,
		"fromAddress":  my.fromAddress,
		"toAddress":    my.toAddress,
		"wei":          my.wei,
		"input":        string(my.inputData),
		"speed":        my.speed,
		"gas_limit":    my.gasLimit,
		"gas_price":    my.gasPrice.String(),
		"error":        my.Error,
	}
	b, _ := json.MarshalIndent(msg, "", "  ")
	return string(b)
}

// GetWei : 전송 금액
func (my GasBoxData) GetWei() string {
	return my.wei
}

// SetWei : 전송 금액 수정
func (my *GasBoxData) SetWei(wei string) {
	my.wei = wei
	if my.isETH == false {
		transferFnSignature := []byte("transfer(address,uint256)")
		hash := sha3.NewLegacyKeccak256()
		hash.Write(transferFnSignature)
		methodID := hash.Sum(nil)[:4]

		toHexAddress := common.HexToAddress(strings.ToLower(my.toAddress))
		paddedAddress := common.LeftPadBytes(toHexAddress.Bytes(), 32)

		amount := new(big.Int)
		amount.SetString(wei, 10)
		paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

		var data []byte
		data = append(data, methodID...)
		data = append(data, paddedAddress...)
		data = append(data, paddedAmount...)

		my.inputData = data
	}
}

// ChangeWeiSubGasForETH : 기존 전송금액에 가스비를 제외시킴.( 전송금액 수정 ) ---이더전송만 가능
func (my *GasBoxData) ChangeWeiSubGasForETH() bool {
	if my.isETH == false {
		return false
	}
	tryWei := my.GetWei()
	gasWei := my.GasWei()
	editWei := jmath.SUB(tryWei, gasWei)
	if jmath.CMP(editWei, 0) <= 0 {
		return false
	}
	my.SetWei(editWei)
	return true
}

// GasBoxDataList :
type GasBoxDataList []GasBoxData

// SpendETH : // [ 이더: 전송량 + 가스비 ] ,  [ 토큰 : 가스비 ]
func (my GasBoxData) SpendETH() string {
	if my.isETH {
		a := WeiToETH(my.wei)
		return jmath.ADD(a, my.GasETH())
	}
	return my.GasETH()
}

// SpendWEI : 전송량 + 가스비
func (my GasBoxData) SpendWEI() string {
	if my.isETH {
		jmath.ADD(my.wei, my.GasWei())
	}
	return my.GasWei()
}

// Limit : LIMIT_PRICE
func (my GasBoxData) Limit() string { return fmt.Sprint(my.gasLimit) }

// SetLimit : LIMIT_PRICE
func (my *GasBoxData) SetLimit(limit string) {
	my.gasLimit = jmath.Uint64(limit)
}

// Price : ETH_price
func (my GasBoxData) Price() string {
	price := big.NewInt(my.gasPrice.Int64())
	v := jmath.VALUE(price)
	return WeiToETH(v)
}

// SetPrice : ETH_price
func (my *GasBoxData) SetPrice(ethPrice string) {
	v := jmath.Int64(ETHToWei(ethPrice))
	my.gasPrice = big.NewInt(v)
}

// GasWei : 가스비
func (my GasBoxData) GasWei() string {
	price := big.NewInt(my.gasPrice.Int64())
	limit := big.NewInt(int64(my.gasLimit))
	fee := price.Mul(price, limit)
	return fee.String()
}

// GasETH : 가스비
func (my GasBoxData) GasETH() string { return WeiToETH(my.GasWei()) }

type GasSnapShot struct {
	Limit      uint64 `bson:"limit" json:"limit"`
	Price      string `bson:"price" json:"price"`
	FeeWei     string `bson:"fee_wei" json:"fee_wei"` // limit * price
	FixedNonce uint64 `bson:"fixedNonce,omitempty" json:"fixedNonce,omitempty"`
}

func (my GasSnapShot) String() string { return dbg.ToJSONString(my) }
func (my GasSnapShot) FeeETH() string { return WeiToETH(my.FeeWei) }

func (my GasSnapShot) priceBig() *big.Int {
	return big.NewInt(jmath.Int64(my.Price))
}

// SnapShot : 가스비 스냅샷 찍음
func (my GasBoxData) SnapShot() GasSnapShot {
	item := GasSnapShot{
		Limit:  my.gasLimit,
		Price:  jmath.VALUE(my.gasPrice.Int64()),
		FeeWei: my.GasWei(),
	}
	return item
}

// ApplySnapShot : 가스비 스냅샷 적용
func (my *GasBoxData) ApplySnapShot(ss GasSnapShot) {
	my.gasLimit = ss.Limit
	my.gasPrice.SetString(ss.Price, 10)
}

// GasBox :
func (my Sender) GasBox(
	tokenAddress,
	fromAddress,
	toAddress string,
	wei string,
	speed GasSpeed,
	isIncludeGas ...bool, //전송금액에 가스비 포함여부 (sendWei = wei - gas) --이더 전송일때만
) GasBoxData {
	tokenAddress = dbg.TrimToLower(tokenAddress)
	fromAddress = dbg.TrimToLower(fromAddress)
	toAddress = dbg.TrimToLower(toAddress)

	isETH := false

	if my.isEtherTag(strings.ToLower(tokenAddress)) {
		isETH = true
	}
	boxdata := GasBoxData{
		isETH:        isETH,
		tokenAddress: tokenAddress,
		fromAddress:  fromAddress,
		toAddress:    toAddress,
		wei:          wei,
		speed:        speed,
		Error:        nil,
	}
	for _, tag := range my.etherTags {
		boxdata.etherTags = append(boxdata.etherTags, tag)
	}

	if isETH {
		fromHexAddress := common.HexToAddress(strings.ToLower(fromAddress))
		toHexAddress := common.HexToAddress(strings.ToLower(toAddress))

		value := new(big.Int)
		value.SetString(wei, 10)
		gasLimit, err := my.client.EstimateGas(context.Background(), ethereum.CallMsg{
			From:  fromHexAddress,
			To:    &toHexAddress,
			Value: value,
			Data:  []byte{},
		})
		if err != nil {
			boxdata.Error = fmt.Errorf("Sender.GasBox : %v", err)
			return boxdata
		}
		//gasLimit := uint64(21000) // in units

		boxdata.gasLimit = gasLimit
		boxdata.gasPrice = ecsaa.SUGGEST_GAS_PRICE(my)

		if boxdata.gasPrice == nil {
			boxdata.Error = errors.New("Sender.GasPrice func nil")
			return boxdata
		}

	} else {
		tokenHexAddress := common.HexToAddress(strings.ToLower(tokenAddress))
		fromHexAddress := common.HexToAddress(strings.ToLower(fromAddress))
		toHexAddress := common.HexToAddress(strings.ToLower(toAddress))

		paddedAddress := common.LeftPadBytes(toHexAddress.Bytes(), 32)

		amount := new(big.Int)
		amount.SetString(wei, 10)
		paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

		transferFnSignature := []byte("transfer(address,uint256)")
		hash := sha3.NewLegacyKeccak256()
		hash.Write(transferFnSignature)
		methodID := hash.Sum(nil)[:4]

		var data []byte
		data = append(data, methodID...)
		data = append(data, paddedAddress...)
		data = append(data, paddedAmount...)

		gasLimit, err := my.client.EstimateGas(context.Background(), ethereum.CallMsg{
			From: fromHexAddress,
			To:   &tokenHexAddress,
			//Value: amount,	//이더value만 들어간다.
			Data: data,
		})
		if err != nil {
			boxdata.Error = fmt.Errorf("Sender.GasBox : %v", err)
			return boxdata
		}

		boxdata.inputData = data
		boxdata.gasLimit = gasLimit
		boxdata.gasPrice = ecsaa.SUGGEST_GAS_PRICE(my)

		// addValue := big.NewInt(EtherClient.TXLimitTokenAddValue)
		// boxdata.gasPrice = boxdata.gasPrice.Add(boxdata.gasPrice, addValue)

	}

	// 전송하려는 금액에서 가스비를 고려해서 뺀다..
	if isETH {
		if dbg.IsTrue2(isIncludeGas...) {
			boxdata.SetWei(jmath.SUB(boxdata.GetWei(), boxdata.GasWei()))
		}
	}

	return boxdata
}

func (my Sender) NonceNumber(address string) (uint64, error) {
	_address := common.HexToAddress(strings.ToLower(address))
	nonce, err := my.client.PendingNonceAt(context.Background(), _address)
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

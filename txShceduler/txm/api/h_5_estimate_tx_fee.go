package api

import (
	"net/http"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/txm/ack"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func init() {
	hEstimateMasterTxFee()
}

type EMT_ITEM struct {
	Symbol        string `json:"symbol"`
	EstimateTxFee string `json:"estimate_tx_fee"`
	IsChainError  bool   `json:"is_chain_error,omitempty"`
	ErrorMessage  string `json:"error_message,omitempty"`
}

func (EMT_ITEM) TagString() []string {
	return []string{
		"symbol", "",
		"estimate_tx_fee", "예측 TX수수료(가스비) -- 0일경우 실패",
		"is_chain_error,omitempty", "예측 실패일경우 true",
		"is_chain_error,omitempty", "실패 사유",
	}
}

type ACK_EMT struct {
	Symbols map[string]EMT_ITEM `json:"symbols"`
}

func (ACK_EMT) TagString() []string {
	return []string{
		"symbols", "심볼별 TX수수료 예측 결과",
		"", "",
	}
}

func hEstimateMasterTxFee() {

	method := chttp.GET
	url := model.V1 + "/estimate/master/tx_fee"
	Doc().Comment("[ 마스터지갑에서 외부 지갑주소로 출금시 트랜잭션 수수료 예측 ]").
		Method(method).URL(url).
		JAckOK(ACK_EMT{}, ACK_EMT{}.TagString()...).
		ETCVAL(EMT_ITEM{}, EMT_ITEM{}.TagString()...).
		JAckError(ack.ChainGasPrice, "네트워크 가스 예측 실패시").
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {

			sender := Sender()
			gasPrice := sender.SUGGEST_GAS_PRICE(ecsx.GasFast)
			if gasPrice.Error() != nil {
				chttp.Fail(w, ack.ChainGasPrice)
				return
			}

			result := ACK_EMT{
				Symbols: map[string]EMT_ITEM{},
			}

			from := inf.Master()
			for _, token := range inf.Config().Tokens {

				item := EMT_ITEM{
					Symbol: token.Symbol,
				}

				limit := uint64(21000)
				var err error
				if !token.IsCoin {
					padBytes := ecsx.PadBytesTransfer(
						"0x0000000000000000000000000000000000000000",
						"1",
					)
					limit, err = sender.XGasLimit(
						padBytes,
						from.Address,
						token.Contract,
						model.ZERO,
					)
					if err != nil {
						limit = 0

						item.IsChainError = true
						item.ErrorMessage = err.Error()
					}
					item.EstimateTxFee = gasPrice.FeeETH(limit)

				} else {
					item.EstimateTxFee = gasPrice.FeeETH(limit)
				}

				result.Symbols[token.Symbol] = item

			} //for

			chttp.OK(w, result)

		},
	)
}

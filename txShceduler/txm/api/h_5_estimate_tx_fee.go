package api

import (
	"context"
	"jtools/cloud/ebcm"
	"jtools/jmath"
	"net/http"
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

	Limit string `json:"limit,omitempty"`
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

			ctx := context.Background()

			sender := Caller()

			gas_price, err := sender.SuggestGasPrice(ctx)
			if err != nil {
				chttp.Fail(w, ack.ChainGasPrice)
				return
			}

			result := ACK_EMT{
				Symbols: map[string]EMT_ITEM{},
			}

			master := inf.Master()
			for _, token := range inf.Config().Tokens {

				item := EMT_ITEM{
					Symbol: token.Symbol,
				}

				limit_text := ""
				if !token.IsCoin {

					padBytes := ebcm.PadByteTransfer(
						ebcm.AddressONE,
						"1",
					)

					limit, err := sender.EstimateGas(
						ctx,
						ebcm.MakeCallMsg(
							master.Address,
							token.Contract,
							model.ZERO,
							padBytes,
						),
					)

					if err != nil {
						limit = 0

						item.IsChainError = true
						item.ErrorMessage = err.Error()
					}

					item.EstimateTxFee = gas_price.EstimateGasFeeETH(limit)
					limit_text = jmath.VALUE(limit)

				} else {
					limit := uint64(21000)
					item.EstimateTxFee = gas_price.EstimateGasFeeETH(limit)
					limit_text = jmath.VALUE(limit)

				}

				item.Limit = limit_text
				result.Symbols[token.Symbol] = item

			} //for

			chttp.OK(w, result)

		},
	)
}

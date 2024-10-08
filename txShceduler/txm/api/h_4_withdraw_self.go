package api

import (
	"context"
	"jcloudnet/itype"
	"jtools/cc"
	"jtools/cloud/ebcm"
	"jtools/cloud/jeth/ecs"
	"jtools/jmath"
	"jtools/mms"
	"net/http"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/database/mongo/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/txm/ack"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func init() {
	hUserWithdrawInfo()
	hUserWithdrawTry()
	hUserWithdrawResult()
}

func Caller() *itype.IClient {
	return itype.New(ecs.RPC_URL(inf.Mainnet()), false, inf.InfuraKey())
}

func Sender() *ebcm.Sender {
	s := Caller()
	return s.EBCMSender(ecs.TxSigner{})
}

func hUserWithdrawInfo() {
	/*
		Comment : 개인지갑에서 직접 출금을 위한 ( 잔액조회 및 가스비 계산)
		Method : POST
		URL : http://host/v1/user/withdraw_info
		Param :
		{
			"from_address" : string			// 회원의 개인지갑 주소
			"symbol" : string				// 심볼 (ETH / GDG) --> 기획상 ETH만 출금
			"price" : string				// 출금 금액 (decimal 제외한 값)
			"to_address" : string			// 수신자 주소
		}
		Response :
		{
			"success" : true,
			"data" : {
				"symbol"` : string	//요청 심볼
				"remain_symbol_price"` : string	// 개인지갑 주소가 보유한 심볼 토큰 가격
				"remain_eth_price"` : string 	// 개인지갑 주소가 보유한 이더 가격
				"estimate_gas_fee"` : string 	// 예상되는 가스 소모량
				"is_chain_error"` : bool		// 가스비 계산 에러 여부 ( true 이면 estimate_gas_fee 값은 "0"이된다.)
			}
		}
		====================================================
		개인지갑이 보유한 금액 보다 큰금액을 전송하려 하면 가스비 계산시 Exception이 뜬다.
	*/
	type CDATA struct {
		FromAddress string `json:"from_address"`
		Symbol      string `json:"symbol"`
		Price       string `json:"price"`
		ToAddress   string `json:"to_address"`
	}
	type RESULT struct {
		Symbol            string `json:"symbol"`
		RemainSymbolPrice string `json:"remain_symbol_price"`
		RemainETHPrice    string `json:"remain_eth_price"`
		EstimateGasFee    string `json:"estimate_gas_fee"`
		IsChainError      bool   `json:"is_chain_error"`
	}
	method := chttp.POST
	url := model.V1 + "/user/withdraw_info"

	Doc().Comment("개인지갑에서 직접 출금을 위한 ( 잔액조회 및 가스비 계산)").
		Method(method).URL(url).
		JParam(CDATA{},
			"from_address", "회원의 개인지갑 주소",
			"symbol", "심볼 (ETH / GDG) --> 기획상 ETH만 출금",
			"price", "출금 금액 (decimal 제외한 값)",
			"to_address", "수신자 주소",
		).
		JAckOK(RESULT{},
			"symbol", "요청 심볼",
			"remain_symbol_price", "개인지갑 주소가 보유한 심볼 토큰 가격",
			"remain_eth_price", "개인지갑 주소가 보유한 이더 가격",
			"estimate_gas_fee", "예상되는 가스 소모량",
			"is_chain_error", `가스비 계산 에러 여부 ( true 이면 estimate_gas_fee 값은 "0"이된다.)`,
			"", "",
		).
		Etc("", `_
			<cc_blue>
			====================================================
			개인지갑이 보유한 금액 보다 큰금액을 전송하려 하면 가스비 계산시 Exception이 뜬다.
			</cc_blue>
		`).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := CDATA{}
			chttp.BindingJSON(req, &cdata)

			model.Trim(&cdata.FromAddress)
			model.Trim(&cdata.ToAddress)
			if jmath.CMP(cdata.Price, 0) <= 0 {
				chttp.Fail(w, ack.UnderZERO)
				return
			}
			if !ebcm.IsAddress(cdata.FromAddress) {
				chttp.Fail(w, ack.InvalidAddress)
				return
			}
			if !ebcm.IsAddress(cdata.ToAddress) {
				chttp.Fail(w, ack.InvalidAddress)
				return
			}

			if !inf.ValidSymbol(cdata.Symbol) {
				chttp.Fail(w, ack.NotFoundSymbol)
				return
			}

			token := inf.Config().Tokens.GetSymbol(cdata.Symbol)
			_ = token

			model.DB(func(db mongo.DATABASE) {
				member := model.LoadMemberAddress(db, cdata.FromAddress)
				if !member.Valid() {
					chttp.Fail(w, ack.NotFoundAddress)
					return
				}

				cloudETH := model.ZERO
				cloudSymbol := model.ZERO
				estimateFEE := model.ZERO

				caller := Caller()
				member_eth_price := caller.GetCoinPrice(member.Address)
				if cdata.Symbol == model.ETH {
					cloudSymbol = member_eth_price
				}
				cloudSymbol = caller.Price(member.Address, token.Contract, token.Decimal)

				if jmath.CMP(member_eth_price, 0) <= 0 {
					chttp.OK(w, RESULT{
						Symbol:            cdata.Symbol,
						RemainSymbolPrice: cloudSymbol,
						RemainETHPrice:    cloudETH,
						EstimateGasFee:    estimateFEE,
					})
					return
				}

				cloudETH = member_eth_price

				var padBytes ebcm.PADBYTES
				sendWEI := model.ZERO
				send_To := ""

				if cdata.Symbol == model.ETH {
					if jmath.CMP(member_eth_price, cdata.Price) <= 0 {
						chttp.OK(w, RESULT{
							Symbol:            cdata.Symbol,
							RemainSymbolPrice: cloudSymbol,
							RemainETHPrice:    cloudETH,
							EstimateGasFee:    estimateFEE,
						})
						return
					}
					member_eth_price = jmath.SUB(member_eth_price, cdata.Price)
					_ = member_eth_price

					padBytes = ebcm.PadByteETH()
					sendWEI = ebcm.TokenToWei(cdata.Price, 18)

					send_To = cdata.ToAddress

				} else {
					tokenPRICE := caller.Price(member.Address, token.Contract, token.Decimal)
					if jmath.CMP(tokenPRICE, cdata.Price) < 0 {
						chttp.OK(w, RESULT{
							Symbol:            cdata.Symbol,
							RemainSymbolPrice: cloudSymbol,
							RemainETHPrice:    cloudETH,
							EstimateGasFee:    estimateFEE,
						})
						return
					}
					padBytes = ebcm.PadByteTransfer(
						cdata.ToAddress,
						ebcm.TokenToWei(cdata.Price, token.Decimal),
					)
					sendWEI = model.ZERO

					send_To = token.Contract

				}

				sender := Sender()

				ctx := context.Background()
				limit, err := sender.EstimateGas(ctx,
					ebcm.MakeCallMsg(
						member.Address,
						send_To,
						sendWEI,
						padBytes,
					),
				)
				if err != nil {
					chttp.OK(w, RESULT{
						Symbol:            cdata.Symbol,
						RemainSymbolPrice: cloudSymbol,
						RemainETHPrice:    cloudETH,
						EstimateGasFee:    estimateFEE,
						IsChainError:      true,
					})
					return
				}

				limit = ebcm.MMA_LimitBuffer(limit)

				gas_price, err := sender.SuggestGasPrice(ctx, true)
				if err != nil {
					chttp.OK(w, RESULT{
						Symbol:            cdata.Symbol,
						RemainSymbolPrice: cloudSymbol,
						RemainETHPrice:    cloudETH,
						EstimateGasFee:    estimateFEE,
						IsChainError:      true,
					})
				}

				gas_price.EstimateGasFeeETH(limit)

				chttp.OK(w,
					RESULT{
						Symbol:            cdata.Symbol,
						RemainSymbolPrice: cloudSymbol,
						RemainETHPrice:    cloudETH,
						EstimateGasFee:    gas_price.EstimateGasFeeETH(limit),
					},
				)

			})

		},
	)
}

func hUserWithdrawTry() {
	/*
		Comment : 개인지갑에서 직접 출금 신청
		Method : POST
		URL : http://host/v1/user/withdraw_try
		Param :
		{
			"from_address" : string			// 회원의 개인지갑 주소
			"symbol" : string				// 심볼 (ETH / GDG) --> 기획상 ETH만 출금
			"price" : string				// 출금 금액 (decimal 제외한 값)
			"to_address" : string			// 수신자 주소
		}
		Response :
		{
			"success" : true,
			"data" : {
				"estimate_gas_fee" : string		// 출금시 예상되는 가스량
				"receipt_code" : string			// 출금 신청 영수증
				"hash" : string					// 출금 트랜젝션 해시값
			}

		}
	*/
	type CDATA struct {
		FromAddress string `json:"from_address"`
		Symbol      string `json:"symbol"`
		Price       string `json:"price"`
		ToAddress   string `json:"to_address"`
	}
	type RESULT struct {
		EstimateGasFee string `json:"estimate_gas_fee"`
		ReceiptCode    string `json:"receipt_code"`
		Hash           string `json:"hash"`
	}

	method := chttp.POST
	url := model.V1 + "/user/withdraw_try"

	Doc().Comment("개인지갑에서 직접 출금 신청").
		Method(method).URL(url).
		JParam(CDATA{},
			"from_address", "회원의 개인지갑 주소",
			"symbol", "심볼 (ETH / GDG) --> 기획상 ETH만 출금",
			"price", "출금 금액 (decimal 제외한 값)",
			"to_address", "수신자 주소",
		).
		JAckOK(RESULT{},
			"estimate_gas_fee", "출금시 예상되는 가스량",
			"receipt_code", `출금 신청 영수증`,
			"hash", `출금 트랜젝션 해시값`,
			"", "",
		).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := CDATA{}
			chttp.BindingJSON(req, &cdata)

			model.Trim(&cdata.FromAddress)
			model.Trim(&cdata.ToAddress)
			if jmath.CMP(cdata.Price, 0) <= 0 {
				chttp.Fail(w, ack.UnderZERO)
				return
			}
			if !ebcm.IsAddress(cdata.FromAddress) {
				chttp.Fail(w, ack.InvalidAddress)
				return
			}
			if !ebcm.IsAddress(cdata.ToAddress) {
				chttp.Fail(w, ack.InvalidAddress)
				return
			}

			if !inf.ValidSymbol(cdata.Symbol) {
				chttp.Fail(w, ack.NotFoundSymbol)
				return
			}

			token := inf.Config().Tokens.GetSymbol(cdata.Symbol)
			_ = token

			model.DB(func(db mongo.DATABASE) {
				member := model.LoadMemberAddress(db, cdata.FromAddress)
				if !member.Valid() {
					chttp.Fail(w, ack.NotFoundAddress)
					return
				}

				if user_tx_ok := model.UserTransactionStart(
					db,
					member.Address,
					url,
					func() bool {

						TAG := dbg.Cat("[SELF_WITHDRAW_API](", member.Address, ") ")
						cc.Cyan(TAG, "UserTransactionStart ------- START")
						defer cc.Cyan(TAG, "UserTransactionStart ------- END")

						userGasETH := Caller().GetCoinPrice(member.Address)
						if jmath.CMP(userGasETH, 0) <= 0 {
							cc.RedItalic(TAG, "GasFee ETH is 0.")
							chttp.Fail(
								w,
								ack.NotEnoughGasPrice,
								dbg.Cat("user_coin_price:", userGasETH),
							)
							return false
						}

						to_address := cdata.ToAddress
						var padBytes ebcm.PADBYTES
						sendWEI := model.ZERO
						if cdata.Symbol == model.ETH {
							if jmath.CMP(userGasETH, cdata.Price) <= 0 {
								cc.RedItalic(TAG, "Less SendTry ETH Price. (", userGasETH, "/", cdata.Price, ")")
								chttp.Fail(w, ack.NotEnoughTargetPrice)
								return false
							}
							userGasETH = jmath.SUB(userGasETH, cdata.Price)

							padBytes = ebcm.PadByteETH()
							sendWEI = ebcm.ETHToWei(cdata.Price)

						} else {
							to_address = token.Contract //////////

							tokenPRICE := Caller().Price(member.Address, token.Contract, token.Decimal)
							if jmath.CMP(tokenPRICE, cdata.Price) < 0 {
								cc.RedItalic(TAG, "Less SendTry ", token.Symbol, " Price. (", tokenPRICE, "/", cdata.Price, ")")
								chttp.Fail(w, ack.NotEnoughTargetPrice)
								return false
							}
							padBytes = ebcm.PadByteTransfer(
								cdata.ToAddress,
								ebcm.TokenToWei(cdata.Price, token.Decimal),
							)

						}

						sender := Sender()

						nonce, err := ebcm.MMA_GetNonce(
							sender,
							member.Address,
							true,
						)
						if err != nil {
							cc.RedItalic(TAG, "MMA_GetNonce : ", err)
							chttp.Fail(w, ack.ChainNonce, err.Error())
							return false
						}

						ctx := context.Background()
						limit, err := sender.EstimateGas(
							ctx,
							ebcm.MakeCallMsg(
								member.Address,
								to_address,
								sendWEI,
								padBytes,
							),
						)
						if err != nil {
							cc.RedItalic(TAG, "EstimateGas : ", err)
							chttp.Fail(w, ack.ChainGasLimit)
							return false
						}

						limit = ebcm.MMA_LimitBuffer(limit)

						gas_price, err := sender.SuggestGasPrice(ctx, true)
						if err != nil {
							cc.RedItalic(TAG, "SuggestGasPrice : ", err)
							chttp.Fail(w, ack.ChainGasPrice)
							return false
						}

						feeETH := gas_price.EstimateGasFeeETH(limit)
						if jmath.CMP(userGasETH, feeETH) < 0 {
							cc.RedItalic(TAG, "Less SendTry ", token.Symbol, " Price. (", userGasETH, "/", feeETH, ")")
							chttp.Fail(w, ack.NotEnoughGasPrice)
							return false
						}

						ntx := sender.NewTransaction(
							nonce,
							to_address,
							sendWEI,
							limit,
							gas_price,
							padBytes,
						)

						stx, err := sender.SignTx(
							ntx,
							member.PrivateKey(),
						)
						if err != nil {
							cc.RedItalic(TAG, "SignTx : ", err)
							chttp.Fail(w, ack.ChainSTX)
							return false
						}

						hash, err := sender.SendTransaction(
							ctx,
							stx,
						)

						if err != nil {
							cc.RedItalic(TAG, "SendTransaction : ", err)
							chttp.Fail(w, ack.ChainSend)
							return false
						}

						nowAt := mms.Now()
						data := model.TxETHWithdraw{
							UID:       member.UID,
							ToAddress: cdata.ToAddress,
							ToPrice:   cdata.Price,

							Hash:     hash,
							GasLimit: jmath.VALUE(limit),
							GasPrice: gas_price.GET_GAS_ETH(),
							Gas:      feeETH,

							Symbol:    cdata.Symbol,
							Decimal:   token.Decimal,
							Timestamp: nowAt,
							YMD:       nowAt.YMD(),

							State: model.TxStatePendingSELF, //22
						}
						receiptCode := data.InsertSELF_DB(db)

						cc.Cyan(TAG, "Hash(", hash, ")")
						r := RESULT{
							ReceiptCode:    receiptCode,
							EstimateGasFee: feeETH,
							Hash:           hash,
						}
						chttp.OK(w, r)

						return true
					},
				); !user_tx_ok {
					chttp.Fail(w, ack.NotYetSendTx)
					return
				}

			})

		},
	)
}

func hUserWithdrawResult() {
	/*
		Comment : 개인지갑 출금신청 결과 확인 요청
		Method : GET
		URL : http://host/v1/user/withdraw_result/[receipt_code]

		Response :
		{
			"success" : true
			"data" : {
				"receipt_code"` : string
				"name"` : string
				"uid"` : long
				"from_address"` : string
				"to_address"` :string
				"hash"` : string
				"symbol"` : string
				"price"` : string
				"gas" : string
				"withdraw_state"` : int  // 22:진행중 , 104:실패 , 200:성공
				"timestamp"`
				"is_send"`
				"send_at"`
			}

		}
	*/
	type RESULT struct {
		ReceiptCode   string        `json:"receipt_code"`
		Name          string        `json:"name"`
		UID           int64         `json:"uid"`
		FromAddress   string        `json:"from_address"`
		ToAddress     string        `json:"to_address"`
		Hash          string        `json:"hash"`
		Symbol        string        `json:"symbol"`
		Price         string        `json:"price"`
		Gas           string        `json:"gas"`
		WithdrawState model.TxState `json:"withdraw_state"`
		Timestamp     mms.MMS       `json:"timestamp"`
		IsSend        bool          `json:"is_send"`
		SendAt        mms.MMS       `json:"send_at"`
	}

	method := chttp.GET
	url := model.V1 + "/user/withdraw_result/:args"

	Doc().Comment("개인지갑 출금신청 결과 확인 요청").
		Method(method).URL(url).
		JAckOK(RESULT{},
			"", "",
			"withdraw_state", `22:진행중 , 104:실패 , 200:성공`,
			"", "",
		).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {

			receiptCode := ps.ByName("args")

			if receiptCode == "" {
				chttp.Fail(w, ack.InvalidReceiptCode)
				return
			}

			model.DB(func(db mongo.DATABASE) {
				witem := model.TxETHWithdraw{}.GetData(db, receiptCode)
				if witem.Valid() {
					member := model.LoadMember(db, witem.UID)
					log := witem.MakeLogWithdraw(member, model.TxStatePending)
					chttp.OK(w, RESULT{
						ReceiptCode:   log.ReceiptCode,
						Name:          log.Name,
						UID:           log.UID,
						FromAddress:   log.Address,
						ToAddress:     log.ToAddress,
						Hash:          log.Hash,
						Symbol:        log.Symbol,
						Price:         log.ToPrice,
						Gas:           log.Gas,
						WithdrawState: log.State,
						Timestamp:     log.Timestamp,
						IsSend:        log.IsSend,
						SendAt:        log.SendAt,
					})

					return
				}

				log := model.LogWithdrawSELF{}.GetData(db, receiptCode)
				if log.Valid() {
					chttp.OK(w, RESULT{
						ReceiptCode:   log.ReceiptCode,
						Name:          log.Name,
						UID:           log.UID,
						FromAddress:   log.Address,
						ToAddress:     log.ToAddress,
						Hash:          log.Hash,
						Symbol:        log.Symbol,
						Price:         log.ToPrice,
						Gas:           log.Gas,
						WithdrawState: log.State,
						Timestamp:     log.Timestamp,
						IsSend:        log.IsSend,
						SendAt:        log.SendAt,
					})
					return
				}

				chttp.Fail(w, ack.InvalidReceiptCode)
			})
		},
	)
}

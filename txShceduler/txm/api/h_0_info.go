package api

import (
	"context"
	"jtools/cloud/ebcm"
	"jtools/jmath"
	"net/http"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/database/mongo/tools/cc"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/txm/ack"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func init() {
	hRoot()
	hVersion(false)
	hInfoGasFee()
	hInfoMaster()
	hInfoMemberName()
	hInfoMemberUID()
	hInfoMemberAddress()
}

func hRoot() {
	hVersion(true)
}

func hVersion(isRoot bool) {
	method := chttp.GET
	url := "/version"
	if isRoot {
		url = "/"
	}

	type RESULT struct {
		Version string `json:"version"`
		Server  string `json:"server"`
		Mainnet bool   `json:"mainnet"`
		IP      string `json:"ip,omitempty"`
	}
	Doc().Comment("[ 버전정보 ] 스케줄러 서버 버전 요청").
		Method(method).URL(url).
		Etc(".", `_
			<cc_blue>응답결과</cc_blue>
			{
				"success" : true,
				"data" : {
					"mainnet" : true,			// 메인넷 / 테스트넷 여부
					"server" : "scheduler",		// 서버 이름 (프로젝트에 따라 변경됩니다.)
					"version" : "2021.04.04"	// 버전 정보
				}
			}
		`).
		Apply()
	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			_ip := ""
			ip, dir, err := chttp.GetIP(req)
			cc.Green("[", dir, "]", ip, err)

			if err != nil {
				_ip = err.Error()
			} else {
				_ip = ip
			}

			chttp.OK(w, RESULT{
				Version: inf.Version(),
				Server:  inf.DBName,
				Mainnet: inf.Mainnet(),
				IP:      _ip,
			})
		},
	)
}

type cGasPriceData struct {
	GasPrice    string `json:"gasPrice"`
	TxGasFeeETH string `json:"gasFeeETH"`
}

func (my *cGasPriceData) Calc(speed ebcm.GasSpeed, limit uint64, price ebcm.GasPrice) {
	gas_wei := jmath.VALUE(price.Gas)
	gas_eth := ebcm.WeiToETH(gas_wei)

	fee_wei := jmath.MUL(gas_wei, limit)

	switch speed {
	default:
		//case ebcm.GasSafeLow:
		my.GasPrice = gas_eth
		my.TxGasFeeETH = ebcm.WeiToETH(fee_wei)

	case ebcm.GasAverage:
		my.GasPrice = gas_eth
		my.TxGasFeeETH = ebcm.WeiToETH(jmath.DOTCUT(jmath.MUL(fee_wei, 1.2), 0))

	case ebcm.GasFast:
		my.GasPrice = gas_eth
		my.TxGasFeeETH = ebcm.WeiToETH(jmath.DOTCUT(jmath.MUL(fee_wei, 1.3), 0))

	case ebcm.GasFastest:
		my.GasPrice = gas_eth
		my.TxGasFeeETH = ebcm.WeiToETH(jmath.DOTCUT(jmath.MUL(fee_wei, 1.5), 0))

	}

}

type cInfoGasFee struct {
	GasLimit uint64 `json:"gasLimit"`

	Low     cGasPriceData `json:"low"`
	Avg     cGasPriceData `json:"avg"`
	Fast    cGasPriceData `json:"fast"`
	Fastest cGasPriceData `json:"fastest"`
}

func (my cInfoGasFee) String() string { return dbg.ToJSONString(my) }

func hInfoGasFee() {
	/*
		Comment : 트랜젝션 가스 수수료 계산
		Method : POST
		URL : http://scheduler.server.org:8080/info/gasfee
		Param :
		{
			"gasLimit" : long 	//가스 리미티드 값 (ETH 전송시 21000 고정)
			( 21000 미만 입력시 21000으로 고정 됨. 이더리움 전송값은 21000 고정 값임.)

			( ex. 이더(ETH) 전송시 : 21000 고정값 )
			( ex. 토큰(GDG기준) 전송시 : 최대 60000 정도 소모)
			( ex. NFT토큰 발행시 : 최대 175000 정도 소모)
		}
		Response :
		{
			"success" : true,
			"data" : {
				"gasLimit": 21000,	//요청한 리미티드 값
				"low": {	//느림
					"gasPrice": "0.00000006",	//가스 가격 (ETH)
					"gasFeeETH": "0.00126"		//전송 수수료 가격 (ETH) --> gasFeeETH = gasLimit * gasPrice
				},
				"avg": {	//평균
					"gasPrice": "0.000000067",
					"gasFeeETH": "0.001407"
				},
				"fast": {	//빠름1  (스케줄러 서버에서 쓰는 값)
					"gasPrice": "0.000000077",
					"gasFeeETH": "0.001617"
				},
				"fastest": {	//빠름2
					"gasPrice": "0.000000089",
					"gasFeeETH": "0.001869"
				}
			}
		}
	*/
	type CDATA struct {
		GasLimit uint64 `json:"gasLimit"`
	}
	method := chttp.POST
	url := "/info/gasfee"
	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := CDATA{}
			chttp.BindingJSON(req, &cdata)

			if cdata.GasLimit < 21000 {
				cdata.GasLimit = 21000
			}

			r := cInfoGasFee{
				GasLimit: cdata.GasLimit,
			}

			f := inf.GetFinder()
			gas_price, _ := f.SuggestGasPrice(context.Background(), true)

			r.Low.Calc(ebcm.GasSafeLow, r.GasLimit, gas_price)
			r.Avg.Calc(ebcm.GasAverage, r.GasLimit, gas_price)
			r.Fast.Calc(ebcm.GasFast, r.GasLimit, gas_price)
			r.Fastest.Calc(ebcm.GasFastest, r.GasLimit, gas_price)

			chttp.OK(w, r)
		},
	)
}

func hInfoMaster() {

	type ownerInfo struct {
		Address string            `json:"address"`
		Token   map[string]string `json:"token"`
	}

	type masterInfo struct {
		Mainnet        bool           `json:"mainnet"`
		MasterAddress  string         `json:"master_address"`
		MasterPrice    model.CoinData `json:"master_price"`
		ChargerAddress string         `json:"charger_address"`
		ChargerPrice   model.CoinData `json:"charger_price"`
		MemberCount    int            `json:"member_count"`
		MasterURL      string         `json:"master_url"`
		ChargerURL     string         `json:"charger_url"`

		IsSupportLockupMode bool       `json:"is_support_lockup_mode"`
		Owner               *ownerInfo `json:"owner,omitempty"`
	}

	method := chttp.GET
	url := model.V1 + "/info/master"
	Doc().Comment("[ 스케줄러 정보 ] 스케줄러 마스터/가스비 지갑 정보").
		Method(method).URL(url).
		Etc(".", `_
			<cc_blue>응답결과</cc_blue>
			{
				"success": true,
				"data": {
					"mainnet": false,	// 메인넷 여부					
					"charger_address": "0x61a671b805a2a9ee6d555c244925a228164cc67f",	// 가스비 주소
					"charger_price": {
						"ETH": "7.96442744"	// 가스비주소 ETH 잔액
					},
					"charger_url": "https://goerli.etherscan.io/address/0x61a671b805a2a9ee6d555c244925a228164cc67f",	// 가스비지갑 이더스캔 주소					
					"master_address": "0xabcc18ad4b268f4c4228ae16111b89839c9a709b",	// 마스터지갑 주소
					"master_price": {
						"ERCT": "199818",		// 마스터 지갑의 ERCT토큰 잔액 (ETH를 제외한 지원 토큰은 프로젝트별로 상이함)
						"ETH": "15.494355036",	// 마스터 지갑의 ETH 잔액
						"USDT": "379200"		// 마스터 지갑의 USDT토큰 잔액 (ETH를 제외한 지원 토큰은 프로젝트별로 상이함)
					},
					"master_url": "https://goerli.etherscan.io/address/0xabcc18ad4b268f4c4228ae16111b89839c9a709b", // 마스터지갑 이더스캔 주소
					"member_count": 3		// 가상계좌(유저 입금계좌) 갯수
				}
			}
		`).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			mainnet := inf.Mainnet()
			masterAddress := inf.Master().Address
			chargerAddress := inf.Charger().Address

			masterPrice := map[string]string{}
			for _, token := range inf.TokenList() {
				masterPrice[token.Symbol] = inf.GetFinder().Price(masterAddress, token.Contract, token.Decimal)
			}

			chargerPrice := map[string]string{}
			chargerPrice[model.ETH] = inf.GetFinder().GetCoinPrice(chargerAddress)

			memberCount := 0
			model.DB(func(db mongo.DATABASE) {
				memberCount, _ = db.C(inf.COLMember).Count()
			})

			result := masterInfo{
				Mainnet:        mainnet,
				MasterAddress:  masterAddress,
				MasterPrice:    masterPrice,
				ChargerAddress: chargerAddress,
				ChargerPrice:   chargerPrice,
				MemberCount:    memberCount,
				MasterURL:      inf.EtherScanAddressURL() + masterAddress,
				ChargerURL:     inf.EtherScanAddressURL() + chargerAddress,
			}

			result.IsSupportLockupMode = inf.IsOnwerTaskMode()

			var owner_info *ownerInfo = nil
			if result.IsSupportLockupMode {
				owner_address := inf.Owner().Address
				owner_info = &ownerInfo{
					Address: owner_address,
					Token:   map[string]string{},
				}
				for _, token := range inf.TokenList() {
					owner_info.Token[token.Symbol] = inf.GetFinder().Price(owner_address, token.Contract, token.Decimal)
				} //for
			}
			result.Owner = owner_info

			chttp.OK(w, result)
		},
	)
}

type cMemberInfo struct {
	UID     int64          `json:"uid"`
	Name    string         `json:"name"`
	Address string         `json:"address"`
	Coin    model.CoinData `json:"coin"`
	URL     string         `json:"url"`
}

func docMemberInfo() string {
	return `
	<cc_blue>< 가입 정보 확인 ></cc_blue>
	가입한 회원의 주소 정보와 실제 코인 잔액을 확인 할수 있다.

	<cc_bold>회원계정으로 요청 ></cc_bold> ` + docURL + `/v1/info/member/name/` + docUser + `
	<cc_bold>발급 ID로 요청 ></cc_bold> ` + docURL + `/v1/info/member/uid/` + docID + `
	<cc_bold>발급 주소로 요청 ></cc_bold> ` + docURL + `/v1/info/member/address/` + docAddress + `

	METHOD : GET

	<cc_purple>성공응답</cc_purple> :
	{
		"success": true,
		"data": {
			"uid" : ` + docID + `,
			"name" : "` + docUser + `",	
			"address" : "` + docAddress + `",
			"coin": {
				"ERCT": "0",
				"ETH": "0.596927409",
				"GDG": "0",
				"USDT": "180"
			},
			"url": "https://goerli.etherscan.io/address/` + docAddress + `"
		}
	}
	`
}

func hInfoMemberName() {
	method := chttp.GET
	url := model.V1 + "/info/member/name/:args"
	Doc().Comment("[ 스케줄러 정보 ] 유저의 가상계좌 정보요청 (회원ID로 검색)").
		Method(method).
		URLS(
			url,
			":args", "test@gmail.com(회원ID)",
		).
		Etc(".", `_
			<cc_blue>응답결과</cc_blue>
			{
				"success": true,
				"data": {
					"address": "0x8738183b11fc107503e782ae4befcbee7a0e2ced",	//회원은 지갑주소 (입금계좌)
					"coin": {
						"ETH": "0.946419098",	//지갑주소가 보유하고 있는 잔액 (본금액은 실제 블록체인 잔액과 차이가 있을수 있습니다.)
						"USDT": "0"
					},
					"name": "test1",	// 회원 ID	(고유식별자)
					"uid": 1001,		// 회원 UID (고유식별자)
					"url": "https://goerli.etherscan.io/address/0x8738183b11fc107503e782ae4befcbee7a0e2ced"	// 이더스캔 주소
				}
			}
		`).
		JResultOK(chttp.AckFormat{}).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			name := ps.ByName("args")
			model.Trim(&name)
			if name == "" {
				chttp.Fail(w, ack.NotFoundName)
				return
			}

			model.DB(func(db mongo.DATABASE) {
				member := model.LoadMemberName(db, name)
				if member.Valid() == false {
					chttp.Fail(w, ack.NotFoundName)
					return
				}

				chttp.OK(w, cMemberInfo{
					UID:     member.UID,
					Name:    member.Name,
					Address: member.Address,
					Coin:    member.Coin,
					URL:     member.EtherScanURL(),
				})

			})
		},
	)
}

func hInfoMemberUID() {
	method := chttp.GET
	url := model.V1 + "/info/member/uid/:args"
	Doc().Comment("[ 스케줄러 정보 ] 유저의 가상계좌 정보요청 (회원 UID로 검색)").
		Method(method).
		URLS(
			url,
			":args", "1001(회원ID)",
		).
		JResultOK(chttp.AckFormat{}).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			uidstring := ps.ByName("args")
			model.Trim(&uidstring)
			if uidstring == "" || jmath.IsNum(uidstring) == false {
				chttp.Fail(w, ack.NotFoundName)
				return
			}
			uid := jmath.Int64(uidstring)

			model.DB(func(db mongo.DATABASE) {
				member := model.LoadMember(db, uid)
				if member.Valid() == false {
					chttp.Fail(w, ack.NotFoundName)
					return
				}

				chttp.OK(w, cMemberInfo{
					UID:     member.UID,
					Name:    member.Name,
					Address: member.Address,
					Coin:    member.Coin,
					URL:     member.EtherScanURL(),
				})

			})
		},
	)
}

func hInfoMemberAddress() {
	method := chttp.GET
	url := model.V1 + "/info/member/address/:args"
	Doc().Comment("[ 스케줄러 정보 ] 유저의 가상계좌 정보요청 (회원 주소로 검색)").
		Method(method).
		URLS(
			url,
			":args", "0x.....",
		).
		JResultOK(chttp.AckFormat{}).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			address := ps.ByName("args")

			address = dbg.TrimToLower(address)
			if address == "" || ebcm.IsAddress(address) == false {
				chttp.Fail(w, ack.NotFoundName)
				return
			}

			model.DB(func(db mongo.DATABASE) {
				member := model.LoadMemberAddress(db, address)
				if member.Valid() == false {
					chttp.Fail(w, ack.NotFoundName)
					return
				}

				chttp.OK(w, cMemberInfo{
					UID:     member.UID,
					Name:    member.Name,
					Address: member.Address,
					Coin:    member.Coin,
					URL:     member.EtherScanURL(),
				})

			})
		},
	)
}

package nft_winners

import (
	"net/http"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/jnet/doc"
	"txscheduler/txm/aadev"
	"txscheduler/txm/inf"
)

var dc doc.Object

// Doc :
func Doc() doc.Object {
	if dc == nil {
		dc = doc.NewObjecter("NFT_WINNERS", "WINNERS NFT API LIST", "api-document last update 2023.05.30")
	}
	return dc
}

// DocEnd :
func DocEnd(classic *chttp.Classic) {
	if dc == nil {
		return
	}
	// dc.Update()
	// dc = nil

	if !inf.Mainnet() {
		classic.SetHandler(
			chttp.GET, "/doc/nft",
			func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
				w.Write(dc.Bytes())
			},
		)
	}

}

const (
	doc_host_url = aadev.DEV_URL
)

func _help_estimate_gas() {
	Doc().Message(`
	<cc_purple>( 가스비 예측 예제 )</cc_purple>

		<요청>
		post `+doc_host_url+`/v1/nft/estimate/eth_price
		{
			"from_address" : "0xe129243a027b25d813aced72fe34b22f5fc4bb20", //회원 주소
			"recipients" : [
				{
					"address" : "0x05b93b0feeb9f60a599ba4b4c76262c22e837579",	//수수료 받을 지갑1
					"price" : "0.01"
				},
				{
					"address" : "0x8CE5bb2013887eD586e6a87211aa126453368b7A",	//수수료 받을 지갑2
					"price" : "0.01"
				}
			]
		}

		<성공 응답>
		{
			"success": true,
			"data": {
				"is_transfer_allow": true, //해당 값이 true 여야만 (ETH 민팅 / 유저간 거래) 가능.
				"from_eth_price": "0.473312127500971272",     // 회원주소의 실제 ETH 잔액
				"estimate_gas_price": "0.064999783656225616", // 예상 가스량
				"estimate_pay_price": "0.084999783656225616"  // 예상 총 비용 = 예상 가스량 + 0.01 + 0.01
			}
		}

		<실패 응답 CASE 1>
		post `+doc_host_url+`/v1/nft/estimate/eth_price
		{
			"from_address" : "0xe129243a027b25d813aced72fe34b22f5fc4bb20",
			"recipients" : [
				{
					"address" : "0x05b93b0feeb9f60a599ba4b4c76262c22e837579",
					"price" : "0.51"
				},
				{
					"address" : "0x8CE5bb2013887eD586e6a87211aa126453368b7A",
					"price" : "0.01"
				}
			]
		}

		[예상 총 비용 = 예상 가스량 + 0.51 + 0.01] 
		위와같이 예상 총 비용을 실제 잔액보다 크게 설정할시에 아래와 같은 실패 응답 반환.

		{
			"success": true,
			"data": {
				"is_transfer_allow": false, 
				"fail_message": "회원의 ETH잔고 부족",
				"from_eth_price": "0.473312127500971272",	// 회원주소의 실제 ETH 잔액
				"estimate_gas_price": "0", 				// 가스비 예측 불가
				"estimate_pay_price": "0"
			}
		}

		< 실패 응답 CASE 2>
		post `+doc_host_url+`/v1/nft/estimate/eth_price
		{
			"from_address" : "0xe129243a027b25d813aced72fe34b22f5fc4bb20",
			"recipients" : [
				{
					"address" : "0x05b93b0feeb9f60a599ba4b4c76262c22e837579",
					"price" : "0.46"
				},
				{
					"address" : "0x8CE5bb2013887eD586e6a87211aa126453368b7A",
					"price" : "0.01"
				}
			]
		}

		[예상 총 비용 = 예상 가스량 + 0.46 + 0.01] 
		가스비용 예측은 되나 총소모 비용이 ETH보유량보다 커서 트랜잭션이 실패할 수 있다!!

		{
			"success": true,
			"data": {
			  "is_transfer_allow": false,
			  "from_eth_price": "0.473312127500971272",
			  "estimate_gas_price": "0.110781376010847636",
			  "estimate_pay_price": "0.580781376010847636"  //from_eth_price 값보다 크다.
			}
		  }

		
	`, doc.Blue)
}

package admin

import (
	"net/http"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/jmath"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/jnet/doc"
	"txscheduler/txm/ack"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func init() {
	hMasterInfo()
	hDepositInfoEDIT()

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
	SeedInfo       string         `json:"seed_info"`

	Symbols []string `json:"symbols"`
}

func hMasterInfo() {
	method := chttp.GET
	url := model.V2 + "/info/master"
	Doc().Comment("[ 정보 ] 마스터 지갑 정보").
		Method(method).URL(url).
		JResultOK(chttp.AckFormat{}).
		ETC(doc.EV(masterInfo{})).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			mainnet := inf.Mainnet()
			masterAddress := inf.Master().Address
			chargerAddress := inf.Charger().Address

			masterPrice := map[string]string{}
			for _, token := range inf.TokenList() {
				finder := ecsx.New(mainnet, inf.InfuraKey())
				wei := finder.Balance2(masterAddress, token.Contract)
				price := ecsx.WeiToToken(wei, token.Decimal)
				masterPrice[token.Symbol] = price
			}

			finder := ecsx.New(mainnet, inf.InfuraKey())
			wei := finder.Balance(chargerAddress)
			chargerPrice := map[string]string{}
			chargerPrice[model.ETH] = ecsx.WeiToToken(wei, "18")

			memberCount := 0
			model.DB(func(db mongo.DATABASE) {
				memberCount, _ = db.C(inf.COLMember).Count()
			})

			chttp.OK(w, masterInfo{
				Mainnet:        mainnet,
				MasterAddress:  masterAddress,
				MasterPrice:    masterPrice,
				ChargerAddress: chargerAddress,
				ChargerPrice:   chargerPrice,
				MemberCount:    memberCount,
				MasterURL:      inf.EtherScanAddressURL() + masterAddress,
				ChargerURL:     inf.EtherScanAddressURL() + chargerAddress,
				SeedInfo:       inf.SeedView(),
				Symbols:        inf.Config().Tokens.SymbolList(),
			})
		},
	)
}

type cDepositInfo struct {
	Coin      model.CoinData `json:"coin"`
	BaseValue string         `json:"base_value"`
}

func (my *cDepositInfo) Valid() bool {
	if my.Coin == nil {
		my.Coin = model.NewCoinData()
	}

	for _, v := range my.Coin {
		if jmath.IsLessZero(v) {
			return false
		}
	}

	if jmath.IsNum(my.BaseValue) == false {
		return false
	}
	if jmath.IsLessZero(my.BaseValue) {
		return false
	}
	return true
}

func hDepositInfoEDIT() {
	method := chttp.POST
	url := model.V2 + "/info/deposit.edit"
	Doc().Comment("[ 정보 ] 디파짓 인포 수량 수정 (ROOT 권한)").
		Method(method).URL(url).
		JParam(cDepositInfo{}).
		JResultOK(chttp.AckFormat{}).
		ETC(doc.EV(masterInfo{})).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			admin := model.GetTokenAdmin(req)
			if admin.IsRoot == false {
				chttp.Fail(w, ack.InvalidRootAdmin)
				return
			}

			cdata := cDepositInfo{}
			chttp.BindingJSON(req, &cdata)
			if cdata.Valid() == false {
				chttp.Fail(w, ack.BadParam)
				return
			}

			model.DB(func(db mongo.DATABASE) {
				info := model.InfoDeposit{}.Get(db)
				info.Coin = cdata.Coin
				info.BaseValue = cdata.BaseValue

				info.UpdateDB(db)

				chttp.OK(w, nil)
			})
		},
	)
}

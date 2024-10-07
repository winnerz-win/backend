package nft_winners

import (
	"net/http"
	"strings"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/jnet/cnet"
	"txscheduler/brix/tools/unix"
	"txscheduler/nft_winners/nwdb"
	"txscheduler/nft_winners/nwtypes"
	"txscheduler/nft_winners/rpc"
	"txscheduler/txm/ack"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func init() {
	hNftSetBaseURI()
}

type AckSetBaseURI struct {
	ReceiptCode nwtypes.RECEIPT_CODE `json:"receipt_code"`
	IsCallback  bool                 `json:"is_callback"`
}

func (AckSetBaseURI) TagString() []string {
	return []string{
		"receipt_code", "새로운 주소",
		"is_callback", "결과 콜백을 받을지 여부",
	}
}

type ReqSetBaseURI struct {
	NewURI     string `json:"new_uri"`
	IsCallback bool   `json:"is_callback"`
}

func (ReqSetBaseURI) TagString() []string {
	return []string{
		"new_uri", "새로운 주소",
		"is_callback", "결과 콜백을 받을지 여부",
	}
}

func hNftSetBaseURI() {
	method := chttp.POST
	url := model.V1 + "/nft/set_base_uri"
	Doc().Comment("[WINNERZ] baseURI 변경 요청").
		Method(method).URL(url).
		JParam(ReqSetBaseURI{}, ReqSetBaseURI{}.TagString()...).
		JAckOK(AckSetBaseURI{}, AckSetBaseURI{}.TagString()...).
		JAckError(ack.NFT_RPC_TIMEOUT).
		JAckError(ack.NFT_SameBaseURI).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := chttp.ParseRequestJson[ReqSetBaseURI](req)
			cdata.NewURI = strings.TrimSpace(cdata.NewURI)

			finder := GetSender()
			if finder == nil {
				chttp.Fail(w, ack.NFT_RPC_TIMEOUT)
				return
			}

			base_uri, err := rpc.ERC721.BaseURI(
				finder, nft_config.Reader(),
			)
			if err != nil {
				chttp.Fail(w, ack.NFT_RPC_TIMEOUT)
				return
			}
			if cdata.NewURI == base_uri {
				chttp.Fail(w, ack.NFT_SameBaseURI)
				return
			}

			model.DB(func(db mongo.DATABASE) {

				receipt_code := nwtypes.GetReceiptCode(nwtypes.RC_SET_BASE_URI, cdata.NewURI)

				master_try := nwtypes.NftMasterTry{
					ReceiptCode: receipt_code,

					DATA_SEQ: nwtypes.MAKE_DATA_SEQ(
						nwtypes.NFT_BASE_URI,
						nwtypes.DataSetBaseURI{
							ReceiptCode: receipt_code,
							NewURI:      cdata.NewURI,
							IsCallback:  cdata.IsCallback,
						},
					),

					TimeTryAt: unix.Now(),
				}

				master_try.InsertDB(db)

				chttp.OK(w, AckSetBaseURI{
					ReceiptCode: receipt_code,
					IsCallback:  cdata.IsCallback,
				})
			})

		},
	)
}

////////////////////////////////////////////////////////////////////////////////

var (
	baseUriC = make(chan any, 1)
)

const (
	URL_NFTS_SET_BASE_URI_CALLBACK = "/v1/nfts/set_base_uri/callback"
)

func waitC_setBaseURI_Callback(db mongo.DATABASE) {

	item := nwtypes.NftSetBaseURIResult{}
	if err := db.C(nwdb.NftSetBaseURIResult).
		Find(mongo.Bson{"is_send": false}).
		Sort("timestamp").
		One(&item); err != nil {
		return
	}

	ack := cnet.POST_JSON_F(
		inf.ClientAddress()+URL_NFTS_SET_BASE_URI_CALLBACK,
		nil,
		item,
	)
	if err := ack.Error(); err != nil {
		dbg.RedItalic("result_callback :", err)
	} else {
		item.SendOK(db, unix.Now())
	}

}

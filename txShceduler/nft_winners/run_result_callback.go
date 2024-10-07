package nft_winners

import (
	"time"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/cnet"
	"txscheduler/brix/tools/runtext"
	"txscheduler/brix/tools/unix"
	"txscheduler/nft_winners/nwdb"
	"txscheduler/nft_winners/nwtypes"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

const (
	URL_NFTS_MINT_CALLBACK = "/v1/nfts/mint/callback"
	URL_NFTS_SALE_CALLBACK = "/v1/nfts/user_sale/callback"
)

func run_result_callback(rtx runtext.Runner) {
	defer dbg.PrintForce("nft_winners.run_result_callback ----------  END")
	<-rtx.WaitStart()
	dbg.PrintForce("nft_winners.run_result_callback ----------  START")

EXIT:
	for {
		select {
		case <-rtx.EndC():
			break EXIT
		default:
		} //select
		time.Sleep(time.Second)

		model.DB(func(db mongo.DATABASE) {
			mongo.IterForeach(
				db.C(nwdb.NftActionResult).
					Find(mongo.Bson{"is_send": false}).
					Sort("insert_at").
					Iter(),
				func(cnt int, item nwtypes.NftActionResult) bool {

					switch item.ResultType {
					case nwtypes.RESULT_MINT:
						ack := cnet.POST_JSON_F(
							inf.ClientAddress()+URL_NFTS_MINT_CALLBACK,
							nil,
							item,
						)
						if err := ack.Error(); err != nil {
							dbg.RedItalic("result_callback :", err)
						} else {
							item.SendOK(db, unix.Now())
						}

					case nwtypes.RESULT_USER_SALE:
						ack := cnet.POST_JSON_F(
							inf.ClientAddress()+URL_NFTS_SALE_CALLBACK,
							nil,
							item,
						)
						if err := ack.Error(); err != nil {
							dbg.RedItalic("result_callback :", err)
						} else {
							item.SendOK(db, unix.Now())
						}
					} //switch

					return false
				},
			)

			waitC_setBaseURI_Callback(db)
		})
	} //for
}

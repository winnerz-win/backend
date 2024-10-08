package model

import (
	"jtools/cc"
	"jtools/unix"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/database/mongo/tools/dbg"
	"txscheduler/brix/tools/database/mongo/tools/jmath"
	"txscheduler/nft_winners/rpc"
	"txscheduler/txm/inf"
)

type AASystemInfo struct {
	AT   unix.Time `bson:"at" json:"at"`
	AKST string    `bson:"akst" json:"akst"`
	AYMD int       `bson:"aymd" json:"aymd"`
	Data mongo.MAP `bson:"data" json:"data"`
}

func (my AASystemInfo) String() string { return dbg.ToJsonString(my) }

var (
	winners_nft = ""
)

func SetNftWinners(nft string) {
	winners_nft = nft
}

func (my AASystemInfo) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.AASystemInfo)
		c.RemoveAll(nil)

		now := unix.Now()
		my.AT = now
		my.AKST = now.KST()
		my.AYMD = now.YMD()

		config := inf.Config()

		finder := inf.GetFinder()
		network := mongo.MAP{
			"host_url": finder.Host(),
			"chain_id": jmath.VALUE(finder.ChainID()),
		}

		my.Data = mongo.MAP{
			"mainnet": config.Mainnet,
			"version": config.Version,
			"network": network,
			//"seed":                      config.Seed,
			"db":                  config.DB,
			"ip_check":            config.IPCheck,
			"client_callback_url": inf.ClientAddress(),
			//"admin_salt":                config.AdminSalt,
			"confirms":                  config.Confirms,
			"is_lock_transfer_by_owner": config.IsLockTransferByOwner,
		}

		if config.IsLockTransferByOwner {
			my.Data["owner"] = inf.Owner().Address
		}
		my.Data["master"] = inf.Master().Address
		my.Data["charger"] = inf.Charger().Address

		my.Data["tokens"] = inf.TokenList()

		my.Data["first_erc20"] = inf.FirstERC20()

		if winners_nft != "" {
			erc721_info := rpc.ERC721Info{}
			rpc.ERC721.ERC721Info(
				inf.GetFinder(), rpc.Reader(winners_nft),
				func(info rpc.ERC721Info) {
					erc721_info = info
				},
			)
			my.Data["winners_nft"] = erc721_info
		}

		c.Insert(my)

		cc.Yellow(my)
	})
}

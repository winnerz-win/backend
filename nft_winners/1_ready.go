package nft_winners

import (
	"context"
	"reflect"
	"time"
	"txscheduler/brix/tools/cloud/ebcm"
	"txscheduler/brix/tools/cloud/ebcm/abi"
	"txscheduler/brix/tools/cloud/jeth/ecs"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/runtext"
	"txscheduler/nft_winners/nwdb"
	"txscheduler/nft_winners/nwtypes"
	"txscheduler/nft_winners/rpc"
	"txscheduler/txm/cloud"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

const (
	limit_real_pow = 1.3

	estimateWaitDu = time.Second * 10

	is_skip_tip_cap = false
)

var (
	LIMIT_TAG_MINT = "mint"
	MINT_LIMIT_MAX = []uint64{210000, 260000} //199711

	LIMIT_TAG_M_TRANS = "mtrans"
	M_TRANS_MAX       = []uint64{130000, 150000}
)

func calc_limit(limit uint64, dot_pow any, limit_tag string) uint64 {
	re_limit := uint64(jmath.Int64(jmath.MUL(limit, dot_pow)))
	if limit_tag != "" {
		switch limit_tag {
		case LIMIT_TAG_MINT:
			if re_limit <= MINT_LIMIT_MAX[0] {
				dbg.RedItalicBG("[MASTER_MINT_GAS_LIMIT_UNDER] :", re_limit, ", limit=", MINT_LIMIT_MAX[1])
				re_limit = MINT_LIMIT_MAX[1]
			}

		case LIMIT_TAG_M_TRANS:
			if re_limit <= M_TRANS_MAX[0] {
				dbg.RedItalicBG("[M_TRANSFER_GAS_LIMIT_UNDER] :", re_limit, ", limit=", M_TRANS_MAX[1])
				re_limit = M_TRANS_MAX[1]
			}

		} //switch

	}

	return re_limit
}

var handle = chttp.PContexts{}

func Ready(classic *chttp.Classic) runtext.Starter {
	dbg.Cyan("nft_winners.Ready ---- START")
	defer dbg.Cyan("nft_winners.Ready ---- READY")

	rtx := runtext.New("nft_winners")

	classic.SetContextHandles(handle)

	check_erc20_approve_all()

	dbg.Cyan("NFT_WINNERS : ", nft_config)

	start_indexing_db()

	cloud.InjectMasterWithdrawProcess(
		proc_master_pending_check,
		proc_master_send_try,
	)
	go run_user_tx_pending(rtx)
	go run_user_tx_try(rtx)
	go run_result_callback(rtx)

	DocEnd(classic)

	return rtx
}

func start_indexing_db() {
	indexing_list := []interface {
		IndexingDB(db mongo.DATABASE)
	}{
		// nwtypes.NftMintTry{},
		// nwtypes.NftMintPending{},

		nwtypes.NftTokenIDS{},
		nwtypes.NftActionResult{},

		nwtypes.NftMasterTry{},
		nwtypes.NftMasterPending{},

		nwtypes.NftUserTry{},
		nwtypes.NftUserPending{},

		nwtypes.NftSetBaseURIResult{},
	}
	model.DB(func(db mongo.DATABASE) {
		for i, v := range indexing_list {
			v.IndexingDB(db)
			dbg.CyanItalic("indexingDB[", i, "/", len(indexing_list), "] ---", reflect.TypeOf(v).Name())
		}

		db.C(nwdb.NftAInfo).RemoveAll(nil)
		db.C(nwdb.NftAInfo).Insert(nft_config)

	})

}

func GetSender() *ebcm.Sender {
	return ecs.New(
		ecs.RPC_URL(inf.Mainnet()),
		inf.InfuraKey(),
	)
}

func check_erc20_approve_all() {
	dbg.Cyan("nft_winners.check_erc20_approve_all ---- START")
	defer dbg.Cyan("nft_winners.check_erc20_approve_all ---- END")

	sender := GetSender()
	if sender == nil {
		dbg.Exit("nft_winners.check_erc20_approve_all : sender is nil")
	}

	if err := rpc.ERC721.Symbol(
		sender, nft_config.Reader(),
		func(_symbol string) {
			nft_config.NftSymbol = _symbol
			dbg.Cyan("NFT_SYMBOL :", _symbol)
		},
	); err != nil {
		dbg.Exit("nft_winners.check_erc20_approve_all : nft_symbol :", err)
	}
	rpc.ERC721.Name(
		sender, nft_config.Reader(),
		func(_name string) {
			nft_config.NftName = _name
			dbg.Cyan("NFT_NAME :", _name)
		},
	)

	nft_config.BaseURI, _ = rpc.ERC721.BaseURI(sender, nft_config.Reader())

	erc20 := inf.TokenList().FirstERC20()
	dbg.Yellow("[check_erc20_approve_all] erc20 :", erc20.Symbol)
	token_reader := rpc.Reader(erc20.Contract)

	is_approve_try := false
	if err := rpc.ERC20.Allowance(
		sender,
		token_reader,
		inf.Master().Address,
		nft_config.NftContract,
		func(amount string) {
			if jmath.CMP(amount, 0) <= 0 {
				is_approve_try = true
			} else {
				dbg.Cyan("ERC20[", erc20.Symbol, "] Allowance :", amount)
			}
		},
	); err != nil {
		dbg.Exit("nft_winners.check_erc20_approve_all : erc20_allowance :", err)
	}

	if is_approve_try {
		dbg.Cyan("NFT_CONTRACT APPROVE TRY")
		exit := func(v ...any) {
			dbg.Exit(v...)
		}
		_ = exit

		tag := "[MASTER.APPROVE]"
		for {
			time.Sleep(time.Second)

			from := inf.Master()
			ctx := context.Background()
			nonce, err := sender.NonceAt(ctx, from.Address)
			if err != nil {
				dbg.RedItalic(tag, "nonce :", err)
				continue
			}
			pending_nonce, _ := sender.PendingNonceAt(ctx, from.Address)
			if nonce != pending_nonce {
				dbg.RedItalic(tag, "master nonce differ (", nonce, "/", pending_nonce, ")")
				continue
			}

			pad_bytes := ebcm.MakePadBytesABI(
				"approve",
				abi.TypeList{
					abi.NewAddress(nft_config.NftContract), //spender
					abi.NewUint256(ebcm.UINT256MAX),        //amount
				},
			)
			gas_limit, err := sender.EstimateGas(
				ctx,
				ebcm.MakeCallMsg(
					from.Address,
					token_reader.Contract(),
					"0",
					pad_bytes,
				),
			)
			if err != nil {
				dbg.RedItalic(tag, "gas_limit :", err)
				continue
			}
			gas_limit = calc_limit(gas_limit, limit_real_pow, "")

			gas_price, err := sender.SuggestGasPrice(ctx, true)
			if err != nil {
				dbg.RedItalic(tag, "gas_price :", err)
				continue
			}

			ntx := sender.NewTransaction(
				nonce,
				token_reader.Contract(),
				"0",
				gas_limit,
				gas_price,
				pad_bytes,
			)
			stx, err := sender.SignTx(
				ntx,
				from.PrivateKey,
			)
			if err != nil {
				dbg.RedItalic(tag, "stx :", err)
				continue
			}

			hash, err := sender.SendTransaction(ctx, stx)
			if err != nil {
				dbg.RedItalic(tag, "send :", err)
				continue
			}
			dbg.Cyan(tag, "send_hash :", hash)
			r := sender.CheckSendTxHashReceiptByHash(
				hash,
				ebcm.SEC_1_HOUR, //1h
				true,
			)
			if r.IsSuccess {
				dbg.Cyan(tag, r)
				break

			} else {
				dbg.RedItalic(tag, r)
			}
		} //for

	}

}

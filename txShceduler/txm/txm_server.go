package txm

import (
	"time"
	"txscheduler/brix/tools"
	"txscheduler/brix/tools/console"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jargs"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/runtext"
	"txscheduler/txm/admin"
	"txscheduler/txm/api"
	"txscheduler/txm/cloud"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
	"txscheduler/txm/nftc"
	"txscheduler/txm/rcc"
	"txscheduler/txm/scv"
)

// Start :
func Start(
	config *inf.IConfig,
	scv_callback_list scv.CallbackList,
	nftCallback *nftc.NftCallback,
) {

	dbg.PrintForce("txm.Start")

	port := 8989
	args := jargs.Args("=")
	inf.SetArgs(args)

	args.Action("port", func(v jargs.ArgValue) {
		port = v.ForceInt()
	})
	_ = port

	mainnet := args.Do("mainnet")

	if config != nil {
		config.Mainnet = mainnet
		inf.SetConfig(config)
	} else {
		inf.InitConfig("config.yaml", mainnet)
	}

	inf.InitMongo("mongo.yaml")

	corsHeader := []string{
		"Access-Control-Allow-Origin",
		model.HeaderAdminToken,
	}

	chttp.SetAckFormat()
	//classic := chttp.NewClassic(port, corsHeader, inf.DBName)

	logger := chttp.NewDayLogger(
		"./log",
		dbg.Cat("log_", config.DB),
		config.DB,
	)
	dbg.SetLogWriteln(logger.Writeln)

	classic := chttp.NewClassicLogger(
		port,
		corsHeader,
		inf.DBName,
		logger,
	)

	model.Ready(
		scv_callback_list,
	)
	rcc.Ready(classic)

	starterList := runtext.StarterList{}

	isNFT := false
	if nftCallback != nil {
		nftc.Ready(nftCallback)
		starterList.Append(
			nftc.Start(classic),
		)
		isNFT = true
	}

	scv_callback_rtx := runtext.New("scv_callback_rtx")
	starterList.Append(
		scv_callback_rtx,
		cloud.Ready(),
		api.Ready(classic),
		admin.Ready(classic),
	)

	start_package_list := scv_callback_list.StartPackageList()
	for _, package_ready_func := range start_package_list {
		starterList.Append(
			package_ready_func(classic),
		)
	} //for

	tools.SetVersion(func() string {
		message := inf.CoreVersion + dbg.ENTER +
			"[" + inf.DBName + "]" + dbg.ENTER +
			dbg.Cat("Mainnet :", config.Mainnet) + dbg.ENTER +
			config.Version + dbg.ENTER +
			inf.SeedView() + dbg.ENTER

		if isNFT {
			message += "[NFT ON]" + dbg.ENTER +
				"NFT Contract :" + nftCallback.NftTokenContract + dbg.ENTER +
				"NFT Owner    :" + nftCallback.NftOwnerAddress + dbg.ENTER +
				"NFT Deposit  :" + nftCallback.NftDepositAddress + dbg.ENTER
		}

		return message
	})

	if args.Do("start") {
		go func() {
			dbg.PrintForce("auto start")
			time.Sleep(time.Millisecond * 10)
			console.TestCommand("start")
		}()
	}

	go scv_callback_list.StartLooper(scv_callback_rtx)

	////////////////////////////////////////////////////////////

	go func() {
		for {
			if !classic.IsServerRun() {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			dbg.Green("-------------------------------------------")
			for i := 0; i < 3; i++ {
				dbg.PrintForce("Wallet-Seed :", config.Seed, config.Mainnet)
			}
			dbg.Green("-------------------------------------------")
			break
		} //for
	}()

	classic.CmdStart(starterList)
}

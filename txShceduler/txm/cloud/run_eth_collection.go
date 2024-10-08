package cloud

import (
	"context"
	"fmt"
	"jtools/cc"
	"jtools/cloud/ebcm"
	"jtools/jmath"
	"jtools/mms"
	"sync"
	"time"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/runtext"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
	"txscheduler/txm/nftc"
)

var (
	cms_once = sync.Once{}
)

func runETHCollection(rtx runtext.Runner) {
	defer dbg.PrintForce("cloud.runETHCollection ----------  END")
	<-rtx.WaitStart()
	dbg.PrintForce("cloud.runETHCollection ----------  START")
	TAG := "[run_eth_collection] "

	sleepNormal := time.Second * 3
	_ = sleepNormal
	txCnt := model.NewTxETHCounter(inf.Mainnet(), "0")
	txCnt.LoadFromDB()

	if nftc.IsRun() {
		ebcm.MMA_MethodID_Append("cloud.CMS", &cms_once, nftc.CMS())
	}

	isview := false
	isStart := false
EXIT:
	for {
		select {
		case <-rtx.EndC():
			break EXIT
		default:
		} //select

		finder := get_finder()
		if finder == nil {
			model.LogError.InsertLog(
				model.ErrorFinderNull,
				"runETHCollection.finder",
			)
			time.Sleep(sleepNormal)
			continue
		}

		// if nftc.IsRun() {
		// 	finder := get_sender_x()
		// 	finder.SetCustomMethods(nftc.CMS()...)
		// }

		ctx := context.Background()
		n, err := finder.BlockNumber(ctx)
		if err != nil {
			cc.Red(TAG, "finder.BlockNumber :", err)
			time.Sleep(sleepNormal)
			continue
		}
		lastNumber := jmath.VALUE(n)

		number := txCnt.Number
		if number == "0" {
			cc.Gray("chain_last_number : ", lastNumber, ", db_number:", txCnt.Number)
			txCnt.Number = lastNumber
			number = lastNumber
		}
		if isview {
			logCollection("find-number :", number)
		}

		if !isStart {
			isStart = true
			if inf.Mainnet() {
				go runETHInternal(
					rtx,
					"foblgate",
					"0x6c0b51971650d28821ce30b15b02b9826a20b129",
					number,
				)
			}
		}

		numberGap := jmath.SUB(lastNumber, number)
		if jmath.CMP(numberGap, inf.Confirms()) < 0 {
			time.Sleep(sleepNormal)
			continue
		}

		data := finder.BlockByNumber(number)
		//data := finder.BlockByNumber(number, false)
		if data == nil {
			cc.Red(TAG, "finder.BlockByNumber dats is null. (", number, ")")
			time.Sleep(sleepNormal)
			continue
		}

		if nftc.IsRun() { //NFT-JOB
			nftc.GetTxList(number, data.TxList, true)
		}

		list := ebcm.TransactionBlockList{}
		txs := data.TxList
		for _, tx := range txs {
			// if tx.IsError {	//에러난것도 일딴 수집한다.
			// 	continue
			// }

			if tx.IsContract {
				token := inf.TokenList().GetContract(tx.ContractAddress)
				if token.Valid() {
					tx.Symbol = token.Symbol
					tx.Decimals = token.Decimal
					list = append(list, tx)
				}
			} else {
				tx.Symbol = model.ETH
				tx.ContractAddress = "eth"
				tx.Decimals = "18"
				list = append(list, tx)
			}
		} //for

		if len(list) > 0 {
			nowAt := mms.Now()
			if err := processTxlist(list, nowAt, false); err != nil {
				model.LogError.InsertLog(
					model.ErrorProcessTxlist,
					dbg.Cat("number[", number, "] :", err.Error()),
				)
			}
		}

		txCnt.Inc(lastNumber)
		if jmath.CMP(numberGap, 10) >= 0 {
			isview = true
			logCollection("number-gap :", numberGap, ", txs:", len(txs))
			fmt.Println()
			time.Sleep(time.Millisecond * 10)
		} else {
			isview = false
			time.Sleep(sleepNormal)
		}

	} //for

}

func processTxlist(list ebcm.TransactionBlockList, nowAt mms.MMS, isInternal bool) error {

	depositList := ETHDepositList{}
	_ = depositList

	err := model.DB(func(db mongo.DATABASE) {

		finder := get_finder()

		for _, tx := range list {

			//마스터지갑으로 들어온 코인
			if tx.To == inf.Master().Address {

				finder.InjectReceipt(&tx, finder.ReceiptByHash(tx.Hash))
				if !tx.IsError {
					innerToMaster(db, tx, nowAt)
				}
				continue
			}

			member := model.LoadMemberAddress(db, tx.To)
			if !member.Valid() {
				continue
			}
			if !isInternal {
				finder.InjectReceipt(&tx, finder.ReceiptByHash(tx.Hash))
			} else {
				finder.InjectReceipt(&tx, finder.ReceiptByHash(tx.Hash))
			}

			// if tx.IsError {
			// 	continue
			// }
			if tx.From == inf.Charger().Address {
				if !tx.IsError {
					model.LockMember(db, tx.To, func(member model.Member) {
						price := ebcm.WeiToToken(tx.Amount, tx.Decimals)
						member.Coin.ADD(tx.Symbol, price)
						member.UpdateDB(db)

						model.CoinSumAdd(db, tx.Symbol, price)

						member.UpdateCoinDB_Legacy(db)
						logCharger("Recv [", member.UID, "]", member.Address, price)
					})
				}
				continue
			}

			newTx := model.TxETHBlock{
				TransactionBlock: tx,
				UID:              member.UID,
				Order:            jmath.Int64(tx.BlockNumber),
				TxState:          model.TxErrorState(tx.IsError),
			}
			if newTx.IsInert(db) != nil {
				continue
			}

			model.LockMember(db, tx.To, func(member model.Member) {
				price := ebcm.WeiToToken(tx.Amount, tx.Decimals)
				if !tx.IsError {
					member.Coin.ADD(tx.Symbol, price)
					member.Deposit.ADD(tx.Symbol, price)
					member.UpdateDB(db)

					model.CoinSumAdd(db, tx.Symbol, price)
					model.CoinDay{}.AddDeposit(db, tx.Symbol, price, nowAt)
				}

				member.UpdateCoinDB_Legacy(db)

				log := model.LogDeposit{
					User:     member.User,
					Hash:     tx.Hash,
					Symbol:   tx.Symbol,
					Contract: tx.ContractAddress,
					Decimal:  tx.Decimals,
					Price:    price,
					From:     tx.From,

					DepositResult: !tx.IsError,

					Timestamp: nowAt,
					YMD:       nowAt.YMD(),
					IsSend:    false,
				}
				log.InsertDB(db)

				if !tx.IsError {
					depositList = append(depositList, ETHDeposit{
						UID:      member.UID,
						Address:  member.Address,
						Symbol:   tx.Symbol,
						Contract: tx.ContractAddress,
						Decimal:  tx.Decimals,
						IsForce:  false,
					})
				}
			})

		} //for
	})

	if err != nil {
		return err
	}
	if len(depositList) > 0 {
		ETHDepositChan <- depositList
	}

	return nil
}

// innerToMaster : 외부에서 마스터지갑으로 들어온것들 처리.
func innerToMaster(db mongo.DATABASE, tx ebcm.TransactionBlock, nowAt mms.MMS) {
	member := model.LoadMemberAddress(db, tx.From)
	if member.Valid() {
		return
	}

	price := ebcm.WeiToToken(tx.Amount, tx.Decimals)
	item := model.LogExMaster{
		Hash:     tx.Hash,
		From:     tx.From,
		Symbol:   tx.Symbol,
		Contract: tx.ContractAddress,
		Decimal:  tx.Decimals,
		Price:    price,

		Timestamp: nowAt,
		YMD:       nowAt.YMD(),
	}
	item.InsertDB(db)
}

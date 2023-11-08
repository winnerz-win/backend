package ecsx

import (
	"sort"

	"txscheduler/brix/tools/jmath"

	"txscheduler/brix/tools/cloudx/ethwallet/EtherScanAPI"
)

// GetTransactions : target(address,contract)  --- 누락 될수 있다.. 안쓰도록 합시다.
func GetTransactions(cfg EtherScanAPI.Config, target string, startblock string, endblock ...string) TransactionBlockList {

	sumlist := TransactionBlockList{}

	ethr := EtherScanAPI.GetEtherTransactions(cfg, target, startblock, endblock...)
	if ethr != nil && ethr.IsSuccess() {
		for _, tx := range ethr.Result {
			newTx := NewTxBlockEx(tx)
			sumlist = append(sumlist, newTx)
		}
	}

	tkr := EtherScanAPI.GetTokenTransactions(cfg, target, startblock, endblock...)
	if tkr != nil && tkr.IsSuccess() {
		for _, tx := range tkr.Result {
			newTx := NewTxBlockEx(tx)
			sumlist = append(sumlist, newTx)
		}
	}

	if cfg.SortOrder() == EtherScanAPI.SortASC {
		sort.Slice(sumlist, func(i, j int) bool {
			return jmath.CompareTo(sumlist[i].BlockNumber, sumlist[j].BlockNumber) < 0
		})
	} else {
		sort.Slice(sumlist, func(i, j int) bool {
			return jmath.CompareTo(sumlist[i].BlockNumber, sumlist[j].BlockNumber) > 0
		})
	}

	return sumlist
}

// GetInternalTransaction :
func GetInternalTransaction(cfg EtherScanAPI.Config, target string, startblock string, endblock ...string) EtherScanAPI.InternalTransactionList {
	data := EtherScanAPI.GetInternalTransaction(cfg, target, startblock, endblock...)
	if data == nil || data.Status != "1" {
		return EtherScanAPI.InternalTransactionList{}
	}
	return data.Result.ToList()
}

package ecsx

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"

	"github.com/ethereum/go-ethereum/ethclient"
)

func NewList(host_list []string, index int, infuraKey string) *Sender {
	sender := newList(host_list, index, infuraKey)
	if sender == nil {
		return nil
	}
	sender.etherTags = append(sender.etherTags,
		etherName1,
		etherName2,
		etherName3,
	)

	//sender.linkGasPriceFunc()
	// sender.gasFunc = SuggestGasPrice
	// sender.SetLinkGasPriceFunc(SuggestGasPrice)
	return sender
}

func newList(host_list []string, index int, infuraKey string) *Sender {
	hosturl := host_list[index]
	url := ""
	infuraKey = strings.TrimSpace(infuraKey)
	if infuraKey != "" {

		url = hosturl + "/" + infuraKey
	} else {
		url = hosturl
	}

	client, err := ethclient.Dial(url)
	if err != nil {
		fmt.Println("Sender.New::Dial :", err)
		return nil
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		fmt.Println("○○○○○○○○○○○○○○○○○○○○○○○○○○○○○○○")
		dbg.Red("newSender.ChainID :", err)
		dbg.Red("host :", hosturl)
		fmt.Println("○○○○○○○○○○○○○○○○○○○○○○○○○○○○○○○")
	}

	return &Sender{
		client:     client,
		mainnet:    true,
		infuraKey:  infuraKey,
		hostURL:    hosturl,
		host_index: 0,
		host_list:  host_list,
		gasURL:     hosturl,
		minGasWei:  "0",

		cacheChainID: jmath.VALUE(chainID.Int64()),
		chainID:      big.NewInt(chainID.Int64()),
		txnType:      TXN_EIP_1559,
	}
}

func (my *Sender) ToggleHost() {
	if len(my.host_list) == 0 {
		return
	}
	infuraKey := my.infuraKey
	index := my.host_index + 1
	if index >= len(my.host_list) {
		index = 0
	}

	hosturl := my.host_list[index]
	url := ""
	infuraKey = strings.TrimSpace(infuraKey)
	if infuraKey != "" {

		url = hosturl + "/" + infuraKey
	} else {
		url = hosturl
	}

	client, err := ethclient.Dial(url)
	if err != nil {
		fmt.Println("Sender.New::Dial :", err)
		return
	}

	fmt.Println("○○○○○○○○○○○○○○○○○○○○○○○○○○○○○○○")
	dbg.Red("change_host :", hosturl)
	fmt.Println("○○○○○○○○○○○○○○○○○○○○○○○○○○○○○○○")

	my.client = client
	my.host_index = index
	my.hostURL = hosturl
	my.gasURL = hosturl

}

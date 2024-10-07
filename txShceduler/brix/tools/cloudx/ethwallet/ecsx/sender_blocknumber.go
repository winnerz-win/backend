package ecsx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
	"txscheduler/brix/tools/jnet/cnet"
)

func (my Sender) BlockNumberRpc() (string, bool) {
	v, err := my.client.BlockNumber(context.Background())
	if err != nil {
		dbg.Red(err)
		return "0", false
	}

	return jmath.VALUE(v), true
}

// BlockNumber : http-post
func (my Sender) BlockNumber() (string, bool) {

	client := cnet.New(my.HostURL())

	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("Accept", "application/json")
	rpcParam := JSONRPCPARAM()
	rpcParam.Method = "eth_blockNumber"

	rpcAck := JSONPRCACK()

	b, _ := json.Marshal(rpcParam)

	resp, err := http.Post(my.HostURL()+"/"+my.InfuraKey(), "application/json", bytes.NewReader(b))
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		fmt.Println("Error :", err)
		return "0", false
	}
	resb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error :", err)
		return "0", false
	}
	if err := json.Unmarshal(resb, &rpcAck); err != nil {
		fmt.Println("Error :", err)
		return "0", false
	}

	return rpcAck.IntString(), true
}

// BlockNumberTry :
func (my Sender) BlockNumberTry(defaultNumber string) string {
	lastNumber, do := my.BlockNumber()
	if do == false {
		return defaultNumber
	}
	return lastNumber
}

func (my Sender) traceTransaction(hash string) {
	client := cnet.New(my.HostURL())

	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("Accept", "application/json")
	rpcParam := JSONRPCPARAM()

	rpcParam.Name = "traceTransaction"
	rpcParam.Method = "debug_traceTransaction"
	rpcParam.Params = append(rpcParam.Params,
		hash,
		struct{}{},
	)

	// rpcParam.Method = "eth_getTransactionByHash"
	// rpcParam.Params = append(rpcParam.Params,
	// 	hash,
	// )

	rpcAck := JSONPRCACK()

	b, _ := json.Marshal(rpcParam)

	resp, err := http.Post(my.HostURL()+"/"+my.InfuraKey(), "application/json", bytes.NewReader(b))
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		fmt.Println("Error :", err)

	}
	resb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error :", err)
	}
	if err := json.Unmarshal(resb, &rpcAck); err != nil {
		fmt.Println("Error :", err)
	}

	//dbg.Purple(string(resb))
	dbg.Green(rpcAck)
}

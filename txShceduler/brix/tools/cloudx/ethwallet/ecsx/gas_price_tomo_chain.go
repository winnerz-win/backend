package ecsx

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net/http"

	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
)

func (my Sender) RpcGasPrice() *big.Int {
	return my.tomoGasPrice()
}

// TomoGasPrice :
func (my Sender) tomoGasPrice() *big.Int {
	url := my.gasURL + "/gasPrice"
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_gasPrice",
		"params":  []string{},
		"id":      73,
	}
	b, _ := json.Marshal(data)
	buf := bytes.NewBuffer(b)

	resp, err := http.Post(url, "application/json", buf)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		dbg.Red(err)
		return nil
	}

	//dbg.Yellow("code :", resp.StatusCode)
	rb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		dbg.Red(err)
		return nil
	}

	result := struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Result  string `json:"result"`
	}{}

	if err := json.Unmarshal(rb, &result); err != nil {
		dbg.Red(err)
		dbg.Cyan(string(rb))
		return nil
	}

	v := jmath.NewBigDecimal(result.Result)
	val := v.ToBigInteger()
	return val
}

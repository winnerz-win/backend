package ecsx

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
)

type txViewer struct {
	*types.Transaction
}

// Map :
func (my txViewer) Map() map[string]interface{} {
	m := map[string]interface{}{}
	m["Hash"] = my.Hash().Hex()
	m["Cost"] = my.Cost().String()
	m["Value"] = my.Value().String()
	m["ChainId"] = my.ChainId().String()
	m["Gas"] = fmt.Sprintf("%v", my.Gas())
	m["GasPrice"] = my.GasPrice().String()
	m["Nonce"] = fmt.Sprintf("%v", my.Nonce())
	//m["ChkNonce"] = fmt.Sprintf("%v", my.CheckNonce())
	m["Data"] = hex.EncodeToString(my.Data())
	m["Size"] = my.Size().TerminalString()
	return m
}

package ecsx

import (
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
)

type cRPCParam struct {
	JSONPRC string        `json:"jsonrpc"`
	Name    string        `json:"name,omitempty"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

func (my cRPCParam) String() string { return dbg.ToJSONString(my) }

// JSONRPCPARAM :
func JSONRPCPARAM() cRPCParam {
	return cRPCParam{
		JSONPRC: "2.0",
		ID:      1,
	}
}

type cRPCAck struct {
	JSONPRC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func (my cRPCAck) String() string { return dbg.ToJSONString(my) }

// JSONPRCACK :
func JSONPRCACK() cRPCAck {
	return cRPCAck{}
}

// IntString :
func (my cRPCAck) IntString() string {
	return jmath.NewBigDecimal(my.Result).ToString()
}

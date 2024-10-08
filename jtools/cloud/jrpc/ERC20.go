package jrpc

import (
	"jtools/cloud/ebcm"
	"jtools/cloud/ebcm/abi"
	"jtools/dbg"
)

type cERC20 struct {
	IContract
}

type EventTransfer struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}
type TokenTransfer struct {
	Addresss string `bson:"addresss" json:"addresss"`
	From     string `bson:"from" json:"from"`
	To       string `bson:"to" json:"to"`
	Value    string `bson:"value" json:"value"`
}

func (my EventTransfer) TokenTransfer(token string) TokenTransfer {
	return TokenTransfer{
		Addresss: dbg.TrimToLower(token),
		From:     my.From,
		To:       my.To,
		Value:    my.Value,
	}
}

type EventApproval struct {
	Owner   string `json:"owner"`
	Spender string `json:"spender"`
	Value   string `json:"value"`
}

var (
	ERC20 = cERC20{
		IContract: NewContract(
			ebcm.MethodIDDataMap(
				ebcm.MethodID(
					"transfer",
					abi.TypeList{
						abi.Address,
						abi.Uint256,
					},
					func(rs abi.RESULT, item *ebcm.TransactionBlock) {
						m := item.NewCustomInputParse()
						m["recipient"] = rs.Address(0)
						m["amount"] = rs.Uint256(1)
					},
				),

				ebcm.MethodID(
					"approve",
					abi.TypeList{
						abi.Address,
						abi.Uint256,
					},
					func(rs abi.RESULT, item *ebcm.TransactionBlock) {
						m := item.NewCustomInputParse()
						m["spender"] = rs.Address(0)
						m["amount"] = rs.Uint256(1)
					},
				),

				ebcm.MethodID(
					"transferFrom",
					abi.TypeList{
						abi.Address,
						abi.Address,
						abi.Uint256,
					},
					func(rs abi.RESULT, item *ebcm.TransactionBlock) {
						m := item.NewCustomInputParse()
						m["sender"] = rs.Address(0)
						m["recipient"] = rs.Address(1)
						m["amount"] = rs.Uint256(2)
					},
				),

				ebcm.MethodID(
					"increaseAllowance",
					abi.TypeList{
						abi.Address,
						abi.Uint256,
					},
					func(rs abi.RESULT, item *ebcm.TransactionBlock) {
						m := item.NewCustomInputParse()
						m["spender"] = rs.Address(0)
						m["amount"] = rs.Uint256(1)
					},
				),

				ebcm.MethodID(
					"decreaseAllowance",
					abi.TypeList{
						abi.Address,
						abi.Uint256,
					},
					func(rs abi.RESULT, item *ebcm.TransactionBlock) {
						m := item.NewCustomInputParse()
						m["spender"] = rs.Address(0)
						m["amount"] = rs.Uint256(1)
					},
				),
			),
			EventMap{
				ebcm.MakeTopicName("event Transfer(address indexed from, address indexed to, uint256 value);"): {
					Name: "Transfer",
					Parse: func(log ebcm.TxLog) interface{} {
						return EventTransfer{
							From:  log.Topics[1].Address(),
							To:    log.Topics[2].Address(),
							Value: log.Data[0].Number(),
						}
					},
				},
				ebcm.MakeTopicName("event Approval(address indexed owner, address indexed spender, uint value);"): {
					Name: "Approval",
					Parse: func(log ebcm.TxLog) interface{} {
						return EventApproval{
							Owner:   log.Topics[1].Address(),
							Spender: log.Topics[2].Address(),
							Value:   log.Data[0].Number(),
						}
					},
				},
			},
		),
	}
)

type ERC20Info struct {
	Address  string `bson:"address" json:"address"`
	Name     string `bson:"name" json:"name"`
	Symbol   string `bson:"symbol" json:"symbol"`
	Decimals uint8  `bson:"decimals" json:"decimals"`
}

func (my ERC20Info) String() string { return dbg.ToJsonString(my) }

func (ERC20Info) TagString() []string {
	return []string{
		"address", "토큰 주소",
		"name", "토큰 이름",
		"symbol", "토큰 심볼",
		"decimals", "토큰 디시멀",
	}
}

func (my cERC20) ERC20Info(caller *ebcm.Sender, reader IReader, f func(erc20_info ERC20Info)) {
	info := ERC20Info{
		Address: reader.Contract(),
	}
	my.Name(caller, reader, func(name string) {
		info.Name = name
	})
	my.Symbol(caller, reader, func(symbol string) {
		info.Symbol = symbol
	})
	my.Decimals(caller, reader, func(decimals uint8) {
		info.Decimals = decimals
	})

	f(info)
}

func (cERC20) Name(caller *ebcm.Sender, reader IReader, f func(name string)) error {
	return abi.Call2(
		caller,
		reader.Contract(),
		abi.Method{
			Name:   "name",
			Params: abi.NewParams(),
			Returns: abi.NewReturns(
				abi.String,
			),
		},
		reader.CallerAddress(),
		func(rs abi.RESULT) {
			f(rs.Text(0))
		},
		caller.IsDebug(),
	)
}

func (cERC20) Symbol(caller *ebcm.Sender, reader IReader, f func(symbol string)) error {
	return abi.Call2(
		caller,
		reader.Contract(),
		abi.Method{
			Name:   "symbol",
			Params: abi.NewParams(),
			Returns: abi.NewReturns(
				abi.String,
			),
		},
		reader.CallerAddress(),
		func(rs abi.RESULT) {
			f(rs.Text(0))
		},
		caller.IsDebug(),
	)
}

func (cERC20) Decimals(caller *ebcm.Sender, reader IReader, f func(decimals uint8)) error {
	return abi.Call2(
		caller,
		reader.Contract(),
		abi.Method{
			Name:   "decimals",
			Params: abi.NewParams(),
			Returns: abi.NewReturns(
				abi.Uint8,
			),
		},
		reader.CallerAddress(),
		func(rs abi.RESULT) {
			f(rs.Uint8(0))
		},
		caller.IsDebug(),
	)
}

func (cERC20) TotalSupply(caller *ebcm.Sender, reader IReader, f func(totalSupply string)) error {
	return abi.Call2(
		caller,
		reader.Contract(),
		abi.Method{
			Name:   "totalSupply",
			Params: abi.NewParams(),
			Returns: abi.NewReturns(
				abi.Uint256,
			),
		},
		reader.CallerAddress(),
		func(rs abi.RESULT) {
			f(rs.Uint256(0))
		},
		caller.IsDebug(),
	)
}

func (cERC20) BalanceOf(caller *ebcm.Sender, reader IReader, account string, f func(amount string)) error {
	return abi.Call2(
		caller,
		reader.Contract(),
		abi.Method{
			Name: "balanceOf",
			Params: abi.NewParams(
				abi.NewAddress(account),
			),
			Returns: abi.NewReturns(
				abi.Uint256,
			),
		},
		reader.CallerAddress(),
		func(rs abi.RESULT) {
			f(rs.Uint256(0))
		},
		caller.IsDebug(),
	)
}

func (cERC20) Allowance(caller *ebcm.Sender, reader IReader, owner, spender string, f func(amount string)) error {
	return abi.Call2(
		caller,
		reader.Contract(),
		abi.Method{
			Name: "allowance",
			Params: abi.NewParams(
				abi.NewAddress(owner),
				abi.NewAddress(spender),
			),
			Returns: abi.NewReturns(
				abi.Uint256,
			),
		},
		reader.CallerAddress(),
		func(rs abi.RESULT) {
			f(rs.Uint256(0))
		},
		caller.IsDebug(),
	)
}

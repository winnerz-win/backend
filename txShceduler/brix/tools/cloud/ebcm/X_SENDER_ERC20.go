package ebcm

import (
	"strings"
	"txscheduler/brix/tools/cloud/ebcm/abi"
	"txscheduler/brix/tools/dbg"
)

type _erc20 string

func (my _erc20) Contract() string { return string(my) }
func IERC20(contract string) _erc20 {
	return _erc20(strings.TrimSpace(contract))
}

func (my _erc20) Name(finder *Sender, f func(name string)) error {
	return abi.Call2(
		finder,
		my.Contract(),
		abi.Method{
			Name:   "name",
			Params: abi.NewParams(),
			Returns: abi.NewReturns(
				abi.String,
			),
		},
		my.Contract(),
		func(rs abi.RESULT) {
			if f != nil {
				f(rs.Text(0))
			}
		},
	)
}
func (my _erc20) Symbol(finder *Sender, f func(symbol string)) error {
	return abi.Call2(
		finder,
		my.Contract(),
		abi.Method{
			Name:   "symbol",
			Params: abi.NewParams(),
			Returns: abi.NewReturns(
				abi.String,
			),
		},
		my.Contract(),
		func(rs abi.RESULT) {
			if f != nil {
				f(rs.Text(0))
			}
		},
	)
}

func (my _erc20) Decimals(finder *Sender, f func(decimals uint8)) error {
	return abi.Call2(
		finder,
		my.Contract(),
		abi.Method{
			Name:   "decimals",
			Params: abi.NewParams(),
			Returns: abi.NewReturns(
				abi.Uint8,
			),
		},
		my.Contract(),
		func(rs abi.RESULT) {
			if f != nil {
				f(rs.Uint8(0))
			}
		},
	)
}

func (my _erc20) TotalSupply(finder *Sender, f func(amount string)) error {
	return abi.Call2(
		finder,
		my.Contract(),
		abi.Method{
			Name:   "totalSupply",
			Params: abi.NewParams(),
			Returns: abi.NewReturns(
				abi.Uint256,
			),
		},
		my.Contract(),
		func(rs abi.RESULT) {
			if f != nil {
				f(rs.Uint256(0))
			}
		},
	)
}

type ERC20Info struct {
	Contract    string `bson:"contract" json:"contract"` //COA
	Name        string `bson:"name" json:"name"`
	Symbol      string `bson:"symbol" json:"symbol"`
	Decimals    uint8  `bson:"decimals" json:"decimals"`
	TotalSupply string `bson:"total_supply" json:"total_supply"`
}

func (my ERC20Info) String() string { return dbg.ToJsonString(my) }

func (my ERC20Info) Price(wei string) string { return WeiToToken(wei, my.Decimals) }

func (my _erc20) Info(finder *Sender, f func(info ERC20Info), isWithoutTotalSupply ...bool) error {
	info := ERC20Info{
		Contract: my.Contract(),
	}
	if err := my.Name(finder, func(name string) { info.Name = name }); err != nil {
		return err
	}
	if err := my.Symbol(finder, func(symbol string) { info.Symbol = symbol }); err != nil {
		return err
	}
	if err := my.Decimals(finder, func(decimals uint8) { info.Decimals = decimals }); err != nil {
		return err
	}
	if !dbg.IsTrue(isWithoutTotalSupply) {
		if err := my.TotalSupply(finder, func(amount string) { info.TotalSupply = amount }); err != nil {
			return err
		}
	}

	if f != nil {
		f(
			info,
		)
	}

	return nil
}

func (my _erc20) BalanceOf(finder *Sender, address string, f func(amount string)) error {
	return abi.Call2(
		finder,
		my.Contract(),
		abi.Method{
			Name: "balanceOf",
			Params: abi.NewParams(
				abi.NewAddress(address),
			),
			Returns: abi.NewReturns(
				abi.Uint256,
			),
		},
		address,
		func(rs abi.RESULT) {
			if f != nil {
				f(rs.Uint(0))
			}
		},
	)
}

func (my _erc20) Allowance(finder *Sender, owner, spender string, f func(amount string)) error {
	return abi.Call2(
		finder,
		my.Contract(),
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
		owner,
		func(rs abi.RESULT) {
			if f != nil {
				f(rs.Uint(0))
			}
		},
	)

}

func (my _erc20) ApproveAll(
	sender *Sender,
	private string,
	spender string,
	f func(r XSendResult),
) error {
	return my.Approve(
		sender,
		private,
		spender,
		UINT256MAX,
		f,
	)
}

func (my _erc20) Approve(
	sender *Sender,
	private string,
	spender, amount string,
	f func(r XSendResult),
) error {
	return sender.XPipe(
		private,
		my.Contract(),
		sender.MakePadBytesABI(
			"approve",
			abi.TypeList{
				abi.NewAddress(spender),
				abi.NewUint256(amount),
			},
		),
		ZERO,
		nil, nil, nil,
		f,
	)
}

func (my _erc20) TransferFrom(
	sender *Sender,
	private string,
	owner, to string,
	amount string,
	f func(r XSendResult),
) error {
	return sender.XPipe(
		private,
		my.Contract(),
		sender.MakePadBytesABI(
			"transferFrom",
			abi.TypeList{
				abi.NewAddress(owner),
				abi.NewAddress(to),
				abi.NewUint256(amount),
			},
		),
		ZERO,
		nil, nil, nil,
		f,
	)
}

func (my _erc20) Transfer(
	sender *Sender,
	private string,
	to string,
	amount string,
	f func(r XSendResult),
) error {
	return sender.XPipe(
		private,
		my.Contract(),
		sender.MakePadBytesABI(
			"transfer",
			abi.TypeList{
				abi.NewAddress(to),
				abi.NewUint256(amount),
			},
		),
		ZERO,
		nil, nil, nil,
		f,
	)
}

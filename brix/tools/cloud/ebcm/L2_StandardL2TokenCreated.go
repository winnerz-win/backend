package ebcm

import (
	"errors"
	"strings"
	"time"
	"txscheduler/brix/tools/cloud/ebcm/abi"
	"txscheduler/brix/tools/dbg"
)

const (
	OptimisticL2TokenFactoryAddress = "0x4200000000000000000000000000000000000012"

	eventStandardL2TokenCreated = "0xceeb8e7d520d7f3b65fc11a262b91066940193b05d4f93df07cfdced0eb551cf"
)

var (
	L2TokenFactory = StandardL2TokenFactory(OptimisticL2TokenFactoryAddress)
)

type StandardL2TokenFactory string

func (my StandardL2TokenFactory) Contract() string { return string(my) }

/*
	contract L2StandardTokenFactory {
		event StandardL2TokenCreated(address indexed _l1Token, address indexed _l2Token);
		//0xceeb8e7d520d7f3b65fc11a262b91066940193b05d4f93df07cfdced0eb551cf

		* @dev Creates an instance of the standard ERC20 token on L2.
		* @param _l1Token Address of the corresponding L1 token.
		* @param _name ERC20 name.
		* @param _symbol ERC20 symbol.

		function createStandardL2Token(
			address _l1Token,
			string memory _name,
			string memory _symbol
		) external {
			require(_l1Token != address(0), "Must provide L1 token address");

			L2StandardERC20 l2Token = new L2StandardERC20(
				Lib_PredeployAddresses.L2_STANDARD_BRIDGE,
				_l1Token,
				_name,
				_symbol
			);

			emit StandardL2TokenCreated(_l1Token, address(l2Token));
		}
	}
*/

func (my StandardL2TokenFactory) createStandardL2Token(
	sender *Sender,
	private string,

	_l1token string,
	_name, _symbol string,
	f func(r XSendResult),
) error {

	if sender.TXNTYPE() != TXN_LEGACY {
		return errors.New("Must TXN_LEGACY only!")
	}

	_l1token = strings.TrimSpace(_l1token)
	_name = strings.TrimSpace(_name)
	_symbol = strings.TrimSpace(_symbol)

	if !IsAddress(_l1token) {
		return errors.New("Must provide L1 token address")
	}
	if f == nil {
		return errors.New("result f is nil")
	}

	return sender.XPipe(
		private,
		my.Contract(),
		MakePadBytesABI(
			"createStandardL2Token",
			abi.TypeList{
				abi.NewAddress(_l1token),
				abi.NewString(_name),
				abi.NewString(_symbol),
			},
		),
		ZERO,
		nil, nil, nil,
		f,
	)
}

type L2TokenPairResult struct {
	L1Token  string `bson:"l1_token" json:"l1_token"`
	TxHash   string `bson:"tx_hash" json:"tx_hash"`
	IsPaird  bool   `bson:"is_paird" json:"is_paird"`
	L2Token  string `bson:"l2_token" json:"l2_token"`
	L2Name   string `bson:"l2_name" json:"l2_name"`
	L2Symbol string `bson:"l2_symbol" json:"l2_symbol"`
}

func (my L2TokenPairResult) String() string { return dbg.ToJsonString(my) }

func (my StandardL2TokenFactory) CreateStandardL2Token(
	sender *Sender,
	private string,

	_l1token string,
	_name, _symbol string,
	f func(r L2TokenPairResult),
) error {
	if f == nil {
		return errors.New("result f is nil")
	}

	_l1token = strings.ToLower(strings.TrimSpace(_l1token))
	_name = strings.TrimSpace(_name)
	_symbol = strings.TrimSpace(_symbol)

	hash := ""
	if err := my.createStandardL2Token(
		sender,
		private,
		_l1token,
		_name, _symbol,
		func(r XSendResult) {
			hash = r.Hash
		},
	); err != nil {
		return err
	}

	r := L2TokenPairResult{
		L1Token:  _l1token,
		TxHash:   hash,
		L2Name:   _name,
		L2Symbol: _symbol,
	}

	chkCnt := 0
EXIT:
	for {
		time.Sleep(time.Second)
		chkCnt++

		tx, pending, err := sender.TransactionByHash(hash)
		_ = pending
		if err != nil {
			f(r)
			return err
		}

		if tx.IsReceiptedByHash {
			for _, log := range tx.Logs {
				if log.Topics.GetName() == eventStandardL2TokenCreated {
					r.IsPaird = log.Topics[1].Address() == _l1token
					r.L2Token = log.Topics[2].Address()
					break EXIT
				}
			}
			r.IsPaird = false
			break EXIT
		} else {
			dbg.PurpleItalic("[", hash, "] Wait :", chkCnt)
		}
	}

	f(r)

	return nil

}

type _IL2StandardERC20 _erc20

func (my _IL2StandardERC20) erc20() _erc20 { return _erc20(my) }

func IL2StandardERC20(contract string) _IL2StandardERC20 {
	return _IL2StandardERC20(IERC20(contract))
}

type StandardERC20Info struct {
	ERC20Info `bson:",inline" json:",inline"`
	L1Token   string `bson:"l1_token" json:"l1_token"`
	L2Bridge  string `bson:"l2_bridge" json:"l2_bridge"`
}

func (my StandardERC20Info) String() string { return dbg.ToJsonString(my) }

func (my StandardERC20Info) Price(wei string) string { return my.ERC20Info.Price(wei) }

func (my _IL2StandardERC20) Info(finder *Sender, f func(info StandardERC20Info), isWithoutTotalSupply ...bool) error {
	l2_info := StandardERC20Info{}

	err := my.erc20().Info(
		finder,
		func(info ERC20Info) {
			l2_info.ERC20Info = info
		},
		isWithoutTotalSupply...,
	)
	if err != nil {
		return err
	}

	my.L1Token(
		finder,
		func(l1_token string) {
			l2_info.L1Token = l1_token
		},
	)

	my.L2Bridge(
		finder,
		func(l2_bridge string) {
			l2_info.L2Bridge = l2_bridge
		},
	)

	if f != nil {
		f(l2_info)
	}
	return nil
}

func (my _IL2StandardERC20) L1Token(finder *Sender, f func(l1_token string)) error {
	return abi.Call2(
		finder,
		my.erc20().Contract(),
		abi.Method{
			Name:   "l1Token",
			Params: abi.NewParams(),
			Returns: abi.NewReturns(
				abi.Address,
			),
		},
		"",
		func(rs abi.RESULT) {
			if f != nil {
				f(rs.Address(0))
			}

		},
	)
}

func (my _IL2StandardERC20) L2Bridge(finder *Sender, f func(l2_bridge string)) error {
	return abi.Call2(
		finder,
		my.erc20().Contract(),
		abi.Method{
			Name:   "l2Bridge",
			Params: abi.NewParams(),
			Returns: abi.NewReturns(
				abi.Address,
			),
		},
		"",
		func(rs abi.RESULT) {
			if f != nil {
				f(rs.Address(0))
			}

		},
	)
}

func (my _IL2StandardERC20) BalanceOf(finder *Sender, address string, f func(amount string)) error {
	return my.erc20().BalanceOf(
		finder,
		address,
		f,
	)
}

func (my _IL2StandardERC20) Allowance(finder *Sender, owner, spender string, f func(amount string)) error {
	return my.erc20().Allowance(
		finder,
		owner,
		spender,
		f,
	)
}

func (my _IL2StandardERC20) ApproveAll(
	sender *Sender,
	private string,
	spender string,
	f func(r XSendResult),
) error {
	return my.erc20().ApproveAll(
		sender,
		private,
		spender,
		f,
	)
}

func (my _IL2StandardERC20) Approve(
	sender *Sender,
	private string,
	spender, amount string,
	f func(r XSendResult),
) error {
	return my.erc20().Approve(
		sender,
		private,
		spender, amount,
		f,
	)
}

func (my _IL2StandardERC20) TransferFrom(
	sender *Sender,
	private string,
	owner, to string,
	amount string,
	f func(r XSendResult),
) error {
	return my.erc20().TransferFrom(
		sender,
		private,
		owner, to,
		amount,
		f,
	)
}

func (my _IL2StandardERC20) Transfer(
	sender *Sender,
	private string,
	to string,
	amount string,
	f func(r XSendResult),
) error {
	return my.erc20().Transfer(
		sender,
		private,
		to,
		amount,
		f,
	)
}

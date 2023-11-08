package rpc

import (
	"strings"
	"txscheduler/brix/tools/cloud/ebcm"
	"txscheduler/brix/tools/cloud/ebcm/abi"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
)

type cERC721 struct {
	IContract
}

type EventTransferNFT struct {
	From    string `bson:"from" json:"from"` //address(0) : mint
	To      string `bson:"to" json:"to"`     //address(0) : burn
	TokenId string `bson:"token_id" json:"token_id"`
}

func (my EventTransferNFT) String() string { return dbg.ToJsonString(my) }

type EventApprovalNFT struct {
	Owner    string `bson:"owner" json:"owner"`
	Approved string `bson:"approved" json:"approved"`
	TokenId  string `bson:"token_id" json:"token_id"`
}

func (my EventApprovalNFT) String() string { return dbg.ToJsonString(my) }

type EventApprovalAllNFT struct {
	Owner    string `bson:"owner" json:"owner"`
	Operator string `bson:"operator" json:"operator"`
	Approved bool   `bson:"approved" json:"approved"`
}

func (my EventApprovalAllNFT) String() string { return dbg.ToJsonString(my) }

var (
	ERC721 = cERC721{
		IContract: newContract(
			ebcm.MethodIDDataMap(
				ebcm.MethodID(
					"setBaseURI",
					abi.TypeList{
						abi.String,
					},
					func(rs abi.RESULT, item *ebcm.TransactionBlock) {
						m := item.NewCustomInputParse()
						m["newURI"] = rs.Text(0)
					},
				),
				ebcm.MethodID(
					"safeTransferFrom",
					abi.TypeList{
						abi.Address,
						abi.Address,
						abi.Uint256,
					},
					func(rs abi.RESULT, item *ebcm.TransactionBlock) {
						m := item.NewCustomInputParse()
						m["from"] = rs.Address(0)
						m["to"] = rs.Address(1)
						m["tokenId"] = rs.Uint256(2)
					},
				),
				ebcm.MethodID(
					"safeTransferFrom",
					abi.TypeList{
						abi.Address,
						abi.Address,
						abi.Uint256,
						abi.Bytes,
					},
					func(rs abi.RESULT, item *ebcm.TransactionBlock) {
						m := item.NewCustomInputParse()
						m["from"] = rs.Address(0)
						m["to"] = rs.Address(1)
						m["tokenId"] = rs.Uint256(2)
						m["data"] = rs.Bytes(3)
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
						m["from"] = rs.Address(0)
						m["to"] = rs.Address(1)
						m["tokenId"] = rs.Uint256(2)
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
						m["to"] = rs.Address(0)
						m["tokenId"] = rs.Uint256(1)
					},
				),

				ebcm.MethodID(
					"setApprovalForAll",
					abi.TypeList{
						abi.Address,
						abi.Bool,
					},
					func(rs abi.RESULT, item *ebcm.TransactionBlock) {
						m := item.NewCustomInputParse()
						m["operator"] = rs.Address(0)
						m["approved"] = rs.Bool(1)
					},
				),

				ebcm.MethodID(
					"setApprovalForAll",
					abi.TypeList{
						abi.Address,
						abi.Bool,
					},
					func(rs abi.RESULT, item *ebcm.TransactionBlock) {
						m := item.NewCustomInputParse()
						m["operator"] = rs.Address(0)
						m["approved"] = rs.Bool(1)
					},
				),

				ebcm.MethodID(
					"mint",
					abi.TypeList{
						abi.Address,
						abi.Uint256,
					},
					func(rs abi.RESULT, item *ebcm.TransactionBlock) {
						m := item.NewCustomInputParse()
						m["to"] = rs.Address(0)
						m["token_id"] = rs.Uint(1)
					},
				),
				ebcm.MethodID(
					"burn",
					abi.TypeList{
						abi.Uint256,
					},
					func(rs abi.RESULT, item *ebcm.TransactionBlock) {
						m := item.NewCustomInputParse()
						m["tokenId"] = rs.Uint256(0)
					},
				),
			),

			MakeEventMap(event_map{
				ebcm.MakeTopicName("event Transfer(address indexed from, address indexed to, uint256 indexed tokenId);"): {
					name: "TransferNFT",
					parse: func(log ebcm.TxLog) interface{} {
						return EventTransferNFT{
							From:    log.Topics[1].Address(),
							To:      log.Topics[2].Address(),
							TokenId: log.Topics[3].Number(),
						}
					},
				},
				ebcm.MakeTopicName("event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId);"): {
					name: "ApprovalNFT",
					parse: func(log ebcm.TxLog) interface{} {
						return EventApprovalNFT{
							Owner:    log.Topics[1].Address(),
							Approved: log.Topics[2].Address(),
							TokenId:  log.Topics[3].Number(),
						}
					},
				},
				ebcm.MakeTopicName("event ApprovalForAll(address indexed owner, address indexed operator, bool approved);"): {
					name: "ApprovalForAllNFT",
					parse: func(log ebcm.TxLog) interface{} {
						return EventApprovalAllNFT{
							Owner:    log.Topics[1].Address(),
							Operator: log.Topics[2].Address(),
							Approved: log.Data[0].Bool(),
						}
					},
				},
			}),
		),
	}
)

type ERC721Info struct {
	Address string `bson:"address" json:"address"`
	Name    string `bson:"name" json:"name"`
	Symbol  string `bson:"symbol" json:"symbol"`
}

func (my ERC721Info) String() string { return dbg.ToJsonString(my) }

func (my cERC721) ERC721Info(
	caller *ebcm.Sender, reader IReader,
	f func(info ERC721Info),
) error {
	info := ERC721Info{
		Address: reader.Contract(),
	}

	if err := my.Name(caller, reader,
		func(_name string) { info.Name = _name },
	); err != nil {
		return err
	}

	if err := my.Symbol(
		caller, reader, func(_symbol string) { info.Symbol = _symbol },
	); err != nil {
		return err
	}

	f(info)
	return nil
}

func (cERC721) Name(
	caller *ebcm.Sender, reader IReader,
	f func(_name string),
) error {
	return caller.Call(
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
		_is_debug_call,
	)
}

func (cERC721) Symbol(
	caller *ebcm.Sender, reader IReader,
	f func(_symbol string),
) error {
	return caller.Call(
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
		_is_debug_call,
	)
}

/**
 * @dev Returns a token ID at a given `index` of all the tokens stored by the contract.
 * Use along with {totalSupply} to enumerate all tokens.
 */
func (cERC721) TokenByIndex(
	caller *ebcm.Sender, reader IReader,
	index string,
	f func(_token_id string),
) error {
	return caller.Call(
		reader.Contract(),
		abi.Method{
			Name: "tokenByIndex",
			Params: abi.NewParams(
				abi.NewUint256(index),
			),
			Returns: abi.NewReturns(
				abi.Uint256,
			),
		},
		reader.CallerAddress(),
		func(rs abi.RESULT) {
			f(rs.Uint256(0))
		},
		_is_debug_call,
	)
}

func (cERC721) TotalSupply(
	caller *ebcm.Sender, reader IReader,
	f func(_total_supply string),
) error {
	return caller.Call(
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
		_is_debug_call,
	)
}

func (cERC721) TokenURI(
	caller *ebcm.Sender, reader IReader,
	tokenId interface{},
	f func(_token_uri string),
) error {
	return caller.Call(
		reader.Contract(),
		abi.Method{
			Name: "tokenURI",
			Params: abi.NewParams(
				abi.NewUint256(tokenId),
			),
			Returns: abi.NewReturns(
				abi.String,
			),
		},
		reader.CallerAddress(),
		func(rs abi.RESULT) {
			f(rs.Text(0))
		},
		_is_debug_call,
	)
}

func (cERC721) BaseURI(
	caller *ebcm.Sender, reader IReader,
) (string, error) {
	base_uri := ""
	err := caller.Call(
		reader.Contract(),
		abi.Method{
			Name:   "baseURI",
			Params: abi.NewParams(),
			Returns: abi.NewReturns(
				abi.String,
			),
		},
		reader.CallerAddress(),
		func(rs abi.RESULT) {
			base_uri = rs.Text(0)
		},
		_is_debug_call,
	)
	return base_uri, err
}

/**
 * @dev Returns a token ID owned by `owner` at a given `index` of its token list.
 * Use along with {balanceOf} to enumerate all of ``owner``'s tokens.
 */
func (cERC721) TokenOfOwnerByIndex(
	caller *ebcm.Sender, reader IReader,
	owner, index string,
	f func(_token_id string),
) error {
	return caller.Call(
		reader.Contract(),
		abi.Method{
			Name: "tokenOfOwnerByIndex",
			Params: abi.NewParams(
				abi.NewAddress(owner),
				abi.NewUint256(index),
			),
			Returns: abi.NewReturns(
				abi.Uint256,
			),
		},
		reader.CallerAddress(),
		func(rs abi.RESULT) {
			f(rs.Uint256(0))
		},
		_is_debug_call,
	)
}

func (cERC721) BalanceOf(
	caller *ebcm.Sender, reader IReader,
	owner string,
	f func(_balance string),
) error {
	return caller.Call(
		reader.Contract(),
		abi.Method{
			Name: "balanceOf",
			Params: abi.NewParams(
				abi.NewAddress(owner),
			),
			Returns: abi.NewReturns(
				abi.Uint256,
			),
		},
		reader.CallerAddress(),
		func(rs abi.RESULT) {
			f(rs.Uint256(0))
		},
		_is_debug_call,
	)
}

func (my cERC721) OwnerTokenAll(
	caller *ebcm.Sender, reader IReader,
	owner string,
) ([]string, error) {
	var re_err error
	list := []string{}
	index := "0"
	for {
		if err := my.TokenOfOwnerByIndex(
			caller, reader,
			owner,
			index,
			func(_token_id string) {
				list = append(list, _token_id)
			},
		); err != nil {
			if !strings.Contains(err.Error(), "index out of bounds") {
				re_err = err
			}
			break
		}
		index = jmath.ADD(index, 1)
	} //for

	return list, re_err
}

func (cERC721) OwnerOf(
	caller *ebcm.Sender, reader IReader,
	tokenId string,
	f func(_owner string),
) error {
	return caller.Call(
		reader.Contract(),
		abi.Method{
			Name: "ownerOf",
			Params: abi.NewParams(
				abi.NewUint256(tokenId),
			),
			Returns: abi.NewReturns(
				abi.Address,
			),
		},
		reader.CallerAddress(),
		func(rs abi.RESULT) {
			f(rs.Address(0))
		},
		_is_debug_call,
	)
}

// IsApprovedForAll : write-func -> setApprovalForAll(operator, approved)
func (cERC721) IsApprovedForAll(
	caller *ebcm.Sender, reader IReader,
	owner string,
	operator string,
) (bool, error) {
	isAllow := false
	err := caller.Call(
		reader.Contract(),
		abi.Method{
			Name: "isApprovedForAll",
			Params: abi.NewParams(
				abi.NewAddress(owner),
				abi.NewAddress(operator),
			),
			Returns: abi.NewReturns(
				abi.Bool,
			),
		},
		reader.CallerAddress(),
		func(rs abi.RESULT) {
			isAllow = rs.Bool(0)
		},
	)

	return isAllow, err
}

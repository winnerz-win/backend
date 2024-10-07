package rpc

import (
	"txscheduler/brix/tools/cloud/ebcm"
	"txscheduler/brix/tools/cloud/ebcm/abi"
)

type cWinnerzNFT struct {
	IContract

	_NFT
}

func init() {
	WinnerzNFT.IContract.cMerge(
		privateNFT,
	)
}

var (
	WinnerzNFT = cWinnerzNFT{
		IContract: newContract(
			ebcm.MethodIDDataMap(
				ebcm.MethodID(
					"finishAdminTransfer",
					abi.TypeList{},
					func(rs abi.RESULT, item *ebcm.TransactionBlock) {
						m := item.NewCustomInputParse()
						m["finishAdminTransfer"] = true
					},
				),
				ebcm.MethodID(
					"adminTransfer",
					abi.TypeList{
						abi.Address,
						abi.Address,
						abi.Uint,
					},
					func(rs abi.RESULT, item *ebcm.TransactionBlock) {
						m := item.NewCustomInputParse()
						m["from"] = rs.Address(0)
						m["to"] = rs.Address(1)
						m["token_id"] = rs.Uint(2)
					},
				),
				ebcm.MethodID(
					"withdrawETH",
					abi.TypeList{
						abi.Uint,
					},
					func(rs abi.RESULT, item *ebcm.TransactionBlock) {
						m := item.NewCustomInputParse()
						m["value"] = rs.Uint(0)
					},
				),
				ebcm.MethodID(
					"withdrawToken",
					abi.TypeList{
						abi.Address,
						abi.Uint,
					},
					func(rs abi.RESULT, item *ebcm.TransactionBlock) {
						m := item.NewCustomInputParse()
						m["token_contract"] = rs.Address(0)
						m["value"] = rs.Uint(1)
					},
				),
				ebcm.MethodID(
					"multiTransferToken",
					abi.TypeList{
						abi.Address,      //token
						abi.AddressArray, //receivers
						abi.Uint256Array, //values
					},
					func(rs abi.RESULT, item *ebcm.TransactionBlock) {
						m := item.NewCustomInputParse()
						m["token_contract"] = rs.Address(0)
						m["receivers"] = rs.AddressArray(1)
						m["values"] = rs.Uint256Array(2)
					},
				),
				ebcm.MethodID(
					"multiTransferETH",
					abi.TypeList{
						abi.AddressArray, //receivers
						abi.Uint256Array, //values
					},
					func(rs abi.RESULT, item *ebcm.TransactionBlock) {
						m := item.NewCustomInputParse()
						m["receivers"] = rs.AddressArray(0)
						m["values"] = rs.Uint256Array(1)
					},
				),
			),
			MakeEventMap(event_map{
				ebcm.MakeTopicName("event MultiTransferETH(address indexed sender, address indexed receiver, uint256 indexed value);"): {
					name: "MultiTransferETH",
					parse: func(log ebcm.TxLog) interface{} {
						return EventMultiTransferETH{
							Sender:   log.Indexed_1().Address(),
							Receiver: log.Indexed_2().Address(),
							Value:    log.Indexed_3().Number(),
						}
					},
				},
				ebcm.MakeTopicName("event MultiTransferToken(address indexed token, address indexed sender, address indexed receiver, uint256 value);"): {
					name: "MultiTransferToken",
					parse: func(log ebcm.TxLog) interface{} {
						return EventMultiTransferToken{
							Token:    log.Indexed_1().Address(),
							Sender:   log.Indexed_2().Address(),
							Receiver: log.Indexed_3().Address(),
							Value:    log.Data[0].Number(),
						}
					},
				},
			}),
		),
	}
)

type EventMultiTransferETH struct {
	Sender   string `bson:"sender" json:"sender"`
	Receiver string `bson:"receiver" json:"receiver"`
	Value    string `bson:"value" json:"value"`
}

type EventMultiTransferToken struct {
	Token    string `bson:"token" json:"token"`
	Sender   string `bson:"sender" json:"sender"`
	Receiver string `bson:"receiver" json:"receiver"`
	Value    string `bson:"value" json:"value"`
}

func (cWinnerzNFT) IsFinishAdminTransfer(
	caller *ebcm.Sender, reader IReader,
) (bool, error) {
	is_finished := false
	err := caller.Call(
		reader.Contract(),
		abi.Method{
			Name:   "isFinishAdminTransfer",
			Params: abi.NewParams(),
			Returns: abi.NewReturns(
				abi.Bool,
			),
		},
		reader.CallerAddress(),
		func(rs abi.RESULT) {
			is_finished = rs.Bool(0)
		},
		_is_debug_call,
	)
	return is_finished, err
}

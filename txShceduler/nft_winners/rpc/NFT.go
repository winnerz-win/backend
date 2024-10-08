package rpc

import (
	"jtools/cloud/ebcm"
	"jtools/cloud/ebcm/abi"
)

type _NFT struct {
	IContract

	cERC721
	dBlockRecord
}

func init() {
	privateNFT.IContract.cMerge(
		ERC721,
	)
}

var (
	privateNFT = _NFT{
		IContract: newContract(
			ebcm.MethodIDDataMap(
				ebcm.MethodID(
					"grantAdminRole",
					abi.TypeList{
						abi.Address,
					},
					func(rs abi.RESULT, item *ebcm.TransactionBlock) {
						m := item.NewCustomInputParse()
						m["to"] = rs.Address(0)
					},
				),
				ebcm.MethodID(
					"setBaseURI",
					abi.TypeList{
						abi.String,
					},
					func(rs abi.RESULT, item *ebcm.TransactionBlock) {
						m := item.NewCustomInputParse()
						m["setBaseURI"] = rs.Text(0)
					},
				),
				// ebcm.MethodID(
				// 	"mint",
				// 	abi.TypeList{
				// 		abi.Address,
				// 		abi.Uint,
				// 	},
				// 	func(rs abi.RESULT, item *ebcm.TransactionBlock) {
				// 		m := item.NewCustomInputParse()
				// 		m["to"] = rs.Address(0)
				// 		m["token_id"] = rs.Uint(1)
				// 	},
				// ),
			),
			MakeEventMap(event_map{}),
		),
	}
)

func (_NFT) BaseURI(
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

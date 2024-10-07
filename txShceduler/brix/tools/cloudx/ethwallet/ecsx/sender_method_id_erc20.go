package ecsx

import "txscheduler/brix/tools/cloudx/ethwallet/abmx"

var (
	MethodDeploy0 = cMethodIDData{
		FuncName:  "deploy",
		ABIString: "0x60...",
		MethodID:  "0x60",
		parseInput: func(data string, item *TransactionBlock) {
			item.ContractAddress = ContractAddressNonce(item.From, item.Nonce)
		},
	}
	MethodDeploy = cMethodIDData{
		FuncName:  "deploy",
		ABIString: "0x60606040...",
		MethodID:  "0x60606040",
		parseInput: func(data string, item *TransactionBlock) {
			item.ContractAddress = ContractAddressNonce(item.From, item.Nonce)
		},
	}
	MethodDeploy2 = cMethodIDData{
		FuncName:  "deploy",
		ABIString: "0x608060...",
		MethodID:  "0x608060",
		parseInput: func(data string, item *TransactionBlock) {
			item.ContractAddress = ContractAddressNonce(item.From, item.Nonce)
		},
	}
	MethodTransferOwnerShip = MakeMethodIDDataABM(
		GetMethodIDABM(
			"transferOwnership",
			abmx.Address,
		),
		func(rs abmx.RESULT, item *TransactionBlock) {
			if !rs.IsError {
				cip := item.NewCustomInputParse()
				cip["to"] = rs.Address(0)
			}
		},
	)

	MethodTransfer = MakeMethodIDDataABM(
		GetMethodIDABM(
			"transfer",
			abmx.Address,
			abmx.Uint256,
		),
		func(rs abmx.RESULT, item *TransactionBlock) {
			if !rs.IsError {
				item.To = rs.Address(0)
				item.Amount = rs.Uint256(1)
			}
		},
	)

	MethodTranferFrom = MakeMethodIDDataABM(
		GetMethodIDABM(
			"transferFrom",
			abmx.Address,
			abmx.Address,
			abmx.Uint256,
		),
		func(rs abmx.RESULT, item *TransactionBlock) {
			if !rs.IsError {
				cip := item.NewCustomInputParse()
				cip["from"] = rs.Address(0)
				cip["to"] = rs.Address(1)
				cip["value"] = rs.Uint256(2)
				item.To = cip.Address("to")
				item.Amount = cip.Number("value")
			}
		},
	)

	MethodApprove = MakeMethodIDDataABM(
		GetMethodIDABM(
			"approve",
			abmx.Address,
			abmx.Uint256,
		),
		func(rs abmx.RESULT, item *TransactionBlock) {
			if !rs.IsError {
				cip := item.NewCustomInputParse()
				cip["to"] = rs.Address(0)
				cip["value"] = rs.Uint256(1)
			}
		},
	)

	MethodAllowance = MakeMethodIDDataABM(
		GetMethodIDABM(
			"allowance",
			abmx.Address,
			abmx.Address,
		),
		func(rs abmx.RESULT, item *TransactionBlock) {
			if !rs.IsError {
				cip := item.NewCustomInputParse()
				cip["owner"] = rs.Address(0)
				cip["spender"] = rs.Address(1)
			}
		},
	)

	MethodIssue = MakeMethodIDDataABM(
		GetMethodIDABM(
			"issue",
			abmx.Uint256,
		),
		func(rs abmx.RESULT, item *TransactionBlock) {
			if !rs.IsError {
				cip := item.NewCustomInputParse()
				cip["value"] = rs.Uint256(0)
			}
		},
	)

	MethodMint = MakeMethodIDDataABM(
		GetMethodIDABM(
			"mint",
			abmx.Uint256,
		),
		func(rs abmx.RESULT, item *TransactionBlock) {
			if !rs.IsError {
				cip := item.NewCustomInputParse()
				cip["value"] = rs.Uint256(0)
			}
		},
	)

	MethodRedeem = MakeMethodIDDataABM(
		GetMethodIDABM(
			"redeem",
			abmx.Uint256,
		),
		func(rs abmx.RESULT, item *TransactionBlock) {
			if !rs.IsError {
				cip := item.NewCustomInputParse()
				cip["value"] = rs.Uint256(0)
			}
		},
	)

	MethodBurn = MakeMethodIDDataABM(
		GetMethodIDABM(
			"burn",
			abmx.Uint256,
		),
		func(rs abmx.RESULT, item *TransactionBlock) {
			if !rs.IsError {
				cip := item.NewCustomInputParse()
				cip["value"] = rs.Uint256(0)
			}
		},
	)
)

var methodERC20s = MethodIDDataList{
	MethodDeploy0,
	MethodDeploy,
	MethodDeploy2,
	MethodTransferOwnerShip,
	MethodTransfer,
	MethodTranferFrom,
	MethodApprove,
	MethodAllowance,
	MethodIssue,
	MethodMint,
	MethodRedeem,
	MethodBurn,
}

package ebcm

import (
	"encoding/hex"
	"jtools/cloud/ebcm/abi"
	"strings"
)

var (
	methodDeploy0 = cMethodIDData{
		FuncName:  "deploy",
		ABIString: "0x60...",
		MethodID:  "0x60",
		parseDeploy: func(input_parser abi.InputParser, rs abi.RESULT, item *TransactionBlock) {
			item.ContractAddress = input_parser.ContractAddressNonce(item.From, item.Nonce)
		},
	}

	methodDeploy1 = cMethodIDData{
		FuncName:  "deploy",
		ABIString: "0x60606040...",
		MethodID:  "0x60606040",
		parseDeploy: func(input_parser abi.InputParser, rs abi.RESULT, item *TransactionBlock) {
			item.ContractAddress = input_parser.ContractAddressNonce(item.From, item.Nonce)
		},
	}

	methodDeploy2 = cMethodIDData{
		FuncName:  "deploy",
		ABIString: "0x608060...",
		MethodID:  "0x608060",
		parseDeploy: func(input_parser abi.InputParser, rs abi.RESULT, item *TransactionBlock) {
			item.ContractAddress = input_parser.ContractAddressNonce(item.From, item.Nonce)
		},
	}

	methodTransferOwnerShip = MethodID(
		"transferOwnership",
		abi.TypeList{
			abi.Address,
		},
		func(rs abi.RESULT, item *TransactionBlock) {
			if !rs.IsError {
				m := item.NewCustomInputParse()
				m["to"] = rs.Address(0)
			}
		},
	)

	methodTransfer = MethodID(
		"transfer",
		abi.TypeList{
			abi.Address,
			abi.Uint256,
		},
		func(rs abi.RESULT, item *TransactionBlock) {
			if !rs.IsError {
				item.To = rs.Address(0)
				item.Amount = rs.Uint256(1)
			}
		},
	)

	methodTranferFrom = MethodID(
		"transferFrom",
		abi.TypeList{
			abi.Address,
			abi.Address,
			abi.Uint256,
		},
		func(rs abi.RESULT, item *TransactionBlock) {
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

	methodApprove = MethodID(
		"approve",
		abi.TypeList{
			abi.Address,
			abi.Uint256,
		},
		func(rs abi.RESULT, item *TransactionBlock) {
			if !rs.IsError {
				cip := item.NewCustomInputParse()
				cip["to"] = rs.Address(0)
				cip["value"] = rs.Uint256(1)
			}
		},
	)

	methodAllowance = MethodID(
		"allowance",
		abi.TypeList{
			abi.Address,
			abi.Address,
		},
		func(rs abi.RESULT, item *TransactionBlock) {
			if !rs.IsError {
				cip := item.NewCustomInputParse()
				cip["owner"] = rs.Address(0)
				cip["spender"] = rs.Address(1)
			}
		},
	)

	methodIssue = MethodID(
		"issue",
		abi.TypeList{
			abi.Uint256,
		},
		func(rs abi.RESULT, item *TransactionBlock) {
			if !rs.IsError {
				cip := item.NewCustomInputParse()
				cip["value"] = rs.Uint256(0)
			}
		},
	)

	methodMint = MethodID(
		"mint",
		abi.TypeList{
			abi.Uint256,
		},
		func(rs abi.RESULT, item *TransactionBlock) {
			if !rs.IsError {
				cip := item.NewCustomInputParse()
				cip["value"] = rs.Uint256(0)
			}
		},
	)

	methodRedeem = MethodID(
		"redeem",
		abi.TypeList{
			abi.Uint256,
		},
		func(rs abi.RESULT, item *TransactionBlock) {
			if !rs.IsError {
				cip := item.NewCustomInputParse()
				cip["value"] = rs.Uint256(0)
			}
		},
	)

	methodBurn = MethodID(
		"burn",
		abi.TypeList{
			abi.Uint256,
		},
		func(rs abi.RESULT, item *TransactionBlock) {
			if !rs.IsError {
				cip := item.NewCustomInputParse()
				cip["value"] = rs.Uint256(0)
			}
		},
	)

	methodERC20s = MethodIDDataList{
		methodDeploy0,
		methodDeploy1,
		methodDeploy2,
		methodTransferOwnerShip,
		methodTransfer,
		methodTranferFrom,
		methodApprove,
		methodAllowance,
		methodIssue,
		methodMint,
		methodRedeem,
		methodBurn,
	}
) //var

func InputDataPure(stringer abi.Bytes32Stringer, data string, typelist abi.TypeList) abi.RESULT {
	v, e := hex.DecodeString(data)
	if e != nil {
		r := abi.RESULT{}
		r.IsError = true
		return r
	}
	return abi.ReceiptDiv(stringer, v, typelist)
}

func (my cMethodIDData) _parseErc20Input(input_parser abi.InputParser, item *TransactionBlock) bool {
	input := item.CustomInput
	if !strings.HasPrefix(input, my.MethodID) {
		return false
	}
	item.ContractMethod = my.FuncName

	if len(input) > len(my.MethodID) {

		if my.parseInputABM != nil {
			data := input[len(my.MethodID):]
			cdata := InputDataPure(
				input_parser,
				data,
				abi.NewReturns(my.params...),
			)
			my.parseInputABM(cdata, item)

		} else if my.parseDeploy != nil {
			data := input[len(my.MethodID):]
			cdata := InputDataPure(
				input_parser,
				data,
				abi.NewReturns(my.params...),
			)
			my.parseDeploy(input_parser, cdata, item)
		}
	}

	return true
}

func CheckMethodERC20(input_parser abi.InputParser, input string, item *TransactionBlock) {
	limitSize := 8 //method name (without 0x)
	getInputSize := len(input)
	isContractCheck := false
	if input != "" && getInputSize >= limitSize {
		isContractCheck = true
	}

	if isContractCheck == false {
		input = strings.ToLower(input)
		item.CustomInput = input

		item.IsContract = false
		return
	}

	input = strings.TrimSpace(input)
	if input != "" {
		item.IsContract = true
		item.ContractAddress = item.To
		item.To = ""
	}

	if item.IsContract {
		if strings.HasPrefix(input, "0x") == false {
			input = "0x" + input
		}
	}
	input = strings.ToLower(input)
	item.CustomInput = input

	isChecked := false
	for _, erc := range methodERC20s {
		isChecked = erc._parseErc20Input(input_parser, item)
		if isChecked {
			break
		}
	}

	if isChecked == false {
		item.ContractMethod = MethodCustomFunction
	}

}

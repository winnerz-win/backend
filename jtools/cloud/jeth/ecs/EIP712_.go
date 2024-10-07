package ecs

import (
	"jtools/cc"
	"jtools/cloud/ebcm"
	"jtools/dbg"

	"jtools/cloud/jeth/sigverify"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// IsVerifyTypedData : ebcm.MakeEIP712Validator(IsVerifyTypedData)
func IsVerifyTypedData(
	domain ebcm.EIP712DomainData,
	struct_data ebcm.IEIP712Struct,
	signer_address string,
	signature string,
) (bool, error) {

	typedData := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": {
				{
					Name: "name",
					Type: "string",
				},
				{
					Name: "version",
					Type: "string",
				},
				{
					Name: "chainId",
					Type: "uint256",
				},
				{
					Name: "verifyingContract",
					Type: "address",
				},
			},
		},
		PrimaryType: struct_data.PrimaryType(),
		Domain: apitypes.TypedDataDomain{
			Name:              domain.Name,
			Version:           domain.Version,
			ChainId:           math.NewHexOrDecimal256(domain.ChainId),
			VerifyingContract: domain.VerifyingContract,
		},
	}

	design := struct_data.Design()

	typed_struct := []apitypes.Type{}
	for _, v := range design.Types() {
		typed_struct = append(typed_struct,
			apitypes.Type{
				Name: v.Name(),
				Type: v.Type(),
			},
		)
	} //for
	typedData.Types[typedData.PrimaryType] = typed_struct

	message := apitypes.TypedDataMessage{}
	for field_name, value := range design.Values() {
		message[field_name] = value
	} //for

	// cc.Yellow(dbg.ToJsonString(typed_struct))
	// cc.Yellow(dbg.ToJsonString(message))

	typedData.Message = message

	cc.Yellow(dbg.ToJsonString(typedData))

	valid, err := sigverify.VerifyTypedDataHexSignatureEx(
		ethcommon.HexToAddress(signer_address),
		typedData,
		signature,
	)

	return valid, err
}

func EIP712Validator() ebcm.EIP712Validator {
	return IsVerifyTypedData
}

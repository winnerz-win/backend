package ecsx

import (
	"crypto/ecdsa"
	"txscheduler/brix/tools/cloudx/ebcmx"
	ebcmABI "txscheduler/brix/tools/cloudx/ebcmx/abix"
	"txscheduler/brix/tools/cloudx/ethwallet/abmx"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx/jwalletx"
	"txscheduler/brix/tools/dbg"

	"github.com/ethereum/go-ethereum/crypto"
)

func EBCMCaller() ebcmABI.Caller {
	caller := abmx.GetEBCM(
		ebcm_InputDataPure,
	)
	return caller
}

func EBCMSignTool() ebcmx.SignTool {

	keccak256HashBytes := func(data []byte) []byte {
		dataHash := crypto.Keccak256Hash(data)
		buf := dataHash.Bytes()
		return buf
	}

	signTool := ebcmx.SignTool{
		HexToECDSA:         crypto.HexToECDSA,
		FromECDSAPub:       crypto.FromECDSAPub,
		Keccak256HashBytes: keccak256HashBytes,
		GetEthereumMessageHash: func(message []byte) []byte {
			const MESSAGE_PREFIX = "\u0019Ethereum Signed Message:\n"
			getEthereumMessagePrefix := func(messageLength int) []byte {
				prefix := MESSAGE_PREFIX + dbg.Cat(messageLength)
				return []byte(prefix)
			}
			prefix := getEthereumMessagePrefix(len(message))
			result := make([]byte, len(prefix)+len(message))
			size := copy(result, prefix)
			copy(result[size:], message)
			return keccak256HashBytes(result)
		},
		MessageV_addVal: 27,
		MessageV_subVal: -27,
		Sign:            crypto.Sign,
		Ecrecover: func(keccak256Hash, sig []byte) (pub []byte, err error) {
			return crypto.Ecrecover(keccak256Hash, sig)
		},
		SigToPub: func(keccak256Hash, sig []byte) (*ecdsa.PublicKey, error) {
			return crypto.SigToPub(keccak256Hash, sig)
		},
		VerifySignature: func(pubkey, digestHash, signature []byte) bool {
			return crypto.VerifySignature(pubkey, digestHash, signature)
		},
	}
	return signTool
}

func getSignMessagePrefix(prefixMsg string, messageLength int) []byte {
	prefix := prefixMsg + dbg.Cat(messageLength)
	return []byte(prefix)
}

func EBCMSignTooler(message_prefix string) ebcmx.SignTool {

	keccak256HashBytes := func(data []byte) []byte {
		dataHash := crypto.Keccak256Hash(data)
		buf := dataHash.Bytes()
		return buf
	}

	signTool := ebcmx.SignTool{
		HexToECDSA:         crypto.HexToECDSA,
		FromECDSAPub:       crypto.FromECDSAPub,
		Keccak256HashBytes: keccak256HashBytes,
		GetEthereumMessageHash: func(message []byte) []byte {
			prefix := getSignMessagePrefix(message_prefix, len(message))
			result := make([]byte, len(prefix)+len(message))
			size := copy(result, prefix)
			copy(result[size:], message)
			return keccak256HashBytes(result)
		},
		MessageV_addVal: 27,
		MessageV_subVal: -27,
		Sign:            crypto.Sign,
		Ecrecover: func(keccak256Hash, sig []byte) (pub []byte, err error) {
			return crypto.Ecrecover(keccak256Hash, sig)
		},
		SigToPub: func(keccak256Hash, sig []byte) (*ecdsa.PublicKey, error) {
			return crypto.SigToPub(keccak256Hash, sig)
		},
		VerifySignature: func(pubkey, digestHash, signature []byte) bool {
			return crypto.VerifySignature(pubkey, digestHash, signature)
		},
	}
	return signTool
}

func EBCMDataItemListParser() ebcmx.DataItemListParser {
	return ebcmx.NewDataItemListParser(EBCM_DataItemList_ParseABI)
}

func EBCMGasStation() ebcmx.GasStation {
	return ebcmx.NewGasStation(ebcm_NewGasStation)
}

func EBCMSener(iSender interface{}) *ebcmx.Sender {
	if iSender == nil {
		return nil
	}
	sender, do := iSender.(*Sender)
	if !do || sender == nil {
		return nil
	}

	signTool := EBCMSignTool()
	ins := ebcmx.Sender{
		ISender:    iSender,
		Caller:     EBCMCaller(),
		SignTool:   &signTool,
		SignTooler: EBCMSignTooler,

		Infomation: sender.String,
		Mainnet:    sender.Mainnet,
		HostURL:    sender.HostURL,
		Key:        sender.Key,

		Balance:        sender.Balance,
		CoinPrice:      sender.CoinPrice,
		TokenBalance:   sender.TokenBalance,
		TokenPrice:     sender.TokenPrice,
		ChainID:        sender.ChainID,
		BlockNumberTry: sender.BlockNumberTry,

		// BlockByNumber:     s.ebcm_BlockByNumber,
		// TransactionByHash: s.ebcm_TransactionByHash,
		ReceiptByHash: sender.ebcm_ReceiptByHash,
		InjectReceipt: sender.ebcm_InjectReceipt,

		MakePadBytesABI: ebcm_MakePadBytesABI,
		MakePadBytes:    ebcm_MakePadBytes,
		MakePadBytes2:   ebcm_MakePadBytes2,

		Token: sender.ebcm_Token,

		PadBytesETH:          ebcm_PadBytesETH,
		TransferPadBytes:     ebcm_TransferPadBytes,
		PadBytesTransfer:     ebcm_PadBytesTransfer,
		PadBytesApprove:      ebcm_PadBytesApprove,
		PadBytesApproveAll:   ebcm_PadBytesApproveAll,
		PadBytesTransferFrom: ebcm_PadBytesTransferFrom,

		XNonce:   sender.XPendingNonceAt,
		XNonceAt: sender.XNonceAt,

		XEstimateFeeETH: sender.ebcm_XEstimateFeeETH,
		XGasLimit:       sender.ebcm_XGasLimit,

		XPipe:                sender.ebcm_XPipe,
		XPipeFixedGAS:        sender.ebcm_XPipeFixedGAS,
		TransferCoin:         sender.ebcm_TransferCoin,
		TransferCoinFixedGAS: sender.ebcm_TransferCoinFixedGAS,

		DelegateCallContract: sender.CallContract,

		IsAddress:            IsAddress,
		ContractAddressNonce: ContractAddressNonce,

		NewSeedI:  jwalletx.EBCM_NewSeedI,
		GetWallet: jwalletx.EBCM_Get,

		DataItemList_ParseABI: EBCM_DataItemList_ParseABI,
		//GetGasResult:          ebcm_NewGasStation,
		GetGasResult: func() ebcmx.GasResult {
			if sender.mainnet {
				return ebcm_NewGasStation()
			}
			return GasResult{val: sender.SuggestGasPrice()}
		},
	}

	ins.SetToggleHost(sender.ToggleHost)

	ins.SetFuncPtr(
		sender.ebcm_BlockByNumber,
		sender.ebcm_TransactionByHash,
		sender.ebcm_TSender,
	)
	return &ins
}

func (my *Sender) ebcm_TSender(contractAddress string) ebcmx.TSender {
	ts := my.TSender(contractAddress)
	ins := ebcmx.TSender{
		ContractAddress: ts.ContractAddress,

		Allowance:    ts.Allowance,
		Approve:      ts.Approve,
		ApproveAll:   ts.ApproveAll,
		TransferFrom: ts.TransferFrom,

		TransferFunction:         ts.ebcm_TransferFunction,
		TransferFunctionFixedGAS: ts.ebcm_TransferFunctionFixedGAS,

		TransferToken:         ts.ebcm_TransferToken,
		TransferTokenFixedGAS: ts.ebcm_TransferTokenFixedGAS,

		Write:         ts.ebcm_Write,
		WriteFixedGAS: ts.ebcm_WriteFixedGAS,
	}
	return ins
}

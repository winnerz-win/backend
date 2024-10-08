package nftc

import (
	"context"
	"jtools/cloud/ebcm"
	"jtools/cloud/ebcm/abi"
	"jtools/cloud/jeth/jwallet"
	"txscheduler/brix/tools/dbg"
	"txscheduler/txm/model"
)

type NFT struct{}

func _nftCall(method abi.Method, caller string, result func(abi.RESULT)) error {
	return Finder().Call(
		nftToken.Contract,
		method,
		caller,
		result,
		debugMode(),
	)
}

func (NFT) StartNumber(f func(string)) error {
	return _nftCall(
		abi.Method{
			Name:   "startNumber",
			Params: abi.NewParams(),
			Returns: abi.NewReturns(
				abi.Uint256,
			),
		},
		nftToken.Address,
		func(rs abi.RESULT) {
			f(rs.Uint256(0))
		},
	)
}

func (NFT) GetBaseURI(f func(string)) error {
	return _nftCall(
		abi.Method{
			Name:   "getBaseURI",
			Params: abi.NewParams(),
			Returns: abi.NewReturns(
				abi.String,
			),
		},
		nftToken.Address,
		func(rs abi.RESULT) {
			f(rs.Text(0))
		},
	)
}

func (NFT) Name(f func(string)) error {
	return _nftCall(
		abi.Method{
			Name:   "name",
			Params: abi.NewParams(),
			Returns: abi.NewReturns(
				abi.String,
			),
		},
		nftToken.Address,
		func(rs abi.RESULT) {
			f(rs.Text(0))
		},
	)
}

func (NFT) Symbol(f func(string)) error {
	return _nftCall(
		abi.Method{
			Name:   "symbol",
			Params: abi.NewParams(),
			Returns: abi.NewReturns(
				abi.String,
			),
		},
		nftToken.Address,
		func(rs abi.RESULT) {
			f(rs.Text(0))
		},
	)
}

func (NFT) TotalSupply(f func(string)) error {
	return _nftCall(
		abi.Method{
			Name:   "totalSupply",
			Params: abi.NewParams(),
			Returns: abi.NewReturns(
				abi.Uint256,
			),
		},
		nftToken.Address,
		func(rs abi.RESULT) {
			f(rs.Uint256(0))
		},
	)
}

func (NFT) BalanceOf(owner string, f func(string)) error {
	return _nftCall(
		abi.Method{
			Name: "balanceOf",
			Params: abi.NewParams(
				abi.NewAddress(owner),
			),
			Returns: abi.NewReturns(
				abi.Uint,
			),
		},
		owner,
		func(rs abi.RESULT) {
			f(rs.Uint(0))
		},
	)
}

func (NFT) DEFAULT_ADMIN_ROLE(f func(string)) error {
	return _nftCall(
		abi.Method{
			Name:   "DEFAULT_ADMIN_ROLE",
			Params: abi.NewParams(),
			Returns: abi.NewReturns(
				abi.Bytes32,
			),
		},
		nftToken.Address,
		func(rs abi.RESULT) {
			f(rs.Bytes32(0))
		},
	)
}

func (NFT) GetApproved(tokenId string, f func(string)) error {
	return _nftCall(
		abi.Method{
			Name: "getApproved",
			Params: abi.NewParams(
				abi.NewUint256(tokenId),
			),
			Returns: abi.NewReturns(
				abi.Address,
			),
		},
		nftToken.Address,
		func(rs abi.RESULT) {
			f(rs.Address(0))
		},
	)
}

func (NFT) GetRoleAdmin(role string, f func(string)) error {
	return _nftCall(
		abi.Method{
			Name: "getRoleAdmin",
			Params: abi.NewParams(
				abi.NewBytes32(role),
			),
			Returns: abi.NewReturns(
				abi.Bytes32,
			),
		},
		nftToken.Address,
		func(rs abi.RESULT) {
			f(rs.Bytes32(0))
		},
	)
}

func (NFT) TokenURI(tokenId string, f func(string)) error {
	return _nftCall(
		abi.Method{
			Name: "tokenURI",
			Params: abi.NewParams(
				abi.NewUint256(tokenId),
			),
			Returns: abi.NewReturns(
				abi.String,
			),
		},
		nftToken.Address,
		func(rs abi.RESULT) {
			f(rs.Text(0))
		},
	)
}

func (NFT) TokenMetaData(tokenId string, f func(model.NftMetaData)) error {
	return _nftCall(
		abi.Method{
			Name: "tokenMetaData",
			Params: abi.NewParams(
				abi.NewUint256(tokenId),
			),
			Returns: abi.NewReturns(
				abi.Uint256,
				abi.String,
				abi.String,
				abi.String,
				abi.String,
			),
		},
		nftToken.Address,
		func(rs abi.RESULT) {
			_tokenType := rs.Uint256(0)
			_tokenName := rs.Text(1)
			_tokenContent := rs.Text(2)
			_tokenDeployer := rs.Text(3)
			_tokenDesc := rs.Text(4)
			meta := model.NftMetaData{
				TokenType: _tokenType,
				Name:      _tokenName,
				Content:   _tokenContent,
				Deployer:  _tokenDeployer,
				Desc:      _tokenDesc,
			}
			f(meta)
		},
	)
}

func (my NFT) Mint(tokenId, to, tokenType string) WriteResult {
	return my.Mint2(tokenId, to, tokenType, nil)

}

func (NFT) Mint2(tokenId, to, tokenType string, try *model.NftBuyTry) WriteResult {

	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "mint",
	}

	data := ebcm.MakePadBytesABI(
		"mint",
		abi.TypeList{
			abi.NewUint256(tokenId),
			abi.NewAddress(to),
			abi.NewUint256(tokenType),
		},
	)

	snap := ebcm.GasSnapShot{}
	from, hash, nonce, err := TransferFuncNTX_Send(
		Sender(),
		nftToken.Contract,

		nftToken.Private,
		data,
		"0",
		&snap,
	)
	if err != nil {
		dbg.Red(err)
	}

	if try != nil {
		try.GasLimit = snap.Limit
		try.GasPrice = snap.Price
		try.GasFeeETH = ebcm.WeiToETH(snap.FeeWei)

		dbg.Purple("gasPrice :", try.GasPrice)
		dbg.Purple("gasLimit :", try.GasLimit)
		dbg.Purple("feeETH :", try.GasFeeETH)
	}

	r.Set(from, hash, nonce, err)
	return r

}

func (NFT) MintMulti(tokenId []string, to, tokenType string) WriteResult {

	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "mintMulti",
	}

	data := ebcm.MakePadBytesABI(
		"mintMulti",
		abi.TypeList{
			abi.NewUint256Array(tokenId),
			abi.NewAddress(to),
			abi.NewUint256(tokenType),
		},
	)

	snap := ebcm.GasSnapShot{}
	from, hash, nonce, err := TransferFuncNTX_Send(
		Sender(),
		nftToken.Contract,

		nftToken.Private,
		data,
		"0",
		&snap,
	)
	if err != nil {
		dbg.Red(err)
	}

	r.Set(from, hash, nonce, err)
	return r
}

func (NFT) MintInfo(tokenId, to, _tokenType, _tokenName, _tokenDesc, _tokenContent, deployer string) WriteResult {

	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "mintInfo",
	}

	data := ebcm.MakePadBytesABI(
		"mintInfo",
		abi.TypeList{
			abi.NewUint256(tokenId),
			abi.NewAddress(to),
			abi.NewUint256(_tokenType),
			abi.NewString(_tokenName),
			abi.NewString(_tokenDesc),
			abi.NewString(_tokenContent),
			abi.NewString(deployer),
		},
	)

	snap := ebcm.GasSnapShot{}
	from, hash, nonce, err := TransferFuncNTX_Send(
		Sender(),
		nftToken.Contract,

		nftToken.Private,
		data,
		"0",
		&snap,
	)
	if err != nil {
		dbg.Red(err)
	}

	r.Set(from, hash, nonce, err)
	return r
}

func (NFT) MintInfoMulti(tokenId []string, to, _tokenType, _tokenName, _tokenDesc, _tokenContent, deployer string) WriteResult {

	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "mintInfoMulti",
	}

	// tc := Sender().TSender(nftToken.Contract)

	// data := ecsx.MakePadBytes(
	// 	"mintInfoMulti(uint256,address,uint256,string,string,string,string)",
	// 	func(pad ecsx.Appender) {
	// 		pos := pad.SetParamHeader(3)
	// 		pos = pad.SetAmountArray(0, pos, tokenId...)
	// 		pad.SetAddress(1, to)
	// 		pad.SetAmount(2, _tokenType)
	// 		pos = pad.SetText(3, pos, _tokenName)
	// 		pos = pad.SetText(4, pos, _tokenDesc)
	// 		pos = pad.SetText(5, pos, _tokenContent)
	// 		pos = pad.SetText(6, pos, deployer)
	// 	},
	// )

	data := ebcm.MakePadBytesABI(
		"mintInfoMulti",
		abi.TypeList{
			abi.NewUint256Array(tokenId),
			abi.NewAddress(to),
			abi.NewUint256(_tokenType),
			abi.NewString(_tokenName),
			abi.NewString(_tokenDesc),
			abi.NewString(_tokenContent),
			abi.NewString(deployer),
		},
	)

	snap := ebcm.GasSnapShot{}
	from, hash, nonce, err := TransferFuncNTX_Send(
		Sender(),
		nftToken.Contract,

		nftToken.Private,
		data,
		"0",
		&snap,
	)
	if err != nil {
		dbg.Red(err)
	}

	r.Set(from, hash, nonce, err)
	return r
}

func (NFT) Burn(privateKey string, tokenId string) WriteResult {
	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "burn",
	}

	// tc := Sender().TSender(nftToken.Contract)
	// data := ecsx.MakePadBytes(
	// 	"burn(uint256)",
	// 	func(pad ecsx.Appender) {
	// 		pad.SetParamHeader(1)
	// 		pad.SetAmount(0, tokenId)
	// 	},
	// )
	data := ebcm.MakePadBytesABI(
		"burn",
		abi.TypeList{
			abi.NewUint256(tokenId),
		},
	)

	snap := ebcm.GasSnapShot{}
	from, hash, nonce, err := TransferFuncNTX_Send(
		Sender(),
		nftToken.Contract,

		nftToken.Private,
		data,
		"0",
		&snap,
	)
	if err != nil {
		dbg.Red(err)
	}

	r.Set(from, hash, nonce, err)
	return r
}

func (NFT) SetBaseURL_Admin(newURI string) WriteResult {
	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "setBaseURI",
	}

	// tc := Sender().TSender(nftToken.Contract)
	// data := ecsx.MakePadBytes(
	// 	"setBaseURI(string)",
	// 	func(pad ecsx.Appender) {
	// 		pos := pad.SetParamHeader(1)
	// 		pad.SetText(0, pos, newURI)
	// 	},
	// )

	data := ebcm.MakePadBytesABI(
		"setBaseURI",
		abi.TypeList{
			abi.NewString(newURI),
		},
	)

	snap := ebcm.GasSnapShot{}
	from, hash, nonce, err := TransferFuncNTX_Send(
		Sender(),
		nftToken.Contract,

		nftToken.Private,
		data,
		"0",
		&snap,
	)
	if err != nil {
		dbg.Red(err)
	}

	r.Set(from, hash, nonce, err)
	return r
}

func (NFT) SetTokenURI_Admin(tokenId, tokenURI string) WriteResult {
	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "setTokenURI",
	}

	// tc := Sender().TSender(nftToken.Contract)
	// data := ecsx.MakePadBytes(
	// 	"setTokenURI(uint256,string)",
	// 	func(pad ecsx.Appender) {
	// 		pos := pad.SetParamHeader(2)
	// 		pad.SetAmount(0, tokenId)
	// 		pad.SetText(1, pos, tokenURI)
	// 	},
	// )

	data := ebcm.MakePadBytesABI(
		"setTokenURI",
		abi.TypeList{
			abi.NewUint256(tokenId),
			abi.NewString(tokenURI),
		},
	)

	snap := ebcm.GasSnapShot{}
	from, hash, nonce, err := TransferFuncNTX_Send(
		Sender(),
		nftToken.Contract,

		nftToken.Private,
		data,
		"0",
		&snap,
	)
	if err != nil {
		dbg.Red(err)
	}

	r.Set(from, hash, nonce, err)
	return r
}
func (NFT) SetTypeName_Admin(nType, nName string) WriteResult {
	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "setTypeName",
	}

	// tc := Sender().TSender(nftToken.Contract)
	// data := ecsx.MakePadBytes(
	// 	"setTypeName(uint256,string)",
	// 	func(pad ecsx.Appender) {
	// 		pos := pad.SetParamHeader(2)
	// 		pad.SetAmount(0, nType)
	// 		pad.SetText(1, pos, nName)
	// 	},
	// )

	data := ebcm.MakePadBytesABI(
		"setTypeName",
		abi.TypeList{
			abi.NewUint256(nType),
			abi.NewString(nName),
		},
	)

	snap := ebcm.GasSnapShot{}
	from, hash, nonce, err := TransferFuncNTX_Send(
		Sender(),
		nftToken.Contract,

		nftToken.Private,
		data,
		"0",
		&snap,
	)
	if err != nil {
		dbg.Red(err)
	}

	r.Set(from, hash, nonce, err)
	return r
}
func (NFT) SetTokenInfo_Admin(tokenId, tokenName, tokenDesc, tokenContent string) WriteResult {
	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "setTokenInfo",
	}

	// tc := Sender().TSender(nftToken.Contract)
	// data := ecsx.MakePadBytes(
	// 	"setTokenInfo(uint256,string,string,string)",
	// 	func(pad ecsx.Appender) {
	// 		pos := pad.SetParamHeader(2)
	// 		pad.SetAmount(0, tokenId)
	// 		pos = pad.SetText(1, pos, tokenName)
	// 		pos = pad.SetText(2, pos, tokenDesc)
	// 		pos = pad.SetText(3, pos, tokenContent)
	// 	},
	// )

	data := ebcm.MakePadBytesABI(
		"setTypeName",
		abi.TypeList{
			abi.NewUint256(tokenId),
			abi.NewString(tokenName),
			abi.NewString(tokenDesc),
			abi.NewString(tokenContent),
		},
	)

	snap := ebcm.GasSnapShot{}
	from, hash, nonce, err := TransferFuncNTX_Send(
		Sender(),
		nftToken.Contract,

		nftToken.Private,
		data,
		"0",
		&snap,
	)
	if err != nil {
		dbg.Red(err)
	}

	r.Set(from, hash, nonce, err)
	return r
}

/*
	interface IAccessControl {
		function hasRole(bytes32 role, address account) external view returns (bool); //권한 여부
		function getRoleAdmin(bytes32 role) external view returns (bytes32);
		function grantRole(bytes32 role, address account) external;		//권한 부여
		function revokeRole(bytes32 role, address account) external;	//권한 회수
		function renounceRole(bytes32 role, address account) external;	//권한 포기
	}
*/

func (NFT) HasRole(role, address string, f func(bool)) error {
	return _nftCall(
		abi.Method{
			Name: "hasRole",
			Params: abi.NewParams(
				abi.NewBytes32(role),
				abi.NewAddress(address),
			),
			Returns: abi.NewReturns(
				abi.Bool,
			),
		},
		nftToken.Address,
		func(rs abi.RESULT) {
			f(rs.Bool(0))
		},
	)
}

func (NFT) GrantRole(role, address string) WriteResult {
	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "grantRole",
	}

	// tc := Sender().TSender(nftToken.Contract)

	// data := ecsx.MakePadBytes(
	// 	"grantRole(bytes32, address)",
	// 	func(pad ecsx.Appender) {
	// 		pad.SetParamHeader(2)
	// 		pad.SetBytes32Hex(0, role)
	// 		pad.SetAddress(0, address)
	// 	},
	// )

	data := ebcm.MakePadBytesABI(
		"grantRole",
		abi.TypeList{
			abi.NewBytes32(role),
			abi.NewAddress(address),
		},
	)

	snap := ebcm.GasSnapShot{}
	from, hash, nonce, err := TransferFuncNTX_Send(
		Sender(),
		nftToken.Contract,

		nftToken.Private,
		data,
		"0",
		&snap,
	)
	if err != nil {
		dbg.Red(err)
	}

	r.Set(from, hash, nonce, err)
	return r
}

func (NFT) RevokeRole(role, address string) WriteResult {
	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "revokeRole",
	}

	// tc := Sender().TSender(nftToken.Contract)

	// data := ecsx.MakePadBytes(
	// 	"revokeRole(bytes32, address)",
	// 	func(pad ecsx.Appender) {
	// 		pad.SetParamHeader(2)
	// 		pad.SetBytes32Hex(0, role)
	// 		pad.SetAddress(0, address)
	// 	},
	// )

	data := ebcm.MakePadBytesABI(
		"revokeRole",
		abi.TypeList{
			abi.NewBytes32(role),
			abi.NewAddress(address),
		},
	)

	snap := ebcm.GasSnapShot{}
	from, hash, nonce, err := TransferFuncNTX_Send(
		Sender(),
		nftToken.Contract,

		nftToken.Private,
		data,
		"0",
		&snap,
	)
	if err != nil {
		dbg.Red(err)
	}

	r.Set(from, hash, nonce, err)
	return r
}

func (NFT) RenounceRole(privateKey string, role, address string) WriteResult {
	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "renounceRole",
	}

	// tc := Sender().TSender(nftToken.Contract)

	// data := ecsx.MakePadBytes(
	// 	"renounceRole(bytes32, address)",
	// 	func(pad ecsx.Appender) {
	// 		pad.SetParamHeader(2)
	// 		pad.SetBytes32Hex(0, role)
	// 		pad.SetAddress(0, address)
	// 	},
	// )

	data := ebcm.MakePadBytesABI(
		"renounceRole",
		abi.TypeList{
			abi.NewBytes32(role),
			abi.NewAddress(address),
		},
	)

	snap := ebcm.GasSnapShot{}
	from, hash, nonce, err := TransferFuncNTX_Send(
		Sender(),
		nftToken.Contract,

		privateKey,
		data,
		"0",
		&snap,
	)
	if err != nil {
		dbg.Red(err)
	}

	r.Set(from, hash, nonce, err)
	return r
}

///////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////

type TRANSFER_TOKEN struct {
	nonce       uint64
	gas_fee_eth string
	stx         ebcm.WrappedTransaction

	/*
		type GasSnapShot struct {
			Limit      uint64 `bson:"limit" json:"limit"`
			Price      string `bson:"price" json:"price"`
			FeeWei     string `bson:"fee_wei" json:"fee_wei"` // limit * price
			FixedNonce uint64 `bson:"fixedNonce,omitempty" json:"fixedNonce,omitempty"`
		}
	*/
	snap ebcm.GasSnapShot
}

func (my TRANSFER_TOKEN) GasFeeETH() string          { return my.gas_fee_eth }
func (my TRANSFER_TOKEN) SnapShot() ebcm.GasSnapShot { return my.snap }

func (NFT) TransferFromNTX(privateKey string, to, tokenId string) (*TRANSFER_TOKEN, error) {
	from, err := jwallet.Get(privateKey)
	if err != nil {
		return nil, err
	}
	from_address := from.Address()

	sender := Sender()
	nonce, err := ebcm.MMA_GetNonce(sender, from_address, true)
	if err != nil {
		return nil, err
	}

	gas_price, err := sender.SuggestGasPrice(context.Background(), true)
	if err != nil {
		return nil, err
	}

	// tc := Sender().TSender(nftToken.Contract)
	// data := ecsx.MakePadBytes(
	// 	"transferFrom(address,address,uint256)",
	// 	func(pad ecsx.Appender) {
	// 		pad.SetParamHeader(3)
	// 		pad.SetAddress(0, w.Address()) //from
	// 		pad.SetAddress(1, to)          //to
	// 		pad.SetAmount(2, tokenId)      //tokenId
	// 	},
	// )
	// ntx, err := tc.TransferFuncNTX(
	// 	privateKey,
	// 	data,
	// 	"0",
	// 	gasSpeed,
	// 	nil,
	// )

	data := ebcm.MakePadBytesABI(
		"transferFrom",
		abi.TypeList{
			abi.NewAddress(from.Address()),
			abi.NewAddress(to),
			abi.NewUint256(tokenId),
		},
	)

	limit, err := sender.EstimateGas(
		context.Background(),
		ebcm.MakeCallMsg(
			from_address,
			to,
			"0",
			data,
		),
	)
	if err != nil {
		return nil, err
	}

	limit = ebcm.MMA_LimitBuffer(limit)

	ntx := sender.NewTransaction(
		nonce,
		nftToken.Contract,
		"0",
		limit,
		gas_price,
		data,
	)
	stx, err := sender.SignTx(ntx, from.PrivateKey())
	if err != nil {
		return nil, err
	}

	// return ntx, err
	tt := &TRANSFER_TOKEN{
		nonce: nonce,

		gas_fee_eth: gas_price.EstimateGasFeeETH(limit),
		stx:         stx,

		snap: ebcm.MakeGasSnapShot(
			nonce,
			limit,
			gas_price,
		),
	}

	return tt, nil
}

func TransferEtherNTX(privateKey string, to, wei string) (*TRANSFER_TOKEN, error) {
	from, err := jwallet.Get(privateKey)
	if err != nil {
		return nil, err
	}
	from_address := from.Address()

	sender := Sender()
	nonce, err := ebcm.MMA_GetNonce(sender, from_address, true)
	if err != nil {
		return nil, err
	}

	gas_price, err := sender.SuggestGasPrice(context.Background(), true)
	if err != nil {
		return nil, err
	}

	data := ebcm.PadByteETH()

	limit, err := sender.EstimateGas(
		context.Background(),
		ebcm.MakeCallMsg(
			from_address,
			to,
			wei,
			data,
		),
	)
	if err != nil {
		return nil, err
	}

	//limit = ebcm.MMA_LimitBuffer(limit)

	ntx := sender.NewTransaction(
		nonce,
		to,
		wei,
		limit,
		gas_price,
		data,
	)

	stx, err := sender.SignTx(ntx, from.PrivateKey())
	if err != nil {
		return nil, err
	}

	tt := &TRANSFER_TOKEN{
		nonce: nonce,

		gas_fee_eth: gas_price.EstimateGasFeeETH(limit),
		stx:         stx,

		snap: ebcm.MakeGasSnapShot(
			nonce,
			limit,
			gas_price,
		),
	}

	return tt, nil
}

func TransferTokenNTX(contract string, privateKey string, to string, token_wei string) (*TRANSFER_TOKEN, error) {

	from, err := jwallet.Get(privateKey)
	if err != nil {
		return nil, err
	}
	from_address := from.Address()

	sender := Sender()
	nonce, err := ebcm.MMA_GetNonce(sender, from_address, true)
	if err != nil {
		return nil, err
	}

	gas_price, err := sender.SuggestGasPrice(context.Background(), true)
	if err != nil {
		return nil, err
	}

	data := ebcm.PadByteTransfer(
		to,
		token_wei,
	)

	limit, err := sender.EstimateGas(
		context.Background(),
		ebcm.MakeCallMsg(
			from_address,
			to,
			"0",
			data,
		),
	)
	if err != nil {
		return nil, err
	}

	limit = ebcm.MMA_LimitBuffer(limit)

	ntx := sender.NewTransaction(
		nonce,
		contract,
		"0",
		limit,
		gas_price,
		data,
	)

	stx, err := sender.SignTx(ntx, from.PrivateKey())
	if err != nil {
		return nil, err
	}

	tt := &TRANSFER_TOKEN{
		nonce: nonce,

		gas_fee_eth: gas_price.EstimateGasFeeETH(limit),
		stx:         stx,

		snap: ebcm.MakeGasSnapShot(
			nonce,
			limit,
			gas_price,
		),
	}

	return tt, nil
}

func TransferNTX_Send(tt *TRANSFER_TOKEN) (string, error) {
	hash, err := Sender().SendTransaction(
		context.Background(),
		tt.stx,
	)
	return hash, err
}

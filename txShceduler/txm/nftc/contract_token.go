package nftc

import (
	"txscheduler/brix/tools/cloudx/ethwallet/abmx"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx/jwalletx"
	"txscheduler/brix/tools/dbg"
	"txscheduler/txm/model"
)

type NFT struct{}

func _nftCall(method abmx.Method, caller string, result func(abmx.RESULT)) error {
	return abmx.Call(
		Sender(),
		nftToken.Contract,
		method,
		caller,
		result,
		debugMode(),
	)
}

func (NFT) StartNumber(f func(string)) error {
	return _nftCall(
		abmx.Method{
			Name:   "startNumber",
			Params: abmx.NewParams(),
			Returns: abmx.NewReturns(
				abmx.Uint256,
			),
		},
		nftToken.Address,
		func(rs abmx.RESULT) {
			f(rs.Uint256(0))
		},
	)
}

func (NFT) GetBaseURI(f func(string)) error {
	return _nftCall(
		abmx.Method{
			Name:   "getBaseURI",
			Params: abmx.NewParams(),
			Returns: abmx.NewReturns(
				abmx.String,
			),
		},
		nftToken.Address,
		func(rs abmx.RESULT) {
			f(rs.Text(0))
		},
	)
}

func (NFT) Name(f func(string)) error {
	return _nftCall(
		abmx.Method{
			Name:   "name",
			Params: abmx.NewParams(),
			Returns: abmx.NewReturns(
				abmx.String,
			),
		},
		nftToken.Address,
		func(rs abmx.RESULT) {
			f(rs.Text(0))
		},
	)
}

func (NFT) Symbol(f func(string)) error {
	return _nftCall(
		abmx.Method{
			Name:   "symbol",
			Params: abmx.NewParams(),
			Returns: abmx.NewReturns(
				abmx.String,
			),
		},
		nftToken.Address,
		func(rs abmx.RESULT) {
			f(rs.Text(0))
		},
	)
}

func (NFT) TotalSupply(f func(string)) error {
	return _nftCall(
		abmx.Method{
			Name:   "totalSupply",
			Params: abmx.NewParams(),
			Returns: abmx.NewReturns(
				abmx.Uint256,
			),
		},
		nftToken.Address,
		func(rs abmx.RESULT) {
			f(rs.Uint256(0))
		},
	)
}

func (NFT) BalanceOf(owner string, f func(string)) error {
	return _nftCall(
		abmx.Method{
			Name: "balanceOf",
			Params: abmx.NewParams(
				abmx.NewAddress(owner),
			),
			Returns: abmx.NewReturns(
				abmx.Uint,
			),
		},
		owner,
		func(rs abmx.RESULT) {
			f(rs.Uint(0))
		},
	)
}

func (NFT) DEFAULT_ADMIN_ROLE(f func(string)) error {
	return _nftCall(
		abmx.Method{
			Name:   "DEFAULT_ADMIN_ROLE",
			Params: abmx.NewParams(),
			Returns: abmx.NewReturns(
				abmx.Bytes32,
			),
		},
		nftToken.Address,
		func(rs abmx.RESULT) {
			f(rs.Bytes32(0))
		},
	)
}

func (NFT) GetApproved(tokenId string, f func(string)) error {
	return _nftCall(
		abmx.Method{
			Name: "getApproved",
			Params: abmx.NewParams(
				abmx.NewUint256(tokenId),
			),
			Returns: abmx.NewReturns(
				abmx.Address,
			),
		},
		nftToken.Address,
		func(rs abmx.RESULT) {
			f(rs.Address(0))
		},
	)
}

func (NFT) GetRoleAdmin(role string, f func(string)) error {
	return _nftCall(
		abmx.Method{
			Name: "getRoleAdmin",
			Params: abmx.NewParams(
				abmx.NewBytes32(role),
			),
			Returns: abmx.NewReturns(
				abmx.Bytes32,
			),
		},
		nftToken.Address,
		func(rs abmx.RESULT) {
			f(rs.Bytes32(0))
		},
	)
}

func (NFT) TokenURI(tokenId string, f func(string)) error {
	return _nftCall(
		abmx.Method{
			Name: "tokenURI",
			Params: abmx.NewParams(
				abmx.NewUint256(tokenId),
			),
			Returns: abmx.NewReturns(
				abmx.String,
			),
		},
		nftToken.Address,
		func(rs abmx.RESULT) {
			f(rs.Text(0))
		},
	)
}

func (NFT) TokenMetaData(tokenId string, f func(model.NftMetaData)) error {
	return _nftCall(
		abmx.Method{
			Name: "tokenMetaData",
			Params: abmx.NewParams(
				abmx.NewUint256(tokenId),
			),
			Returns: abmx.NewReturns(
				abmx.Uint256,
				abmx.String,
				abmx.String,
				abmx.String,
				abmx.String,
			),
		},
		nftToken.Address,
		func(rs abmx.RESULT) {
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

func (NFT) Mint(tokenId, to, tokenType string) WriteResult {

	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "mint",
	}

	tc := Sender().TSender(nftToken.Contract)

	data := ecsx.MakePadBytes(
		"mint(uint256,address,uint256)",
		func(pad ecsx.Appender) {
			pad.SetParamHeader(3)
			pad.SetAmount(0, tokenId)
			pad.SetAddress(1, to)
			pad.SetAmount(2, tokenType)
		},
	)

	ntx, err := tc.TransferFuncNTX(
		nftToken.Private,
		data,
		"0",
		ecsx.GasFast,
		nil,
	)
	if err != nil {
		dbg.Red(err)
		r.Set(ntx.From(), "", 0, err)
		return r
	}

	h, n, e := tc.TransferFuncSEND(ntx)
	r.Set(ntx.From(), h, n, e)
	return r
}

func (NFT) Mint2(tokenId, to, tokenType string, try *model.NftBuyTry) WriteResult {

	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "mint",
	}

	data := ecsx.MakePadBytes(
		"mint(uint256,address,uint256)",
		func(pad ecsx.Appender) {
			pad.SetParamHeader(3)
			pad.SetAmount(0, tokenId)
			pad.SetAddress(1, to)
			pad.SetAmount(2, tokenType)
		},
	)

	sender := Sender()

	nonce, err := sender.XPendingNonceAt(nftToken.Address)
	if err != nil {
		dbg.Red(err)
		r.Set(nftToken.Address, "", nonce, err)
		return r
	}

	gasLimit, err := sender.XGasLimit(
		data,
		nftToken.Address,
		nftToken.Contract,
		"0",
	)
	if err != nil {
		dbg.Red(err)
		r.Set(nftToken.Address, "", nonce, err)
		return r
	}
	gasPrice := sender.SUGGEST_GAS_PRICE(ecsx.GasFast)
	if gasPrice.Error() != nil {
		dbg.Red(err)
		r.Set(nftToken.Address, "", nonce, gasPrice.Error())
		return r
	}

	ntx := sender.XNTX(
		data,
		nftToken.Contract,
		"0",
		nonce,
		gasLimit,
		gasPrice,
	)
	if ntx.Error() != nil {
		dbg.Red(err)
		r.Set(nftToken.Address, "", nonce, ntx.Error())
		return r
	}

	stx := sender.XSTX(
		nftToken.Private,
		ntx,
	)
	if stx.Error() != nil {
		dbg.Red(err)
		r.Set(nftToken.Address, "", nonce, stx.Error())
		return r
	}

	if err := sender.XSend(stx); err != nil {
		dbg.Red(err)
		r.Set(nftToken.Address, "", nonce, err)
		return r
	}
	r.Set(
		nftToken.Address,
		stx.Hash(),
		nonce,
		nil,
	)

	try.GasLimit = gasLimit
	try.GasPrice = gasPrice.ETH()
	try.GasFeeETH = gasPrice.FeeETH(gasLimit)

	dbg.Purple("gasPrice :", gasPrice.WEI())
	dbg.Purple("gasLimit :", gasLimit)
	dbg.Purple("feeETH :", gasPrice.FeeETH(gasLimit))

	return r
}

func (NFT) MintMulti(tokenId []string, to, tokenType string) WriteResult {

	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "mintMulti",
	}

	tc := Sender().TSender(nftToken.Contract)

	data := ecsx.MakePadBytes(
		"mintMulti(uint256[],address,uint256)",
		func(pad ecsx.Appender) {
			pos := pad.SetParamHeader(3)
			pad.SetAmountArray(0, pos, tokenId...)
			pad.SetAddress(1, to)
			pad.SetAmount(2, tokenType)
		},
	)

	ntx, err := tc.TransferFuncNTX(
		nftToken.Private,
		data,
		"0",
		ecsx.GasFast,
		nil,
	)
	if err != nil {
		dbg.Red(err)
		r.Set(ntx.From(), "", 0, err)
		return r
	}

	h, n, e := tc.TransferFuncSEND(ntx)
	r.Set(ntx.From(), h, n, e)
	return r
}

func (NFT) MintInfo(tokenId, to, _tokenType, _tokenName, _tokenDesc, _tokenContent, deployer string) WriteResult {

	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "mintInfo",
	}

	tc := Sender().TSender(nftToken.Contract)

	data := ecsx.MakePadBytes(
		"mintInfo(uint256,address,uint256,string,string,string,string)",
		func(pad ecsx.Appender) {
			pos := pad.SetParamHeader(3)
			pad.SetAmount(0, tokenId)
			pad.SetAddress(1, to)
			pad.SetAmount(2, _tokenType)
			pos = pad.SetText(3, pos, _tokenName)
			pos = pad.SetText(4, pos, _tokenDesc)
			pos = pad.SetText(5, pos, _tokenContent)
			pos = pad.SetText(6, pos, deployer)
		},
	)

	ntx, err := tc.TransferFuncNTX(
		nftToken.Private,
		data,
		"0",
		ecsx.GasFast,
		nil,
	)
	if err != nil {
		dbg.Red(err)
		r.Set(ntx.From(), "", 0, err)
		return r
	}

	h, n, e := tc.TransferFuncSEND(ntx)
	r.Set(ntx.From(), h, n, e)
	return r
}

func (NFT) MintInfoMulti(tokenId []string, to, _tokenType, _tokenName, _tokenDesc, _tokenContent, deployer string) WriteResult {

	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "mintInfoMulti",
	}

	tc := Sender().TSender(nftToken.Contract)

	data := ecsx.MakePadBytes(
		"mintInfoMulti(uint256,address,uint256,string,string,string,string)",
		func(pad ecsx.Appender) {
			pos := pad.SetParamHeader(3)
			pos = pad.SetAmountArray(0, pos, tokenId...)
			pad.SetAddress(1, to)
			pad.SetAmount(2, _tokenType)
			pos = pad.SetText(3, pos, _tokenName)
			pos = pad.SetText(4, pos, _tokenDesc)
			pos = pad.SetText(5, pos, _tokenContent)
			pos = pad.SetText(6, pos, deployer)
		},
	)

	ntx, err := tc.TransferFuncNTX(
		nftToken.Private,
		data,
		"0",
		ecsx.GasFast,
		nil,
	)
	if err != nil {
		dbg.Red(err)
		r.Set(ntx.From(), "", 0, err)
		return r
	}

	h, n, e := tc.TransferFuncSEND(ntx)
	r.Set(ntx.From(), h, n, e)
	return r
}

func (NFT) Burn(privateKey string, tokenId string) WriteResult {
	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "burn",
	}

	tc := Sender().TSender(nftToken.Contract)
	data := ecsx.MakePadBytes(
		"burn(uint256)",
		func(pad ecsx.Appender) {
			pad.SetParamHeader(1)
			pad.SetAmount(0, tokenId)
		},
	)
	ntx, err := tc.TransferFuncNTX(
		privateKey,
		data,
		"0",
		ecsx.GasFast,
		nil,
	)
	if err != nil {
		dbg.Red(err)
		r.Set(ntx.From(), "", 0, err)
		return r
	}

	h, n, e := tc.TransferFuncSEND(ntx)
	r.Set(ntx.From(), h, n, e)
	return r
}

func (NFT) SetBaseURL_Admin(newURI string) WriteResult {
	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "setBaseURI",
	}

	tc := Sender().TSender(nftToken.Contract)
	data := ecsx.MakePadBytes(
		"setBaseURI(string)",
		func(pad ecsx.Appender) {
			pos := pad.SetParamHeader(1)
			pad.SetText(0, pos, newURI)
		},
	)
	ntx, err := tc.TransferFuncNTX(
		nftToken.Private,
		data,
		"0",
		ecsx.GasFast,
		nil,
	)
	if err != nil {
		dbg.Red(err)
		r.Set(ntx.From(), "", 0, err)
		return r
	}

	h, n, e := tc.TransferFuncSEND(ntx)
	r.Set(ntx.From(), h, n, e)
	return r
}

func (NFT) SetTokenURI_Admin(tokenId, tokenURI string) WriteResult {
	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "setTokenURI",
	}

	tc := Sender().TSender(nftToken.Contract)
	data := ecsx.MakePadBytes(
		"setTokenURI(uint256,string)",
		func(pad ecsx.Appender) {
			pos := pad.SetParamHeader(2)
			pad.SetAmount(0, tokenId)
			pad.SetText(1, pos, tokenURI)
		},
	)
	ntx, err := tc.TransferFuncNTX(
		nftToken.Private,
		data,
		"0",
		ecsx.GasFast,
		nil,
	)
	if err != nil {
		dbg.Red(err)
		r.Set(ntx.From(), "", 0, err)
		return r
	}

	h, n, e := tc.TransferFuncSEND(ntx)
	r.Set(ntx.From(), h, n, e)
	return r
}
func (NFT) SetTypeName_Admin(nType, nName string) WriteResult {
	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "setTypeName",
	}

	tc := Sender().TSender(nftToken.Contract)
	data := ecsx.MakePadBytes(
		"setTypeName(uint256,string)",
		func(pad ecsx.Appender) {
			pos := pad.SetParamHeader(2)
			pad.SetAmount(0, nType)
			pad.SetText(1, pos, nName)
		},
	)
	ntx, err := tc.TransferFuncNTX(
		nftToken.Private,
		data,
		"0",
		ecsx.GasFast,
		nil,
	)
	if err != nil {
		dbg.Red(err)
		r.Set(ntx.From(), "", 0, err)
		return r
	}

	h, n, e := tc.TransferFuncSEND(ntx)
	r.Set(ntx.From(), h, n, e)
	return r
}
func (NFT) SetTokenInfo_Admin(tokenId, tokenName, tokenDesc, tokenContent string) WriteResult {
	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "setTokenInfo",
	}

	tc := Sender().TSender(nftToken.Contract)
	data := ecsx.MakePadBytes(
		"setTokenInfo(uint256,string,string,string)",
		func(pad ecsx.Appender) {
			pos := pad.SetParamHeader(2)
			pad.SetAmount(0, tokenId)
			pos = pad.SetText(1, pos, tokenName)
			pos = pad.SetText(2, pos, tokenDesc)
			pos = pad.SetText(3, pos, tokenContent)
		},
	)
	ntx, err := tc.TransferFuncNTX(
		nftToken.Private,
		data,
		"0",
		ecsx.GasFast,
		nil,
	)
	if err != nil {
		dbg.Red(err)
		r.Set(ntx.From(), "", 0, err)
		return r
	}

	h, n, e := tc.TransferFuncSEND(ntx)
	r.Set(ntx.From(), h, n, e)
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
		abmx.Method{
			Name: "hasRole",
			Params: abmx.NewParams(
				abmx.NewBytes32(role),
				abmx.NewAddress(address),
			),
			Returns: abmx.NewReturns(
				abmx.Bool,
			),
		},
		nftToken.Address,
		func(rs abmx.RESULT) {
			f(rs.Bool(0))
		},
	)
}

func (NFT) GrantRole(role, address string) WriteResult {
	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "grantRole",
	}

	tc := Sender().TSender(nftToken.Contract)

	data := ecsx.MakePadBytes(
		"grantRole(bytes32, address)",
		func(pad ecsx.Appender) {
			pad.SetParamHeader(2)
			pad.SetBytes32Hex(0, role)
			pad.SetAddress(0, address)
		},
	)

	ntx, err := tc.TransferFuncNTX(
		nftToken.Private,
		data,
		"0",
		ecsx.GasFast,
		nil,
	)
	if err != nil {
		dbg.Red(err)
		r.Set(ntx.From(), "", 0, err)
		return r
	}

	h, n, e := tc.TransferFuncSEND(ntx)
	r.Set(ntx.From(), h, n, e)
	return r
}

func (NFT) RevokeRole(role, address string) WriteResult {
	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "revokeRole",
	}

	tc := Sender().TSender(nftToken.Contract)

	data := ecsx.MakePadBytes(
		"revokeRole(bytes32, address)",
		func(pad ecsx.Appender) {
			pad.SetParamHeader(2)
			pad.SetBytes32Hex(0, role)
			pad.SetAddress(0, address)
		},
	)

	ntx, err := tc.TransferFuncNTX(
		nftToken.Private,
		data,
		"0",
		ecsx.GasFast,
		nil,
	)
	if err != nil {
		dbg.Red(err)
		r.Set(ntx.From(), "", 0, err)
		return r
	}

	h, n, e := tc.TransferFuncSEND(ntx)
	r.Set(ntx.From(), h, n, e)
	return r
}

func (NFT) RenounceRole(privateKey string, role, address string) WriteResult {
	r := WriteResult{
		Constract: nftToken.Contract,
		FuncName:  "renounceRole",
	}

	tc := Sender().TSender(nftToken.Contract)

	data := ecsx.MakePadBytes(
		"renounceRole(bytes32, address)",
		func(pad ecsx.Appender) {
			pad.SetParamHeader(2)
			pad.SetBytes32Hex(0, role)
			pad.SetAddress(0, address)
		},
	)

	ntx, err := tc.TransferFuncNTX(
		privateKey,
		data,
		"0",
		ecsx.GasFast,
		nil,
	)
	if err != nil {
		dbg.Red(err)
		r.Set(ntx.From(), "", 0, err)
		return r
	}

	h, n, e := tc.TransferFuncSEND(ntx)
	r.Set(ntx.From(), h, n, e)
	return r
}

func (NFT) TransferFromNTX(privateKey string, to, tokenId string) (*ecsx.NTX, error) {
	w, err := jwalletx.Get(privateKey)
	if err != nil {
		return nil, err
	}

	tc := Sender().TSender(nftToken.Contract)
	data := ecsx.MakePadBytes(
		"transferFrom(address,address,uint256)",
		func(pad ecsx.Appender) {
			pad.SetParamHeader(3)
			pad.SetAddress(0, w.Address()) //from
			pad.SetAddress(1, to)          //to
			pad.SetAmount(2, tokenId)      //tokenId
		},
	)
	ntx, err := tc.TransferFuncNTX(
		privateKey,
		data,
		"0",
		gasSpeed,
		nil,
	)

	return ntx, err
}

func (NFT) TransferFromSEND(ntx *ecsx.NTX) (string, error) {
	tc := Sender().TSender(nftToken.Contract)
	h, n, e := tc.TransferFuncSEND(ntx)
	_ = n
	return h, e
}

package ecsx

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"txscheduler/brix/tools/cloudx/ebcmx"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsaa"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx/jwalletx"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// XPendingNonceAt : PendingNonceAt
func (my *Sender) XPendingNonceAt(hexAddress string) (uint64, error) {
	address := common.HexToAddress(hexAddress)
	nonce, err := my.client.PendingNonceAt(context.Background(), address)
	return nonce, err
}

// XNonceAt : NonceAt ( pendingNonce )
func (my *Sender) XNonceAt(hexAddress string) (uint64, error) {
	address := common.HexToAddress(hexAddress)
	nonce, err := my.client.NonceAt(context.Background(), address, nil)
	//my.client.PendingTransactionCount(context.Background())

	return nonce, err
}

func (my *Sender) XGasLimit(
	paddedData PADBYTES,
	fromAddress, toAddress string,
	ethWEIs ...string,
) (uint64, error) {
	ethWEI := "0"
	if len(ethWEIs) > 0 {
		ethWEI = ethWEIs[0]
	}
	from := common.HexToAddress(fromAddress)
	to := common.HexToAddress(toAddress)
	var ethValue *big.Int
	if !jmath.IsUnderZero(ethWEI) {
		ethValue = new(big.Int)
		ethValue.SetString(jmath.VALUE(ethWEI), 10)
	}

	gasLimit, err := my.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  from,
		To:    &to,
		Value: ethValue,
		Data:  paddedData.Bytes(),
	})

	return gasLimit, err
}

func (my *Sender) ebcm_XGasLimit(
	paddedData ebcmx.PADBYTES,
	fromAddress, toAddress string,
	wei string,
) (uint64, error) {
	from := common.HexToAddress(fromAddress)
	to := common.HexToAddress(toAddress)
	var ethValue *big.Int
	if !jmath.IsUnderZero(wei) {
		ethValue = new(big.Int)
		ethValue.SetString(jmath.VALUE(wei), 10)
	}

	data := paddedData.Bytes()
	dbg.Cyan(paddedData.Hex())

	gasLimit, err := my.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  from,
		To:    &to,
		Value: ethValue,
		Data:  data,
	})

	return gasLimit, err
}

type XGasPrice struct {
	val    *big.Int
	xError error
	xWEI   string
}

func (my XGasPrice) GetGWEI() string {
	return ebcmx.WEIToGWEI(jmath.VALUE(my.val)).String()
}

func (my XGasPrice) Valid() bool { return my.val != nil }

func (my XGasPrice) String() string {
	if my.xError != nil || my.val == nil {
		return my.xError.Error()
	}
	tag := dbg.Cat(
		"< gasPirce >", dbg.ENTER,
		" ETH  :", my.ETH(), dbg.ENTER,
		" GWEI :", ETH(my.ETH()).ToGWEI(), dbg.ENTER,
	)
	return tag
}

func (my XGasPrice) Clone() XGasPrice {
	clone := XGasPrice{
		val:    big.NewInt(my.val.Int64()),
		xError: my.xError,
		xWEI:   my.xWEI,
	}
	return clone
}
func (my XGasPrice) Cmp(dst XGasPrice) int {
	return jmath.CMP(my.val, dst.val)
}

func NewXGasPriceETH(ethVal string) XGasPrice {
	return NewXGasPrice(ETHToWei(ethVal))
}

func NewXGasPrice(wei string) XGasPrice {
	v := big.NewInt(0)
	v.SetString(wei, 10)
	item := XGasPrice{
		val:  v,
		xWEI: wei,
	}
	return item
}

func (my XGasPrice) Error() error {
	return my.xError
}

func (my XGasPrice) WEI() string {
	if my.val == nil {
		return "0"
	}
	return jmath.VALUE(my.val)
}
func (my XGasPrice) ETH() string {
	return WeiToETH(my.WEI())
}

func (my XGasPrice) GWEI() string {
	return WEIToGWEI(WEI(my.WEI())).String()
}

func (my *XGasPrice) AddWEI(weiVal interface{}) {
	addWei := jmath.VALUE(weiVal)
	//addWei := ETHToWei(val)

	org := my.WEI()
	newVal := jmath.ADD(org, addWei)

	my.val.SetString(newVal, 10)
}

func (my *XGasPrice) AddGWEI(gwei string) {
	my.AddWEI(GWEI(gwei).ToWEI().String())
}
func (my *XGasPrice) SetGWEI(gwei string) {
	my.val = jmath.New(GWEI(gwei).ToWEI().String()).ToBigInteger()
}

func (my XGasPrice) FeeWEI(limit uint64) string {
	return jmath.MUL(my.WEI(), limit)
}
func (my XGasPrice) FeeETH(limit uint64) string {
	wei := my.FeeWEI(limit)
	return WeiToETH(wei)
}

func (my *Sender) SUGGEST_GAS_PRICE(speed GasSpeed) XGasPrice {
	gas := ecsaa.SUGGEST_GAS_PRICE(my)
	if gas == nil {
		return XGasPrice{xError: errors.New("[XGasPrice] gasValue is nil")}
	}
	item := XGasPrice{val: gas}
	item.xWEI = item.WEI()
	return item
}
func (my *Sender) ebcm_XGasPrice(speed ebcmx.GasSpeed) ebcmx.XGasPrice {
	item := my.SUGGEST_GAS_PRICE(GasSpeed(speed))
	return &item
}

type xNTX struct {
	tx       *types.Transaction
	gasPrice *big.Int
	xError   error
}

func (my xNTX) Error() error {
	return my.xError
}

func (my *Sender) XNTX(
	paddedData PADBYTES,
	toAddress string,
	ethWei string,
	nonce uint64,
	gasLimit uint64,
	gasPrice XGasPrice,
) xNTX {
	to := common.HexToAddress(toAddress)

	var ethValue *big.Int
	if !jmath.IsUnderZero(ethWei) {
		ethValue = new(big.Int)
		ethValue.SetString(jmath.VALUE(ethWei), 10)
	}

	var tx *types.Transaction
	switch my.txnType {
	case TXN_EIP_1559:
		tip_cap := ecsaa.SUGGEST_TIP_PRICE(my)

		tx = types.NewTx(&types.DynamicFeeTx{
			ChainID:   my.chainID,
			Nonce:     nonce,
			GasTipCap: tip_cap,
			GasFeeCap: gasPrice.val,
			Gas:       gasLimit,
			To:        &to,
			Value:     ethValue,
			Data:      paddedData.Bytes(),
		})

	default:
		tx = types.NewTransaction(
			nonce,
			to,
			ethValue,
			gasLimit,
			gasPrice.val,
			paddedData.Bytes(),
		)
	}

	item := xNTX{
		tx:       tx,
		gasPrice: gasPrice.val,
	}
	if tx == nil {
		item.xError = errors.New("[XNTX] is nil")
	}

	return item
}

func (my *Sender) XNTX_FixedGAS(
	paddedData PADBYTES,
	toAddress string,
	ethWei string,
	nonce uint64,
	gasLimit uint64,
	gasPair ...*big.Int,
) xNTX {
	to := common.HexToAddress(toAddress)

	var ethValue *big.Int
	if !jmath.IsUnderZero(ethWei) {
		ethValue = new(big.Int)
		ethValue.SetString(jmath.VALUE(ethWei), 10)
	}

	var tx *types.Transaction

	var left *big.Int
	var right *big.Int

	left = gasPair[0]
	if len(gasPair) > 1 {
		right = gasPair[1]
	} else {
		right = jmath.New(left).ToBigInteger()
	}

	switch my.txnType {
	case TXN_EIP_1559:

		tx = types.NewTx(&types.DynamicFeeTx{
			ChainID:   my.chainID,
			Nonce:     nonce,
			GasTipCap: left,
			GasFeeCap: right,
			Gas:       gasLimit,
			To:        &to,
			Value:     ethValue,
			Data:      paddedData.Bytes(),
		})

	default:
		tx = types.NewTransaction(
			nonce,
			to,
			ethValue,
			gasLimit,
			gasPair[0],
			paddedData.Bytes(),
		)
	}

	item := xNTX{
		tx:       tx,
		gasPrice: right,
	}
	if tx == nil {
		item.xError = errors.New("[XNTX] is nil")
	}

	return item
}

type xSTX struct {
	tx     *types.Transaction
	xError error
}

func (my xSTX) Error() error {
	return my.xError
}
func (my xSTX) Hash() string {
	if my.tx == nil {
		return ""
	}
	return my.tx.Hash().Hex()
}

func (my *Sender) XSTX(privateKeyString string, ntx xNTX) xSTX {
	privatekey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		return xSTX{xError: fmt.Errorf("[XSTX] HexToECDSA : %v", err)}
	}

	chainID, err := my.client.NetworkID(context.Background())
	if err != nil {
		return xSTX{xError: fmt.Errorf("[XSTX] NetworkID : %v", err)}
	}

	var signer types.Signer
	if ntx.tx.Type() == 2 {
		signer = types.NewLondonSigner(chainID)
	} else {
		signer = types.NewEIP155Signer(chainID)
	}

	signedTx, err := types.SignTx(ntx.tx, signer, privatekey)
	if err != nil {
		return xSTX{xError: fmt.Errorf("[XSTX] SignTx : %v", err)}
	}

	return xSTX{tx: signedTx}
}

func (my *Sender) XSend(stx xSTX) error {
	return my.client.SendTransaction(context.Background(), stx.tx)
}

///////////////////////////////////////////////////////////////////////////

type XSendResult struct {
	From        string `bson:"from" json:"from"`
	To          string `bson:"to" json:"to"`
	PadBytesHex string `bson:"padBytesHex" json:"padBytesHex"` // hexString
	Nonce       uint64 `bson:"nonce" json:"nonce"`
	GasLimit    uint64 `bson:"gasLimit" json:"gasLimit"`
	GasPrice    string `bson:"gasPrice" json:"gasPrice"` // wei
	Hash        string `bson:"hash" json:"hash"`
}

func (my XSendResult) GetPadBytes() PadBytes   { return PadBytesFromHex(my.PadBytesHex) }
func (my XSendResult) GetXGasPrice() XGasPrice { return NewXGasPrice(my.GasPrice) }

func (my XSendResult) String() string { return dbg.ToJSONString(my) }

func (my *Sender) XPipe(
	fromPrivate string,
	toAddress string,
	padBytes PadBytes,
	wei string,
	speed GasSpeed,
	limitCB func(gasLimit uint64) uint64,
	nonceCB func(nonce uint64) uint64,
	gaspCB func(gasPrice XGasPrice) XGasPrice,
	resultCB func(r XSendResult),
) error {
	if resultCB == nil {
		return errors.New("resultCB is nil")
	}

	from, err := jwalletx.Get(fromPrivate)
	if err != nil {
		return err
	}

	toAddress = dbg.TrimToLower(toAddress)
	gasLimit, err := my.XGasLimit(
		padBytes,
		from.Address(),
		toAddress,
		wei,
	)
	if err != nil {
		return err
	}
	if limitCB != nil {
		gasLimit = limitCB(gasLimit)
	}
	nonce, err := my.XPendingNonceAt(from.Address())
	if err != nil {
		return err
	}
	if nonceCB != nil {
		nonce = nonceCB(nonce)
	}
	gasPrice := my.SUGGEST_GAS_PRICE(speed)
	if gaspCB != nil {
		gasPrice = gaspCB(gasPrice)
	}

	ntx := my.XNTX(
		padBytes,
		toAddress,
		wei,
		nonce,
		gasLimit,
		gasPrice,
	)
	if err := ntx.Error(); err != nil {
		return err
	}

	stx := my.XSTX(
		from.PrivateKey(),
		ntx,
	)

	if err := stx.Error(); err != nil {
		return err
	}

	if err := my.XSend(stx); err != nil {
		return err
	}

	result := XSendResult{
		From:        from.Address(),
		To:          toAddress,
		PadBytesHex: padBytes.Hex(),
		Nonce:       nonce,
		GasLimit:    gasLimit,
		GasPrice:    gasPrice.WEI(),
		Hash:        stx.Hash(),
	}
	resultCB(result)

	return nil
}

func (my *Sender) ebcm_XEstimateFeeETH(
	fromPrivate string,
	toAddress string,
	ipadBytes ebcmx.PADBYTES,
	wei string,
	speed ebcmx.GasSpeed,
) string {
	from, err := jwalletx.Get(fromPrivate)
	if err != nil {
		return "0"
	}
	//padBytes := ipadBytes.(PadBytes)

	toAddress = dbg.TrimToLower(toAddress)
	gasLimit, err := my.XGasLimit(
		ipadBytes,
		from.Address(),
		toAddress,
		wei,
	)
	if err != nil {
		return "0"
	}

	ebcmGasPrice := my.ebcm_XGasPrice(speed)
	return ebcmGasPrice.FeeETH(gasLimit)

}

func (my *Sender) ebcm_XPipe(
	fromPrivate string,
	toAddress string,
	ipadBytes ebcmx.PADBYTES,
	wei string,
	speed ebcmx.GasSpeed,
	limitCB func(gasLimit uint64) uint64,
	nonceCB func(nonce uint64) uint64,
	gaspCB func(gasPrice ebcmx.XGasPrice) ebcmx.XGasPrice,
	resultCB func(r ebcmx.XSendResult),
) error {
	if resultCB == nil {
		return errors.New("resultCB is nil")
	}

	from, err := jwalletx.Get(fromPrivate)
	if err != nil {
		return err
	}
	//padBytes := ipadBytes.(PadBytes)

	toAddress = dbg.TrimToLower(toAddress)
	gasLimit, err := my.XGasLimit(
		ipadBytes,
		from.Address(),
		toAddress,
		wei,
	)
	if err != nil {
		return err
	}
	if limitCB != nil {
		gasLimit = limitCB(gasLimit)
	}
	nonce, err := my.XPendingNonceAt(from.Address())
	if err != nil {
		return err
	}
	if nonceCB != nil {
		nonce = nonceCB(nonce)
	}
	ebcmGasPrice := my.ebcm_XGasPrice(speed)
	if gaspCB != nil {
		ebcmGasPrice = gaspCB(ebcmGasPrice)
	}

	gasPrice := *(ebcmGasPrice.(*XGasPrice))

	ntx := my.XNTX(
		ipadBytes,
		toAddress,
		wei,
		nonce,
		gasLimit,
		gasPrice,
	)
	if err := ntx.Error(); err != nil {
		return err
	}

	stx := my.XSTX(
		from.PrivateKey(),
		ntx,
	)

	if err := stx.Error(); err != nil {
		return err
	}

	if err := my.XSend(stx); err != nil {
		return err
	}

	result := XSendResult{
		From:        from.Address(),
		To:          toAddress,
		PadBytesHex: ipadBytes.Hex(),
		Nonce:       nonce,
		GasLimit:    gasLimit,
		GasPrice:    gasPrice.WEI(),
		Hash:        stx.Hash(),
	}

	r := ebcmx.XSendResult{}
	dbg.ChangeStruct(result, &r)
	resultCB(r)

	return nil
}

func (my *Sender) ebcm_TransferCoin(
	fromPrivate string,
	toAddress string,
	wei string,
	speed ebcmx.GasSpeed,
	limitCB func(gasLimit uint64) uint64,
	nonceCB func(nonce uint64) uint64,
	gaspCB func(gasPrice ebcmx.XGasPrice) ebcmx.XGasPrice,
	resultCB func(r ebcmx.XSendResult),
) error {
	return my.ebcm_XPipe(
		fromPrivate,
		toAddress,
		PadBytesETH(),
		wei,
		speed,
		limitCB,
		nonceCB,
		gaspCB,
		resultCB,
	)
}

func (my *Sender) XPipeFixedGAS(
	fromPrivate string,
	toAddress string,
	padBytes PadBytes,
	wei string,
	speed GasSpeed,
	limitCB func(gasLimit uint64) uint64,
	nonceCB func(nonce uint64) uint64,
	gasPair []*big.Int,
	resultCB func(r XSendResult),
) error {
	if resultCB == nil {
		return errors.New("resultCB is nil")
	}
	from, err := jwalletx.Get(fromPrivate)
	if err != nil {
		return err
	}

	toAddress = dbg.TrimToLower(toAddress)
	gasLimit, err := my.XGasLimit(
		padBytes,
		from.Address(),
		toAddress,
		wei,
	)
	if err != nil {
		return err
	}
	if limitCB != nil {
		gasLimit = limitCB(gasLimit)
	}
	nonce, err := my.XPendingNonceAt(from.Address())
	if err != nil {
		return err
	}
	if nonceCB != nil {
		nonce = nonceCB(nonce)
	}

	ntx := my.XNTX_FixedGAS(
		padBytes,
		toAddress,
		wei,
		nonce,
		gasLimit,
		gasPair...,
	)

	if err := ntx.Error(); err != nil {
		return err
	}

	stx := my.XSTX(
		from.PrivateKey(),
		ntx,
	)

	if err := stx.Error(); err != nil {
		return err
	}

	if err := my.XSend(stx); err != nil {
		return err
	}

	gaspriceWEI := jmath.VALUE(gasPair[0])
	if len(gasPair) > 1 {
		gaspriceWEI = jmath.VALUE(gasPair[1])
	}

	result := XSendResult{
		From:        from.Address(),
		To:          toAddress,
		PadBytesHex: padBytes.Hex(),
		Nonce:       nonce,
		GasLimit:    gasLimit,
		GasPrice:    gaspriceWEI,
		Hash:        stx.Hash(),
	}
	resultCB(result)

	return nil
}

func (my *Sender) ebcm_XPipeFixedGAS(
	fromPrivate string,
	toAddress string,
	ipadBytes ebcmx.PADBYTES,
	wei string,
	limitCB func(gasLimit uint64) uint64,
	nonceCB func(nonce uint64) uint64,
	gasPair []*big.Int,
	txFeeWeiAllow func(feeWEI string) bool, //isContinue
	resultCB func(r ebcmx.XSendResult),
) error {
	if resultCB == nil {
		return errors.New("resultCB is nil")
	}

	from, err := jwalletx.Get(fromPrivate)
	if err != nil {
		return err
	}
	//padBytes := ipadBytes.(PadBytes)

	toAddress = dbg.TrimToLower(toAddress)
	gasLimit, err := my.XGasLimit(
		ipadBytes,
		from.Address(),
		toAddress,
		wei,
	)
	if err != nil {
		return err
	}
	if limitCB != nil {
		gasLimit = limitCB(gasLimit)
	}
	nonce, err := my.XPendingNonceAt(from.Address())
	if err != nil {
		return err
	}
	if nonceCB != nil {
		nonce = nonceCB(nonce)
	}

	ntx := my.XNTX_FixedGAS(
		ipadBytes,
		toAddress,
		wei,
		nonce,
		gasLimit,
		gasPair...,
	)
	if err := ntx.Error(); err != nil {
		return err
	}

	if txFeeWeiAllow != nil {
		txFeeWei := jmath.MUL(gasLimit, ntx.gasPrice)
		isAllow := txFeeWeiAllow(txFeeWei)
		if !isAllow {
			return ebcmx.TxCancelUserFee
		}
	}

	stx := my.XSTX(
		from.PrivateKey(),
		ntx,
	)

	if err := stx.Error(); err != nil {
		return err
	}

	if err := my.XSend(stx); err != nil {
		return err
	}

	gaspriceWEI := jmath.VALUE(gasPair[0])
	if len(gasPair) > 1 {
		gaspriceWEI = jmath.VALUE(gasPair[1])
	}
	result := XSendResult{
		From:        from.Address(),
		To:          toAddress,
		PadBytesHex: ipadBytes.Hex(),
		Nonce:       nonce,
		GasLimit:    gasLimit,
		GasPrice:    gaspriceWEI,
		Hash:        stx.Hash(),
	}

	r := ebcmx.XSendResult{}
	dbg.ChangeStruct(result, &r)
	resultCB(r)

	return nil
}

func (my *Sender) ebcm_TransferCoinFixedGAS(
	fromPrivate string,
	toAddress string,
	wei string,
	limitCB func(gasLimit uint64) uint64,
	nonceCB func(nonce uint64) uint64,
	gasPair []*big.Int,
	txFeeWeiAllow func(feeWEI string) bool,
	resultCB func(r ebcmx.XSendResult),
) error {
	return my.ebcm_XPipeFixedGAS(
		fromPrivate,
		toAddress,
		PadBytesETH(),
		wei,
		limitCB,
		nonceCB,
		gasPair,
		txFeeWeiAllow,
		resultCB,
	)
}

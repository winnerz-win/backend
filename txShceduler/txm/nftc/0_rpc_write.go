package nftc

import (
	"context"
	"jtools/cloud/ebcm"
	"jtools/cloud/jeth/jwallet"
)

func TransferFuncNTX_Send(
	sender *ebcm.Sender,
	contract string,

	fromPrivate string,
	data ebcm.PADBYTES,
	wei string,
	snap *ebcm.GasSnapShot,
) (string, string, uint64, error) { //from.address, hash , nonce , err
	from, err := jwallet.Get(fromPrivate)
	if err != nil {
		return "", "", 0, err
	}
	from_address := from.Address()

	nonce, err := ebcm.MMA_GetNonce(sender, from.Address(), true)
	if err != nil {
		return from_address, "", 0, err
	}

	gas_price, err := sender.SuggestGasPrice(context.Background(), true)
	if err != nil {
		return from_address, "", 0, err
	}

	limit, err := sender.EstimateGas(
		context.Background(),
		ebcm.MakeCallMsg(
			from.Address(),
			contract,
			wei,
			data,
		),
	)
	if err != nil {
		return from_address, "", 0, err
	}

	limit = ebcm.MMA_LimitBuffer(limit)

	if snap != nil {
		snap.Limit = limit
		snap.Price = gas_price.GET_GAS_ETH()
		snap.FeeWei = gas_price.EstimateGasFeeWEI(limit)
	}

	ntx := sender.NewTransaction(
		nonce,
		contract,
		wei,
		limit,
		gas_price,
		data,
	)

	stx, err := sender.SignTx(ntx, fromPrivate)
	if err != nil {
		return from_address, "", 0, err
	}
	hash, err := sender.SendTransaction(context.Background(), stx)
	if err != nil {
		return from_address, "", 0, err
	}

	return from_address, hash, nonce, nil

}

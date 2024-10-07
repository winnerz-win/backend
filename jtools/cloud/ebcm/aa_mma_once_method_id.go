package ebcm

import (
	"context"
	"jtools/cc"
	"jtools/dbg"
	"jtools/jmath"
	"sync"
	"time"
)

func MMA_MethodID_Append(
	TAG string,
	once_mu *sync.Once,
	sl MethodIDDataList,
) {
	once_mu.Do(func() {
		cc.Yellow("#########  ebcm.MMA_MethodID_Append[", TAG, "] #######")
		methodERC20s = append(methodERC20s, sl...)
	})
}

// MMA_LimitBuffer : jmath.Uint64(jmath.MUL(limit, 1.3))
func MMA_LimitBuffer(limit uint64) uint64 {
	return jmath.Uint64(jmath.MUL(limit, 1.3))
}

// MMA_LimitBuffer_MasterOut : jmath.Uint64(jmath.MUL(limit, 2))
func MMA_LimitBuffer_MasterOut(limit uint64) uint64 {
	return jmath.Uint64(jmath.MUL(limit, 1.5))
}

// MMA_LimitBufferMax : 1.8 배수
func MMA_LimitBufferMax(limit uint64) uint64 {
	return jmath.Uint64(jmath.MUL(limit, 1.8))
}

func MMA_LimitBufferCustom(limit uint64, mv any) uint64 {
	return jmath.Uint64(jmath.MUL(limit, mv))
}

func MMA_GetNonce(
	caller *Sender,
	account string,
	is_differ_break ...bool,
) (uint64, error) {
	ctx := context.Background()
	for {
		nonce, err := caller.NonceAt(ctx, account)
		if err != nil {
			return 0, err
		}

		pending, err := caller.PendingNonceAt(ctx, account)
		if err != nil {
			return 0, err
		}
		if nonce != pending {
			if dbg.IsTrue(is_differ_break) {
				return 0, dbg.Error("ebcm.MMA_GetNonce[", nonce, "/", pending, "] is Differ.")
			}
			time.Sleep(time.Second)
			continue
		}

		return nonce, nil
	} //for
}

/////////////////////////////////////////////////////////////////

type GasSnapShot struct {
	Limit      uint64 `bson:"limit" json:"limit"`
	Price      string `bson:"price" json:"price"`
	FeeWei     string `bson:"fee_wei" json:"fee_wei"` // limit * price
	FixedNonce uint64 `bson:"fixedNonce,omitempty" json:"fixedNonce,omitempty"`
}

func MakeGasSnapShot(
	nonce uint64,
	limit uint64,
	gas_price GasPrice,
) GasSnapShot {
	snap := GasSnapShot{
		FixedNonce: nonce,
		Limit:      limit,
		FeeWei:     gas_price.EstimateGasFeeWEI(limit),
	}
	return snap
}

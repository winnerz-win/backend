package ecsx

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"txscheduler/brix/tools/jmath"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// STXList :
type STXList []*STX

// Append :
func (my *STXList) Append(stx *STX) {
	*my = append(*my, stx)
}

// Clear :
func (my *STXList) Clear() {
	*my = (*my)[:0]
}

// STX :
type STX struct {
	client    *ethclient.Client
	tx        *types.Transaction
	isSend    bool
	isPending bool

	gasFeeWei string //트랜젝션 수수료

	from        string
	blockNumber string

	UserData TxUserData
}

// GasFeeWei : 가스 수수료 (트랜젝션 발생시만)
func (my STX) GasFeeWei() string {
	return my.gasFeeWei
}

// Transaction :
func (my STX) Transaction() *types.Transaction {
	return my.tx
}

// BlockNumber :
func (my STX) BlockNumber() string {
	//my.Receipt().GetData().BlockNumber
	return my.blockNumber
}

// String :
func (my STX) String() string {
	view := txViewer{my.tx}.Map()
	view["isSend"] = my.isSend
	view["isPending"] = my.isPending
	b, _ := json.MarshalIndent(view, "", "  ")
	return string(b)
}

// Send :
func (my *STX) Send() error {
	if my.isSend {
		return errors.New("STX - already send.")
	}
	my.isSend = true
	err := my.client.SendTransaction(context.Background(), my.tx)
	if err != nil {
		return fmt.Errorf("STX - %v", err)
	}
	my.isPending = true
	return nil
}

// IsPending :
func (my STX) IsPending() bool {
	return my.isPending
}

// Hash : Send Transaction Hash
func (my STX) Hash() string {
	return my.tx.Hash().Hex()
}

// Receipt :
func (my STX) Receipt() *Receipt {
	sender := Sender{
		client: my.client,
	}
	return sender.Receipt(my.Hash())
}

// IsTokenTx : 토큰 트랙잭션 여부 반환
func (my STX) IsTokenTx() bool {
	buffer := hex.EncodeToString(my.tx.Data())
	return buffer != ""
}

// ContractAddress :
func (my STX) ContractAddress() string {
	buffer := hex.EncodeToString(my.tx.Data())
	if buffer == "" {
		return "n"
	}
	return strings.ToLower(my.tx.To().Hex())
}

// From :
func (my STX) From() string {
	return my.from
}

// To : 토큰 전송일 경우 to 는 컨트랙트 주소이다.
func (my STX) To() string {
	return strings.ToLower(my.tx.To().Hex())
}

// ContractTo : contract to address
func (my STX) ContractTo() string {
	buffer := hex.EncodeToString(my.tx.Data())
	if buffer == "" {
		return "n"
	}
	buffer = buffer[32:]
	buffer = buffer[:40]
	return "0x" + strings.ToLower(buffer)
}

// TokenValue :
func (my STX) TokenValue() string {
	buffer := hex.EncodeToString(my.tx.Data())
	if buffer == "" {
		return "0"
	}
	buffer = "0x" + buffer[32+40+46:]
	return jmath.NewBigDecimal(buffer).ToString()
}

// ETHValue :
func (my STX) ETHValue() string {
	return my.tx.Value().String()
}

func (my STX) Wei() string { return my.ETHValue() }

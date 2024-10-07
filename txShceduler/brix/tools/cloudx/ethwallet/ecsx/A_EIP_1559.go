package ecsx

type TXNTYPE int

const (
	TXN_LEGACY   = 0
	TXN_EIP_1559 = 2
)

func (my TXNTYPE) String() string {
	msg := ""
	switch my {
	case TXN_LEGACY:
		msg = "Legacy"
	case TXN_EIP_1559:
		msg = "EIP-1559"
	}
	return msg
}

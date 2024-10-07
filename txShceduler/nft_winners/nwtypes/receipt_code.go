package nwtypes

import (
	"strings"
	"txscheduler/brix/tools/crypt"
	"txscheduler/brix/tools/dbg"
)

type ReceiptKind string

func (my ReceiptKind) String() string { return string(my) }

const (
	RC_MINT_FREE  = ReceiptKind("nft_mint_free_")
	RC_MINT_TOKEN = ReceiptKind("nft_mint_token_")
	RC_MINT_COIN  = ReceiptKind("nft_mint_coin_")
	RC_SALE_COIN  = ReceiptKind("nft_sale_coin_")

	RC_SET_BASE_URI = ReceiptKind("nft_set_base_uri_")
)

////////////////////////////////////////////////////////////////////////////

type RECEIPT_CODE string

func (my RECEIPT_CODE) String() string { return string(my) }
func (my RECEIPT_CODE) Valid() bool {
	str := my.String()
	if strings.HasPrefix(str, string(RC_MINT_FREE)) {
		return true
	}
	if strings.HasPrefix(str, string(RC_MINT_TOKEN)) {
		return true
	}
	if strings.HasPrefix(str, string(RC_MINT_COIN)) {
		return true
	}
	if strings.HasPrefix(str, string(RC_SALE_COIN)) {
		return true
	}

	if strings.HasPrefix(str, string(RC_SET_BASE_URI)) {
		return true
	}

	return false
}

func GetReceiptCode(kind ReceiptKind, token_id string) RECEIPT_CODE {
	prefix := dbg.Cat(kind.String(), "[", token_id, "]_")
	uuid := crypt.MakeUID256()

	//cut := len(uuid)/2 + len(prefix)
	code := prefix + uuid[:len(uuid)-4]

	return RECEIPT_CODE(code)
}

////////////////////////////////////////////////////////////////////////////

type RESULT_TYPE string

func (my RESULT_TYPE) String() string { return string(my) }

const (
	RESULT_MINT      = RESULT_TYPE("mint")
	RESULT_USER_SALE = RESULT_TYPE("user_sale")
)

func (my RECEIPT_CODE) ResultType() RESULT_TYPE {
	if strings.Contains(my.String(), "mint") {
		return RESULT_MINT
	}
	if strings.Contains(my.String(), "sale") {
		return RESULT_USER_SALE
	}
	return ""
}

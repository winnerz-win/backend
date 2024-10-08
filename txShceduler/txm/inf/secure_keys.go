package inf

import (
	"fmt"
	"jtools/cloud/jeth/jwallet"
	"strings"
)

// KeyPair :
type KeyPair struct {
	Mainnet    bool   `yaml:"mainnet" json:"mainnet"`
	PrivateKey string `yaml:"privatekey" json:"privatekey"`
	Address    string `yaml:"address" json:"address"`
}

func (my *KeyPair) Refactory() {
	my.PrivateKey = strings.TrimSpace(my.PrivateKey)
	my.Address = strings.TrimSpace(my.Address)
}

// KeyPairList :
type KeyPairList []KeyPair

func (my KeyPair) String() string { return "dbg.To" }

// valid :
func (my KeyPair) valid() bool {
	//jwallet.New()
	w, err := jwallet.Get(my.PrivateKey)
	if err != nil {
		return false
	}
	isValid := w.Address() == my.Address
	if isValid == false {
		panic(fmt.Sprintf("inf.secure_keys : %v", my.Address))
	}
	return isValid
}

/*
블록체인 입출금 서버 구축
블록체인 입출금 서버 중요 지갑 정보 제공 기능
중요지갑 : 시스템 운영 지갑 (마스터지갑 , 가스비지갑)

회원가입시 유저 입금전용 가상계좌(코인 입금용 주소) 발급 기능
유저의 가상계좌로 코인 입금시 입금 금액 전송 기능
유저의 코인(암호화화폐) 출금 요청 기능
현재 출금 대기중인 건수 요청 기능
출금 요청한 코인의 출금결과(성공,실패) 전송기능


*/

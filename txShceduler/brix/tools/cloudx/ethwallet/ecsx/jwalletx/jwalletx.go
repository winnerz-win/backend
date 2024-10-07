package jwalletx

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"txscheduler/brix/tools/cloudx/ebcmx"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

// Wallet :
type Wallet struct {
	hexPrivatekey string //대문자
	hexAddress    string //소문자
}

// SWallet : Wallet of Seed
type SWallet struct {
	*Wallet
	text string
	seq  interface{}
}

// IWallet :
type IWallet interface {
	Index() interface{}
	String() string
	PrivateKey() string
	Address() string
	CompareAddress(cmpAddress string) bool
}

// Index :
func (my SWallet) Index() interface{} {
	return my.seq
}

func (my Wallet) Index() interface{} {
	return 0
}

// String :
func (my Wallet) String() string {
	m := map[string]interface{}{
		"privateKey": my.hexPrivatekey,
		"address":    my.hexAddress,
	}
	b, _ := json.MarshalIndent(m, "", "  ")
	return string(b)
}

// PrivateKey :
func (my Wallet) PrivateKey() string {
	return strings.TrimSpace(my.hexPrivatekey)
}

// Address :
func (my Wallet) Address() string {
	return strings.TrimSpace(my.hexAddress)
}

// CompareAddress :
func (my Wallet) CompareAddress(cmpAddress string) bool {
	cmpAddress = strings.ToLower(cmpAddress)
	return my.hexAddress == cmpAddress
}

func (my *Wallet) genAddress() error {
	privatekey, err := crypto.HexToECDSA(my.hexPrivatekey)
	if err != nil {
		return err
	}
	publickey := privatekey.Public()
	ecdsakey := publickey.(*ecdsa.PublicKey)

	publickeyBytes := crypto.FromECDSAPub(ecdsakey)

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publickeyBytes[1:])

	address := fmt.Sprintf("%v", hexutil.Encode(hash.Sum(nil)[12:]))

	my.hexAddress = strings.ToLower(address)

	return nil
}

// New : Generate
func New() *Wallet {
	privateKey, _ := crypto.GenerateKey()
	privateBytes := crypto.FromECDSA(privateKey)
	wallet := &Wallet{
		hexPrivatekey: strings.ToUpper(hexutil.Encode(privateBytes)[2:]),
	}
	wallet.genAddress()
	return wallet
}

// NewSeed :
func NewSeed(text string, seq interface{}) *SWallet {
	sText := fmt.Sprintf("%v%v", text, seq)
	hash := sha256.Sum256([]byte(sText))

	encHash := hex.EncodeToString(hash[:])
	reader := bytes.NewBuffer([]byte(encHash))

	privateKey, _ := ecdsa.GenerateKey(crypto.S256(), reader)
	privateBytes := crypto.FromECDSA(privateKey)
	wallet := &Wallet{
		hexPrivatekey: strings.ToUpper(hexutil.Encode(privateBytes)[2:]),
	}
	wallet.genAddress()
	my := &SWallet{
		Wallet: wallet,
		text:   text,
		seq:    seq,
	}
	return my
}

// NewSeedI :
func NewSeedI(text string, seq interface{}) IWallet {
	return NewSeed(text, seq)
}

// Get :
func Get(hexPrivate string) (*Wallet, error) {
	hexPrivate = strings.ToUpper(hexPrivate)
	wallet := &Wallet{
		hexPrivatekey: hexPrivate,
	}
	err := wallet.genAddress()
	return wallet, err
}

func EBCM_NewSeedI(text string, seq interface{}) ebcmx.IWallet {
	return NewSeedI(text, seq)
}
func EBCM_Get(hexPrivate string) (ebcmx.IWallet, error) {
	return Get(hexPrivate)
}

package EtherClient

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"regexp"
	"strings"

	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/crypto"
)

//EthWallet :
type EthWallet struct {
	privatekey string

	address string
}

//PrivateKeyString :
func (my EthWallet) PrivateKeyString() string {
	return my.privatekey
}

//PrivateKey :
func (my EthWallet) PrivateKey() *ecdsa.PrivateKey {
	pk, _ := crypto.HexToECDSA(my.privatekey)
	return pk
}

//PublicKey :
func (my EthWallet) PublicKey() *ecdsa.PublicKey {
	pk := my.PrivateKey()
	pubkey := pk.Public()
	pubECDSA, _ := pubkey.(*ecdsa.PublicKey)
	return pubECDSA
}

//PublicKeyString :
func (my EthWallet) PublicKeyString() string {
	pubECDSA := my.PublicKey()
	publicBytes := crypto.FromECDSAPub(pubECDSA)
	return hexutil.Encode(publicBytes)[4:]
}

//AddressString :
func (my *EthWallet) AddressString() string {
	if my.address != "" {
		return my.address
	}
	pubECDSA := my.PublicKey()
	publicBytes := crypto.FromECDSAPub(pubECDSA)

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicBytes[1:])

	address := fmt.Sprintf("%v", hexutil.Encode(hash.Sum(nil)[12:]))
	my.address = strings.ToLower(address)
	return my.address

}

func (my EthWallet) ToString() string {
	return `privateKey : ` + my.PrivateKeyString() + `
publicKey  : ` + my.PublicKeyString() + `
address    : ` + my.AddressString()
}

//GenerateEthWallet :
func GenerateEthWallet() EthWallet {
	privatekey, _ := crypto.GenerateKey()
	pkbytes := crypto.FromECDSA(privatekey)

	return EthWallet{
		privatekey: hexutil.Encode(pkbytes)[2:],
	}
}

//NewEthWallet :
func NewEthWallet(privatekey string) EthWallet {
	privatekey = strings.ToLower(privatekey)
	return EthWallet{privatekey: privatekey}
}

//GenerateWallet : Wallet Test Code
func GenerateWallet() {
	privatekey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	pkbytes := crypto.FromECDSA(privatekey)
	fmt.Println("privatekey :", hexutil.Encode(pkbytes)[2:])

	publickey := privatekey.Public()
	pubECDSA, ok := publickey.(*ecdsa.PublicKey)
	if ok == false {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	pubbytes := crypto.FromECDSAPub(pubECDSA)
	fmt.Println("publickey :", hexutil.Encode(pubbytes)[4:])

	address := crypto.PubkeyToAddress(*pubECDSA).Hex()
	fmt.Println("address :", address)

	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubbytes[1:])

	fmt.Println("address :", hexutil.Encode(hash.Sum(nil)[12:]))

	fmt.Println("------------------------------------------")
	ew := EthWallet{
		privatekey: hexutil.Encode(pkbytes)[2:],
	}
	fmt.Println("privatekey :", ew.PrivateKeyString())
	fmt.Println("publickey  :", ew.PublicKeyString())
	fmt.Println("address    :", ew.AddressString())

	//valid := util.IsValidAddress(ew.AddressString())
	valid := IsValidAddress(ew.AddressString())
	fmt.Println(valid)
}

// IsValidAddress validate hex address
func IsValidAddress(iaddress interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	switch v := iaddress.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}

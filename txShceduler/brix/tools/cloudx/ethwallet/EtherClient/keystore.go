package EtherClient

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

//NewKeyStore :
func NewKeyStore(path string, password string) (accounts.Account, error) {
	ks := keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)

	account, err := ks.NewAccount(password)
	if err == nil {
		fmt.Println(account.Address.Hex())
	}
	return account, err
}

//LoadKeyStore :
func LoadKeyStore(path, fullName string, password string) (accounts.Account, error) {
	ks := keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)
	jsonBytes, err := ioutil.ReadFile(fullName)
	if err != nil {
		return accounts.Account{}, err
	}

	account, err := ks.Import(jsonBytes, password, password)
	if err == nil {
		fmt.Println(account.Address.Hex())

		if err := os.Remove(fullName); err != nil {
			fmt.Println("LoadKeyStore@Remove", err)
		}

	} else {
		fmt.Println("Import Error :", err)
	}
	return account, err
}

package inf

import (
	"jtools/cloud/ebcm"
	"strings"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
)

type TokenInfo struct {
	Mainnet  bool   `yaml:"mainnet" json:"mainnet" bson:"mainnet"`
	Contract string `yaml:"contract" json:"contract" bson:"contract"`
	Symbol   string `yaml:"symbol" json:"symbol" bson:"symbol"`
	Decimal  string `yaml:"decimal" json:"decimal" bson:"decimal"`
	IsCoin   bool   `yaml:"-" json:"is_coin" bson:"is_coin"`
}

func (my TokenInfo) String() string { return dbg.ToJSONString(my) }

func (my TokenInfo) Valid() bool { return my.Symbol != "" }

func (my TokenInfo) Wei(price string) string {
	return ebcm.TokenToWei(price, my.Decimal)
}
func (my TokenInfo) Price(wei string) string {
	return ebcm.WeiToToken(wei, my.Decimal)
}

type TokenInfoList []TokenInfo

func (my TokenInfoList) FirstERC20() TokenInfo {
	for _, v := range my {
		if !v.IsCoin {
			return v
		}
	}
	return TokenInfo{}
}

func (my *TokenInfo) Refactory(finder *ebcm.Sender) {
	my.Contract = dbg.TrimToLower(my.Contract)

	if finder == nil { //COIN
		my.IsCoin = true
		my.Symbol = strings.TrimSpace(my.Symbol)
		my.Decimal = strings.TrimSpace(my.Decimal)
		return
	}

	my.IsCoin = false
	token := GetERC20(finder, my.Contract)
	my.Symbol = token.Symbol()
	my.Decimal = token.Decimals()

}
func (my *TokenInfoList) Refactory() {
	for i := 0; i < len(*my); i++ {
		(*my)[i].Contract = dbg.TrimToLower((*my)[i].Contract)
	}
}

func (my TokenInfoList) GetContract(contract string) TokenInfo {
	for _, v := range my {
		if v.Contract == contract {
			return v
		}
	}
	return TokenInfo{}
}

func (my TokenInfoList) GetSymbol(symbol string) TokenInfo {
	for _, v := range my {
		if v.Symbol == symbol {
			return v
		}
	}
	return TokenInfo{}
}

func (my TokenInfoList) GetList() TokenInfoList {
	return my
}

func (my TokenInfoList) SymbolList() []string {
	list := []string{}
	for _, v := range my {
		list = append(list, v.Symbol)
	}
	return list
}

func (TokenInfo) IndexingDB() {
	DB().Action(config.DB, func(db mongo.DATABASE) {
		c := db.C(TXETHTokenEx)
		c.EnsureIndex(mongo.SingleIndex("contract", "1", true))
		c.EnsureIndex(mongo.SingleIndex("mainnet", "1", false))

		list := TokenInfoList{}
		c.Find(nil).All(&list)
		if len(list) > 0 {
			config.mu.Lock()
			for _, newToken := range list {
				if newToken.Mainnet != config.Mainnet {
					continue
				}

				isAdd := true
				for _, token := range config.Tokens {
					if token.Contract == newToken.Contract {
						isAdd = false
						break
					}
				}
				if isAdd {
					config.Tokens = append(config.Tokens, newToken)
				}
			} //for
			config.mu.Unlock()
		}
	})
}

func (my TokenInfo) insertDB(db mongo.DATABASE) error {
	return db.C(TXETHTokenEx).Insert(my)
}

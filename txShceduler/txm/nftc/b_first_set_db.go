package nftc

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func firstSetDB(db mongo.DATABASE) {
	isFirstSet := false

	aset := model.NftAset{}
	db.C(inf.NFTASET).Find(aset.Selector()).One(&aset)

	if aset.NFTContract != nftToken.Contract {
		isFirstSet = true
	}

	if isFirstSet {
		for _, colName := range inf.NFTList() {
			db.C(colName).RemoveAll(nil)
		} //for

		_startNumber = model.ZERO
		NFT{}.StartNumber(func(number string) {
			_startNumber = number
		})
		if _startNumber == model.ZERO {
			emsg := "[nftc] firstSetDB:: nft_token_contract : startNumber is ZERO"
			panic(emsg)
		}

		NFT{}.Name(func(n string) {
			aset.NftName = n
		})
		NFT{}.Symbol(func(n string) {
			aset.NftSymbol = n
		})

		aset.Number = _startNumber
		aset.NFTContract = nftToken.Contract
		aset.IsEnd = false
		aset.Selector()
		db.C(inf.NFTASET).Insert(aset)

	}

	NFT{}.GetBaseURI(func(v string) {
		aset.UpdateBaseURI(db, v)
	})
}

package ecsx

import "txscheduler/brix/tools/dbg"

//TxUserData : NTX , STX 에 유저가 추가하고 싶은 데이타 넣기
type TxUserData struct {
	item map[string]interface{}
}

func newTxUserData() TxUserData {
	return TxUserData{
		item: map[string]interface{}{},
	}
}

func (my *TxUserData) Set(key string, val interface{}) {
	my.item[key] = val
}

func (my *TxUserData) Parse(key string, p interface{}) {
	if v, do := my.item[key]; do {
		dbg.ChangeStruct(v, p)
	} else {
		dbg.Red("TxUserData not foudn key :", key)
	}
}

func (my *TxUserData) Get(key string) interface{} {
	if v, do := my.item[key]; do {
		return v
	}
	return nil
}

func (my *TxUserData) Clone() TxUserData {
	clone := newTxUserData()
	for k, v := range my.item {
		clone.item[k] = v
	}
	return clone
}

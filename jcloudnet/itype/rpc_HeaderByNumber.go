package itype

import "jtools/cloud/ebcm"

func (my IClient) HeaderByNumber(number any) *ebcm.BlockHeader {
	data := my.BlockByNumberSimple(number)
	if data == nil {
		return nil
	}
	return ebcm.NewBlockHeader(data, my.isKlay)
}

package txsigner

import (
	"jcloudnet/itype"
	"jtools/dbg"
)

type Option struct {
	InfraTag string `json:"infra_tag"`
	Signer   itype.TxSigner
}

func (my Option) String() string { return dbg.ToJsonString(my) }

package txsigner

import (
	"jcloudnet/itype"
	"jtools/jnet/chttp"
	"net/http"
)

func init() {
	hVersion()
	hSignTx()
}

func hVersion() {
	type RESULT struct {
		TITLE    string `json:"title"`
		InfraTag string `json:"infra_tag"`
	}

	method := chttp.GET
	url := "/version"

	handle.Add(
		method, url,
		func(w chttp.ResponseWriter, req *http.Request, ps chttp.Params) {
			result := RESULT{
				TITLE:    config.TITLE,
				InfraTag: opt.InfraTag,
			}
			chttp.OK(w, result)
		},
	)
}

func hSignTx() {
	method := chttp.POST
	url := "/sign_tx"
	handle.Add(
		method, url,
		func(w chttp.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := chttp.BindingStruct[itype.ReqTx](req)

			signer := opt.Signer
			ntx := signer.NewTransaction(
				cdata.TXNTYPE,
				cdata.Nonce(),
				cdata.To,
				cdata.Value,
				cdata.Limit(),
				cdata.GasPrice(),
				cdata.Data,
			)

			stx, err := signer.SignTx(
				cdata.ChainID(),
				ntx,
				cdata.PrivateKey(),
			)
			if err != nil {
				chttp.Fail(w, ERROR_SignTx, err)
				return
			}

			raw, err := stx.MarshalBinary()
			if err != nil {
				chttp.Fail(w, ERROR_MarshalBinary, err)
				return
			}

			hash := signer.GetHash(stx)
			result := itype.AckTx{
				Hash: hash,
				Raw:  raw,
			}
			chttp.OK(w, result)
		},
	)
}

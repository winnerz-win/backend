package abmx

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"golang.org/x/crypto/sha3"
)

func PadBytes(pureName string, ps AbiParams) []byte {

	pns := []string{}
	for _, v := range ps {
		pns = append(pns, v.p.String())
	}
	methodName := strings.TrimSpace(pureName) + "(" + strings.Join(pns, ",") + ")"
	//dbg.Purple("InputBytesName :", methodName)

	getMethod := func() map[string]abi.Method {
		m := abi.NewMethod(
			methodName,
			methodName,
			abi.Function,
			"",
			false,
			false,
			ps.getArgument(),
			[]abi.Argument{},
		)
		return map[string]abi.Method{
			methodName: m,
		}
	}

	abiSafe := abi.ABI{
		Methods: getMethod(),
	}
	r, err := abiSafe.Pack(methodName, ps.getParames()...)
	if err == nil {
		transferFnSignature := []byte(methodName) // ----- "transfer(address,uint256)"
		hash := sha3.NewLegacyKeccak256()
		hash.Write(transferFnSignature)
		pad := hash.Sum(nil)[:4]
		r[0] = pad[0]
		r[1] = pad[1]
		r[2] = pad[2]
		r[3] = pad[3]
	}
	return r

}

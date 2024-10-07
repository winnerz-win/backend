package cert

import (
	"fmt"
	"io/ioutil"
	"os"

	"txscheduler/brix/tools/dbg"
)

//CopyTarget :
func CopyTarget(rootPath, name, targetPath string) {
	copyTarget(rootPath, name, targetPath)
}

func copyTarget(rootPath, name, targetPath string) {
	nowPath := rootPath
	cert, err := ioutil.ReadFile(fmt.Sprintf("%v\\%v_cert.pem", nowPath, name))
	if err != nil {
		dbg.Red(err)
		return
	}
	dbg.Purple(string(cert))

	key, err := ioutil.ReadFile(fmt.Sprintf("%v\\%v_key.pem", nowPath, name))
	if err != nil {
		dbg.Red(err)
		return
	}
	dbg.Purple(string(key))

	req, err := ioutil.ReadFile(fmt.Sprintf("%v\\%v_req.csr", nowPath, name))
	if err != nil {
		dbg.Red(err)
		return
	}
	dbg.Purple(string(key))

	if file, err := os.Create(fmt.Sprintf("%v\\%v_cert.pem", targetPath, name)); err != nil {
		dbg.Red(err)
		return
	} else if _, err := file.Write(cert); err != nil {
		dbg.Red(err)
		return
	}
	if file, err := os.Create(fmt.Sprintf("%v\\%v_key.pem", targetPath, name)); err != nil {
		dbg.Red(err)
		return
	} else if _, err := file.Write(key); err != nil {
		dbg.Red(err)
		return
	}
	if file, err := os.Create(fmt.Sprintf("%v\\%v_req.csr", targetPath, name)); err != nil {
		dbg.Red(err)
		return
	} else if _, err := file.Write(req); err != nil {
		dbg.Red(err)
		return
	}
	dbg.Green("copy success")
}

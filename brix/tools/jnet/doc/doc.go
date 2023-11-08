package doc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"txscheduler/brix/tools/jnet/chttp"
)

//DocumentData :
type DocumentData struct {
	Href     string
	Size     string
	Weight   string
	Color    string
	Text     string
	onlyText bool
}

//DocStringList :
type DocStringList []DocumentData

//Count :
func (my DocStringList) Count() int {
	cnt := 0
	for _, v := range my {
		if v.onlyText == false {
			cnt++
		}
	} //for
	return cnt
}

var (
	urlPath    string
	docTItle   string
	docVersion string
	doclist    = DocStringList{}
)

//Ready : http://address:65530/doc/path ,
func Ready(path, titleString, versionString string) {
	urlPath = path
	docTItle = titleString
	docVersion = versionString
}

//Comment :
func Comment(comment string) *Doc {
	return newDoc(comment)
}

//Update : API-Doc 서버로 전송
func Update(isLocal ...bool) {
	update(urlPath, docTItle, docVersion, doclist, isLocal...)
}

//update :
func update(urlPath, docTItle, docVersion string, doclist DocStringList, isLocal ...bool) {
	if chttp.IPDO(localIP) == false {
		fmt.Println("[doc.Update] REAL-SERVER-RUN SKIP")
		return
	}
	address := realAddress + urlPath
	if len(isLocal) > 0 && isLocal[0] == true {
		address = localAddress + urlPath
	}
	_update(address, docTItle, docVersion, doclist)
}

func updateCustom(address, urlPath, docTItle, docVersion string, doclist DocStringList) {
	address = address + urlPath
	_update(address, docTItle, docVersion, doclist)
}

func _update(address, docTItle, docVersion string, doclist DocStringList) {

	reader := bytes.NewReader(HTMLBytes(urlPath, docTItle, docVersion, doclist))
	request, err := http.NewRequest("POST", address, reader)
	if err != nil {
		fmt.Println("[doc.Update].NewRequest", err)
		return
	}
	client := http.DefaultClient
	client.Transport = &http.Transport{
		Dial: (&net.Dialer{
			KeepAlive: 600 * time.Second,
		}).Dial,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		DisableKeepAlives:   true,
	}
	client.Timeout = time.Second * 3

	res, err := client.Do(request)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		fmt.Println("[doc.Update].Do", err)
		return
	}
	if res == nil {
		fmt.Println("[doc.Update].request.Body is nil")
		return
	}

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("[doc.Update].ReadAll", err)
		return
	}

	fmt.Println("--- doc.Update -------------------------------------------------------------------------")
	fmt.Println("; code :", res.StatusCode)
	fmt.Println("; body :", string(buf))
	fmt.Println("----------------------------------------------------------------------------------------")
}

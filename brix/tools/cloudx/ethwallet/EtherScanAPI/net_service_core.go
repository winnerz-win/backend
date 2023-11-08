package EtherScanAPI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"txscheduler/brix/tools/dbg"
)

const (
	get  = 0
	post = 1
)

//cFailCallbackData :
type cFailCallbackData interface {
	GetOrignalClass() interface{}
}

type cService struct {
	addr               string
	url                string
	method             int
	contentType        string
	requestData        interface{}
	resultCallbackData interface{}
	failCallbackData   cFailCallbackData
	resultErr          error

	fixTimeCallback func()
}

func (my *cService) setURL(v string) *cService {
	my.url = v
	return my
}
func (my cService) getURL() string {
	if my.url != "" {
		return my.addr + my.url
	}
	return my.addr
}
func (my cService) getParamsURL() string {
	url := my.getURL()
	return fmt.Sprintf("%v%v", url, my.requestData)
}

//Execute :
func (my *cService) Execute(isLog ...bool) *cService {
	doLog := false
	if len(isLog) > 0 {
		doLog = true
	}
	switch my.method {
	case get:
		url := my.getParamsURL()
		if doLog {
			dbg.Yellow("[GET]", url)
		}
		//dbg.Yellow("[GET]", url)
		if resp, err := http.Get(url); err == nil {
			defer resp.Body.Close()
			if body, err := ioutil.ReadAll(resp.Body); err == nil {
				if doLog {
					dbg.Yellow(string(body))
				}
				//dbg.Yellow(string(body))

				err := json.Unmarshal(body, my.resultCallbackData)
				if err != nil {
					if my.failCallbackData != nil {
						err = json.Unmarshal(body, my.failCallbackData)
						if err == nil {
							my.resultCallbackData = my.failCallbackData.GetOrignalClass()
						} else {
							my.resultErr = fmt.Errorf("[err:%v], Body:%v, Data:%v", err, string(body), my.requestData)
						}
					} else {
						my.resultErr = fmt.Errorf("[err:%v], Body:%v, Data:%v", err, string(body), my.requestData)
					}
				}
			} else {
				my.resultErr = fmt.Errorf("[err:%v], Body:%v, Data:%v", err, string(body), my.requestData)
			}
		} else {
			my.resultErr = err
		}

	case post:
		jdata, _ := json.Marshal(my.requestData)
		body := bytes.NewBuffer(jdata)
		url := my.getURL()
		if doLog {
			dbg.Yellow("[post]", url)
			dbg.Yellow("[body]", string(jdata))
		}

		if resp, err := http.Post(url, "application/json", body); err == nil {
			defer resp.Body.Close()
			if body, err := ioutil.ReadAll(resp.Body); err == nil {
				if doLog {
					dbg.Yellow(string(body))
				}
				err := json.Unmarshal(body, my.resultCallbackData)
				if err != nil {
					my.resultErr = fmt.Errorf("[%v]%v", string(body), err)
				}
			} else {
				my.resultErr = err
			}
		} else {
			my.resultErr = err
		}
	} //switch
	if my.fixTimeCallback != nil {
		my.fixTimeCallback()
	}
	return my
}

//Body :
func (my *cService) Body() interface{} {
	if my.resultErr != nil {
		if isSkipError == false {
			dbg.Red(my.getURL(), ", error:", my.resultErr)
		}
		return nil
	}
	return my.resultCallbackData
}

// ETHScanAPI :
type ETHScanAPI struct {
	*cService
	apikey    string
	sortOrder string //"asc"
}

//getService :
func getService(cfg Config) *ETHScanAPI {
	ns := &ETHScanAPI{
		cService: &cService{
			addr:               cfg.ConnenctURL(),
			url:                "",
			contentType:        "application/json",
			requestData:        nil,
			resultCallbackData: nil,
			resultErr:          nil,
			fixTimeCallback:    cfg.FixTime,
		},
		sortOrder: cfg.SortOrder(),
		apikey:    cfg.Apikey(),
	}
	return ns
}

func (my *ETHScanAPI) makeGetParams(params ...interface{}) *ETHScanAPI {

	size := len(params)

	msg := ""
	for i := 0; i < size; i += 2 {
		msg = fmt.Sprintf("%v%v", msg, params[i])
		if i+1 < size {
			msg = fmt.Sprintf("%v=%v", msg, params[i+1])
		}
		if i+2 < size {
			msg = fmt.Sprintf("%v&", msg)
		}
	} //for

	my.requestData = msg

	return my
}

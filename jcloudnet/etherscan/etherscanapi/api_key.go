package etherscanapi

import (
	"fmt"
	"jtools/dbg"
	"strings"
	"sync"
	"time"
)

var (
	FLAG_NETWORK_MAINNET    = false
	ETHEREUM_ETHERSCAN_MAIN = "http://api.etherscan.io"
	ETHEREUM_ETHERSCAN_TEST = "https://api-goerli.etherscan.io"

	ETHEREUM_ETHERSCAN_APIKEY = "YOURApiKeyToken"

	isSkipError = false
)

// SkipError :
func SkipError(flag ...bool) {
	isSkipError = dbg.IsTrue(flag)
}

//
/////////////////////////////////////////////////////////////////////////////////
//

// Config :
type Config interface {
	Mainnet() bool
	ConnenctURL() string
	Apikey() string
	SnapShotKey() string
	FixTime()

	SortAsc()
	SortDesc()
	SortOrder() string
}

const (
	//TESTKey : 1/3 sec
	TESTKey = "YOURApiKeyToken"
	//LIMITCOUNTFORSEC :
	LIMITCOUNTFORSEC = time.Millisecond * 3100 //3초에 1회 제한

	endblockNumber = "9999999999999"
	//EndBlockNumber :
	EndBlockNumber = endblockNumber

	SortASC  = "asc"
	SortDESC = "desc"
)

type testKeyConfig struct {
	callFlag bool
	callTime time.Time
	reqKeyC  chan chan string
	keyLock  chan struct{}
	counter  int64
}

var testCfg = testKeyConfig{
	callFlag: false,
	callTime: time.Now(),
	reqKeyC:  make(chan chan string),
	keyLock:  make(chan struct{}),
	counter:  0,
}

// func init() {
// 	go testCfg.run()
// }

var (
	testKeyOnce sync.Once
)

func (my *testKeyConfig) run() {
	fmt.Println("EtherScanAPI ::: testKeyConfig Once run from getAPIKey()")
	for req := range my.reqKeyC {
		if my.callFlag == true {
			du := time.Now().Sub(my.callTime)
			//dbg.Purple("du :", du)
			if du < LIMITCOUNTFORSEC {
				sleepValue := LIMITCOUNTFORSEC - du
				//dbg.Yellow("TestKey sleep :", sleepValue)
				time.Sleep(sleepValue)

			} else {
				//dbg.Yellow("3Sec Over")
			}
		}
		my.callFlag = true
		my.counter++
		//dbg.Yellow("call :", my.counter)

		req <- TESTKey
		<-my.keyLock
	} //for
}

func (my *testKeyConfig) getAPIKey() string {
	testKeyOnce.Do(func() {
		go testCfg.run()
	})
	keyC := make(chan string)
	my.reqKeyC <- keyC
	return <-keyC
}

type esKeyConfig struct {
	_mainnet   bool
	_url       string
	_apikey    string
	_isTestKey bool
	_sortOrder string
}

func (my *esKeyConfig) SortAsc() {
	my._sortOrder = SortASC
}
func (my *esKeyConfig) SortDesc() {
	my._sortOrder = SortDESC
}
func (my *esKeyConfig) SortOrder() string {
	return my._sortOrder
}

func (my *esKeyConfig) Mainnet() bool {
	return my._mainnet
}

func (my *esKeyConfig) ConnenctURL() string {
	// if my._mainnet {
	// 	return "http://api.etherscan.io"
	// }
	// return "https://api-goerli.etherscan.io"
	return my._url
}
func (my *esKeyConfig) Apikey() string {
	if my._isTestKey == true {
		return testCfg.getAPIKey()
	} else {
		//dbg.Yellow("allow user key")
	}
	return my._apikey
}

func (my *esKeyConfig) FixTime() {
	if my._isTestKey == true {
		testCfg.callTime = time.Now()
		testCfg.keyLock <- struct{}{}
	}
}

func (my *esKeyConfig) SnapShotKey() string {
	return my._apikey
}

// NewConfig :
func NewConfig(mainnet bool, apiKeys ...string) Config {
	cfg := &esKeyConfig{
		_mainnet:   mainnet,
		_isTestKey: true,
		_apikey:    TESTKey,
	}
	if cfg._mainnet {
		cfg._url = ETHEREUM_ETHERSCAN_MAIN
	} else {
		cfg._url = ETHEREUM_ETHERSCAN_TEST
	}
	cfg.SortAsc()

	if len(apiKeys) > 0 {
		cfg._isTestKey = false
		cfg._apikey = strings.TrimSpace(apiKeys[0])
		if cfg._apikey == TESTKey {
			cfg._isTestKey = true
		}
	}
	return cfg
}

func NewCustomConfig(mainnet bool, baseURL string, apiKeys ...string) Config {
	cfg := &esKeyConfig{
		_mainnet:   mainnet,
		_url:       baseURL,
		_isTestKey: true,
		_apikey:    TESTKey,
	}
	cfg.SortAsc()

	if len(apiKeys) > 0 {
		cfg._isTestKey = false
		cfg._apikey = strings.TrimSpace(apiKeys[0])
		if cfg._apikey == TESTKey {
			cfg._isTestKey = true
		}
	}
	return cfg
}

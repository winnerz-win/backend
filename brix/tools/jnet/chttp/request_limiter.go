package chttp

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"txscheduler/brix/tools/console"
	"txscheduler/brix/tools/dbg"
)

//hVisitor :
type hVisitor struct {
	sync.RWMutex
	manager  *limitManager
	ip       string
	inCnt    int
	lastSeen time.Time
}

func (my *hVisitor) Last() time.Time {
	defer my.RUnlock()
	my.RLock()
	return my.lastSeen
}

func (my *hVisitor) checkClear() {
	//defer fmt.Println(my.ip, " ------------ clear")
	for {
		time.Sleep(my.manager.clearDuration)
		if time.Now().Sub(my.Last()) >= 0 {
			my.manager.cleanUp(my.ip, my.inCnt)
			break
		}
	} //for
}
func (my *hVisitor) next() bool {
	defer my.Unlock()
	my.Lock()
	if my.inCnt >= my.manager.allowCount {
		my.inCnt++
		return false
	}
	my.inCnt++
	return true
}

type limitManager struct {
	allowCount    int
	clearDuration time.Duration
	blackCount    int

	visitors map[string]*hVisitor
	mu       sync.RWMutex

	blocks  map[string]int
	blockMu sync.RWMutex
}

//LimitHandler :
func LimitHandler(allowCount int, clearDuration time.Duration, blackCount int) HandlerFunc {
	if blackCount < allowCount {
		blackCount = allowCount * 2
	}
	manager := &limitManager{
		visitors:      map[string]*hVisitor{},
		blocks:        map[string]int{},
		allowCount:    allowCount,
		clearDuration: clearDuration,
		blackCount:    blackCount,
	}
	console.SetCmd(console.Commands{
		{
			Cmd: "ddos",
			Help: `
ddos view          ----  ( caching / block visitors )
ddos block.clear   ----  ( block visitors clear )
`,
			Work: func(done chan<- bool, ps []string) {
				defer console.DoneC(done)
				defer console.Atap()
				console.Atap()

				switch ps[0] {
				case "view":
					info := manager.Info()
					console.Log(info.ToString())

				case "block.clear":
					manager.clearBlockers()
					console.Log("clear blockers success.")
				} //switch

			},
		},
	})
	return manager.handler
}

type limitInfo struct {
	VCnt    int
	Blocker map[string]int
}

func (my limitInfo) ToString() string {
	return dbg.ToJSONString(my)
}

func (my *limitManager) Info() limitInfo {
	info := limitInfo{
		Blocker: map[string]int{},
	}
	my.mu.RLock()
	info.VCnt = len(my.visitors)
	my.mu.RUnlock()

	my.blockMu.RLock()
	for k, v := range my.blocks {
		info.Blocker[k] = v
	}
	my.blockMu.RUnlock()
	return info
}

func (my *limitManager) clearBlockers() {
	defer my.blockMu.Unlock()
	my.blockMu.Lock()
	my.blocks = map[string]int{}
}

func (my *limitManager) getVisitor(ip string) *hVisitor {
	my.mu.RLock()
	if visitor, do := my.visitors[ip]; do {
		my.mu.RUnlock()
		return visitor
	}
	my.mu.RUnlock()

	my.mu.Lock()
	newVisitor := &hVisitor{
		manager:  my,
		ip:       ip,
		lastSeen: time.Now(),
	}
	my.visitors[ip] = newVisitor
	go newVisitor.checkClear()
	my.mu.Unlock()

	return newVisitor
}

func (my *limitManager) cleanUp(ip string, callCnt int) {
	defer my.mu.Unlock()
	my.mu.Lock()
	delete(my.visitors, ip)

	if callCnt >= my.blackCount {
		defer my.blockMu.Unlock()
		my.blockMu.Lock()
		my.blocks[ip] = callCnt
	}
}
func (my *limitManager) checkBlockVisitor(ip string) bool {
	defer my.blockMu.RUnlock()
	my.blockMu.RLock()
	if v, do := my.blocks[ip]; do == true {
		fmt.Println("DDOS Attacker IP :", v)
		return true
	}
	return false
}

func (my *limitManager) handler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if my.checkBlockVisitor(ip) {
		http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
		return
	}

	visitor := my.getVisitor(ip)
	if visitor.next() == false {
		http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
		return
	}
	next(w, r)
}

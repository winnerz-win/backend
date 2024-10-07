package console

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
)

const (
	notSupportd = "This command is not supported."
	baseIntag   = ">"
)

type cConsole struct {
	log_func func(a ...interface{})
	err_func func(a ...interface{})

	titleName string
	isStart   bool
	mu        sync.RWMutex

	cmds CMDList
}

func (my *cConsole) run() {
	my.initHelp()
	reader := bufio.NewReader(os.Stdin)
	for {

		//fmt.Print(fmt.Sprintf("%v%v", my.titleName, baseIntag))
		fmt.Print(my.titleName, baseIntag)

		line, _ := reader.ReadString('\n')
		buf := []byte(line)
		switch runtime.GOOS {
		case "windows":
			if len(buf) < 2 {
				continue
			}
			buf = buf[:len(buf)-2]
		default:
			if len(buf) < 1 {
				continue
			}
			buf = buf[:len(buf)-1]
		} //switch
		line = strings.TrimSpace(string(buf))
		name := strings.Split(line, " ")[0]
		if v, do := checkHelp(line); do {
			viewMainHelp(my, v)
			continue
		}
		line = strings.TrimSpace(line[len(name):])

		ok := parseCmd(
			my,
			my.cmds,
			name,
			line,
		)
		if !ok {
			if line != "" {
				my.Error(notSupportd)
			}
		}
	} //for
}

func parseCmd(my *cConsole, cmds CMDList, name, line string) bool {
	defer func() {
		if e := recover(); e != nil {
			my.Error("------ console.Panic -------")
			my.Error(e)
		}
	}()

	ok := false
	isHelp := true
	for _, cmd := range cmds {
		if cmd.Name == name {
			if _, do := checkHelp(line); do {
				viewCmdHelp(my, cmd, true)
				return true
			}

			sublist := cmd.CMDS
			if len(sublist) == 0 {
				ps := strings.Split(line, " ")
				ps = DIV(
					cmd.Pair,
					ps,
					cmd.PairWordCut,
				)
				cmd.Action(ps)
				ok = true
			} else {
				name = strings.Split(line, " ")[0]
				line = strings.TrimSpace(line[len(name):])
				ok = parseCmd(
					my,
					sublist,
					name,
					line,
				)
			}
			isHelp = false
			break
		}
	} //for
	if isHelp {
		//cc.Yellow("help")
	}
	return ok
}

func DIV(pair []string, ps []string, isCut ...bool) []string {
	sl := []string{}
	ss := ps
	spot := -1
	depth := 0
	for i, v := range ss {

		if len(pair) > 0 {
			if strings.HasPrefix(v, pair[0]) {
				depth++
				if spot == -1 {
					spot = i
				}
			}
			if strings.HasSuffix(v, pair[1]) {
				depth--
				if depth == 0 && spot != -1 {
					tt := ss[spot : i+1]
					word := strings.Join(tt, " ")
					if len(isCut) > 0 && isCut[0] {
						word = word[1:]
						word = word[:len(word)-1]
					}
					if word != "" {
						sl = append(sl, word)
					}
					depth = 0
					spot = -1
					continue
				}
			}
		}

		if spot == -1 {
			if v != "" {
				sl = append(sl, v)
			}

		}
	}
	return sl
}

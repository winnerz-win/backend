package console

import (
	"fmt"
	"strings"
)

const (
	ob = "<" //open block
	cb = ">" //close block
	bb = "<>"
)

// getArgs : [] <--- text block
func getArgs(cmd ICmd, ps []string) ([]string, error) {
	args := []string{}
	comments := ""

	_tagOpen := ""
	_tagClose := ""
	customCount := 0
	customFlag := false
	if argTags := cmd.ArgTags(); len(argTags) > 1 {
		_tagOpen = argTags[0]
		_tagClose = argTags[1]
		customFlag = true
	}

	isCmt := false
	blockCount := 0
	_, _ = isCmt, blockCount

	for _, v := range ps {
		v = strings.Trim(v, " ")
		if v == "" {
			continue
		}

		if customFlag && blockCount == 0 {
			if customCount == 0 {
				if strings.HasPrefix(v, _tagOpen) == true {
					customCount = 1
					if strings.HasSuffix(v, _tagClose) == true {
						customCount = 0
						args = append(args, strings.Trim(v, " "))
					} else {
						comments = v
					}
					continue
				} else {
					// args = append(args, v)
					// fmt.Println("bc", v)
					// continue
				}
			} else {
				customCount += strings.Count(v, _tagOpen)
				customCount -= strings.Count(v, _tagClose)
				if customCount <= 0 {
					comments = fmt.Sprintf("%v %v", comments, v)
					args = append(args, comments)
					customCount = 0
					//fmt.Println("bc", comments)

				} else {
					comments = fmt.Sprintf("%v %v", comments, v)
				}
				continue
			}
		} //if

		if customCount == 0 {
			if blockCount == 0 {
				if strings.HasPrefix(v, ob) == true {
					if strings.HasSuffix(v, cb) == true {
						v = v[1:]        //cut [
						v = v[:len(v)-1] //cut ]
						if tr := strings.Trim(v, " "); tr != "" {
							args = append(args, v)
						}
					} else {
						v = v[1:]
						comments = v
						blockCount = 1
					}
				} else {
					args = append(args, v)
					//fmt.Println("cc", v)
					continue
				}
			} else {
				blockCount += strings.Count(v, ob)
				blockCount -= strings.Count(v, cb)
				if blockCount <= 0 {
					v = v[:len(v)-1]
					comments = fmt.Sprintf("%v %v", comments, v)
					if tr := strings.Trim(comments, " "); tr != "" {
						args = append(args, comments)
						//fmt.Println("cc", comments)
					}
					blockCount = 0
					continue
				} else {
					comments = fmt.Sprintf("%v %v", comments, v)
				}
			}
		} //if

	} //for

	if blockCount > 0 {
		return nil, fmt.Errorf("blockCount is error : %v", blockCount)
	}
	if customFlag && customCount > 0 {
		return nil, fmt.Errorf("customCount is error : %v", customCount)
	}

	// fmt.Println("-----------------------------------------------------------------------")
	// for _, v := range args {
	// 	fmt.Print(v, " ◆ ")
	// }
	// fmt.Println()
	// fmt.Println("-----------------------------------------------------------------------")

	return args, nil
}

package dbg

import (
	"errors"
	"fmt"
	"strings"
)

//MakeError :
func MakeError(a ...interface{}) error {
	msg := "[dbg.MakeError]"
	for i, c := range a {
		if i == 0 {
			msg = fmt.Sprintf("%v", c)
		} else {
			msg = fmt.Sprintf("%v : %v", msg, c)
		}
	} //for
	return errors.New(msg)
}

func Error(a ...interface{}) error {
	msg := []string{}
	for _, v := range a {
		msg = append(msg, fmt.Sprintf("%v", v))
	} //for
	return errors.New(strings.Join(msg, " "))
}

package jargs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type cArgs struct {
	startPath string
	items     []string
}

type ARGS interface {
	RootPath() string
	Next(key string, f func(val string)) bool
	Do(key string) bool
	Split(key_sep, sep string, f func(val string)) bool
}

func (my cArgs) String() string {
	msg := fmt.Sprintln("< cArgs >")
	msg += fmt.Sprintln("startPath :", my.startPath)
	for i, v := range my.items {
		msg += fmt.Sprintln("[", i, "] ", v)
	}
	return msg
}

func New() ARGS {
	args := cArgs{
		startPath: os.Args[0],
		items:     os.Args[1:],
	}
	return args
}

func Test(sl ...any) ARGS {
	items := []string{}
	for _, v := range sl {
		items = append(items, fmt.Sprintf("%v", v))
	}

	args := cArgs{
		startPath: "test",
		items:     items,
	}
	return args
}

func _trim_toLower(v string) string {
	return strings.ToLower(strings.TrimSpace(v))
}

func (my cArgs) RootPath() string {
	dir, _ := filepath.Abs(filepath.Dir(my.startPath))
	return dir
}

func (my cArgs) Next(key string, f func(val string)) bool {
	key = _trim_toLower(key)
	size := len(my.items)
	for i, v := range my.items {
		if _trim_toLower(v) == key {
			if i+1 < size {
				f(my.items[i+1])
			}
			return true
		}
	}
	return false
}

func (my cArgs) Do(key string) bool {
	key = _trim_toLower(key)
	for _, v := range my.items {
		if _trim_toLower(v) == key {
			return true
		}
	}

	return false
}

func (my cArgs) Split(key_sep, sep string, f func(val string)) bool {
	key_sep = _trim_toLower(key_sep) + sep
	for _, v := range my.items {
		lo := _trim_toLower(v)
		if strings.Contains(lo, key_sep) {
			cut_size := len(key_sep)
			val := v[cut_size:]
			f(val)
			return true
		}
	} //for
	return false
}

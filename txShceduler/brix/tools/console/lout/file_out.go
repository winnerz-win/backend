package lout

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	Tag = "[F]"

	outDir = "console_log"
)

var (
	//IsFileLog :
	isFileLog bool

	logTitle   string
	logMessage string
)

func replaceFileFilter(str string) string {
	flist := map[string]string{
		"\\": "",
		"/":  "",
		":":  ";",
		"*":  "",
		"?":  "",
		"\"": "",
		"<":  "",
		">":  "",
		"|":  "",
	}
	for k, v := range flist {
		str = strings.ReplaceAll(str, k, v)
	}
	str = strings.TrimSpace(str)
	return str
}

//Set :
func Set(on bool) {
	isFileLog = on
	if on {
		os.MkdirAll(outDir, os.ModeDir)
	}
}

//Do :
func Do() bool {
	return isFileLog
}

//RemoteWrite :
func RemoteWrite(title, msg string) {
	if isFileLog == false {
		return
	}

	fileName := fmt.Sprintf(outDir+"/console_%v.txt", replaceFileFilter(title))
	ioutil.WriteFile(fileName, []byte(msg), os.ModePerm)
	fmt.Println("write.file :", fileName)
}

//Write :
func Write() {
	if isFileLog == false {
		return
	}
	fileName := fmt.Sprintf(outDir+"/console_%v.txt", logTitle)
	ioutil.WriteFile(fileName, []byte(logMessage), os.ModePerm)
	fmt.Println("write.file :", fileName)

	logTitle = ""
	logMessage = ""
}

//Clear :
func Clear() {
	logTitle = ""
	logMessage = ""
}

//Title :
func Title(line string) {
	if isFileLog == false {
		return
	}
	logTitle = replaceFileFilter(line)
}

//Log :
func Log(a ...interface{}) {
	if isFileLog == false {
		return
	}
	for _, v := range a {
		logMessage += fmt.Sprintf("%v", v)
	}
	logMessage += "\n"
}

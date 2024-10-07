package jpath

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//NowPath :
func NowPath() string {
	if runtime.GOOS != "windows" {
		return NowPathLinux()
	}

	ex, err := os.Getwd() //C:\Users\[UserName]\AppData\Local\Temp\
	if err != nil {
		fmt.Println("NowPath Error :", err)
	}
	//dbg.Purple(ex)
	//exPath := filepath.Dir(ex)

	//dbg.Blue(os.Getwd())

	return ex
}

func NowPathLinux() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir

}

//MakeFolder :
func MakeFolder(fullPath string) {
	os.MkdirAll(fullPath, os.ModeDir)
}

//dirPath :
type dirPath struct {
	pathlist []string
}

//NewDirPath :
func NewDirPath(path string) *dirPath {
	return &dirPath{strings.Split(path, "\\")}
}

//NowDirPath :
func NowDirPath() *dirPath {
	path := NowPath()
	return NewDirPath(path)
}

//ToString :
func (my dirPath) ToString() string {
	path := ""
	for _, v := range my.pathlist {
		path += v + "\\"
	}
	return path[:len(path)-1]
}

//String :
func (my dirPath) String() string {
	return my.ToString()
}

//Up :
func (my dirPath) Up(depth int) *dirPath {
	if depth >= len(my.pathlist) {
		return NewDirPath(my.pathlist[0])
	}

	counter := len(my.pathlist) - depth
	path := ""
	for _, v := range my.pathlist {
		path += v + "\\"
		counter--
		if counter == 0 {
			break
		}
	}
	return NewDirPath(path[:len(path)-1])
}

//Add :
func (my *dirPath) Add(add string) *dirPath {
	newPath := NewDirPath(my.String())
	newPath.pathlist = append(newPath.pathlist, add)
	return newPath
}

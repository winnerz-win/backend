package jpath

import (
	"os"
	"path/filepath"
)

func NowPath() string {

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir

	// if runtime.GOOS != "windows" {
	// 	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	// 	return dir
	// }
	// ex, _ := os.Getwd()
	// return ex
}

func IsPathExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

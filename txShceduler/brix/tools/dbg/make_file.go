package dbg

import (
	"fmt"
	"os"
)

//FileWriteMakeDir : path , name(filename.ext) , []byte
func FileWriteMakeDir(path, name string, buf []byte) error {
	os.MkdirAll(path, os.ModeDir)
	fileName := fmt.Sprintf("%v\\%v", path, name)
	fp, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer fp.Close()
	if _, err := fp.Write(buf); err != nil {
		return err
	}
	return nil
}

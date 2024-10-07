package jzip

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"io"
)

//ZIP :
func ZIP(name string, data []byte) ([]byte, error) {
	compressBuf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(compressBuf)
	f, err := zipWriter.Create(name)
	if err != nil {
		return nil, err
	}
	f.Write([]byte(data))
	if err := zipWriter.Close(); err != nil {
		return nil, err
	}
	return compressBuf.Bytes(), nil
}

//UNZIP :
func UNZIP(compBuf []byte) ([]byte, error) {
	readBuf := new(bytes.Buffer)
	r, err := zip.NewReader(bytes.NewReader(compBuf), int64(len(compBuf)))
	if err != nil {
		return nil, err
	}
	for _, f := range r.File {
		reader, err := f.Open()
		if err != nil {
			return nil, err
		}
		io.Copy(readBuf, reader)
		reader.Close()
	} //for
	return readBuf.Bytes(), nil
}

//ZIP64 :
func ZIP64(name string, data []byte) (string, error) {
	if cmp, err := ZIP(name, data); err != nil {
		return "", err
	} else {
		return base64.StdEncoding.EncodeToString(cmp), nil
	}
}

//UNZIP64 :
func UNZIP64(b64 string) ([]byte, error) {
	if cmp, err := base64.StdEncoding.DecodeString(b64); err != nil {
		return nil, err
	} else {
		return UNZIP(cmp)
	}
}

//ZIPHex :
func ZIPHex(name string, data []byte) (string, error) {
	if cmp, err := ZIP(name, data); err != nil {
		return "", err
	} else {
		return hex.EncodeToString(cmp), nil
	}
}

//UNZIPHex :
func UNZIPHex(hexString string) ([]byte, error) {
	if cmp, err := hex.DecodeString(hexString); err != nil {
		return nil, err
	} else {
		return UNZIP(cmp)
	}
}

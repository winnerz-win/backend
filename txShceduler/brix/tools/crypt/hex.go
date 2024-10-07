package crypt

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"log"
	"time"
)

// CurrentTime local time
func CurrentTime() int64 {
	return time.Now().Unix()
}

// CurrentTimeUTC UTC time
func CurrentTimeUTC() int64 {
	// loc, _ := time.LoadLocation("UTC")
	// return time.Now().In(loc).Unix()
	return time.Now().UTC().Unix()
}

// Bytes2Hex returns the hexadecimal encoding of d.
func Bytes2Hex(d []byte) string {
	return hex.EncodeToString(d)
}

// Hex2Bytes returns the bytes represented by the hexadecimal string str.
func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

// IntToHex converts an int64 to a byte array
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// ReverseBytes reverses a byte array
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

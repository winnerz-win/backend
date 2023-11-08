package ebcmx

import (
	"encoding/hex"
	"strings"

	"golang.org/x/crypto/sha3"
)

/*
	https://steemit.com/kr/@anpigon/ethereum-1
*/

type elem55 struct {
	val   int
	text  string
	upper string
}

var (
	eip55map = map[byte]elem55{
		'0': {0, "0", "0"},
		'1': {1, "1", "1"},
		'2': {2, "2", "2"},
		'3': {3, "3", "3"},
		'4': {4, "4", "4"},
		'5': {5, "5", "5"},
		'6': {6, "6", "6"},
		'7': {7, "7", "7"},
		'8': {8, "8", "8"},
		'9': {9, "9", "9"},
		'a': {10, "a", "A"},
		'b': {11, "b", "B"},
		'c': {12, "c", "C"},
		'd': {13, "d", "D"},
		'e': {14, "e", "E"},
		'f': {15, "f", "F"},
	}
)

func EIP55(address string) string {
	address = strings.ToLower(address)
	address = strings.TrimPrefix(address, "0x")
	if len(address) != 40 {
		return "" //invalid address format
	}

	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(address))
	hashText := hex.EncodeToString(hash.Sum(nil))

	eip55Address := ""
	for i := range address {
		addressChar := address[i]
		if elem, do := eip55map[addressChar]; !do {
			return "" //invalid address format
		} else {
			if elem.val >= 10 { //a,b,c,d,e,f
				hashChar := hashText[i]
				cv := eip55map[hashChar]
				if cv.val >= 8 {
					eip55Address += elem.upper
					continue
				}
			}
			eip55Address += elem.text
		}

	}
	return "0x" + eip55Address
}

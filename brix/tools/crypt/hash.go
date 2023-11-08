package crypt

import (
	"crypto/sha256"
)

const (
	HashLength = 32
)

// Hash : pw를 db에 저장할때 씀
type Hash [HashLength]byte

// ToString :
func (h Hash) ToString() string {
	return Bytes2Hex(h[:])
}

// FromString :
func (h *Hash) FromString(hex string) {
	data := Hex2Bytes(hex)
	copy(h[:], data[:HashLength])
}

// GetPwHash :
func GetPwHash(pw string) string {
	hash := DoubleHash([]byte(pw))
	return hash.ToString()
}

// DoubleHash :
func DoubleHash(data []byte) Hash {
	once := sha256.Sum256(data)
	ReverseBytes(once[:])
	return Hash(sha256.Sum256(once[:]))
}

// MultiHash :
func MultiHash(data []byte, count int) Hash {

	if count <= 0 {
		count = 1
	}

	var h Hash
	s := sha256.Sum256(data)
	ReverseBytes(s[:])
	copy(h[:], s[:HashLength])
	//fmt.Println(h)

	for i := 1; i < count; i++ {
		h = sha256.Sum256(h[:])
		ReverseBytes(h[:])
		//fmt.Println(h)
	} //for

	return h
}

package abi

import (
	"encoding/hex"
	"strings"

	"golang.org/x/crypto/sha3"
)

const (
	// AddressLength is the expected length of the address
	_AddressLength = 20
)

type addressHash [_AddressLength]byte

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (a *addressHash) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-_AddressLength:]
	}
	copy(a[_AddressLength-len(b):], b)
}
func (a addressHash) Hex() string { return string(a.checksumHex()) }

func (a *addressHash) checksumHex() []byte {
	buf := a.hex()

	// compute checksum
	sha := sha3.NewLegacyKeccak256()
	sha.Write(buf[2:])
	hash := sha.Sum(nil)
	for i := 2; i < len(buf); i++ {
		hashByte := hash[(i-2)/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if buf[i] > '9' && hashByte > 7 {
			buf[i] -= 32
		}
	}
	return buf[:]
}

func (a addressHash) hex() []byte {
	var buf [len(a)*2 + 2]byte
	copy(buf[:2], "0x")
	hex.Encode(buf[2:], a[:])
	return buf[:]
}

///////////////////////////////

type byteToAddressHexer struct{}

func (my byteToAddressHexer) BytesToAddressHex(data32 []byte) string {
	var a addressHash
	a.SetBytes(data32)
	return strings.ToLower(a.Hex())
}

func ByteToAddressHexer() byteToAddressHexer {
	return byteToAddressHexer{}
}

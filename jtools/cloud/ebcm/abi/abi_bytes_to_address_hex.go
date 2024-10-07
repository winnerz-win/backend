package abi

import (
	"encoding/hex"
	"hash"
	"jtools/cloud/ebcm/abi/rlp"
	"strings"

	"golang.org/x/crypto/sha3"
)

const (
	// AddressLength is the expected length of the address
	_AddressLength = 20
)

type AddressHASH [_AddressLength]byte

// BytesToAddress returns Address with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToAddress(b []byte) AddressHASH {
	var a AddressHASH
	a.SetBytes(b)
	return a
}

// has0xPrefix validates str begins with '0x' or '0X'.
func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// Hex2Bytes returns the bytes represented by the hexadecimal string str.
func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

// FromHex returns the bytes represented by the hexadecimal string s.
// s may be prefixed with "0x".
func FromHex(s string) []byte {
	if has0xPrefix(s) {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}

// HexToAddress returns Address with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToAddress(s string) AddressHASH { return BytesToAddress(FromHex(s)) }

// Bytes gets the string representation of the underlying address.
func (a AddressHASH) Bytes() []byte { return a[:] }

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (a *AddressHASH) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-_AddressLength:]
	}
	copy(a[_AddressLength-len(b):], b)
}
func (a AddressHASH) Hex() string {
	return string(a.checksumHex())
}

// String implements fmt.Stringer.
func (a AddressHASH) String() string {
	return a.Hex()
}

func (a *AddressHASH) checksumHex() []byte {
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

func (a AddressHASH) hex() []byte {
	var buf [len(a)*2 + 2]byte
	copy(buf[:2], "0x")
	hex.Encode(buf[2:], a[:])
	return buf[:]
}

////////////////////////////////////////////////////////
// byteToAddressHexer
////////////////////////////////////////////////////////

type byteToAddressHexer struct{}

func (my byteToAddressHexer) BytesToAddressHex(data32 []byte) string {
	var a AddressHASH
	a.SetBytes(data32)
	return strings.ToLower(a.Hex())
}

func (my byteToAddressHexer) ContractAddressNonce(from string, nonce uint64) string {
	return ContractAddressNonce(from, nonce)
}

func GetInputParser() InputParser {
	return byteToAddressHexer{}
}

////////////////////////////////////////////////////////

func ByteToAddressHexer() byteToAddressHexer {
	return byteToAddressHexer{}
}

func ContractAddressNonce(from string, nonce uint64) string {
	if !strings.HasPrefix(from, "0x") {
		return ""
	} else {
		if EIP55(from) == "" {
			return ""
		}
	}

	v := CreateAddress(
		HexToAddress(from),
		nonce,
	)
	return strings.ToLower(v.String())
}

/////////////////////////////////////////////////////////////////////////////////////////////////////

// KeccakState wraps sha3.state. In addition to the usual hash methods, it also supports
// Read to get a variable amount of data from the hash state. Read is faster than Sum
// because it doesn't copy the internal state, but also modifies the internal state.
type KeccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

// NewKeccakState creates a new KeccakState
func NewKeccakState() KeccakState {
	return sha3.NewLegacyKeccak256().(KeccakState)
}

// Keccak256Hash calculates and returns the Keccak256 hash of the input data,
// converting it to an internal Hash data structure.
func Keccak256Hash(data ...[]byte) (h AddressHASH) {
	d := NewKeccakState()
	for _, b := range data {
		d.Write(b)
	}
	d.Read(h[:])
	return h
}

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	b := make([]byte, 32)
	d := NewKeccakState()
	for _, b := range data {
		d.Write(b)
	}
	d.Read(b)
	return b
}

// CreateAddress creates an ethereum address given the bytes and the nonce
func CreateAddress(b AddressHASH, nonce uint64) AddressHASH {
	data, _ := rlp.EncodeToBytes([]interface{}{b, nonce})
	return BytesToAddress(Keccak256(data)[12:])
}

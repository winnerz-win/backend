package ebcm

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
)

const hexutil_uintBits = 32 << (uint64(^uint(0)) >> 63)

// Errors
var (
	ErrEmptyString   = &hexutil_decError{"empty hex string"}
	ErrSyntax        = &hexutil_decError{"invalid hex string"}
	ErrMissingPrefix = &hexutil_decError{"hex string without 0x prefix"}
	ErrOddLength     = &hexutil_decError{"hex string of odd length"}
	ErrEmptyNumber   = &hexutil_decError{"hex string \"0x\""}
	ErrLeadingZero   = &hexutil_decError{"hex number with leading zero digits"}
	ErrUint64Range   = &hexutil_decError{"hex number > 64 bits"}
	ErrUintRange     = &hexutil_decError{fmt.Sprintf("hex number > %d bits", hexutil_uintBits)}
	ErrBig256Range   = &hexutil_decError{"hex number > 256 bits"}
)

type hexutil_decError struct{ msg string }

func (err hexutil_decError) Error() string { return err.msg }

// Hexutil_Decode decodes a hex string with 0x prefix.
func Hexutil_Decode(input string) ([]byte, error) {
	if len(input) == 0 {
		return nil, ErrEmptyString
	}
	if !hexutil_has0xPrefix(input) {
		return nil, ErrMissingPrefix
	}
	b, err := hex.DecodeString(input[2:])
	if err != nil {
		err = hexutil_mapError(err)
	}
	return b, err
}

// Hexutil_MustDecode decodes a hex string with 0x prefix. It panics for invalid input.
func Hexutil_MustDecode(input string) []byte {
	dec, err := Hexutil_Decode(input)
	if err != nil {
		panic(err)
	}
	return dec
}

// Hexutil_Encode encodes b as a hex string with 0x prefix.
func Hexutil_Encode(b []byte) string {
	enc := make([]byte, len(b)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], b)
	return string(enc)
}

// Hexutil_DecodeUint64 decodes a hex string with 0x prefix as a quantity.
func Hexutil_DecodeUint64(input string) (uint64, error) {
	raw, err := hexutil_checkNumber(input)
	if err != nil {
		return 0, err
	}
	dec, err := strconv.ParseUint(raw, 16, 64)
	if err != nil {
		err = hexutil_mapError(err)
	}
	return dec, err
}

// Hexutil_MustDecodeUint64 decodes a hex string with 0x prefix as a quantity.
// It panics for invalid input.
func Hexutil_MustDecodeUint64(input string) uint64 {
	dec, err := Hexutil_DecodeUint64(input)
	if err != nil {
		panic(err)
	}
	return dec
}

// Hexutil_EncodeUint64 encodes i as a hex string with 0x prefix.
func Hexutil_EncodeUint64(i uint64) string {
	enc := make([]byte, 2, 10)
	copy(enc, "0x")
	return string(strconv.AppendUint(enc, i, 16))
}

var bigWordNibbles int

func init() {
	// This is a weird way to compute the number of nibbles required for big.Word.
	// The usual way would be to use constant arithmetic but go vet can't handle that.
	b, _ := new(big.Int).SetString("FFFFFFFFFF", 16)
	switch len(b.Bits()) {
	case 1:
		bigWordNibbles = 16
	case 2:
		bigWordNibbles = 8
	default:
		panic("weird big.Word size")
	}
}

// Hexutil_DecodeBig decodes a hex string with 0x prefix as a quantity.
// Numbers larger than 256 bits are not accepted.
func Hexutil_DecodeBig(input string) (*big.Int, error) {
	raw, err := hexutil_checkNumber(input)
	if err != nil {
		return nil, err
	}
	if len(raw) > 64 {
		return nil, ErrBig256Range
	}
	words := make([]big.Word, len(raw)/bigWordNibbles+1)
	end := len(raw)
	for i := range words {
		start := end - bigWordNibbles
		if start < 0 {
			start = 0
		}
		for ri := start; ri < end; ri++ {
			nib := hexutil_decodeNibble(raw[ri])
			if nib == hexutil_badNibble {
				return nil, ErrSyntax
			}
			words[i] *= 16
			words[i] += big.Word(nib)
		}
		end = start
	}
	dec := new(big.Int).SetBits(words)
	return dec, nil
}

// Hexutil_MustDecodeBig decodes a hex string with 0x prefix as a quantity.
// It panics for invalid input.
func Hexutil_MustDecodeBig(input string) *big.Int {
	dec, err := Hexutil_DecodeBig(input)
	if err != nil {
		panic(err)
	}
	return dec
}

// Hexutil_EncodeBig encodes bigint as a hex string with 0x prefix.
// The sign of the integer is ignored.
func Hexutil_EncodeBig(bigint *big.Int) string {
	nbits := bigint.BitLen()
	if nbits == 0 {
		return "0x0"
	}
	return fmt.Sprintf("%#x", bigint)
}

func hexutil_has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

func hexutil_checkNumber(input string) (raw string, err error) {
	if len(input) == 0 {
		return "", ErrEmptyString
	}
	if !hexutil_has0xPrefix(input) {
		return "", ErrMissingPrefix
	}
	input = input[2:]
	if len(input) == 0 {
		return "", ErrEmptyNumber
	}
	if len(input) > 1 && input[0] == '0' {
		return "", ErrLeadingZero
	}
	return input, nil
}

const hexutil_badNibble = ^uint64(0)

func hexutil_decodeNibble(in byte) uint64 {
	switch {
	case in >= '0' && in <= '9':
		return uint64(in - '0')
	case in >= 'A' && in <= 'F':
		return uint64(in - 'A' + 10)
	case in >= 'a' && in <= 'f':
		return uint64(in - 'a' + 10)
	default:
		return hexutil_badNibble
	}
}

func hexutil_mapError(err error) error {
	if err, ok := err.(*strconv.NumError); ok {
		switch err.Err {
		case strconv.ErrRange:
			return ErrUint64Range
		case strconv.ErrSyntax:
			return ErrSyntax
		}
	}
	if _, ok := err.(hex.InvalidByteError); ok {
		return ErrSyntax
	}
	if err == hex.ErrLength {
		return ErrOddLength
	}
	return err
}

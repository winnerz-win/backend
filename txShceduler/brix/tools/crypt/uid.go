package crypt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

//MakeInt64 :
func MakeInt64() int64 {
	v := int64(MakeUInt64())
	if v < 0 {
		v *= -1
	}
	return v
}

//MakeUInt32 :
func MakeUInt32() uint32 {
	rndBuf := make([]byte, 256)
	rand.Read(rndBuf)
	h := sha256.Sum256(rndBuf)
	return binary.BigEndian.Uint32(h[:])
}

//MakeUInt64 :
func MakeUInt64() uint64 {
	rndBuf := make([]byte, 256)
	rand.Read(rndBuf)
	h := sha256.Sum256(rndBuf)
	return binary.BigEndian.Uint64(h[:])
}

//MakeUInt64String :
func MakeUInt64String() string {
	return fmt.Sprintf("%v", MakeUInt64())
}

// MakeUID :
func MakeUID(length int) string {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	h := sha256.Sum256(randomBytes)
	return hex.EncodeToString(h[:])
	// uid := base64.StdEncoding.EncodeToString(randomBytes)
	// uid = strings.Replace(uid, "+", "^", -1)
	// return uid
}

// MakeUID256 :
func MakeUID256() string {
	return MakeUID(256)
}

// MakeUID128 :
func MakeUID128() string {
	return MakeUID(128)
}

// MakeUID64 :
func MakeUID64() string {
	return MakeUID(64)
}

//MakeUID32 :
func MakeUID32() string {
	return MakeUID(32)
}

func MakeUIDString(text string) string {
	h := sha256.Sum256([]byte(text))
	v := hex.EncodeToString(h[:])
	return v
}

// MakePin :
func MakePin(clen ...int) string {
	size := 6
	if len(clen) > 0 && clen[0] > 0 {
		size = clen[0]
	}
	cnt := size * 2
	rb := make([]byte, cnt)
	_, err := rand.Read(rb)
	if err != nil {
		panic(err)
	}

	pin := ""
	for i := 0; i < cnt; i += 2 {
		pin += fmt.Sprint((rb[i] + rb[i+1]) % 10)
	}
	return pin
	// p1 := (rb[0] + rb[1]) % 10
	// p2 := (rb[2] + rb[3]) % 10
	// p3 := (rb[4] + rb[5]) % 10
	// p4 := (rb[6] + rb[7]) % 10
	// p5 := (rb[8] + rb[9]) % 10
	// p6 := (rb[10] + rb[11]) % 10
	// return fmt.Sprintf("%v%v%v%v%v%v", p1, p2, p3, p4, p5, p6)
}

var mcA0 = []string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J",
	"K", "L", "M", "N", "O", "P", "Q", "R", "S", "T",
	"U", "V", "W", "X", "Y", "Z",
}
var mcA1 = []string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J",
	"K", "L", "M", "N", "O", "P", "Q", "R", "S", "T",
	"U", "V", "W", "X", "Y", "Z", "a", "b", "c", "d",
	"e", "f", "g", "h", "i", "j", "k", "m", "n", "p",
	"q", "r", "x", "t", "u", "v", "w", "x", "y", "z",
	//"l", "o",
}

func makeCode(max int, cs []string) string {
	size := max * 2
	rb := make([]byte, size)
	_, err := rand.Read(rb)
	if err != nil {
		panic(err)
	}
	mod := byte(36)

	code := ""
	for i := 0; i < size; i += 2 {
		code += cs[int((rb[i]+rb[i+1])%mod)]
	}
	return code
}

//MakeCode : [0~9,A~Z] default cnt : 6
func MakeCode(cnt ...int) string {
	max := 6
	if len(cnt) > 0 {
		max = cnt[0]
	}
	return makeCode(max, mcA0)

}

//MakeCode2 : [0~9,A~Z] X(l,o)default cnt : 6
func MakeCode2(cnt ...int) string {
	max := 6
	if len(cnt) > 0 {
		max = cnt[0]
	}
	return makeCode(max, mcA1)
}

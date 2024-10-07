package pwd

import (
	"crypto/sha512"
	"encoding/hex"

	"txscheduler/brix/tools/crypt"
	"txscheduler/brix/tools/dbg"

	"golang.org/x/crypto/pbkdf2"
)

const (
	//IterCount :
	IterCount = 512
	//KeySize :
	KeySize = 64
)

var saltValue = "default_salt_value"

// Salt :
func SaltVale() string { return saltValue }

// InitPWD : default_salt_value
func InitPWD(salt string) {
	saltValue = salt
	dbg.PrintForce("admin-Salt :", saltValue)
}

// Hex : salt( saltValue )
func Hex(pwd string) string {
	dbg.Yellow("pwd.Hex(", pwd, ")")
	dbg.Yellow("saltValue :", saltValue)
	sault := saltValue
	dk := pbkdf2.Key([]byte(pwd), []byte(sault), IterCount, KeySize, sha512.New)
	hexString := hex.EncodeToString(dk)

	return hexString
}

// Salt :
func Salt(salt, pwd string) string {
	dk := pbkdf2.Key([]byte(pwd), []byte(salt), IterCount, KeySize, sha512.New)
	hexString := hex.EncodeToString(dk)

	return hexString
}

// MakeDBPWD : dbPWD , db_salt
func MakeDBPWD(client_pwd string) (string, string) {
	mergePwdSalt := func(client_pwd string) (string, string) {
		db_salt := crypt.MakeUID128()[:20]
		return db_salt + client_pwd, db_salt
	}
	sum_pwd, db_salt := mergePwdSalt(client_pwd)
	dbPwd := Salt(db_salt, sum_pwd)
	return dbPwd, db_salt
}

// VerifyPWD :
func VerifyPWD(db_salt, pass string) string {
	sum_pwd := db_salt + pass
	password := Salt(db_salt, sum_pwd)
	return password
}

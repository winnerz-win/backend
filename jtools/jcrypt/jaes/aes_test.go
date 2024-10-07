package jaes

import (
	"fmt"
	"testing"
)

func Test_OPENSSL_AES(t_ *testing.T) {
	aeskey := New("11")

	enc, _ := aeskey.EncryptBytesString([]byte("hellow world!!"))
	fmt.Println(enc)

	aeskey2 := New("11")
	text, _ := aeskey2.DecryptStringString(enc)
	fmt.Println(text)

}

func Test_DEC(t_ *testing.T) {
	enc := "U2FsdGVkX1+BdzSXQoLg8mLUOyVgymHfm9UGmJ+1lEYZUPmwh21BR0CyuzLq7XJE"
	enc = "U2FsdGVkX1+ykmvGFIzZ31WmXYKBsSbC39nH/9lfHVZAc5geEo6tWX+4uKmOG/zl"
	aeskey := New("test_openssl_key_1234")
	text, _ := aeskey.DecryptStringString(enc)
	fmt.Println(text)
}

func Test_Openssl(t *testing.T) {
	/*
		javascript	- CBC mode
		golang		- CFB mode
		https://stackoverflow.com/questions/55308051/aes-cbc-javascript-cryptojs-encrypt-golang-decrypt
		https://medium.com/@thanat.arp/encrypt-decrypt-aes256-cbc-shell-script-golang-node-js-ffb675a05669
	*/
	// key := "scserver_wallet_!@#_2019ddfsgdfg"
	// b64 := "U2FsdGVkX19ltWj0tM6IdcTsVsqalS2Ym5xyCtXtoAyxFMK5yRAurohsr2jXTrMlZa5UnrU+OAYOg+6nIyhMzrmkUKPFOJtE2sUpTNS5GmE="
	// text := "aes암호화_예제_입니다._~_씨부럴탱탱././"

	// //buf, _ := base64.StdEncoding.DecodeString(b64)

	// o := openssl.New()

	// dec, err := o.DecryptBytes(key, []byte(b64), openssl.DigestMD5Sum)
	// if err != nil {
	// 	fmt.Printf("An error occurred: %s\n", err)
	// }

	// result := string(dec)
	// fmt.Printf("Decrypted text: %s\n", result)

	// if common.BytesCompair([]byte(text), dec) {
	// 	fmt.Println("result is same.")
	// } else {
	// 	fmt.Println("result is differ.")
	// }

	// enc, err := o.EncryptBytes(key, []byte(text), openssl.DigestMD5Sum)
	// if err != nil {
	// 	fmt.Println("Encrypt :", err)
	// }
	// encB64 := string(enc)
	// fmt.Println("[encB64]", encB64)
	// if encB64 == b64 {
	// 	fmt.Println("enc.b64 is same.")
	// } else {
	// 	fmt.Println("enc.b64 is differ.")
	// }
}

package api

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"jtools/cc"
	"jtools/dbg"
	"jtools/unix"
	"net/http"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/jnet/doc"
	"txscheduler/txm/ack"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

const (
	aes256cbc = CBC("---------- AES 256 CBC KEY ----------------")
)

func init() {
	hInfoPrivateKey()
}

type CBC string

func (my CBC) key() []byte {
	hash := sha256.New()
	hash.Write([]byte(string(my)))
	key := hash.Sum(nil)
	return key
}
func (my CBC) hKey() string { return hex.EncodeToString(my.key()) }

func (CBC) pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func (CBC) unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, fmt.Errorf("unpad error. This could happen when incorrect encryption key is used")
	}

	return src[:(length - unpadding)], nil
}

func (my CBC) Encrypt(key []byte, plaintext []byte, is_debug ...bool) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintext = my.pad(plaintext)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	if dbg.IsTrue(is_debug) {
		cc.Gray("encrypt-iv :", hex.EncodeToString(iv))
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (my CBC) Decrypt(key []byte, cryptoText string, is_debug ...bool) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	decodedCryptoText, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	if len(decodedCryptoText) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := decodedCryptoText[:aes.BlockSize]
	if dbg.IsTrue(is_debug) {
		cc.Gray("decrypt-iv :", hex.EncodeToString(iv))
	}
	decodedCryptoText = decodedCryptoText[aes.BlockSize:] //16byte padding

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decodedCryptoText, decodedCryptoText)

	if dbg.IsTrue(is_debug) {
		cc.Gray("decodedCryptoText :", string(decodedCryptoText))
	}

	decryptedText, err := my.unpad(decodedCryptoText)
	if err != nil {
		return "", fmt.Errorf("unpad error. This could happen when incorrect encryption key is used")
	}

	return string(decryptedText), nil
}

/////////////////////////////////////////////////////////////////////////

func hInfoPrivateKey() {
	method := chttp.GET
	url := model.V1 + "/info/pk/name/:args"
	Doc().Comment("[ ★★★ SECRET ★★★ ] 유저 가상계좌의 개인키 요청 (회원ID로 검색)").
		Method(method).
		URLS(
			url,
			":args", "dd237eb2-4cd8-4f39-831e-c396ff3fbc7b(회원ID)",
		).
		Etc(".", `_
			<cc_blue>응답결과</cc_blue>
			{
				"success": true,
				"data": {
					"uid": 1001,
					"address": "0xf811b879b9f4f24b411a92ebd10dfb7e79c4a200",
					"name": "dd237eb2-4cd8-4f39-831e-c396ff3fbc7b",
					"cbc_key": "ndBY5mbA . . . 13WMaBegYt",  //암호화된 유저 지갑주소의 개인키.
					"debug": {	//개발 서버에서만 debug 필드 존재함.
						"secret_key": "160bccf48b93a22f700b2ac4f8dcb1f86983eb9aa432108ab8a4fe546e1dc1cb", //CBC 시크릿키 (hex)
						"private_key": "6637643 . . . A949F63" //실제 유저 개인키 (cbc_key을 Decode하면 해당 값이 나와야 합니다.)
					}
				}
			}
		`).
		JResultOK(chttp.AckFormat{}).
		Apply(doc.Red)

	type DEBUG struct {
		SecretKey  string `json:"secret_key"`
		PrivateKey string `json:"private_key"`
	}
	type RESULT struct {
		UID     int64  `json:"uid"`
		Address string `json:"address"`
		Name    string `json:"name"`
		CBC_Key string `json:"cbc_key"`
		Debug   *DEBUG `json:"debug,omitempty"`
	}

	secret_key_buf := aes256cbc.key()
	secret_key_hex := aes256cbc.hKey()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			name := ps.ByName("args")
			model.Trim(&name)
			if name == "" {
				chttp.Fail(w, ack.NotFoundName)
				return
			}

			model.DB(func(db mongo.DATABASE) {
				member := model.LoadMemberName(db, name)
				if !member.Valid() {
					chttp.Fail(w, ack.NotFoundName)
					return
				}

				privateKeyText := member.PrivateKey()

				cbc_key, err := aes256cbc.Encrypt(
					secret_key_buf,
					[]byte(privateKeyText),
				)
				if err != nil {
					chttp.Fail(w, ack.DBJob, dbg.Cat("aes256cbc.Endrypt Fail :", err))
					return
				}

				result := RESULT{
					UID:     member.UID,
					Address: member.Address,
					Name:    member.Name,
					CBC_Key: cbc_key,
				}

				if !inf.Config().Mainnet {
					debug := DEBUG{
						SecretKey:  secret_key_hex,
						PrivateKey: privateKeyText,
					}
					result.Debug = &debug
				}

				request_host, req_dir, _ := chttp.GetIP(req)

				model.LogCBC.Set(
					"cbc",
					"uid", member.UID,
					"address", member.Address,
					"name", member.Name,
					"cbc_key", cbc_key,
					"req_host", dbg.Cat("[", req_dir, "]", request_host),
					"req_url_host", req.URL.Host,
					"req_time", unix.Now().KST(),
				)

				chttp.OK(w, result)

			})
		},
	)
}

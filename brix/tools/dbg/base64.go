package dbg

import "encoding/base64"

func Base64DecodeString(b64 string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(b64)
}

func Base64Encoding(buf []byte) string {
	return base64.StdEncoding.EncodeToString(buf)
}

package jrsa

import (
	"fmt"
	"testing"
)

func Test_RSA_MYSELF(t_ *testing.T) {
	plainText := "가나다라마바사"
	_ = plainText

	private, public := GenerateKeys(2048)
	fmt.Println(private.ToString())
	fmt.Println(public.ToString())

	c64, _ := public.EncBase64([]byte(plainText))
	decText, _ := private.DecBase64String(c64)
	fmt.Println(decText)
}

func Test_RSA_PEM(t_ *testing.T) {
	private, _ := ToPrivateKeyString(privatePEM)
	public, _ := ToPublicKeyString(publicPEM)
	_ = public

	text, _ := private.DecBase64String(chiperB64)
	fmt.Println(text)

	chiper642, _ := public.EncBase64String(text)
	text2, _ := private.DecBase64String(chiper642)
	fmt.Println(text2)

}

var privatePEM = `-----BEGIN RSA PRIVATE KEY-----            
MIICXQIBAAKBgQDLZL7JkIAxVOLqtEZr+SjJZYqS59ARp3/Z/l5Q6AOWryyrCJI4            
A/4RjFiZ2lN2a08WGgZuYrM+rk7mXmxQdAX1aUTXk9xxZX2YrO9g2TMxpJzrdDD6            
AKbDSqz9TLx2mBuLFKr+dJXX4FSg6GQF9ZAlP/PxuriRMQ5kbnqGPiMCvwIDAQAB            
AoGAMlaw0XouAAeeUbBkbXyxF4dGEK3G1Ve7UNyfwy5pFPYt+/aXGb4DN5ygoRNj            
7L8KR9IRHWjYK/9AD8v2ysKsZmuXNt23ojkKHq5wirSCpO2vb244ApMQAlvZYdoA            
eUaYjMyq9RLBsPXh8yLAMrLGU6Yxsv5evRjrHnStW8LcImECQQDwnypVThxZP5d5            
vxDHV4xFvS0PFucR05TIEtpLvGIGQ8bzol9jz+A4irpIRHnyRfydyJQziS7Svp9E            
sSlWiOaVAkEA2GR8XZPTXTwtotF4dSSqu0sdERUCyShTAK4/vRjGviH0qITGPk+Z            
eCQp2ZV0WsY84/6mMRviiBn92+J7hS5TAwJAASRQOB1pxwalOl+svbVtpfsS1qp+            
KDh/0T89p/RZ5ru1mvxfRYL8BmiqH6OrjHnGjB0ijugMv9VFvja1AoMdzQJBALjo            
1SUZpunq/Iw/NxHS7Vnyi7oHHERMgvD39VtfCqV6WpiOLOEeH+R78o8NmUngUDP7            
bIRWcbMfksAMvsRFm4UCQQC3suGb43vEPOv8dGaIX7DxRCuvbgI+3mRn07bZe//i            
axiQXphcMJArJB4bpxl+Sdq59miHPhHpTIYGN9fVHYr+            
-----END RSA PRIVATE KEY-----`

var publicPEM = `-----BEGIN PUBLIC KEY-----            
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDLZL7JkIAxVOLqtEZr+SjJZYqS            
59ARp3/Z/l5Q6AOWryyrCJI4A/4RjFiZ2lN2a08WGgZuYrM+rk7mXmxQdAX1aUTX            
k9xxZX2YrO9g2TMxpJzrdDD6AKbDSqz9TLx2mBuLFKr+dJXX4FSg6GQF9ZAlP/Px            
uriRMQ5kbnqGPiMCvwIDAQAB            
-----END PUBLIC KEY-----`

var chiperB64 = "ZEZnw8r5a9Oaog2ytibEPkeIcFHvfQdzqymDnlnbDbsOFJ3KVpTo4EYwi1cZ/3eXgsIea7tCa7ABpZQZ3E18ycbJ7WaVZ+uoG3+YU0GN4i1eYXa04j+2OdJr9QQIRIVtFcZoBhZ1ZkqpL6P+W2VvVH0OM2McmWGTUMGP98LbBuw="

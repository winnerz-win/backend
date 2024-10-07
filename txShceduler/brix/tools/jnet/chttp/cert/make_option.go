package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"strings"
	"time"

	"jtools/mms"
	"txscheduler/brix/tools/jpath"
)

/*
	Option {
		Country:            "ko",
		Organization:       "Golang Organization",
		OrganizationalUnit: "Brickstream",
		CommonName:         "Go Https for Web",
		SerialNumberString: GetSerialNumber(),
		StartYMD:           mms.FromYMD(20200101).YMD(),
		EndYMD:             mms.FromYMD(21200101).YMD(),
		IP:                 []string{"127.0.0.1"},
		DNSNames:           []string{},
		FileName:           "ssl",
		RootPath:           jpath.NowPath() + "\\default_ssl",
	}
*/
type Option struct {
	Country            string
	Organization       string
	OrganizationalUnit string
	CommonName         string
	SerialNumberString string
	StartYMD           int
	EndYMD             int
	IP                 []string
	DNSNames           []string
	FileName           string
	RootPath           string
}

func (my Option) getIPList() []net.IP {
	list := []net.IP{}
	for _, ip := range my.IP {
		list = append(list, net.ParseIP(ip))
	}
	return list
}

// Year :
const Year = time.Hour * 24 * 365

// GetSerialNumber :
func GetSerialNumber() string {
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)
	return serialNumber.String()
}

// DefaultOption :
func DefaultOption() *Option {
	createAt := mms.Now()
	return &Option{
		Country:            "ko",
		Organization:       "Golang Organization",
		OrganizationalUnit: "Brickstream",
		CommonName:         "Go Https for Web",
		SerialNumberString: GetSerialNumber(),
		StartYMD:           createAt.YMD(),
		EndYMD:             createAt.Add(Year * 100).YMD(),
		IP:                 []string{"127.0.0.1"},
		DNSNames:           []string{},
		FileName:           "ssl",
		RootPath:           jpath.NowPath() + "\\default_ssl",
	}
}

func (my *Option) SetExpiredYear(year int) {
	createAt := mms.Now()
	my.StartYMD = createAt.YMD()
	my.EndYMD = createAt.Add(Year * time.Duration(year)).YMD()
}
func (my *Option) SetIP(ip string) {
	ip = strings.TrimSpace(ip)
	isAdd := true
	for _, v := range my.IP {
		if v == ip {
			isAdd = false
			break
		}
	}
	if isAdd {
		my.IP = append(my.IP, ip)
	}
}
func (my *Option) SetDNS(dns string) {
	dns = strings.TrimSpace(dns)
	isAdd := true
	for _, v := range my.DNSNames {
		if v == dns {
			isAdd = false
			break
		}
	}
	if isAdd {
		my.DNSNames = append(my.DNSNames, dns)
	}
}
func (my *Option) ClearRootPath() {
	my.RootPath = "."
}

// MakeOption :
func MakeOption(opt *Option) {
	if opt == nil || opt.FileName == "" {
		opt = DefaultOption()
	}

	serialNumber := big.NewInt(0)
	serialNumber, _ = serialNumber.SetString(opt.SerialNumberString, 10)

	subject := pkix.Name{
		Country:            []string{opt.Country},
		Organization:       []string{opt.Organization},
		OrganizationalUnit: []string{opt.OrganizationalUnit},
		CommonName:         opt.CommonName,
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    mms.FromYMD(opt.StartYMD).ToTime(),
		NotAfter:     mms.FromYMD(opt.EndYMD).ToTime(),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  opt.getIPList(),
		DNSNames:     opt.DNSNames,
	}

	os.MkdirAll(opt.RootPath, os.ModeDir)

	pk, _ := rsa.GenerateKey(rand.Reader, 2048)
	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)

	//
	type basicConstraints struct {
		IsCA       bool `asn1:"optional"`
		MaxPathLen int  `asn1:"optional,default:-1"`
	}
	val, _ := asn1.Marshal(basicConstraints{true, -1})
	creq := x509.CertificateRequest{
		Subject:            subject,
		SignatureAlgorithm: x509.SHA512WithRSA,
		ExtraExtensions: []pkix.Extension{
			{
				Id:       asn1.ObjectIdentifier{2, 5, 29, 19},
				Value:    val,
				Critical: true,
			},
		},
		DNSNames: opt.DNSNames,
	}
	csr, _ := x509.CreateCertificateRequest(rand.Reader, &creq, pk)
	csrtOut, _ := os.Create(opt.RootPath + "/" + opt.FileName + "_req.csr")
	pem.Encode(csrtOut, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csr})
	csrtOut.Close()

	certOut, _ := os.Create(opt.RootPath + "/" + opt.FileName + "_cert.pem")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, _ := os.Create(opt.RootPath + "/" + opt.FileName + "_key.pem")
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	keyOut.Close()

	//pem.Encode(nil, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&pk.PublicKey)})

}

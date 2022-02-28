package cert

import (
	"LiteKube/pkg/certificate"
	"crypto/rsa"
	"crypto/x509"
	"time"
)

var defaultCA = certificate.CAConfig{
	OrganizationName: []string{"LITEKUBE"},
	FileName:         "ca",
	SubjectKeyId:     []byte{1, 2, 3, 4, 5},
	NotAfter:         time.Now().AddDate(10, 0, 0),
	KeyLength:        2048,
}

var defaultServer = certificate.ServerConfig{
	OrganizationName: []string{"LITEKUBE"},
	FileName:         "server",
	SubjectKeyId:     []byte{1, 2, 3, 4, 6},
	Hosts:            []string{"localhost", "127.0.0.1"},
	NotAfter:         time.Now().AddDate(10, 0, 0),
}

var defaultClient = certificate.ClientConfig{
	OrganizationName: []string{"LITEKUBE"},
	FileName:         "client",
	SubjectKeyId:     []byte{1, 2, 3, 4, 7},
	NotAfter:         time.Now().AddDate(10, 0, 0),
}

func CreateCACert(foldPath string) error {
	certificate.GenerateCA(foldPath, defaultCA)
	return nil
}

func CreateServerCert(foldPath string, caCert *x509.Certificate, caKey *rsa.PrivateKey) error {
	certificate.GenerateServerCert(foldPath, defaultServer, caCert, caKey)
	return nil
}

func CreateClientCert(foldPath string, caCert *x509.Certificate, caKey *rsa.PrivateKey) error {
	certificate.GenerateClientCert(foldPath, defaultClient, caCert, caKey)
	return nil
}

// return cert_base64, key_base64, true/false
func CreateClientCertBase64(caCert *x509.Certificate, caKey *rsa.PrivateKey) ([]byte, []byte, error) {
	return certificate.GenerateClientCertBase64(defaultClient, caCert, caKey)
}

package certificate

import (
	"LiteKube/pkg/common"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"

	"k8s.io/klog/v2"
)

type CAConfig struct {
	OrganizationName []string
	FileName         string
	SubjectKeyId     []byte
	NotAfter         time.Time
	KeyLength        int
}

type ServerConfig struct {
	OrganizationName []string
	FileName         string
	SubjectKeyId     []byte
	Hosts            []string
	NotAfter         time.Time
}

type ClientConfig struct {
	OrganizationName []string
	FileName         string
	SubjectKeyId     []byte
	NotAfter         time.Time
}

func GenerateCA(foldPath string, config CAConfig) error {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Organization: config.OrganizationName,
		},
		NotBefore:             time.Now(),
		NotAfter:              config.NotAfter,
		SubjectKeyId:          config.SubjectKeyId,
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	privCa, _ := rsa.GenerateKey(rand.Reader, config.KeyLength)
	return CreateCertificateFile(foldPath, config.FileName, ca, privCa, ca, nil)
}

func GenerateServerCert(foldPath string, config ServerConfig, caCert *x509.Certificate, caKey *rsa.PrivateKey) error {
	server := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: config.OrganizationName,
		},
		NotBefore:    time.Now(),
		NotAfter:     config.NotAfter,
		SubjectKeyId: config.SubjectKeyId,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	for _, h := range config.Hosts {
		if ip := net.ParseIP(h); ip != nil {
			server.IPAddresses = append(server.IPAddresses, ip)
		} else {
			server.DNSNames = append(server.DNSNames, h)
		}
	}

	privSer, _ := rsa.GenerateKey(rand.Reader, caKey.N.BitLen())
	return CreateCertificateFile(foldPath, config.FileName, server, privSer, caCert, caKey)
}

func GenerateClientCert(foldPath string, config ClientConfig, caCert *x509.Certificate, caKey *rsa.PrivateKey) error {
	client := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: config.OrganizationName,
		},
		NotBefore:    time.Now(),
		NotAfter:     config.NotAfter,
		SubjectKeyId: config.SubjectKeyId,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	privCli, _ := rsa.GenerateKey(rand.Reader, caKey.N.BitLen())
	return CreateCertificateFile(foldPath, config.FileName, client, privCli, caCert, caKey)
}

func GenerateClientCertBase64(foldPath string, config ClientConfig, caCert *x509.Certificate, caKey *rsa.PrivateKey) ([]byte, []byte, error) {
	client := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: config.OrganizationName,
		},
		NotBefore:    time.Now(),
		NotAfter:     config.NotAfter,
		SubjectKeyId: config.SubjectKeyId,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	privCli, _ := rsa.GenerateKey(rand.Reader, caKey.N.BitLen())
	return CreateCertificateBase64(client, privCli, caCert, caKey)
}

// func CreateCertificateFile(name string, cert *x509.Certificate, key *rsa.PrivateKey, caCert *x509.Certificate, caKey *rsa.PrivateKey) {
// 	// value copy
// 	priv := key

// 	pub := &priv.PublicKey

// 	privPm := priv
// 	if caKey != nil {
// 		privPm = caKey
// 	}

// 	ca_b, err := x509.CreateCertificate(rand.Reader, cert, caCert, pub, privPm)
// 	if err != nil {
// 		log.Println("create failed", err)
// 		return
// 	}
// 	ca_f := name + ".pem"
// 	log.Println("write to pem", ca_f)
// 	var certificate = &pem.Block{Type: "CERTIFICATE",
// 		Headers: map[string]string{},
// 		Bytes:   ca_b}
// 	ca_b64 := pem.EncodeToMemory(certificate)
// 	ioutil.WriteFile(ca_f, ca_b64, 0777)

// 	priv_f := name + ".key"
// 	priv_b := x509.MarshalPKCS1PrivateKey(priv)
// 	log.Println("write to key", priv_f)
// 	ioutil.WriteFile(priv_f, priv_b, 0777)
// 	var privateKey = &pem.Block{Type: "PRIVATE KEY",
// 		Headers: map[string]string{},
// 		Bytes:   priv_b}
// 	priv_b64 := pem.EncodeToMemory(privateKey)
// 	ioutil.WriteFile(priv_f, priv_b64, 0777)
// }

func CreateCertificateFile(foldPath string, name string, cert *x509.Certificate, key *rsa.PrivateKey, caCert *x509.Certificate, caKey *rsa.PrivateKey) error {
	// calculate file-name
	certFileName := name + ".pem"
	privateKeyFileName := name + "-key.pem"
	if !common.IsZero(foldPath) {
		// ensure kubelet-client-cert-config exists
		if err := os.MkdirAll(foldPath, os.ModePerm); err != nil {
			klog.Errorf("fail to create fold path: %s", foldPath)
			return err
		}
		certFileName = filepath.Join(foldPath, certFileName)
		privateKeyFileName = filepath.Join(foldPath, privateKeyFileName)
	}

	certBase64, privateKeyBase64, err := CreateCertificateBase64(cert, key, caCert, caKey)
	if err != nil {
		return err
	}

	// create certificate file
	klog.Infof("create X.509 Certificate: %s", certFileName)
	ioutil.WriteFile(certFileName, certBase64, 0777)

	// create Private-key file
	klog.Infof("create X.509 Certificate key: %s", privateKeyFileName)

	ioutil.WriteFile(privateKeyFileName, privateKeyBase64, 0777)
	return nil
}

func CreateCertificateBase64(cert *x509.Certificate, key *rsa.PrivateKey, caCert *x509.Certificate, caKey *rsa.PrivateKey) ([]byte, []byte, error) {
	// value copy
	privateKey := key

	publicKey := &privateKey.PublicKey

	privateKeyCA := privateKey
	if caKey != nil {
		privateKeyCA = caKey
	}

	// create certificate
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, caCert, publicKey, privateKeyCA)
	if err != nil {
		return nil, nil, err
	}

	var certificate = &pem.Block{
		Type:    "CERTIFICATE",
		Headers: map[string]string{},
		Bytes:   certBytes,
	}

	certBase64 := pem.EncodeToMemory(certificate)

	// create Private-key file
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	var pk = &pem.Block{
		Type:    "PRIVATE KEY",
		Headers: map[string]string{},
		Bytes:   privateKeyBytes,
	}
	privateKeyBase64 := pem.EncodeToMemory(pk)

	return certBase64, privateKeyBase64, nil
}

func GetPrivateKeyLen(PriKey []byte) (int, error) {
	if PriKey == nil {
		return 0, errors.New("input arguments error")
	}

	block, _ := pem.Decode(PriKey)
	if block == nil {
		return 0, fmt.Errorf("RSA private Key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return 0, err
	}

	return priv.N.BitLen(), nil
}

func ReadPrivateKeyFromFile(filepath string) *rsa.PrivateKey {
	if _, err := os.Stat(filepath); err != nil {
		return nil
	}

	key, err := ioutil.ReadFile(filepath)
	if err != nil {
		klog.Errorf("fail to read private-key file: %s, error: %s", filepath, err.Error())
		return nil
	}

	block, _ := pem.Decode(key)
	if block == nil {
		klog.Errorf("fail to decode private-key pem file: %s", filepath)
		return nil
	}
	// parse to pkcs8 format
	// privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	// parse to pkcs1 format

	// parse to public key if file is public-key
	// publicKey, _ := x509.ParsePKIXPublicKey(key)

	privateKey, p_err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if p_err != nil {
		klog.Errorf("fail to Parse private-key pem file: %s, error: %s", filepath, p_err.Error())
		return nil
	}

	return privateKey
}

func ReadCertificateFromFile(filepath string) *x509.Certificate {
	if _, err := os.Stat(filepath); err != nil {
		return nil
	}

	pemBlock, err := ioutil.ReadFile(filepath)
	if err != nil {
		klog.Errorf("fail to read certificate file: %s, error: %s", filepath, err.Error())
		return nil
	}

	block, _ := pem.Decode(pemBlock)
	if block == nil {
		klog.Errorf("fail to decode certificate pem file: %s", filepath)
		return nil
	}

	x509Cert, x_err := x509.ParseCertificate(block.Bytes)
	if x_err != nil {
		klog.Errorf("fail to Parse certificate pem file: %s, error: %s", filepath, x_err.Error())
		return nil
	}
	return x509Cert
}

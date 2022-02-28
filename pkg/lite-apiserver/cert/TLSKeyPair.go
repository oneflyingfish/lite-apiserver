package cert

import (
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/klog/v2"
)

type TLSKeyPair struct {
	certPath string
	keyPath  string
	Valid    bool
}

func NewTLSKeyPair() *TLSKeyPair {
	return &TLSKeyPair{
		certPath: "",
		keyPath:  "",
		Valid:    false,
	}
}

func CreateTLSKeyPair(certPath string, keyPath string) (*TLSKeyPair, error) {
	tlsKeyPair := NewTLSKeyPair()
	err := tlsKeyPair.SetTLSKeyPair(certPath, keyPath)
	return tlsKeyPair, err
}

func (opt *TLSKeyPair) GetTLSKeyPair() (string, string, bool) {
	if opt.Valid {
		return opt.certPath, opt.keyPath, true
	} else {
		return "", "", false
	}
}

// Set ABs path to TLSKeyPair, if there are one pair, opt.Valid==true
func (opt *TLSKeyPair) SetTLSKeyPair(certPath string, keyPath string) error {
	path, err := filepath.Abs(certPath)
	if err != nil {
		klog.Errorf("fail to translate %s to absolute path", certPath)
		return err
	} else {
		opt.certPath = path
	}

	path, err = filepath.Abs(keyPath)
	if err != nil {
		klog.Errorf("fail to translate %s to absolute path", keyPath)
		return err
	} else {
		opt.keyPath = path
	}

	if err := ValidateTLSPair(certPath, keyPath); err != nil {
		opt.Valid = false
		return err
	}

	opt.Valid = true
	return nil
}

// Load X.509 pair from go-map. Format: {"cert": "$CERT_PATH", "key", "$KEY_PATH"}
func (opt *TLSKeyPair) LoadFromMap(config map[string]string) error {
	certPath, c_ok := config["cert"]
	if !c_ok {
		return fmt.Errorf("loss key: `cert`")
	}

	keyPath, k_ok := config["key"]
	if !k_ok {
		return fmt.Errorf("loss key: `key`")
	}

	return opt.SetTLSKeyPair(certPath, keyPath)
}

// check if certificate and key are one pair
func ValidateTLSPair(certPath string, keyPath string) error {
	// validate Certificate exist
	if _, err := os.Stat(certPath); err != nil {
		klog.Errorf("invalid X.509 Certificate path: %s", certPath)
		return fmt.Errorf("invalid X.509 Certificate path: %s", certPath)
	}

	// Validate private-key exist
	if _, err := os.Stat(keyPath); err != nil {
		klog.Errorf("invalid private-key path: %s", keyPath)
		return fmt.Errorf("invalid private-key path: %s", keyPath)
	}

	// validate pair for certificate and key
	if _, err := tls.LoadX509KeyPair(certPath, keyPath); err != nil {
		return err
	}

	return nil
}

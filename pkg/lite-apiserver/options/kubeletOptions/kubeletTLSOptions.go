package kubeletOptions

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

var (
	tlsJsonPath string
)

type KubeletTLSKeyPair struct {
	//CaPath   string `json:"ca"`
	CertPath string `json:"cert"`
	KeyPath  string `json:"key"`
}

func NewKubeletTLSKeyPair() *KubeletTLSKeyPair {
	return &KubeletTLSKeyPair{}
}

func (opt *KubeletTLSKeyPair) LoadFromJson(jsonPath *string) error {
	if jsonPath == nil {
		jsonPath = &tlsJsonPath
	}

	if len(*jsonPath) <= 0 {
		err := "loss json config for X.509 Certificate information to kubelet ,you can set by \"--kubelet-client-cert-config=$path\""
		klog.Error(err)
		return fmt.Errorf(err)
	}

	// load json
	bytes, err := ioutil.ReadFile(*jsonPath)
	if err != nil {
		klog.Errorf("fail to load %s while load X.509 Certificate information to kubelet", jsonPath)
		return err
	}

	// unmarshal json
	if err := json.Unmarshal(bytes, opt); err != nil {
		klog.Errorf("fail to unmarshal json while load X.509 Certificate information to kubelet", jsonPath)
		return err
	}

	if err := opt.validate(); err != nil {
		klog.Errorf("fail to load X.509 Certificate information to kubelet from json")
		return err
	}
	klog.Info("success to load X.509 Certificate information to kubelet from json")

	return nil
}

// check if certificate and key are one pair
func (opt *KubeletTLSKeyPair) validate() error {
	var e error

	// validate Certificate
	if err := FileExist(opt.CertPath); err != nil {
		klog.Error("Validate Kubelet TLS config: ", err.Error())
		e = err
	}

	// Validate private Key
	if err := FileExist(opt.KeyPath); err != nil {
		klog.Error("Validate Kubelet TLS config: ", err.Error())
		e = err
	}

	// // Validate CA
	// if err := FileExist(opt.caPath); err != nil {
	// 	klog.Error("Validate Kubelet TLS config: ", err.Error())
	// 	e = err
	// }

	if e != nil {
		return e
	}

	// validate pair for certificate and key
	_, err := tls.LoadX509KeyPair(opt.CertPath, opt.KeyPath)
	if err != nil {
		klog.Error("Validate Kubelet TLS pair for certificate and key: ", err.Error())
		return err
	}

	return nil
}

func (opt *KubeletTLSKeyPair) AddFlagsTo(fs *pflag.FlagSet) {
	fs.StringVar(&tlsJsonPath, "kubelet-client-cert-config", "", "path to config store the X.509 Certificate information to kubelet")
}

func FileExist(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}

	if os.IsNotExist(err) {
		return fmt.Errorf("%s is not exist", path)
	}

	return fmt.Errorf("unknow error for %s, row error string: %s", path, err.Error())
}

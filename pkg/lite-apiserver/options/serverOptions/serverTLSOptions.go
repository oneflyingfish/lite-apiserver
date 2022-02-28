package serverOptions

// import (
// 	"LiteKube/pkg/common"
// 	"crypto/tls"
// 	"fmt"
// 	"io/ioutil"
// 	"os"

// 	"gopkg.in/yaml.v2"
// 	"k8s.io/klog/v2"
// )

// type ServerTLSKeyPair struct {
// 	CACertPath string `yaml:"cacert"`
// 	CAKeyPath  string `yaml:"cakey"`
// }

// func NewServerTLSKeyPair() *ServerTLSKeyPair {
// 	return &ServerTLSKeyPair{}
// }

// func (opt *ServerTLSKeyPair) LoadFromConfig(configPath *string) error {
// 	if configPath == nil || len(*configPath) <= 0 {
// 		err := "loss config for X.509 Certificate information for lite-apiserver ,you can set by \"--tls-configpath=$path\""
// 		klog.Error(err)
// 		return fmt.Errorf(err)
// 	}

// 	// load config
// 	bytes, err := ioutil.ReadFile(*configPath)
// 	if err != nil {
// 		klog.Errorf("fail to load %s while load X.509 Certificate information for lite-apiserver", *configPath)
// 		return err
// 	}

// 	// unmarshal config
// 	if err := yaml.Unmarshal(bytes, opt); err != nil {
// 		klog.Errorf("fail to unmarshal config while load X.509 Certificate information for lite-apiserver", *configPath)
// 		return err
// 	}

// 	// to absolute path
// 	if err := common.AbsPath(&opt.CACertPath); err != nil {
// 		klog.Errorf("fail to translate %s to absolute path", opt.CACertPath)
// 		return err
// 	}

// 	if err := common.AbsPath(&opt.CAKeyPath); err != nil {
// 		klog.Errorf("fail to translate %s to absolute path", opt.CAKeyPath)
// 		return err
// 	}

// 	if err := opt.validate(); err != nil {
// 		klog.Errorf("fail to load X.509 Certificate information for lite-apiserver from config")
// 		return err
// 	}

// 	klog.Info("success to load X.509 Certificate information for lite-apiserver from config")

// 	return nil
// }

// // check if certificate and key are one pair
// func (opt *ServerTLSKeyPair) validate() error {
// 	var e error

// 	// validate CA Certificate
// 	if err := FileExist(opt.CACertPath); err != nil {
// 		klog.Error("Validate lite-apiserver CA TLS config: ", err.Error())
// 		e = err
// 	}

// 	// Validate CA private Key
// 	if err := FileExist(opt.CAKeyPath); err != nil {
// 		klog.Error("Validate lite-apiserver CA TLS config: ", err.Error())
// 		e = err
// 	}

// 	if e != nil {
// 		return e
// 	}

// 	// validate pair for CA certificate and key
// 	_, err := tls.LoadX509KeyPair(opt.CACertPath, opt.CAKeyPath)
// 	if err != nil {
// 		klog.Error("Validate lite-apiserver CA TLS pair for certificate and key: ", err.Error())
// 		return err
// 	}

// 	return nil
// }

// func FileExist(path string) error {
// 	_, err := os.Stat(path)
// 	if err == nil {
// 		return nil
// 	}

// 	if os.IsNotExist(err) {
// 		return fmt.Errorf("%s is not exist", path)
// 	}

// 	return fmt.Errorf("unknow error for %s, row error string: %s", path, err.Error())
// }

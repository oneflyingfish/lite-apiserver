package serverOptions

import (
	"LiteKube/pkg/common"
	"LiteKube/pkg/lite-apiserver/cert"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

var (
	serverConfigPath string
	defaultValue     ServerOption = ServerOption{
		CATLSKeyPair:     nil,
		ServerTLSKeyPair: nil,
		Debug:            false,
		Hostname:         "127.0.0.1",
		Port:             13500,
		InsecurePort:     -1,
		TLSStoreFold:     "",
		SyncDuration:     10,
	}
)

type ServerOption struct {
	CATLSKeyPair     *cert.TLSKeyPair `json:"-"`
	ServerTLSKeyPair *cert.TLSKeyPair `json:"-"`
	Debug            bool             `json:"-"`
	Hostname         string           `yaml:"hostname"`
	Port             int              `yaml:"port"`
	InsecurePort     int              `yaml:"insecure-port"`
	TLSStoreFold     string           `yaml:"tls-store-fold"`
	SyncDuration     int              `yaml:"syncduration"`
}

func NewServerOptions() *ServerOption {
	return &ServerOption{
		CATLSKeyPair:     cert.NewTLSKeyPair(),
		ServerTLSKeyPair: cert.NewTLSKeyPair(),
	}
}

func (opt *ServerOption) AddFlagsTo(fs *pflag.FlagSet) {
	fs.StringVar(&serverConfigPath, "config", "", "config for lite-apiserver (lower priority to flags)")
	fs.StringVar(&opt.Hostname, "hostname", "", fmt.Sprintf("hostname of lite-apiserver (default: %s)", defaultValue.Hostname))
	fs.IntVar(&opt.Port, "port", 0, fmt.Sprintf("https port of lite-apiserver (default: %d)", defaultValue.Port))
	fs.IntVar(&opt.InsecurePort, "insecure-port", 0, fmt.Sprintf("http port of lite-apiserver, not secure, set -1 to disable (default: %d)", defaultValue.InsecurePort))
	fs.StringVar(&opt.TLSStoreFold, "tls-store-fold", "", fmt.Sprintf("fold path to store CA and server X.509 files for lite-apiserver, which contains {ca, server}.{pem, -key.pem} (default: \"%s\")", defaultValue.TLSStoreFold))
	fs.IntVar(&opt.SyncDuration, "syncduration", 0, fmt.Sprintf("max time for one-request last (default: %d)", defaultValue.SyncDuration))
	fs.BoolVar(&opt.Debug, "debug", false, fmt.Sprintf("enable debug or not, this value is not allow to set with config-file (default: %s)", strconv.FormatBool(defaultValue.Debug)))
}

func (opt *ServerOption) LoadServerConfig() error {
	opt_file := &ServerOption{
		CATLSKeyPair: nil,
		Debug:        false,
	}

	if len(serverConfigPath) > 0 {
		// load config
		bytes, err := ioutil.ReadFile(serverConfigPath)
		if err != nil {
			klog.Warningf("fail to read %s for config, process skip directly", serverConfigPath)
			goto SKIP
		}

		// unmarshal config
		if err := yaml.Unmarshal(bytes, opt_file); err != nil {
			klog.Warningf("fail to unmarshal %s for config, process skip directly", serverConfigPath)
			goto SKIP
		}
	}

SKIP:

	opt.MergeConfig(opt_file)

	// load whole X509 config
	if err := opt.LoadX509(); err != nil {
		return err
	}

	return nil
}

func (opt *ServerOption) MergeConfig(opt_file *ServerOption) error {
	// serverConfigPath to absolute path
	if err := common.AbsPath(&serverConfigPath); err != nil {
		klog.Errorf("fail to translate %s to absolute path", serverConfigPath)
		return err
	}

	// merge config-file to flags
	common.Merge(opt, opt_file, &defaultValue, "Hostname")
	common.Merge(opt, opt_file, &defaultValue, "Port")
	common.Merge(opt, opt_file, &defaultValue, "InsecurePort")
	common.Merge(opt, opt_file, &defaultValue, "TLSStoreFold")
	common.Merge(opt, opt_file, &defaultValue, "SyncDuration")
	common.Merge(opt, opt_file, &defaultValue, "Debug")

	// CATLSConfigPath to absolute path
	if err := common.AbsPath(&opt.TLSStoreFold); err != nil {
		klog.Errorf("fail to translate %s to absolute path", opt.TLSStoreFold)
		return err
	}

	return nil
}

func (opt *ServerOption) PrintArgs() error {
	klog.Infof("--debug=%s ", strconv.FormatBool(opt.Debug))
	klog.Infof("--config=%s ", serverConfigPath)
	klog.Infof("--hostname=%s ", opt.Hostname)
	klog.Infof("--port=%d ", opt.Port)
	klog.Infof("--insecure-port=%d ", opt.InsecurePort)
	klog.Infof("--tls-store-fold=%s", opt.TLSStoreFold)
	klog.Infof("--syncduration=%d", opt.SyncDuration)
	return nil
}

func (opt *ServerOption) LoadX509() error {
	if err := opt.LoadCAX509(); err != nil {
		return err
	}

	if err := opt.LoadServerX509(); err != nil {
		return err
	}
	//return opt.CATLSKeyPair.LoadFromConfig(&opt.TLSConfigPath)
	return nil
}

func (opt *ServerOption) LoadCAX509() error {
	if err := opt.tryLoadCAX509(); err != nil {
		klog.Warningf("X.509 CA Certificate for lite-apiserver are lost, we will try to generate one.")
		cert.CreateCACert(opt.TLSStoreFold)

		if e := opt.tryLoadCAX509(); e != nil {
			klog.Warningf("fail to create CA file for lite-apiserver in %s", opt.TLSStoreFold)
			return err
		}

		klog.Info("Success to create new CA Certificate for lite-apiserver")
	}

	klog.Info("Success to Load CA Certificate for lite-apiserver")
	return nil
}

func (opt *ServerOption) LoadServerX509() error {
	if err := opt.tryLoadServerX509(); err != nil {
		klog.Warningf("X.509 Server Certificate for lite-apiserver are lost, we will try to generate one.")

		// caKey := rsa.PrivateKey.Load()
		caCert, caKey, valid := opt.CATLSKeyPair.GetTLSKeyPairCertificate()
		if !valid {
			klog.Error("fail to read CA X.509 for lite-apiserver when create server-certificate")
			return fmt.Errorf("fail to read CA X.509 for lite-apiserver when create server-certificate")
		}

		if err := cert.CreateServerCert(opt.TLSStoreFold, caCert, caKey); err != nil {
			klog.Errorf("fail to create Server file for lite-apiserver in %s", opt.TLSStoreFold)
			return err
		}

		if e := opt.tryLoadServerX509(); e != nil {
			klog.Errorf("Server Certificate file for lite-apiserver is not useful in %s", opt.TLSStoreFold)
			return err
		}

		klog.Info("Success to create new Server Certificate for lite-apiserver")
	}

	klog.Info("Success to load Server Certificate for lite-apiserver")
	return nil
}

func (opt *ServerOption) tryLoadCAX509() error {
	certPath := filepath.Join(opt.TLSStoreFold, "ca.pem")
	keyPath := filepath.Join(opt.TLSStoreFold, "ca-key.pem")

	if _, err := os.Stat(certPath); err != nil {
		return fmt.Errorf("bad ca certificate")
	}

	if _, err := os.Stat(keyPath); err != nil {
		return fmt.Errorf("bad ca private key")
	}

	if err := opt.CATLSKeyPair.SetTLSKeyPair(certPath, keyPath); err != nil {
		return err
	}
	return nil
}

func (opt *ServerOption) tryLoadServerX509() error {
	certPath := filepath.Join(opt.TLSStoreFold, "server.pem")
	keyPath := filepath.Join(opt.TLSStoreFold, "server-key.pem")

	if _, err := os.Stat(certPath); err != nil {
		return fmt.Errorf("bad server certificate")
	}

	if _, err := os.Stat(keyPath); err != nil {
		return fmt.Errorf("bad server private key")
	}

	if err := opt.ServerTLSKeyPair.SetTLSKeyPair(certPath, keyPath); err != nil {
		return err
	}
	return nil
}

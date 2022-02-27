package serverOptions

import (
	"LiteKube/pkg/common"
	"fmt"
	"io/ioutil"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

var (
	serverConfigPath string
	defaultValue     ServerOption = ServerOption{
		CATLSKeyPair:    nil,
		Hostname:        "127.0.0.1",
		Port:            13500,
		InsecurePort:    0,
		CATLSConfigPath: "",
		SyncDuration:    10,
	}
)

type ServerOption struct {
	CATLSKeyPair    *ServerTLSKeyPair
	Hostname        string `yaml:"hostname"`
	Port            int    `yaml:"port"`
	InsecurePort    int    `yaml:"insecure-port"`
	CATLSConfigPath string `yaml:"ca-tls-configpath"`
	SyncDuration    int    `yaml:"syncduration"`
}

func NewServerOptions() *ServerOption {
	return &ServerOption{
		CATLSKeyPair: NewServerTLSKeyPair(),
	}
}

func (opt *ServerOption) AddFlagsTo(fs *pflag.FlagSet) {
	fs.StringVar(&serverConfigPath, "config", "", "config for lite-apiserver (lower priority to flags)")
	fs.StringVar(&opt.Hostname, "hostname", "", fmt.Sprintf("hostname of lite-apiserver (default: %s)", defaultValue.Hostname))
	fs.IntVar(&opt.Port, "port", 0, fmt.Sprintf("https port of lite-apiserver (default: %d)", defaultValue.Port))
	fs.IntVar(&opt.InsecurePort, "insecure-port", 0, fmt.Sprintf("http port of lite-apiserver, not secure, set 0 to disable (default: %d)", defaultValue.InsecurePort))
	fs.StringVar(&opt.CATLSConfigPath, "ca-tls-configpath", "", fmt.Sprintf("path to config store the X.509 Certificate information for lite-apiserver (default: \"%s\")", defaultValue.CATLSConfigPath))
	fs.IntVar(&opt.SyncDuration, "--syncduration", 0, fmt.Sprintf("max time for one-request last (default: %d)", defaultValue.SyncDuration))
}

func (opt *ServerOption) LoadServerConfig() error {
	opt_file := &ServerOption{
		CATLSKeyPair: nil,
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
	common.Merge(opt, opt_file, &defaultValue, "CATLSConfigPath")
	common.Merge(opt, opt_file, &defaultValue, "SyncDuration")

	// CATLSConfigPath to absolute path
	if err := common.AbsPath(&opt.CATLSConfigPath); err != nil {
		klog.Errorf("fail to translate %s to absolute path", opt.CATLSConfigPath)
		return err
	}

	return nil
}

func (opt *ServerOption) PrintArgs() error {
	klog.Infof("--config=%s ", serverConfigPath)
	klog.Infof("--hostname=%s ", opt.Hostname)
	klog.Infof("--port=%d ", opt.Port)
	klog.Infof("--insecure-port=%d ", opt.InsecurePort)
	klog.Infof("--ca-tls-configpath=%s", opt.CATLSConfigPath)
	klog.Infof("--syncduration=%d", opt.SyncDuration)
	return nil
}

func (opt *ServerOption) LoadX509() error {
	return opt.CATLSKeyPair.LoadFromConfig(&opt.CATLSConfigPath)
}

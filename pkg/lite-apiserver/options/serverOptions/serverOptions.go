package serverOptions

import (
	"LiteKube/pkg/common"
	"encoding/json"
	"io/ioutil"

	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

var (
	serverConfigPath string
	defaultValue     ServerOption = ServerOption{
		Hostname: "127.0.0.1",
		Port:     6500,
	}
)

type ServerOption struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

func NewServerOptions() *ServerOption {
	return &ServerOption{}
}

func (opt *ServerOption) AddFlagsTo(fs *pflag.FlagSet) {
	fs.StringVar(&serverConfigPath, "config", "", "json-config for lite-apiserver (lower priority to flags)")
	fs.StringVar(&opt.Hostname, "hostname", defaultValue.Hostname, "hostname of kubelet")
	fs.IntVar(&opt.Port, "port", defaultValue.Port, "port of kubelet")
}

func (opt *ServerOption) LoadServerConfig() error {
	if len(serverConfigPath) <= 0 {
		return nil
	}

	opt_file := &ServerOption{}

	// load json
	bytes, err := ioutil.ReadFile(serverConfigPath)
	if err != nil {
		klog.Warningf("fail to read %s for kubelet-config, process skip directly", serverConfigPath)
		return nil
	}

	// unmarshal json
	if err := json.Unmarshal(bytes, opt_file); err != nil {
		klog.Warningf("fail to unmarshal %s for kubelet-config, process skip directly", serverConfigPath)
		return nil
	}

	opt.MergeConfig(opt_file)

	return nil
}

func (opt *ServerOption) MergeConfig(opt_file *ServerOption) {
	// hostname
	if (common.IsZero(opt.Hostname) || opt.Hostname == defaultValue.Hostname) && !common.IsZero(opt_file.Hostname) {
		opt.Hostname = opt_file.Hostname
	}

	// port
	if (common.IsZero(opt.Port) || opt.Port == defaultValue.Port) && !common.IsZero(opt_file.Port) {
		opt.Port = opt_file.Port
	}
}

func (opt *ServerOption) PrintArgs() error {
	klog.Infof("--config=%s", serverConfigPath)
	klog.Infof("--hostname=%s", opt.Hostname)
	klog.Infof("--port=%d", opt.Port)
	return nil
}

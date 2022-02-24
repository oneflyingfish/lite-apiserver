package kubeletOptions

import (
	"encoding/json"
	"io/ioutil"

	"LiteKube/pkg/common"

	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

var (
	kubeletConfigPath string
	defaultValue      KubeletOption = KubeletOption{
		TLSKeyPair:  nil,
		Hostname:    "127.0.0.1",
		HealthzPort: 10248,
		Port:        10250,
	}
)

type KubeletOption struct {
	TLSKeyPair  *KubeletTLSKeyPair
	Hostname    string `json:"kubelet-hostname"`
	HealthzPort int    `json:"kubelet-healthzport"`
	Port        int    `json:"kubelet-port"`
}

func NewKubeletOption() *KubeletOption {
	return &KubeletOption{
		TLSKeyPair: NewKubeletTLSKeyPair(),
	}
}

func (opt *KubeletOption) AddFlagsTo(fs *pflag.FlagSet) {
	fs.StringVar(&kubeletConfigPath, "kubelet-config", "", "json-config for kubelet (lower priority to flags)")
	fs.StringVar(&opt.Hostname, "kubelet-hostname", defaultValue.Hostname, "hostname of kubelet")
	fs.IntVar(&opt.HealthzPort, "kubelet-healthzport", defaultValue.HealthzPort, "healthz port of kubelet")
	fs.IntVar(&opt.Port, "kubelet-port", defaultValue.Port, "port of kubelet")

	// add flags of load X.509 json config for connection to kubelet.
	opt.TLSKeyPair.AddFlagsTo(fs)
}

func (opt *KubeletOption) LoadKubeletConfig() error {
	// load whole X509 config
	if err := opt.LoadX509(); err != nil {
		return err
	}

	if len(kubeletConfigPath) <= 0 {
		return nil
	}

	opt_file := &KubeletOption{
		TLSKeyPair: nil,
	}

	// load json
	bytes, err := ioutil.ReadFile(kubeletConfigPath)
	if err != nil {
		klog.Warningf("fail to read %s for kubelet-config, process skip directly", kubeletConfigPath)
		return nil
	}

	// unmarshal json
	if err := json.Unmarshal(bytes, opt_file); err != nil {
		klog.Warningf("fail to unmarshal %s for kubelet-config, process skip directly", kubeletConfigPath)
		return nil
	}

	opt.MergeConfig(opt_file)
	return nil
}

func (opt *KubeletOption) MergeConfig(opt_file *KubeletOption) {
	// kubelet-hostname
	if (common.IsZero(opt.Hostname) || opt.Hostname == defaultValue.Hostname) && !common.IsZero(opt_file.Hostname) {
		opt.Hostname = opt_file.Hostname
	}

	// kubelet-healthzPort
	if (common.IsZero(opt.HealthzPort) || opt.HealthzPort == defaultValue.HealthzPort) && !common.IsZero(opt_file.HealthzPort) {
		opt.HealthzPort = opt_file.HealthzPort
	}

	// kubelet-port
	if (common.IsZero(opt.Port) || opt.Port == defaultValue.Port) && !common.IsZero(opt_file.Port) {
		opt.Port = opt_file.Port
	}
}

func (opt *KubeletOption) LoadX509() error {
	return opt.TLSKeyPair.LoadFromJson(nil)
}

func (opt *KubeletOption) PrintArgs() error {
	klog.Infof("--kubelet-config=%s", kubeletConfigPath)
	klog.Infof("--kubelet-hostname=%s", opt.Hostname)
	klog.Infof("--kubelet-healthzport=%d", opt.HealthzPort)
	klog.Infof("--kubelet-port=%d", opt.Port)
	return nil
}

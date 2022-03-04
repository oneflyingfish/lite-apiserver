package kubeletOptions

import (
	"LiteKube/pkg/lite-apiserver/cert"
	"fmt"
	"io/ioutil"
	"os"

	"LiteKube/pkg/common"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

var (
	kubeletConfigPath string
	defaultValue      KubeletOption = KubeletOption{
		TLSKeyPair:    nil,
		Hostname:      "127.0.0.1",
		HealthzPort:   10248,
		Port:          10250,
		TLSConfigPath: "",
		ManifestsFold: "",
	}
)

type KubeletOption struct {
	TLSKeyPair    *cert.TLSKeyPair `json:"-"`
	Hostname      string           `yaml:"kubelet-hostname"`
	HealthzPort   int              `yaml:"kubelet-healthzport"`
	Port          int              `yaml:"kubelet-port"`
	ManifestsFold string           `yaml:"kubelet-pod-manifest-path"`
	TLSConfigPath string           `yaml:"kubelet-client-cert-config"`
}

func NewKubeletOption() *KubeletOption {
	return &KubeletOption{
		TLSKeyPair: cert.NewTLSKeyPair(),
	}
}

func (opt *KubeletOption) AddFlagsTo(fs *pflag.FlagSet) {
	fs.StringVar(&kubeletConfigPath, "kubelet-config", "", "config for kubelet (lower priority to flags)")
	fs.StringVar(&opt.Hostname, "kubelet-hostname", "", fmt.Sprintf("hostname of kubelet (default: %s)", defaultValue.Hostname))
	fs.IntVar(&opt.HealthzPort, "kubelet-healthzport", 0, fmt.Sprintf("healthz port of kubelet (default: %d)", defaultValue.HealthzPort))
	fs.IntVar(&opt.Port, "kubelet-port", 0, fmt.Sprintf("port of kubelet (default: %d)", defaultValue.Port))
	fs.StringVar(&opt.TLSConfigPath, "kubelet-client-cert-config", "", fmt.Sprintf("path to config store the X.509 Certificate information to kubelet (default: \"%s\")", defaultValue.TLSConfigPath))
	fs.StringVar(&opt.ManifestsFold, "kubelet-pod-manifest-path", "", fmt.Sprintf("same value with kubelet --pod-manifest-path， set \"\" will disable lite-apiserver with alpha version (default: \"%s\")", defaultValue.TLSConfigPath))
}

func (opt *KubeletOption) LoadKubeletConfig() error {
	opt_file := &KubeletOption{
		TLSKeyPair: nil,
	}

	if len(kubeletConfigPath) > 0 {
		// load yaml
		bytes, err := ioutil.ReadFile(kubeletConfigPath)
		if err != nil {
			klog.Warningf("fail to read %s for kubelet-config, process skip directly", kubeletConfigPath)
			goto SKIP
		}

		// unmarshal yaml
		if err := yaml.Unmarshal(bytes, opt_file); err != nil {
			klog.Warningf("fail to unmarshal %s for kubelet-config, process skip directly", kubeletConfigPath)
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

func (opt *KubeletOption) MergeConfig(opt_file *KubeletOption) error {
	// kubeletConfigPath to absolute path
	if err := common.AbsPath(&kubeletConfigPath); err != nil {
		klog.Errorf("fail to translate %s to absolute path", kubeletConfigPath)
		return err
	}

	// merge config-file to flags
	common.Merge(opt, opt_file, &defaultValue, "Hostname")
	common.Merge(opt, opt_file, &defaultValue, "HealthzPort")
	common.Merge(opt, opt_file, &defaultValue, "Port")
	common.Merge(opt, opt_file, &defaultValue, "TLSConfigPath")
	common.Merge(opt, opt_file, &defaultValue, "ManifestsFold")

	if common.IsZero(opt.ManifestsFold) {
		klog.Warning("It seems that you have forgotten to set --kubelet-pod-manifest-path, lite-Apiserver will not work in alpha version")
	} else {
		if err := common.AbsPath(&opt.ManifestsFold); err != nil {
			klog.Errorf("fail to translate %s to absolute path", opt.ManifestsFold)
			return err
		}

		// ensure kubelet-client-cert-config exists
		if err := os.MkdirAll(opt.ManifestsFold, os.ModePerm); err != nil {
			klog.Errorf("fail to create fold path: %s", opt.ManifestsFold)
			return err
		}
	}

	// TLSConfigPath to absolute path
	if err := common.AbsPath(&opt.TLSConfigPath); err != nil {
		klog.Errorf("fail to translate %s to absolute path", opt.TLSConfigPath)
		return err
	}
	return nil
}

func (opt *KubeletOption) LoadX509() error {
	if len(opt.TLSConfigPath) <= 0 {
		err := "loss config for X.509 Certificate information to kubelet ,you can set by \"--kubelet-client-cert-config=$path\""
		klog.Error(err)
		return fmt.Errorf(err)
	}

	// load config
	bytes, err := ioutil.ReadFile(opt.TLSConfigPath)
	if err != nil {
		klog.Errorf("fail to load %s while load X.509 Certificate information to kubelet", opt.TLSConfigPath)
		return err
	}

	data := make(map[string]string)

	// unmarshal config
	if err := yaml.Unmarshal(bytes, &data); err != nil {
		klog.Errorf("fail to unmarshal config while load X.509 Certificate information to kubelet", opt.TLSConfigPath)
		return err
	}

	if err := opt.TLSKeyPair.LoadFromMap(data); err != nil {
		klog.Errorf("fail to load X.509 Certificate information to kubelet from config file")
	}

	klog.Info("success to load X.509 Certificate information to kubelet from config file")

	return nil
}

func (opt *KubeletOption) PrintArgs() error {
	klog.Infof("--kubelet-config=%s", kubeletConfigPath)
	klog.Infof("--kubelet-hostname=%s", opt.Hostname)
	klog.Infof("--kubelet-healthzport=%d", opt.HealthzPort)
	klog.Infof("--kubelet-port=%d", opt.Port)
	klog.Infof("--kubelet-client-cert-config=%s", opt.TLSConfigPath)
	klog.Infof("--kubelet-pod-manifest-path=%s", opt.ManifestsFold)
	return nil
}

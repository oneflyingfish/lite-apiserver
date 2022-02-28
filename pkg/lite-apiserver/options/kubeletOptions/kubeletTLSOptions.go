package kubeletOptions

// type KubeletTLSKeyPair struct {
// 	//CaPath   string `yaml:"ca"`
// 	CertPath string `yaml:"cert"`
// 	KeyPath  string `yaml:"key"`
// }

// func NewKubeletTLSKeyPair() *KubeletTLSKeyPair {
// 	return &KubeletTLSKeyPair{}
// }

// func (opt *KubeletTLSKeyPair) LoadFromConfig(configPath *string) error {

// 	if configPath == nil || len(*configPath) <= 0 {
// 		err := "loss config for X.509 Certificate information to kubelet ,you can set by \"--kubelet-client-cert-config=$path\""
// 		klog.Error(err)
// 		return fmt.Errorf(err)
// 	}

// 	// load config
// 	bytes, err := ioutil.ReadFile(*configPath)
// 	if err != nil {
// 		klog.Errorf("fail to load %s while load X.509 Certificate information to kubelet", *configPath)
// 		return err
// 	}

// 	// unmarshal config
// 	if err := yaml.Unmarshal(bytes, opt); err != nil {
// 		klog.Errorf("fail to unmarshal config while load X.509 Certificate information to kubelet", *configPath)
// 		return err
// 	}

// 	// to absolute path
// 	if err := common.AbsPath(&opt.CertPath); err != nil {
// 		klog.Errorf("fail to translate %s to absolute path", opt.CertPath)
// 		return err
// 	}

// 	if err := common.AbsPath(&opt.KeyPath); err != nil {
// 		klog.Errorf("fail to translate %s to absolute path", opt.KeyPath)
// 		return err
// 	}

// 	if err := opt.validate(); err != nil {
// 		klog.Errorf("fail to load X.509 Certificate information to kubelet from config file")
// 		return err
// 	}
// 	klog.Info("success to load X.509 Certificate information to kubelet from config file")

// 	return nil
// }

// // check if certificate and key are one pair
// func (opt *KubeletTLSKeyPair) validate() error {
// 	var e error

// 	// validate Certificate
// 	if _, err := os.Stat(opt.CertPath); err != nil {
// 		klog.Error("Validate Kubelet TLS config: ", err.Error())
// 		e = err
// 	}

// 	// Validate private Key
// 	if _, err := os.Stat(opt.KeyPath); err != nil {
// 		klog.Error("Validate Kubelet TLS config: ", err.Error())
// 		e = err
// 	}

// 	// // Validate CA
// 	// if err := FileExist(opt.caPath); err != nil {
// 	// 	klog.Error("Validate Kubelet TLS config: ", err.Error())
// 	// 	e = err
// 	// }

// 	if e != nil {
// 		return e
// 	}

// 	// validate pair for certificate and key
// 	_, err := tls.LoadX509KeyPair(opt.CertPath, opt.KeyPath)
// 	if err != nil {
// 		klog.Error("Validate Kubelet TLS pair for certificate and key: ", err.Error())
// 		return err
// 	}

// 	return nil
// }

package options

import (
	"LiteKube/pkg/lite-apiserver/options/kubeApiserverOptions"
	"LiteKube/pkg/lite-apiserver/options/kubeletOptions"
	"LiteKube/pkg/lite-apiserver/options/serverOptions"

	cliflag "k8s.io/component-base/cli/flag"
)

type ServerRunOption struct {
	apiserverOptions *kubeApiserverOptions.KubeApiserverOption
	kubeletOption    *kubeletOptions.KubeletOption
	serverOption     *serverOptions.ServerOption
}

func NewServerRunOption() *ServerRunOption {
	return &ServerRunOption{
		apiserverOptions: kubeApiserverOptions.NewKubeApiserverOption(),
		kubeletOption:    kubeletOptions.NewKubeletOption(),
		serverOption:     serverOptions.NewServerOptions(),
	}
}

func (opt *ServerRunOption) GetNamedFlagsSet() (fsSet cliflag.NamedFlagSets) {
	opt.serverOption.AddFlagsTo(fsSet.FlagSet("lite-apiserver"))
	opt.apiserverOptions.AddFlagsTo(fsSet.FlagSet("kube-apiserver"))
	opt.kubeletOption.AddFlagsTo(fsSet.FlagSet("kubelet"))
	return
}

// load config from disk-file and merge with flags
func (opt *ServerRunOption) LoadConfig() []error {
	var errors []error

	if err := opt.serverOption.LoadServerConfig(); err != nil {
		errors = append(errors, err)
	}

	if err := opt.kubeletOption.LoadKubeletConfig(); err != nil {
		errors = append(errors, err)
	}

	if err := opt.apiserverOptions.LoadKubeApiserverConfig(); err != nil {
		errors = append(errors, err)
	}

	// print all args
	opt.PrintArgs()

	var new_errors []error
	for _, item := range errors {
		if item != nil {
			new_errors = append(new_errors, item)
		}
	}

	if len(new_errors) <= 0 {
		return new_errors
	}

	return nil
}

// Complete set default ServerRunOptions. It should be called after flags parsed.
func (opt *ServerRunOption) Complete() error {
	return nil
}

func (opt *ServerRunOption) PrintArgs() error {
	var err_ error

	if err := opt.serverOption.PrintArgs(); err != nil {
		err_ = err
	}

	if err := opt.kubeletOption.PrintArgs(); err != nil {
		err_ = err
	}

	if err := opt.apiserverOptions.PrintArgs(); err != nil {
		err_ = err
	}

	return err_
}

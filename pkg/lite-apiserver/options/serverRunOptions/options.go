package serverRunOptions

import (
	"LiteKube/pkg/lite-apiserver/options/kubeApiserverOptions"
	"LiteKube/pkg/lite-apiserver/options/kubeletOptions"
	"LiteKube/pkg/lite-apiserver/options/serverOptions"

	verflag "LiteKube/pkg/version/varflag"

	cliflag "k8s.io/component-base/cli/flag"
)

type ServerRunOption struct {
	ApiserverOptions *kubeApiserverOptions.KubeApiserverOption
	KubeletOption    *kubeletOptions.KubeletOption
	ServerOption     *serverOptions.ServerOption
}

func NewServerRunOption() *ServerRunOption {
	return &ServerRunOption{
		ApiserverOptions: kubeApiserverOptions.NewKubeApiserverOption(),
		KubeletOption:    kubeletOptions.NewKubeletOption(),
		ServerOption:     serverOptions.NewServerOptions(),
	}
}

func (opt *ServerRunOption) GetNamedFlagsSet() (fsSet cliflag.NamedFlagSets) {
	opt.ServerOption.AddFlagsTo(fsSet.FlagSet("lite-apiserver"))
	opt.ApiserverOptions.AddFlagsTo(fsSet.FlagSet("kube-apiserver"))
	opt.KubeletOption.AddFlagsTo(fsSet.FlagSet("kubelet"))
	verflag.AddFlagsTo(fsSet.FlagSet("others"))
	return
}

// load config from disk-file and merge with flags
func (opt *ServerRunOption) LoadConfig() []error {
	var errors []error

	if err := opt.ServerOption.LoadServerConfig(); err != nil {
		errors = append(errors, err)
	}

	if err := opt.KubeletOption.LoadKubeletConfig(); err != nil {
		errors = append(errors, err)
	}

	if err := opt.ApiserverOptions.LoadKubeApiserverConfig(); err != nil {
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

	if err := opt.ServerOption.PrintArgs(); err != nil {
		err_ = err
	}

	if err := opt.KubeletOption.PrintArgs(); err != nil {
		err_ = err
	}

	if err := opt.ApiserverOptions.PrintArgs(); err != nil {
		err_ = err
	}

	return err_
}

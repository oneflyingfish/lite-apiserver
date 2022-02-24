package kubeApiserverOptions

import (
	"github.com/spf13/pflag"
)

type KubeApiserverOption struct {
}

func NewKubeApiserverOption() *KubeApiserverOption {
	return &KubeApiserverOption{}
}

func (opt *KubeApiserverOption) AddFlagsTo(fs *pflag.FlagSet) {
	// do nothing now, because we don't connect to kube-apiserver at this time.
}

func (opt *KubeApiserverOption) LoadKubeApiserverConfig() error {
	return nil
}

func (opt *KubeApiserverOption) PrintArgs() error {
	return nil
}

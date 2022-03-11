package StorgeOperator

import k8sv1 "k8s.io/api/core/v1"

type StorgeOperator interface {
	HasPod(k8sv1.Pod) (bool, error)
	ValidatePod(k8sv1.Pod) (bool, error)
	ReadPodList() ([]k8sv1.Pod, error)
	ReadNamespaces() ([]string, error)
	ReadNames() ([]string, error)
	ReadPorts() ([]int, error)
	ReadPod(string, string) (k8sv1.Pod, error)
	WritePod(k8sv1.Pod) error
	DeletePod(namespace string, name string) error
}

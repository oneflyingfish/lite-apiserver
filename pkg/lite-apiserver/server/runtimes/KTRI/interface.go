package KTRI

import (
	"LiteKube/pkg/lite-apiserver/server/runtimes/KTRI/apis"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	k8sv1 "k8s.io/api/core/v1"
)

func Log(logpath string) (string, error) {
	resp, err := apis.GetRequest(filepath.Join(Log_Path, logpath))
	if err != nil {
		return fmt.Sprintf("error: %s", err.Error()), err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("error: %s", err.Error()), err
	}

	return string(data), nil
}

func ReadPod(namespace string, podName string) (*k8sv1.Pod, error) {
	pods, err := ReadPodList()
	if err != nil {
		return nil, err
	}

	for _, pod := range pods.Items {
		if pod.Namespace == namespace {
			return &pod, nil
		}
	}
	return nil, nil
}

func ReadPodList() (*k8sv1.PodList, error) {
	resp, err := apis.GetRequest(PODLIST)

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var podList k8sv1.PodList
	err = json.Unmarshal(data, &podList)
	if err != nil {
		return nil, err
	} else {
		return &podList, nil
	}
}

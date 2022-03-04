package describe

import (
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type StatusInfo struct {
	Kind       string      `json:"kind"`
	ApiVersion string      `json:"apiVersion"`
	Metadata   interface{} `json:"metadata"`
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Reason     string      `json:"reason"`
	Details    interface{} `json:"details"`
	Code       int         `json:"code"`
}

func (status StatusInfo) Complete() StatusInfo {
	if len(status.Kind) < 1 {
		status.Kind = "Status"
	}

	if len(status.ApiVersion) < 1 {
		status.ApiVersion = "v1"
	}

	if status.Metadata == nil {
		status.Metadata = struct{}{}
	}

	if len(status.Status) < 1 {
		status.Status = "Failure"
	}

	if status.Details == nil {
		status.Details = struct{}{}
	}

	return status
}

func LoadJson(news string) StatusInfo {
	status := StatusInfo{}
	json.Unmarshal([]byte(news), &status)
	return status.Complete()
}

func (status StatusInfo) ToJson() string {
	bytes, err := json.Marshal(status)
	if err != nil {
		return string(bytes)
	} else {
		return ""
	}
}

func (status StatusInfo) ToBytes() []byte {
	fmt.Println(status)
	bytes, err := json.Marshal(status)
	if err != nil {
		return bytes
	} else {
		klog.Error(err.Error())
		return nil
	}
}

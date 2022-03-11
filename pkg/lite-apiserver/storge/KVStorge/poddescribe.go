package KVStorge

import "encoding/json"

type PodDescribe struct {
	Ports     []int32 `json:"ports"`
	Name      string  `json:"name"`
	Namespace string  `json:"namespace"`
	Hash      string  `json:"hash"`
	PodNow    []byte  `json:"PodNow"`
}

func (pd PodDescribe) Marshal() ([]byte, error) {
	return json.Marshal(pd)
}

func (pd PodDescribe) UnMarshal(data []byte) error {
	return json.Unmarshal(data, &pd)
}

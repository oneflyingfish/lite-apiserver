package KVStorge

import (
	"encoding/json"
)

type KeyRecord struct {
	keys []string `json:"keys"`
}

const dbKey = "@@ALL@@"

func (kr KeyRecord) Marshal() ([]byte, error) {
	return json.Marshal(kr)
}

func (kr KeyRecord) UnMarshal(data []byte) error {
	return json.Unmarshal(data, &kr)
}

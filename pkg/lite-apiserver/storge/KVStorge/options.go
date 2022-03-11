package KVStorge

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"
	k8sv1 "k8s.io/api/core/v1"
)

const (
	storgeFileName = "~/litekube-lite-apiserver.db"
)

var podsMemoryCache []k8sv1.Pod

type DBWriter struct {
	StorgeFileName string
	KR             KeyRecord
}

func NewDBWriter(dbpath string) *DBWriter {
	if len(dbpath) < 1 {
		return &DBWriter{
			StorgeFileName: storgeFileName,
		}
	} else {

		err := os.MkdirAll(dbpath, os.ModePerm)
		if err != nil {
			return nil
		}

		return &DBWriter{
			StorgeFileName: filepath.Join(dbpath, "litekube-lite-apiserver.db"),
			KR:             KeyRecord{keys: nil},
		}
	}
}

func mergeKey(namespace string, name string) string {
	return namespace + "@" + name
}

func getKey(pod k8sv1.Pod) string {
	return mergeKey(pod.GetNamespace(), pod.GetName())
}

func (dw *DBWriter) HasPod(pod k8sv1.Pod) (bool, error) {
	return dw.has(mergeKey(pod.GetNamespace(), pod.GetName()))
}

func (dw *DBWriter) ValidatePod(pod k8sv1.Pod) (bool, error) {

}

func (dw *DBWriter) ReadPodList() ([]k8sv1.Pod, error) {

}

func (dw *DBWriter) ReadNamespaces() ([]string, error) {

}

func (dw *DBWriter) ReadNames() ([]string, error) {

}

func (dw *DBWriter) ReadPorts() ([]int, error) {

}

func (dw *DBWriter) ReadPod(namespace string, name string) (k8sv1.Pod, error) {

}

func (dw *DBWriter) DeletePod(namespace string, name string) error {
	//keys := dw.read(mergeKey(namespace, name))
}

func (dw *DBWriter) UpdatePod(namespace string, name string, pod k8sv1.Pod) error {
	data, err := dw.read(mergeKey(namespace, name))
	if err == leveldb.ErrNotFound {
		return ErrNotFound
	} else if err != nil {
		return err
	}

	pd := PodDescribe{}
	if err:= pd.UnMarshal(data);err!=nil{
		return err
	}
	newhash := GetHashString(pod)
	if pd.Hashstatus==newhash
}

func (dw *DBWriter) CreatePod(pod k8sv1.Pod) error {
	hostports := make([]int32, 0)
	for _, container := range append(pod.Spec.InitContainers, pod.Spec.Containers...) {
		for _, ports := range container.Ports {
			if ports.HostPort > 0 {
				hostports = append(hostports, ports.HostPort)
			}
		}
	}

	hashValue, err := GetHashString(pod)
	if err != nil {
		return err
	}

	podDescribe := PodDescribe{
		Name:       pod.GetName(),
		Namespace:  pod.GetNamespace(),
		Ports:      hostports,
		Hash:       hashValue,
		Hashstatus: hashValue,
	}

	// write pod describe to DB
	pdBytes, err := podDescribe.Marshal()
	if err != nil {
		return err
	}

	err = dw.write(getKey(pod), pdBytes)
	if err != nil {
		return err
	}

	podBytes, err := pod.Marshal()
	if err != nil {
		return err
	}

	err = dw.write(hashValue, podBytes)
	if err != nil {
		return err
	}

	return dw.AddKey(getKey(pod))
}

func GetHashString(pod k8sv1.Pod) (string, error) {
	data, err := pod.Marshal()
	if err != nil {
		return "", err
	}

	hashBytes := sha1.Sum(data)
	return hex.EncodeToString(hashBytes[:]), nil
}

func (dw *DBWriter) AddKey(newKey string) error {
	kr := dw.KR

	// load data from DB
	datas, err := dw.read(dbKey)
	if err == leveldb.ErrNotFound {
		kr.keys = make([]string, 0)
	} else if err == nil {
		err := kr.UnMarshal(datas)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	// adjest repeat
	for _, key := range kr.keys {
		if key == newKey {
			return nil
		}
	}

	// write to DB
	kr.keys = append(kr.keys, newKey)
	keyBytes, err := kr.Marshal()
	if err != nil {
		return err
	}
	return dw.write(dbKey, keyBytes)
}

func (dw *DBWriter) DeleteKey(oldKey string) error {
	kr := dw.KR

	// load data from DB
	datas, err := dw.read(dbKey)
	if err == leveldb.ErrNotFound {
		kr.keys = make([]string, 0)
	} else if err == nil {
		err := kr.UnMarshal(datas)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	if len(kr.keys) <= 0 {
		return nil
	}

	flag := false
	newKeys := make([]string, len(kr.keys))
	// adjest repeat
	for _, key := range kr.keys {
		if key != oldKey {
			newKeys = append(newKeys, key)
		} else {
			flag = true
		}
	}

	// no change
	if !flag {
		return nil
	}

	kr.keys = newKeys
	keyBytes, err := kr.Marshal()
	if err != nil {
		return err
	}
	return dw.write(dbKey, keyBytes)
}

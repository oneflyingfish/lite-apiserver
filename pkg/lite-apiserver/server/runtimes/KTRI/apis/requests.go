package apis

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"path/filepath"

	options "LiteKube/pkg/lite-apiserver/options/kubeletOptions"

	"k8s.io/klog/v2"
)

type KubeletRuntime struct {
	*options.KubeletOption
	// BackendTimeout int

	// runtime args
	RequestClient *http.Client
	ServerUrl     string
}

func CreateKubeletRuntime(opt *options.KubeletOption) (*KubeletRuntime, error) {
	client := CreateHttpsClient(opt)
	if client == nil {
		return nil, fmt.Errorf("fail to create https-client")
	}

	return &KubeletRuntime{
		opt,
		client,
		fmt.Sprintf("https://%s:%d", opt.Hostname, opt.Port),
	}, nil
}

func InitKubelet(kt *KubeletRuntime) error {
	APIKubelet = kt
	return nil
}

func GetKubelet() *KubeletRuntime {
	return APIKubelet
}

// set APIKubelet by kubeletOptions
func CreateHttpsClient(opt *options.KubeletOption) *http.Client {
	if opt == nil {
		return nil
	}

	certFile, keyFile, isValid := opt.TLSKeyPair.GetTLSKeyPair()
	if !isValid {
		return nil
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		klog.Errorf("fail to load client certificate for kubelet while create https-client, error tip: %s", err.Error())
		return nil
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			Certificates:       []tls.Certificate{cert},
		},
	}

	return &http.Client{Transport: tr}
}

func (kr *KubeletRuntime) DoGenericRequest(req *http.Request) (*http.Response, error) {
	if kr == nil {
		return nil, fmt.Errorf("kubelet options is nil")
	}

	resp, err := (*(kr.RequestClient)).Do(req)

	return resp, err
}

func DoGenericRequest(req *http.Request) (*http.Response, error) {
	return APIKubelet.DoGenericRequest(req)
}

func (kr *KubeletRuntime) GetRequest(path string) (*http.Response, error) {
	if kr == nil {
		return nil, fmt.Errorf("kubelet options is nil")
	}

	req, _ := http.NewRequest(GET, filepath.Join(kr.ServerUrl, path), nil)

	return kr.DoGenericRequest(req)
}

func GetRequest(path string) (*http.Response, error) {
	return APIKubelet.GetRequest(path)
}

func (kr *KubeletRuntime) PutRequest(path string, bodyData []byte) (*http.Response, error) {
	if kr == nil {
		return nil, fmt.Errorf("kubelet options is nil")
	}

	req, _ := http.NewRequest(PUT, filepath.Join(kr.ServerUrl, path), bytes.NewBuffer(bodyData))
	req.Header.Set("Content-Type", "text/plain")
	return kr.DoGenericRequest(req)
}

func PutRequest(path string, bodyData []byte) (*http.Response, error) {
	return APIKubelet.PutRequest(path, bodyData)
}

func (kr *KubeletRuntime) PostRequest(client *http.Client, path string, bodyData []byte) (*http.Response, error) {
	if kr == nil {
		return nil, fmt.Errorf("kubelet options is nil")
	}

	req, _ := http.NewRequest(POST, filepath.Join(kr.ServerUrl, path), bytes.NewBuffer(bodyData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return kr.DoGenericRequest(req)
}

func PostRequest(client *http.Client, path string, bodyData []byte) (*http.Response, error) {
	return APIKubelet.PostRequest(client, path, bodyData)
}

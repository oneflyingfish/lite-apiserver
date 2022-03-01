package tls

import (
	"LiteKube/pkg/lite-apiserver/cert"
	"LiteKube/pkg/lite-apiserver/server/runtimes/ReturnFormat"
	"fmt"
	"net/http"
)

func TLSHandleFunc(caTLSKeyPair *cert.TLSKeyPair, port int, checkPrivilege bool) func(w http.ResponseWriter, r *http.Request) (int, error) {
	return func(w http.ResponseWriter, r *http.Request) (status int, e error) {
		if caTLSKeyPair == nil && checkPrivilege {
			return http.StatusMethodNotAllowed, fmt.Errorf("this work is not allowed by http")
		} else {
			caCert, caKey, ok := caTLSKeyPair.GetTLSKeyPairCertificate()
			if !ok {
				return http.StatusInternalServerError, fmt.Errorf("fail to load CA informations")
			}

			clientCertBase64, clientKeyBase64, err := cert.CreateClientCertBase64(caCert, caKey)
			if err != nil {
				return http.StatusInternalServerError, fmt.Errorf("error occured while generate certificate for client, tips: %s", err)
			}

			info := TLSInfo{
				CACert:     string(caTLSKeyPair.GetCertBase64()),
				ClientCert: string(clientCertBase64),
				ClientKey:  string(clientKeyBase64),
				Port:       port,
			}

			fmt.Fprint(w, TLSReturn(info, r.URL.Query().Get("format") != ReturnFormat.Json))
			return http.StatusOK, nil
		}
	}
}

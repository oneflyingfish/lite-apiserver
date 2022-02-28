// some urls for debug

package ServerHandlers

import (
	"LiteKube/pkg/lite-apiserver/cert"
	"fmt"
	"net/http"
	"strings"

	"k8s.io/klog/v2"
)

var TLSReturnString string = "Create client.pem and client-key.pem on your disk with the following command: \ncat >ca.pem<<EOF\n%sEOF\n\ncat >client.pem<<EOF\n%sEOF\n\ncat >client-key.pem<<EOF\n%sEOF\n\nTips: you can run `openssl pkcs12 -export -clcerts -in client.pem -inkey client-key.pem -out client.p12` in cmd.exe to create a certificate for windows"
var prefixUrl string = "/debug"

type HandleFunc func(w http.ResponseWriter, r *http.Request) error

type DebugHandle struct {
	prefixUrl string
	handles   map[string]HandleFunc
	caKeyPair *cert.TLSKeyPair
}

func NewDebugHandle(caTLSKeyPair *cert.TLSKeyPair) DebugHandle {
	return DebugHandle{
		prefixUrl: prefixUrl,
		handles:   createDebugHandle(caTLSKeyPair),
		caKeyPair: caTLSKeyPair,
	}
}

func createDebugHandle(caTLSKeyPair *cert.TLSKeyPair) map[string]HandleFunc {
	handles := make(map[string]HandleFunc)
	handles["/hello"] = func(w http.ResponseWriter, r *http.Request) error {
		fmt.Fprintf(w, "Hello to see you, LiteKube is here!\n")
		return nil
	}

	handles["/about"] = func(w http.ResponseWriter, r *http.Request) error {
		http.RedirectHandler("https://github.com/kubesys/LiteKube", http.StatusTemporaryRedirect).ServeHTTP(w, r)
		return nil
	}

	handles["/tls"] = func(w http.ResponseWriter, r *http.Request) error {
		if caTLSKeyPair == nil {
			w.WriteHeader(http.StatusMethodNotAllowed) // http status: 400
			fmt.Fprintf(w, "this work is not allowed by http\n")
		} else {
			caCert, caKey, ok := caTLSKeyPair.GetTLSKeyPairCertificate()
			if !ok {
				return fmt.Errorf("fail to load CA informations")
			}

			clientCertBase64, clientKeyBase64, err := cert.CreateClientCertBase64(caCert, caKey)
			if err != nil {
				return fmt.Errorf("error occured while generate certificate for client, tips: %s", err)
			}

			fmt.Fprintf(w, TLSReturnString, caTLSKeyPair.GetCertBase64(), clientCertBase64, clientKeyBase64)
		}
		return nil
	}

	return handles
}

func (handle DebugHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer w.Header().Set("WWW-Authenticate", `Basic realm="mydomain"`)

	klog.Infof("Get request %s", r.URL.Path)

	if !strings.HasPrefix(r.URL.Path, handle.prefixUrl) {
		w.WriteHeader(http.StatusNotFound) // http status: 404
		return
	}

	deal, ok := handle.handles[r.URL.Path[len(handle.prefixUrl):]]
	if !ok || !strings.HasPrefix(r.URL.Path, handle.prefixUrl) {
		w.WriteHeader(http.StatusNotFound) // http status: 404
		fmt.Fprintf(w, "page is not found\n")
	} else {
		if err := deal(w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError) // http status: 500
			fmt.Fprintf(w, "some errors accure while deal with your request for %s, error: %s\n", r.URL.Path, err.Error())
		}
	}
}

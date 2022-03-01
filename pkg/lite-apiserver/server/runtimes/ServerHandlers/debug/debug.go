package debug

import (
	"LiteKube/pkg/common"
	"LiteKube/pkg/lite-apiserver/cert"
	"LiteKube/pkg/lite-apiserver/server/runtimes/ReturnFormat"
	"LiteKube/pkg/lite-apiserver/server/runtimes/ServerHandlers/tls"
	"fmt"
	"net/http"
	"strings"

	"k8s.io/klog/v2"
)

//var TLSReturnString string = "Create client.pem and client-key.pem on your disk with the following command: \ncat >ca.pem<<EOF\n%sEOF\n\ncat >client.pem<<EOF\n%sEOF\n\ncat >client-key.pem<<EOF\n%sEOF\n\nTips: you can run `openssl pkcs12 -export -clcerts -in client.pem -inkey client-key.pem -out client.p12` in cmd.exe to create a certificate for windows"
var prefixUrl string = "/debug"

type HandleFunc func(w http.ResponseWriter, r *http.Request) (int, error)

type DebugHandle struct {
	prefixUrl string
	handles   map[string]HandleFunc
	caKeyPair *cert.TLSKeyPair
}

func NewDebugHandle(caTLSKeyPair *cert.TLSKeyPair, port int) DebugHandle {
	return DebugHandle{
		prefixUrl: prefixUrl,
		handles:   createDebugHandle(caTLSKeyPair, port),
		caKeyPair: caTLSKeyPair,
	}
}

func createDebugHandle(caTLSKeyPair *cert.TLSKeyPair, port int) map[string]HandleFunc {
	handles := make(map[string]HandleFunc)
	handles["/hello"] = func(w http.ResponseWriter, r *http.Request) (int, error) {
		fmt.Fprintf(w, "Hello to see you, LiteKube is here!\n")
		return http.StatusOK, nil
	}

	handles["/health"] = func(w http.ResponseWriter, r *http.Request) (int, error) {
		fmt.Fprintf(w, "ok\n")
		return http.StatusOK, nil
	}

	handles["/about"] = func(w http.ResponseWriter, r *http.Request) (int, error) {
		http.RedirectHandler("https://github.com/kubesys/LiteKube", http.StatusTemporaryRedirect).ServeHTTP(w, r)
		return http.StatusTemporaryRedirect, nil
	}

	handles["/tls"] = tls.TLSHandleFunc(caTLSKeyPair, port, false)

	return handles
}

func (handle DebugHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	klog.Infof("Get request %s", r.URL.Path)
	if !strings.HasPrefix(r.URL.Path, handle.prefixUrl) {
		w.WriteHeader(http.StatusNotFound) // http status: 404
		return
	}

	deal, ok := handle.handles[r.URL.Path[len(handle.prefixUrl):]]
	if !ok || !strings.HasPrefix(r.URL.Path, handle.prefixUrl) {
		w.WriteHeader(http.StatusNotFound) // http status: 404
		fmt.Fprint(w, common.ErrorString("Page is not found", r.URL.Query().Get("format") == ReturnFormat.Raw))
	} else {
		statusCode, err := deal(w, r)
		if statusCode != http.StatusOK {
			w.WriteHeader(statusCode)
		}
		if err != nil {
			fmt.Fprint(w, common.ErrorString(err.Error(), r.URL.Query().Get("format") == ReturnFormat.Raw))
		}
	}
}

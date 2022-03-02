package tls

import (
	"LiteKube/pkg/lite-apiserver/cert"
	"fmt"
	"net/http"

	"k8s.io/klog/v2"

	"github.com/emicklei/go-restful"
)

func RegisterWebService(container *restful.Container, caKeyPair *cert.TLSKeyPair, port int, checkPrivilege bool) {
	ws := new(restful.WebService)
	ws.Path("/tls")

	ws.Route(ws.GET("").To(TLSResponse(caKeyPair, port, checkPrivilege)))
	container.Add(ws)
}

func TLSResponse(caTLSKeyPair *cert.TLSKeyPair, port int, checkPrivilege bool) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		klog.Info("one request for https certificate")

		if (caTLSKeyPair == nil || request.Request.TLS == nil) && checkPrivilege {
			response.AddHeader("Content-Type", "text/plain")
			response.WriteErrorString(http.StatusInternalServerError, "this work is not allowed by HTTP, you can access by https, which means you may need to seek admin for privilege.")
			return
		} else {
			caCert, caKey, ok := caTLSKeyPair.GetTLSKeyPairCertificate()
			if !ok {
				response.AddHeader("Content-Type", "text/plain")
				response.WriteErrorString(http.StatusInternalServerError, "fail to load CA informations")
				return
			}

			clientCertBase64, clientKeyBase64, err := cert.CreateClientCertBase64(caCert, caKey)
			if err != nil {
				response.AddHeader("Content-Type", "text/plain")
				response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("error occured while generate certificate for client, tips: %s", err))
				return
			}

			info := TLSInfo{
				CACert:     string(caTLSKeyPair.GetCertBase64()),
				ClientCert: string(clientCertBase64),
				ClientKey:  string(clientKeyBase64),
				Port:       port,
			}

			if request.QueryParameter("format") != "json" {
				response.AddHeader("Content-Type", "text/html")
				response.Write([]byte(TLSReturn(info, true)))
			} else {
				response.WriteHeaderAndJson(http.StatusOK, info, "application/json")
			}

			klog.Info("success to return https certificate")
		}
	}
}

// func TLSHandleFunc(caTLSKeyPair *cert.TLSKeyPair, port int, checkPrivilege bool) func(w http.ResponseWriter, r *http.Request) (int, error) {
// 	return func(w http.ResponseWriter, r *http.Request) (status int, e error) {
// 		if caTLSKeyPair == nil && checkPrivilege {
// 			return http.StatusMethodNotAllowed, fmt.Errorf("this work is not allowed by http")
// 		} else {
// 			caCert, caKey, ok := caTLSKeyPair.GetTLSKeyPairCertificate()
// 			if !ok {
// 				return http.StatusInternalServerError, fmt.Errorf("fail to load CA informations")
// 			}

// 			clientCertBase64, clientKeyBase64, err := cert.CreateClientCertBase64(caCert, caKey)
// 			if err != nil {
// 				return http.StatusInternalServerError, fmt.Errorf("error occured while generate certificate for client, tips: %s", err)
// 			}

// 			info := TLSInfo{
// 				CACert:     string(caTLSKeyPair.GetCertBase64()),
// 				ClientCert: string(clientCertBase64),
// 				ClientKey:  string(clientKeyBase64),
// 				Port:       port,
// 			}

// 			fmt.Fprint(w, TLSReturn(info, r.URL.Query().Get("format") != ReturnFormat.Json))
// 			return http.StatusOK, nil
// 		}
// 	}
// }

package debug

import (
	"LiteKube/pkg/common"
	"LiteKube/pkg/lite-apiserver/cert"
	"LiteKube/pkg/lite-apiserver/server/runtimes/ServerHandlers/tls"
	"net/http"

	"github.com/emicklei/go-restful"
)

//var TLSReturnString string = "Create client.pem and client-key.pem on your disk with the following command: \ncat >ca.pem<<EOF\n%sEOF\n\ncat >client.pem<<EOF\n%sEOF\n\ncat >client-key.pem<<EOF\n%sEOF\n\nTips: you can run `openssl pkcs12 -export -clcerts -in client.pem -inkey client-key.pem -out client.p12` in cmd.exe to create a certificate for windows"
//var prefixUrl string = "/debug"

type HandleFunc func(w http.ResponseWriter, r *http.Request) (int, error)

type DebugHandle struct {
	//prefixUrl string
	//handles   map[string]HandleFunc
	caKeyPair *cert.TLSKeyPair
	port      int
}

func NewDebugHandle(caTLSKeyPair *cert.TLSKeyPair, port int) DebugHandle {
	return DebugHandle{
		//prefixUrl: prefixUrl,
		//handles:   createDebugHandle(caTLSKeyPair, port),
		caKeyPair: caTLSKeyPair,
		port:      port,
	}
}

func (handle DebugHandle) HealthGet(request *restful.Request, response *restful.Response) {
	response.AddHeader("Content-Type", "text/plain")
	response.Write([]byte("ok\n"))
}

func (handle DebugHandle) HelloGet(request *restful.Request, response *restful.Response) {
	response.AddHeader("Content-Type", "text/plain")
	response.Write([]byte("Hello to see you, LiteKube is here!\n"))
}

func (handle DebugHandle) TLSGet(request *restful.Request, response *restful.Response) {
	tls.TLSResponse(handle.caKeyPair, handle.port, false)(request, response)
}

func (handle DebugHandle) AboutGet(request *restful.Request, response *restful.Response) {
	http.RedirectHandler("https://github.com/kubesys/LiteKube", http.StatusTemporaryRedirect).ServeHTTP(response.ResponseWriter, request.Request)
}

func (handle DebugHandle) RegisterWebService(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/debug")

	ws.Route(ws.GET("").To(common.HelpWebService(ws, map[string]interface{}{"message": "/debug/... will only be enable while --debug=true"})))

	ws.Route(ws.GET("/health").To(handle.HealthGet))
	ws.Route(ws.GET("/hello").To(handle.HelloGet))
	ws.Route(ws.GET("/tls").To(handle.TLSGet))
	// Doc("describe: Get X.509 certificate for https access").
	// Param(ws.PathParameter("format", fmt.Sprintf("%s,%s", ReturnFormat.Json, ReturnFormat.Raw)).DataType("string")).
	// Do(handle.returns200, handle.returns500))
	ws.Route(ws.GET("/about").To(handle.AboutGet))

	container.Add(ws)
}

// func (handle DebugHandle) returns200(b *restful.RouteBuilder) {
// 	klog.Info("---------222222-------")
// 	b.Returns(http.StatusOK, "OK111111111111111111111", "sucess")
// }

// func (handle DebugHandle) returns500(b *restful.RouteBuilder) {
// 	klog.Info("--------55555--------")
// 	b.Returns(http.StatusInternalServerError, "bad params", nil)
// }

// func createDebugHandle(caTLSKeyPair *cert.TLSKeyPair, port int) map[string]HandleFunc {
// 	handles := make(map[string]HandleFunc)
// 	handles["/hello"] = func(w http.ResponseWriter, r *http.Request) (int, error) {
// 		fmt.Fprintf(w, "Hello to see you, LiteKube is here!\n")
// 		return http.StatusOK, nil
// 	}

// 	handles["/health"] = func(w http.ResponseWriter, r *http.Request) (int, error) {
// 		fmt.Fprintf(w, "ok\n")
// 		return http.StatusOK, nil
// 	}

// 	handles["/about"] = func(w http.ResponseWriter, r *http.Request) (int, error) {
// 		http.RedirectHandler("https://github.com/kubesys/LiteKube", http.StatusTemporaryRedirect).ServeHTTP(w, r)
// 		return http.StatusTemporaryRedirect, nil
// 	}

// 	handles["/tls"] = tls.TLSHandleFunc(caTLSKeyPair, port, false)

// 	return handles
// }

// func (handle DebugHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	klog.Infof("Get request %s", r.URL.Path)
// 	if !strings.HasPrefix(r.URL.Path, handle.prefixUrl) {
// 		w.WriteHeader(http.StatusNotFound) // http status: 404
// 		return
// 	}

// 	deal, ok := handle.handles[r.URL.Path[len(handle.prefixUrl):]]
// 	if !ok || !strings.HasPrefix(r.URL.Path, handle.prefixUrl) {
// 		w.WriteHeader(http.StatusNotFound) // http status: 404
// 		fmt.Fprint(w, common.ErrorString("Page is not found", r.URL.Query().Get("format") == ReturnFormat.Raw))
// 	} else {
// 		statusCode, err := deal(w, r)
// 		if statusCode != http.StatusOK {
// 			w.WriteHeader(statusCode)
// 		}
// 		if err != nil {
// 			fmt.Fprint(w, common.ErrorString(err.Error(), r.URL.Query().Get("format") == ReturnFormat.Raw))
// 		}
// 	}
// }

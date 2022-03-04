package debug

import (
	"LiteKube/pkg/lite-apiserver/cert"
	"LiteKube/pkg/lite-apiserver/server/runtimes/ServerHandlers/tls"
	"LiteKube/pkg/restfulenhance"
	"net/http"

	"github.com/emicklei/go-restful"
)

type DebugHandle struct {
	caKeyPair *cert.TLSKeyPair
	port      int
}

func NewDebugHandle(caTLSKeyPair *cert.TLSKeyPair, port int) DebugHandle {
	return DebugHandle{
		caKeyPair: caTLSKeyPair,
		port:      port,
	}
}

func (handle DebugHandle) RegisterWebService(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/debug")

	ws.Route(ws.GET("").To(restfulenhance.HelpWebService(ws, map[string]interface{}{"message": "/debug/... will only be enable while --debug=true"})))

	ws.Route(ws.GET("/health").To(handle.HealthGet))
	ws.Route(ws.GET("/hello").To(handle.HelloGet))
	ws.Route(ws.GET("/tls").To(handle.TLSGet))
	ws.Route(ws.GET("/about").To(handle.AboutGet))

	container.Add(ws)
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

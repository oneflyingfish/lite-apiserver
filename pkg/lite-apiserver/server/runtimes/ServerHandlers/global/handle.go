package global

import (
	"LiteKube/pkg/lite-apiserver/cert"
	"LiteKube/pkg/lite-apiserver/server/runtimes/ServerHandlers/tls"
	"LiteKube/pkg/restfulenhance"
	litekubeVersion "LiteKube/pkg/version"
	"net/http"

	"github.com/emicklei/go-restful"
)

func RegisterWebService(container *restful.Container, caTLSKeyPair *cert.TLSKeyPair, port int) {
	ws := new(restful.WebService)
	ws.Path("")

	ws.Route(ws.GET("").To(restfulenhance.HelpContainer(container, nil)))
	ws.Route(ws.GET("/healthz").To(healthz))
	ws.Route(ws.GET("/version").To(version))
	ws.Route(ws.GET("/tls").To(tls.TLSResponse(caTLSKeyPair, port, true)))
	container.Add(ws)
}

func healthz(request *restful.Request, response *restful.Response) {
	// adjust if kubelet ok
	response.AddHeader("Content-Type", "text/plain")
	response.Write([]byte("ok\n"))
}

func version(request *restful.Request, response *restful.Response) {
	response.WriteHeaderAndJson(http.StatusOK, litekubeVersion.Get(), "application/json")
}

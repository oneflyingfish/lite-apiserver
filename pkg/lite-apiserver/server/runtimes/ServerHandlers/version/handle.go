package version

import (
	litekubeVersion "LiteKube/pkg/version"
	"net/http"

	"github.com/emicklei/go-restful"
)

func RegisterWebService(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/version")

	ws.Route(ws.GET("").To(responseRootGet))
	container.Add(ws)
}



package ServerHandlers

import (
	"LiteKube/pkg/common"

	"github.com/emicklei/go-restful"
)

func GlobalRegisterWebService(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("")

	ws.Route(ws.GET("").To(common.HelpContainer(container, nil)))
	container.Add(ws)
}

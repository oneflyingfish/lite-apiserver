package api

import (
	v1 "LiteKube/pkg/lite-apiserver/server/runtimes/ServerHandlers/api/v1"
	"LiteKube/pkg/restfulenhance"
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
)

type Address struct {
	ClientCIDR    string `json:"clientCIDR"`
	ServerAddress string `json:"serverAddress"`
}

type API struct {
	Kind           string                         `json:"kind"`
	Versions       []string                       `json:"versions"`
	ServerAddress  *Address                       `json:"serverAddressByClientCIDRs"`
	WebServiceNode *restfulenhance.WebServiceNode `json:"-"`
}

func NewAPI(serverHostname string, serverPort int) *API {
	return &API{
		Kind:     "APIVersions",
		Versions: make([]string, 0),
		ServerAddress: &Address{
			ClientCIDR:    "0.0.0.0/0",
			ServerAddress: fmt.Sprintf("%s:%d", serverHostname, serverPort),
		},
		WebServiceNode: restfulenhance.NewWebServiceNode(nil, "/api"), // init "/api" as root node
	}
}

func (api *API) AddRoutes() {
	api.WebServiceNode.Ws.Route(api.WebServiceNode.Ws.GET("").To(api.list))

}

func (api *API) list(request *restful.Request, response *restful.Response) {
	response.WriteHeaderAndJson(http.StatusOK, api, "application/json")
}

func (api *API) AddChildWebService() {
	wsNode := restfulenhance.NewWebServiceNode(api.WebServiceNode, "/v1")
	v1.NewAPIV1(wsNode).InitCurrentWebService()
}

func (api *API) InitCurrentWebService() {
	api.AddRoutes()
	api.AddChildWebService()

	for _, node := range api.WebServiceNode.ChildNodes {
		if node != nil && len(node.Path) > 1 {
			api.Versions = append(api.Versions, node.Path[1:])
		}
	}
}

func (api *API) RegisteredWebServices(container *restful.Container) {
	api.InitCurrentWebService()
	api.WebServiceNode.RegisterToContainer(container)
}

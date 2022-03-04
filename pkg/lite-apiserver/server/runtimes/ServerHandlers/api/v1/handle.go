package v1

import (
	"LiteKube/pkg/lite-apiserver/describe"
	"LiteKube/pkg/lite-apiserver/server/runtimes/ServerHandlers/api/v1/namespaces"
	"LiteKube/pkg/lite-apiserver/server/runtimes/ServerHandlers/api/v1/namespaces/pods"
	"LiteKube/pkg/restfulenhance"
	"net/http"

	"github.com/emicklei/go-restful"
)

type APIV1 struct {
	Kind           string                         `json:"kind"`
	GroupVersion   string                         `json:"groupVersion"`
	Resources      []describe.Resource            `json:"resources"`
	WebServiceNode *restfulenhance.WebServiceNode `json:"-"`
}

func NewAPIV1(node *restfulenhance.WebServiceNode) *APIV1 {
	return &APIV1{
		Kind:           "APIResourceList",
		GroupVersion:   "v1",
		Resources:      make([]describe.Resource, 0),
		WebServiceNode: node,
	}
}

func (api *APIV1) AddRoutes() {
	api.WebServiceNode.Ws.Route(api.WebServiceNode.Ws.GET("").To(api.list))
	api.WebServiceNode.Ws.Route(api.WebServiceNode.Ws.GET("/pods").To(api.list))
}

func (api *APIV1) list(request *restful.Request, response *restful.Response) {
	response.WriteHeaderAndJson(http.StatusOK, api, "application/json")
}

func (api *APIV1) AddChildWebService() {
	wsNode := restfulenhance.NewWebServiceNode(api.WebServiceNode, "/namespaces")
	namespaces.NewNamespaces(wsNode).InitCurrentWebService()

	api.Resources = append(api.Resources, namespaces.Describe)
	api.Resources = append(api.Resources, pods.Describe)
}

func (api *APIV1) InitCurrentWebService() {
	api.AddRoutes()
	api.AddChildWebService()
}

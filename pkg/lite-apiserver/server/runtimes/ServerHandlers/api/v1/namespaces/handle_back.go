package namespaces

// import (
// 	"LiteKube/pkg/lite-apiserver/describe"
// 	"LiteKube/pkg/restfulenhance"
// 	"net/http"

// 	"github.com/emicklei/go-restful"
// )

// var Describe describe.Resource = describe.Resource{
// 	Name:         "namespaces",
// 	SingularName: "",
// 	Namespaced:   false,
// 	Kind:         "Namespace",
// 	Verbs: []string{
// 		"list",
// 	},
// }

// type Namespaces struct {
// 	Kind           string                         `json:"kind"`
// 	ApiVersion     string                         `json:"apiVersion"`
// 	Items          []describe.NamespaceItem       `json:"items"`
// 	WebServiceNode *restfulenhance.WebServiceNode `json:"-"`
// }

// func NewAPIV1(node *restfulenhance.WebServiceNode) *Namespaces {
// 	return &Namespaces{
// 		Kind:           "NamespaceList",
// 		ApiVersion:     "v1",
// 		Items:          make([]describe.NamespaceItem, 0),
// 		WebServiceNode: node,
// 	}
// }

// func (ns *Namespaces) AddRoutes() {
// 	ns.WebServiceNode.Ws.Route(ns.WebServiceNode.Ws.GET("{}").To(ns.list))
// }

// func (ns *Namespaces) list(request *restful.Request, response *restful.Response) {
// 	response.WriteHeaderAndJson(http.StatusOK, ns, "application/json")
// }

// func (ns *Namespaces) AddChildWebService() {
// 	//wsNode := restfulenhance.NewWebServiceNode(ns.WebServiceNode, "/pods")

// 	// api.Resources = append(api.Resources, namespaces.Describe)
// 	// api.Resources = append(api.Resources, pods.Describe)
// }

// func (ns *Namespaces) InitCurrentWebService() {
// 	ns.AddRoutes()
// 	ns.AddChildWebService()
// }

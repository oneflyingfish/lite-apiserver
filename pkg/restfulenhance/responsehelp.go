package restfulenhance

import (
	"net/http"
	"sort"

	"github.com/emicklei/go-restful"
)

func HelpContainer(container *restful.Container, addition map[string]interface{}) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		jsonMap := make(map[string]interface{})
		paths := make([]string, 0)
		for _, ws := range container.RegisteredWebServices() {
			for _, route := range ws.Routes() {
				if len(route.Path) > 0 && route.Path != "/" {
					paths = append(paths, route.Path)
				}
			}
		}
		sort.Slice(paths, func(i, j int) bool {
			return paths[i] < paths[j]
		})
		jsonMap["paths"] = paths

		for key, value := range addition {
			jsonMap[key] = value
		}

		response.WriteHeaderAndJson(http.StatusOK, jsonMap, "application/json")
	}
}

func HelpWebService(ws *restful.WebService, addition map[string]interface{}) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		jsonMap := make(map[string]interface{})
		paths := make([]string, 0)
		for _, route := range ws.Routes() {
			if len(route.Path) > 0 && route.Path != "/" {
				paths = append(paths, route.Path)
			}
		}

		sort.Slice(paths, func(i, j int) bool {
			return paths[i] < paths[j]
		})
		jsonMap["paths"] = paths

		for key, value := range addition {
			jsonMap[key] = value
		}
		response.WriteHeaderAndJson(http.StatusOK, jsonMap, "application/json")
	}
}

func HelpWebServiceNode(node *WebServiceNode, addition map[string]interface{}) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		jsonMap := make(map[string]interface{})

		// get all paths
		paths := GetWebServiceNodePaths(node)
		sort.Slice(paths, func(i, j int) bool {
			return paths[i] < paths[j]
		})

		jsonMap["paths"] = paths[1:] // delete root path from paths

		for key, value := range addition {
			jsonMap[key] = value
		}
		response.WriteHeaderAndJson(http.StatusOK, jsonMap, "application/json")
	}
}

func GetWebServiceNodePaths(node *WebServiceNode) []string {
	// invalid node
	if node == nil || len(node.Path) < 2 {
		return make([]string, 0)
	}

	paths := make([]string, 1)

	for _, route := range node.Ws.Routes() {
		paths = append(paths, route.Path)
	}

	for _, node := range node.ChildNodes {
		paths = append(paths, GetWebServiceNodePaths(node)...)
	}

	return paths
}

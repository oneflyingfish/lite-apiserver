package common

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

		jsonMap["paths"] = paths

		for key, value := range addition {
			jsonMap[key] = value
		}
		response.WriteHeaderAndJson(http.StatusOK, jsonMap, "application/json")
	}
}

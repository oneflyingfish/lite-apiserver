package restfulenhance

import (
	"path/filepath"

	"github.com/emicklei/go-restful"
)

type WebServiceNode struct {
	Path       string
	Ws         *restful.WebService
	ParentNode *WebServiceNode
	ChildNodes []*WebServiceNode
}

// parentNode==nil means a root node
func NewWebServiceNode(parentNode *WebServiceNode, path string) *WebServiceNode {
	if len(path) < 2 || path[0] != '/' {
		return nil
	}

	wsNode := &WebServiceNode{
		Path:       path,
		Ws:         new(restful.WebService),
		ParentNode: parentNode,
		ChildNodes: make([]*WebServiceNode, 0),
	}

	if parentNode == nil {
		wsNode.Ws.Path(path) // root node
	} else {
		wsNode.Ws.Path(wsNode.CompletePath())
		parentNode.ChildNodes = append(parentNode.ChildNodes, wsNode)
	}

	return wsNode
}

// Get path from root-node to current node
func (wsNode *WebServiceNode) CompletePath() string {
	path := wsNode.Path

	pNode := wsNode.ParentNode
	for pNode != nil {
		path = filepath.Join(pNode.Path, path)
		pNode = pNode.ParentNode
	}
	return path
}

// register one WebserviceTree to restful.Container by the root-node
func (wsNode *WebServiceNode) RegisterToContainer(container *restful.Container) {
	// invalid node
	if wsNode == nil || len(wsNode.Path) < 2 {
		return
	}

	container.Add(wsNode.Ws)

	for _, node := range wsNode.ChildNodes {
		node.RegisterToContainer(container)
	}
}

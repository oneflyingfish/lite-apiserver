package pods

import "LiteKube/pkg/lite-apiserver/describe"

var Describe describe.Resource = describe.Resource{
	Name:         "pods",
	SingularName: "",
	Namespaced:   true,
	Kind:         "Pod",
	Verbs: []string{
		"create",
		"delete",
		"deletecollection",
		"get",
		"list",
		"patch",
		"update",
		"watch",
	},
}

package schema

import (
	"github.com/rancher/norman/types"
	m "github.com/rancher/norman/types/mapper"
	"github.com/rancher/types/apis/project.cattle.io/v3"
	"github.com/rancher/types/factory"
	"github.com/rancher/types/mapper"
	"k8s.io/api/core/v1"
)

var (
	Version = types.APIVersion{
		Version:          "v3",
		Group:            "cluster.cattle.io",
		Path:             "/v3/cluster",
		SubContext:       true,
		SubContextSchema: "/v3/schemas/cluster",
	}

	Schemas = factory.Schemas(&Version).
		Init(namespaceTypes).
		Init(nodeTypes).
		Init(volumeTypes)
)

func namespaceTypes(schemas *types.Schemas) *types.Schemas {
	return schemas.
		AddMapperForType(&Version, v1.NamespaceSpec{},
			&m.Drop{Field: "finalizers"},
		).
		AddMapperForType(&Version, v1.Namespace{},
			&m.AnnotationField{Field: "description"},
			&m.AnnotationField{Field: "projectId"},
			&m.Drop{Field: "status"},
		).
		MustImport(&Version, v1.Namespace{}, struct {
			Description string `json:"description"`
			ProjectID   string `norman:"type=reference[/v3/schemas/project]"`
		}{})
}

func nodeTypes(schemas *types.Schemas) *types.Schemas {
	return NodeTypes(&Version, schemas)
}

func NodeTypes(version *types.APIVersion, schemas *types.Schemas) *types.Schemas {
	return schemas.
		AddMapperForType(version, v1.NodeStatus{},
			&mapper.NodeAddressMapper{},
			&mapper.OSInfo{},
			&m.Drop{Field: "addresses"},
			&m.Drop{Field: "daemonEndpoints"},
			&m.Drop{Field: "images"},
			&m.Drop{Field: "nodeInfo"},
			&m.Move{From: "conditions", To: "nodeConditions"},
			&m.Drop{Field: "phase"},
			&m.SliceToMap{Field: "volumesAttached", Key: "devicePath"},
		).
		AddMapperForType(version, v1.NodeSpec{},
			&m.Drop{Field: "externalID"},
			&m.Drop{Field: "configSource"},
			&m.Move{From: "providerID", To: "providerId"},
			&m.Move{From: "podCIDR", To: "podCidr"},
			m.Access{Fields: map[string]string{
				"podCidr":       "r",
				"providerId":    "r",
				"taints":        "ru",
				"unschedulable": "ru",
			}}).
		AddMapperForType(version, v1.Node{},
			&m.AnnotationField{Field: "description"},
			&m.AnnotationField{Field: "publicEndpoints", List: true},
			&m.Embed{Field: "status"},
		).
		MustImport(version, v1.NodeStatus{}, struct {
			IPAddress string
			Hostname  string
			Info      NodeInfo
		}{}).
		MustImport(version, v3.PublicEndpoint{}).
		MustImport(version, v1.Node{}, struct {
			Description     string `json:"description"`
			PublicEndpoints string `json:"publicEndpoints" norman:"type=array[publicEndpoint],nocreate,noupdate"`
		}{})
}

func volumeTypes(schemas *types.Schemas) *types.Schemas {
	return schemas.
		MustImport(&Version, v1.PersistentVolume{})
}

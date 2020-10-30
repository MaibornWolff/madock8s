package adapter

import (
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
)

type UInstance struct {
	Name      string
	Address   string
	Config    map[string]string
	BaseURL   string
	JSONURL   string
	Endpoints map[string][]Endpoint
}

type Endpoint struct {
	Path   string
	Method string
	Op     *spec.Operation
}

func GetSwaggerSpec(jsonURL string) (map[string][]Endpoint, error) {
	doc, err := loads.Spec(jsonURL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load spec")
	}

	return ExtractEndpoints(doc), nil
}

func ExtractEndpoints(doc *loads.Document) map[string][]Endpoint {
	endpoints := make(map[string][]Endpoint)
	for path, pathItem := range doc.Spec().Paths.Paths {
		methods := extractMethodsFrom(pathItem)
		endpoints[path] = methods
	}
	return endpoints
}

func extractMethodsFrom(pathItem spec.PathItem) []Endpoint {
	methods := []Endpoint{}
	if pathItem.Get != nil {
		methods = append(methods, Endpoint{Method: "GET", Op: pathItem.Get})
	}
	if pathItem.Put != nil {
		methods = append(methods, Endpoint{Method: "PUT", Op: pathItem.Put})
	}
	if pathItem.Post != nil {
		methods = append(methods, Endpoint{Method: "POST", Op: pathItem.Post})
	}
	if pathItem.Delete != nil {
		methods = append(methods, Endpoint{Method: "DELETE", Op: pathItem.Delete})
	}
	if pathItem.Options != nil {
		methods = append(methods, Endpoint{Method: "OPTIONS", Op: pathItem.Options})
	}
	if pathItem.Head != nil {
		methods = append(methods, Endpoint{Method: "HEAD", Op: pathItem.Head})
	}
	if pathItem.Patch != nil {
		methods = append(methods, Endpoint{Method: "PATCH", Op: pathItem.Patch})
	}
	return methods
}

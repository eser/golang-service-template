package openapi

import (
	"github.com/eser/go-service/pkg/bliss/httpservice"
	"github.com/getkin/kin-openapi/openapi3"
)

func GenerateOpenApiSpec(identity *ApiIdentity, routes *httpservice.Router) any {
	spec := &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:      identity.name,
			Version:    identity.version,
			Extensions: map[string]interface{}{},
		},
		Components: &openapi3.Components{
			Schemas:    make(openapi3.Schemas),
			Extensions: map[string]interface{}{},
		},
		Paths:      &openapi3.Paths{},
		Extensions: map[string]interface{}{},
	}

	for _, route := range routes.Routes {
		operation := &openapi3.Operation{}
		operation.OperationID = route.Spec.OperationId
		operation.Summary = route.Spec.Summary
		operation.Description = route.Spec.Description
		operation.Tags = route.Spec.Tags
		operation.Deprecated = route.Spec.Deprecated

		path := &openapi3.PathItem{}
		path.SetOperation(route.Pattern.Method, operation)

		spec.Paths.Set(route.Pattern.Path, path)
	}

	return spec
}

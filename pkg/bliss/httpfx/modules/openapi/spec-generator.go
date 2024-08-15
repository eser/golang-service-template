package openapi

import (
	"github.com/eser/go-service/pkg/bliss/httpfx"
	"github.com/getkin/kin-openapi/openapi3"
)

func GenerateOpenApiSpec(identity *ApiIdentity, routes httpfx.Router) any {
	spec := &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:      identity.name,
			Version:    identity.version,
			Extensions: map[string]any{},
		},
		Components: &openapi3.Components{
			Schemas:    make(openapi3.Schemas),
			Extensions: map[string]any{},
		},
		Paths:      &openapi3.Paths{},
		Extensions: map[string]any{},
	}

	for _, route := range routes.GetRoutes() {
		operation := &openapi3.Operation{}
		operation.OperationID = route.Spec.OperationId
		operation.Summary = route.Spec.Summary
		operation.Description = route.Spec.Description
		operation.Tags = route.Spec.Tags
		operation.Deprecated = route.Spec.Deprecated

		operation.Responses = &openapi3.Responses{}

		for _, response := range route.Spec.Responses {
			description := ""

			operation.AddResponse(response.StatusCode, &openapi3.Response{
				Description: &description,
				Content:     openapi3.Content{},
				Headers:     openapi3.Headers{},
				Links:       openapi3.Links{},
			})
		}

		path := &openapi3.PathItem{}
		path.SetOperation(route.Pattern.Method, operation)

		spec.Paths.Set(route.Pattern.Path, path)
	}

	return spec
}

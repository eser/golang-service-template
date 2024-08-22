package httpfx

import (
	"net/http"

	"github.com/eser/go-service/pkg/bliss/httpfx/uris"
)

type RouteOpenApiSpecRequest struct {
	Model any
}

type RouteOpenApiSpecResponse struct {
	StatusCode int
	HasModel   bool
	Model      any
}

type RouteOpenApiSpec struct {
	OperationId string
	Summary     string
	Description string
	Tags        []string
	Deprecated  bool

	Requests  []RouteOpenApiSpecRequest
	Responses []RouteOpenApiSpecResponse
}

type Route struct {
	Pattern        *uris.Pattern
	Handlers       []Handler
	MuxHandlerFunc func(http.ResponseWriter, *http.Request)

	Spec RouteOpenApiSpec
}

func (r *Route) HasOperationId(operationId string) *Route {
	r.Spec.OperationId = operationId

	return r
}

func (r *Route) HasSummary(summary string) *Route {
	r.Spec.Summary = summary

	return r
}

func (r *Route) HasDescription(description string) *Route {
	r.Spec.Description = description

	return r
}

func (r *Route) HasTags(tags ...string) *Route {
	r.Spec.Tags = tags

	return r
}

func (r *Route) IsDeprecated() *Route {
	r.Spec.Deprecated = true

	return r
}

func (r *Route) HasPathParameter(name string, description string) *Route {
	return r
}

func (r *Route) HasQueryParameter(name string, description string) *Route {
	return r
}

func (r *Route) HasRequestModel(model any) *Route {
	r.Spec.Requests = append(r.Spec.Requests, RouteOpenApiSpecRequest{
		Model: model,
	})

	return r
}

func (r *Route) HasResponse(statusCode int) *Route {
	r.Spec.Responses = append(r.Spec.Responses, RouteOpenApiSpecResponse{
		StatusCode: statusCode,
		HasModel:   false,
		Model:      nil,
	})

	return r
}

func (r *Route) HasResponseModel(statusCode int, model any) *Route {
	r.Spec.Responses = append(r.Spec.Responses, RouteOpenApiSpecResponse{
		StatusCode: statusCode,
		HasModel:   true,
		Model:      model,
	})

	return r
}

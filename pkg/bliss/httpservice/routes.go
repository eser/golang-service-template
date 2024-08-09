package httpservice

import (
	"net/http"

	"github.com/eser/go-service/pkg/bliss/httpservice/uris"
)

type RouteOpenApiSpec struct {
	OperationId string
	Summary     string
	Description string
	Tags        []string
	Deprecated  bool
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
	return r
}

func (r *Route) HasResponseModel(statusCode int, model any) *Route {
	return r
}

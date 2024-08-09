package httpservice

import "net/http"

type RouteOpenApiSpec struct {
	Summary     string
	Description string
}

type Route struct {
	Pattern        string
	Handlers       []Handler
	MuxHandlerFunc func(http.ResponseWriter, *http.Request)

	Spec RouteOpenApiSpec
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
	return r
}

func (r *Route) HasPathParameter(name string, description string) *Route {
	return r
}

func (r *Route) HasResponseModel(statusCode int, model any) *Route {
	return r
}

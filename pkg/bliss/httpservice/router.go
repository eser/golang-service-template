package httpservice

import (
	"fmt"
	"net/http"

	"github.com/eser/go-service/pkg/bliss/lib"
)

type Router struct {
	Mux      *http.ServeMux
	Handlers []Handler

	Path string
}

type Route struct {
	Pattern  string
	Handlers []Handler

	Summary     string
	Description string
}

func NewRouter(path string) *Router {
	mux := http.NewServeMux()

	return &Router{Mux: mux, Path: path}
}

func (r *Router) Group(path string) *Router {
	return NewRouter(r.Path + path)
}

func (r *Router) Use(handlers ...Handler) {
	r.Handlers = append(r.Handlers, handlers...)
}

func (r *Router) Route(pattern string, handlers ...Handler) *Route {
	route := &Route{Pattern: pattern, Handlers: handlers}

	// TODO r.Path+route.Pattern
	r.Mux.HandleFunc(route.Pattern, func(responseWriter http.ResponseWriter, req *http.Request) {
		routeHandlers := lib.CreateCopy(r.Handlers, route.Handlers)

		ctx := &Context{
			Request:        req,
			ResponseWriter: responseWriter,
			routeDef:       route,
			handlers:       routeHandlers,
		}

		result := routeHandlers[0](ctx)

		responseWriter.WriteHeader(result.StatusCode)
		_, err := responseWriter.Write(result.Body)
		if err != nil {
			fmt.Println("error writing response body: %w", err)
		}
	})

	return route
}

func (r *Route) WithSummary(summary string) *Route {
	r.Summary = summary

	return r
}

func (r *Route) WithDescription(description string) *Route {
	r.Description = description

	return r
}

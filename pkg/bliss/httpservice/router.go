package httpservice

import (
	"net/http"

	"github.com/eser/go-service/pkg/bliss/lib"
)

type Router struct {
	Mux *http.ServeMux

	Path string
}

type Route struct {
	Pattern  string
	Handlers []Handler

	Summary     string
	Description string
}

// func CreateMiddlewareStack(xs ...httpservice.Middleware) httpservice.Middleware {
// 	return func(next http.Handler) http.Handler {
// 		for i := len(xs) - 1; i > 0; i-- {
// 			x := xs[i]
// 			next = x(next)
// 		}

// 		return next
// 	}
// }

func NewRouter(path string) *Router {
	mux := http.NewServeMux()

	return &Router{Mux: mux, Path: path}
}

func (r *Router) Group(path string) *Router {
	return NewRouter(r.Path + path)
}

func (r *Router) Use(handlers ...Handler) {
}

func (r *Router) Route(pattern string, handlers ...Handler) *Route {
	// TODO prepend middlewares
	routeHandlers := lib.CreateCopy(handlers)

	route := &Route{Pattern: pattern, Handlers: routeHandlers}

	// TODO r.Path+route.Pattern
	r.Mux.HandleFunc(route.Pattern, func(responseWriter http.ResponseWriter, req *http.Request) {
		ctx := &Context{
			Request:        req,
			ResponseWriter: responseWriter,
			routeDef:       route,
			handlers:       routeHandlers,
		}

		result := routeHandlers[0](ctx)

		responseWriter.WriteHeader(result.StatusCode)
		responseWriter.Write(result.Body)
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

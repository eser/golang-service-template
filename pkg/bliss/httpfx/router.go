package httpfx

import (
	"fmt"
	"net/http"

	"github.com/eser/go-service/pkg/bliss/httpfx/uris"
	"github.com/eser/go-service/pkg/bliss/lib"
)

type Router struct {
	Mux      *http.ServeMux
	Handlers []Handler
	Routes   []*Route

	Path string
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
	parsed, err := uris.ParsePattern(pattern)
	if err != nil {
		panic(err)
	}

	// parsed.method

	route := &Route{Pattern: parsed, Handlers: handlers}
	route.MuxHandlerFunc = func(responseWriter http.ResponseWriter, req *http.Request) {
		routeHandlers := lib.ArraysCopy(r.Handlers, route.Handlers)

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
	}

	// TODO r.Path+route.Pattern
	r.Mux.HandleFunc(route.Pattern.Str, route.MuxHandlerFunc)

	r.Routes = append(r.Routes, route)

	return route
}

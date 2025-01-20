package httpfx

import (
	"fmt"
	"net/http"

	"github.com/eser/go-service/pkg/bliss/httpfx/uris"
	"github.com/eser/go-service/pkg/bliss/lib"
)

type Router struct {
	mux  *http.ServeMux
	path string

	handlers []Handler
	routes   []*Route
}

func NewRouter(path string) *Router {
	mux := http.NewServeMux()

	return &Router{
		mux:  mux,
		path: path,

		handlers: make([]Handler, 0),
		routes:   make([]*Route, 0),
	}
}

func (r *Router) GetMux() *http.ServeMux {
	return r.mux
}

func (r *Router) GetPath() string {
	return r.path
}

func (r *Router) GetHandlers() []Handler {
	return r.handlers
}

func (r *Router) GetRoutes() []*Route {
	return r.routes
}

func (r *Router) Group(path string) *Router {
	return NewRouter(r.path + path)
}

func (r *Router) Use(handlers ...Handler) {
	r.handlers = append(r.handlers, handlers...)
}

func (r *Router) Route(pattern string, handlers ...Handler) *Route {
	parsed, err := uris.ParsePattern(pattern)
	if err != nil {
		panic(err)
	}

	// parsed.method

	route := &Route{Pattern: parsed, Handlers: handlers} //nolint:exhaustruct
	route.MuxHandlerFunc = func(responseWriter http.ResponseWriter, req *http.Request) {
		routeHandlers := lib.ArraysCopy(r.handlers, route.Handlers)

		ctx := &Context{
			Request:        req,
			ResponseWriter: responseWriter,

			Results: Results{},

			routeDef: route,
			handlers: routeHandlers,
			index:    0,
		}

		result := routeHandlers[0](ctx)

		responseWriter.WriteHeader(result.StatusCode())

		_, err := responseWriter.Write(result.Body())
		if err != nil {
			// TODO(@eser) replace it with logger
			fmt.Println("error writing response body: %w", err) //nolint:forbidigo
		}
	}

	// TODO(@eser) r.Path+route.Pattern
	r.mux.HandleFunc(route.Pattern.Str, route.MuxHandlerFunc)

	r.routes = append(r.routes, route)

	return route
}

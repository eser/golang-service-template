package httpfx

import (
	"fmt"
	"net/http"

	"github.com/eser/go-service/pkg/bliss/httpfx/uris"
	"github.com/eser/go-service/pkg/bliss/lib"
)

type Router interface {
	GetMux() *http.ServeMux
	GetPath() string
	GetHandlers() []Handler
	GetRoutes() []*Route

	Group(path string) Router
	Use(handlers ...Handler)
	Route(pattern string, handlers ...Handler) *Route
}

type RouterImpl struct {
	mux  *http.ServeMux
	path string

	handlers []Handler
	routes   []*Route
}

var _ Router = (*RouterImpl)(nil)

func NewRouter(path string) *RouterImpl {
	mux := http.NewServeMux()

	return &RouterImpl{
		mux:  mux,
		path: path,

		handlers: []Handler{},
		routes:   []*Route{},
	}
}

func (r *RouterImpl) GetMux() *http.ServeMux {
	return r.mux
}

func (r *RouterImpl) GetPath() string {
	return r.path
}

func (r *RouterImpl) GetHandlers() []Handler {
	return r.handlers
}

func (r *RouterImpl) GetRoutes() []*Route {
	return r.routes
}

func (r *RouterImpl) Group(path string) Router { //nolint:ireturn
	return NewRouter(r.path + path)
}

func (r *RouterImpl) Use(handlers ...Handler) {
	r.handlers = append(r.handlers, handlers...)
}

func (r *RouterImpl) Route(pattern string, handlers ...Handler) *Route {
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
			fmt.Println("error writing response body: %w", err)
		}
	}

	// TODO(@eser) r.Path+route.Pattern
	r.mux.HandleFunc(route.Pattern.Str, route.MuxHandlerFunc)

	r.routes = append(r.routes, route)

	return route
}

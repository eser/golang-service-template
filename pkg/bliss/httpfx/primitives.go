package httpfx

import "net/http"

type (
	Handler      func(*Context) Result
	HandlerChain []Handler
	Middleware   func() Handler

	HttpService struct {
		Server *http.Server

		Config *Config

		Routes Router
	}
)

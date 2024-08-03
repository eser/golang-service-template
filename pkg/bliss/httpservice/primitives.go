package httpservice

type (
	Handler      func(*Context) Result
	HandlerChain []Handler
	Middleware   func() Handler
)

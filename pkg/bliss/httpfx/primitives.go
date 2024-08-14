package httpfx

type (
	Handler      func(*Context) Response
	HandlerChain []Handler
	Middleware   func() Handler
)

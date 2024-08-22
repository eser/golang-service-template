package httpfx

type (
	Handler      func(*Context) ResponseResult
	HandlerChain []Handler
	Middleware   func() Handler
)

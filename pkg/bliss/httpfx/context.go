package httpfx

import (
	"net/http"
)

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter

	Results Results

	//nolint:dupword
	// Params  Params
	// Errors  errorMsgs |or|

	routeDef *Route
	handlers HandlerChain
	index    int
	// isAborted bool
}

func (c *Context) Next() Response {
	c.index++

	for c.index < len(c.handlers) {
		if c.handlers[c.index] == nil {
			c.index++

			continue
		}

		return c.handlers[c.index](c)
	}

	return c.Results.Ok()
}

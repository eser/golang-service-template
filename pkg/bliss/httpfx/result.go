package httpfx

import (
	"github.com/eser/go-service/pkg/bliss/results"
)

type Result struct { //nolint:errname
	InnerRedirectToUri string
	results.Result

	InnerBody []byte

	InnerStatusCode int
}

func (r Result) StatusCode() int {
	return r.InnerStatusCode
}

func (r Result) Body() []byte {
	return r.InnerBody
}

func (r Result) RedirectToUri() string {
	return r.InnerRedirectToUri
}

func (r Result) WithStatusCode(statusCode int) Result {
	r.InnerStatusCode = statusCode

	return r
}

func (r Result) WithBody(body string) Result {
	r.InnerBody = []byte(body)

	return r
}

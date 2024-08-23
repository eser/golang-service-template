package httpfx

import (
	"github.com/eser/go-service/pkg/bliss/results"
)

type Result interface {
	results.Result

	StatusCode() int
	Body() []byte

	RedirectToUri() string
}

type ResultImpl struct {
	results.ResultImpl

	InnerStatusCode int
	InnerBody       []byte

	InnerRedirectToUri string
}

var _ Result = (*ResultImpl)(nil)

func (r ResultImpl) StatusCode() int {
	return r.InnerStatusCode
}

func (r ResultImpl) Body() []byte {
	return r.InnerBody
}

func (r ResultImpl) RedirectToUri() string {
	return r.InnerRedirectToUri
}

func (r ResultImpl) WithStatusCode(statusCode int) ResultImpl {
	r.InnerStatusCode = statusCode

	return r
}

func (r ResultImpl) WithBody(body string) ResultImpl {
	r.InnerBody = []byte(body)

	return r
}

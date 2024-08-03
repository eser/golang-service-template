package httpservice

import "net/http"

type Result struct {
	StatusCode int
	Body       []byte

	RedirectToUri string
}

func (r Result) WithStatusCode(statusCode int) Result {
	r.StatusCode = statusCode

	return r
}

func (r Result) WithBody(body string) Result {
	r.Body = []byte(body)

	return r
}

type Results struct{}

func (r *Results) Ok() Result {
	return Result{
		StatusCode: http.StatusNoContent,
		Body:       []byte{},
	}
}

func (r *Results) Bytes(body []byte) Result {
	return Result{
		StatusCode: http.StatusOK,
		Body:       body,
	}
}

func (r *Results) PlainText(body string) Result {
	return Result{
		StatusCode: http.StatusOK,
		Body:       []byte(body),
	}
}

func (r *Results) Json(body string) Result {
	return Result{
		StatusCode: http.StatusOK,
		Body:       []byte(body),
	}
}

func (r *Results) Redirect(uri string) Result {
	return Result{
		StatusCode:    http.StatusTemporaryRedirect,
		Body:          []byte{},
		RedirectToUri: uri,
	}
}

func (r *Results) NotFound() Result {
	return Result{
		StatusCode: http.StatusNotFound,
		Body:       []byte("Not Found"),
	}
}

func (r *Results) BadRequest() Result {
	return Result{
		StatusCode: http.StatusBadRequest,
		Body:       []byte("Bad Request"),
	}
}

func (r *Results) Error(statusCode int, body string) Result {
	return Result{
		StatusCode: statusCode,
		Body:       []byte(body),
	}
}

func (r *Results) Abort() Result {
	// TODO
	return Result{}
}

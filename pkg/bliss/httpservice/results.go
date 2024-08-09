package httpservice

import (
	"encoding/json"
	"net/http"
)

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

func (r *Results) Json(body any) Result {
	encoded, err := json.Marshal(body)
	if err != nil {
		// TODO(@eser): Log error
		return r.Error(
			http.StatusInternalServerError,
			"Failed to encode JSON",
		)
	}

	return Result{
		StatusCode: http.StatusOK,
		Body:       encoded,
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

func (r *Results) Error(statusCode int, message string) Result {
	return Result{
		StatusCode: statusCode,
		Body:       []byte(message),
	}
}

func (r *Results) Abort() Result {
	// TODO
	return Result{}
}

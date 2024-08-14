package httpfx

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	StatusCode int
	Body       []byte

	RedirectToUri string
}

func (r Response) WithStatusCode(statusCode int) Response {
	r.StatusCode = statusCode

	return r
}

func (r Response) WithBody(body string) Response {
	r.Body = []byte(body)

	return r
}

type Results struct{}

func (r *Results) Ok() Response {
	return Response{
		StatusCode: http.StatusNoContent,
		Body:       []byte{},
	}
}

func (r *Results) Bytes(body []byte) Response {
	return Response{
		StatusCode: http.StatusOK,
		Body:       body,
	}
}

func (r *Results) PlainText(body string) Response {
	return Response{
		StatusCode: http.StatusOK,
		Body:       []byte(body),
	}
}

func (r *Results) Json(body any) Response {
	encoded, err := json.Marshal(body)
	if err != nil {
		// TODO(@eser): Log error
		return r.Error(
			http.StatusInternalServerError,
			"Failed to encode JSON",
		)
	}

	return Response{
		StatusCode: http.StatusOK,
		Body:       encoded,
	}
}

func (r *Results) Redirect(uri string) Response {
	return Response{
		StatusCode:    http.StatusTemporaryRedirect,
		Body:          []byte{},
		RedirectToUri: uri,
	}
}

func (r *Results) NotFound() Response {
	return Response{
		StatusCode: http.StatusNotFound,
		Body:       []byte("Not Found"),
	}
}

func (r *Results) BadRequest() Response {
	return Response{
		StatusCode: http.StatusBadRequest,
		Body:       []byte("Bad Request"),
	}
}

func (r *Results) Error(statusCode int, message string) Response {
	return Response{
		StatusCode: statusCode,
		Body:       []byte(message),
	}
}

func (r *Results) Abort() Response {
	// TODO
	return Response{}
}

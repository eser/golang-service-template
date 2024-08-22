package httpfx

import (
	"encoding/json"
	"net/http"

	"github.com/eser/go-service/pkg/bliss/results"
)

type ResponseResult struct {
	results.Result

	StatusCode int
	Body       []byte

	RedirectToUri string
}

func (r ResponseResult) WithStatusCode(statusCode int) ResponseResult {
	r.StatusCode = statusCode

	return r
}

func (r ResponseResult) WithBody(body string) ResponseResult {
	r.Body = []byte(body)

	return r
}

type Results struct{}

func (r *Results) Ok() ResponseResult {
	return ResponseResult{ //nolint:exhaustruct
		StatusCode: http.StatusNoContent,
		Body:       []byte{},
	}
}

func (r *Results) Bytes(body []byte) ResponseResult {
	return ResponseResult{ //nolint:exhaustruct
		StatusCode: http.StatusOK,
		Body:       body,
	}
}

func (r *Results) PlainText(body string) ResponseResult {
	return ResponseResult{ //nolint:exhaustruct
		StatusCode: http.StatusOK,
		Body:       []byte(body),
	}
}

func (r *Results) Json(body any) ResponseResult {
	encoded, err := json.Marshal(body)
	if err != nil {
		// TODO(@eser): Log error
		return r.Error(
			http.StatusInternalServerError,
			"Failed to encode JSON",
		)
	}

	return ResponseResult{ //nolint:exhaustruct
		StatusCode: http.StatusOK,
		Body:       encoded,
	}
}

func (r *Results) Redirect(uri string) ResponseResult {
	return ResponseResult{
		StatusCode:    http.StatusTemporaryRedirect,
		Body:          []byte{},
		RedirectToUri: uri,
	}
}

func (r *Results) NotFound() ResponseResult {
	return ResponseResult{ //nolint:exhaustruct
		StatusCode: http.StatusNotFound,
		Body:       []byte("Not Found"),
	}
}

func (r *Results) Unauthorized(body string) ResponseResult {
	return ResponseResult{ //nolint:exhaustruct
		StatusCode: http.StatusUnauthorized,
		Body:       []byte(body),
	}
}

func (r *Results) BadRequest() ResponseResult {
	return ResponseResult{ //nolint:exhaustruct
		StatusCode: http.StatusBadRequest,
		Body:       []byte("Bad Request"),
	}
}

func (r *Results) Error(statusCode int, message string) ResponseResult {
	return ResponseResult{ //nolint:exhaustruct
		StatusCode: statusCode,
		Body:       []byte(message),
	}
}

func (r *Results) Abort() ResponseResult {
	// TODO(@eser) implement this
	return ResponseResult{ //nolint:exhaustruct
		StatusCode: http.StatusNotImplemented,
		Body:       []byte("Not Implemented"),
	}
}

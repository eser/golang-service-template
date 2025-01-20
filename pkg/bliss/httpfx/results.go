package httpfx

import (
	"encoding/json"
	"net/http"
)

type Results struct{}

func (r *Results) Ok() Result {
	return Result{ //nolint:exhaustruct
		InnerStatusCode: http.StatusNoContent,
		InnerBody:       make([]byte, 0),
	}
}

func (r *Results) Bytes(body []byte) Result {
	return Result{ //nolint:exhaustruct
		InnerStatusCode: http.StatusOK,
		InnerBody:       body,
	}
}

func (r *Results) PlainText(body []byte) Result {
	return Result{ //nolint:exhaustruct
		InnerStatusCode: http.StatusOK,
		InnerBody:       body,
	}
}

func (r *Results) Json(body any) Result {
	encoded, err := json.Marshal(body)
	if err != nil {
		// TODO(@eser) Log error
		return r.Error(
			http.StatusInternalServerError,
			[]byte("Failed to encode JSON"),
		)
	}

	return Result{ //nolint:exhaustruct
		InnerStatusCode: http.StatusOK,
		InnerBody:       encoded,
	}
}

func (r *Results) Redirect(uri string) Result {
	return Result{ //nolint:exhaustruct
		InnerStatusCode:    http.StatusTemporaryRedirect,
		InnerBody:          make([]byte, 0),
		InnerRedirectToUri: uri,
	}
}

func (r *Results) NotFound() Result {
	return Result{ //nolint:exhaustruct
		InnerStatusCode: http.StatusNotFound,
		InnerBody:       []byte("Not Found"),
	}
}

func (r *Results) Unauthorized(body []byte) Result {
	return Result{ //nolint:exhaustruct
		InnerStatusCode: http.StatusUnauthorized,
		InnerBody:       body,
	}
}

func (r *Results) BadRequest() Result {
	return Result{ //nolint:exhaustruct
		InnerStatusCode: http.StatusBadRequest,
		InnerBody:       []byte("Bad Request"),
	}
}

func (r *Results) Error(statusCode int, message []byte) Result {
	return Result{ //nolint:exhaustruct
		InnerStatusCode: statusCode,
		InnerBody:       message,
	}
}

func (r *Results) Abort() Result {
	// TODO(@eser) implement this
	return Result{ //nolint:exhaustruct
		InnerStatusCode: http.StatusNotImplemented,
		InnerBody:       []byte("Not Implemented"),
	}
}

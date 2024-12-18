package httpfx

import (
	"encoding/json"
	"net/http"
)

type Results struct{}

func (r *Results) Ok() ResultImpl {
	return ResultImpl{ //nolint:exhaustruct
		InnerStatusCode: http.StatusNoContent,
		InnerBody:       make([]byte, 0),
	}
}

func (r *Results) Bytes(body []byte) ResultImpl {
	return ResultImpl{ //nolint:exhaustruct
		InnerStatusCode: http.StatusOK,
		InnerBody:       body,
	}
}

func (r *Results) PlainText(body []byte) ResultImpl {
	return ResultImpl{ //nolint:exhaustruct
		InnerStatusCode: http.StatusOK,
		InnerBody:       body,
	}
}

func (r *Results) Json(body any) ResultImpl {
	encoded, err := json.Marshal(body)
	if err != nil {
		// TODO(@eser) Log error
		return r.Error(
			http.StatusInternalServerError,
			[]byte("Failed to encode JSON"),
		)
	}

	return ResultImpl{ //nolint:exhaustruct
		InnerStatusCode: http.StatusOK,
		InnerBody:       encoded,
	}
}

func (r *Results) Redirect(uri string) ResultImpl {
	return ResultImpl{ //nolint:exhaustruct
		InnerStatusCode:    http.StatusTemporaryRedirect,
		InnerBody:          make([]byte, 0),
		InnerRedirectToUri: uri,
	}
}

func (r *Results) NotFound() ResultImpl {
	return ResultImpl{ //nolint:exhaustruct
		InnerStatusCode: http.StatusNotFound,
		InnerBody:       []byte("Not Found"),
	}
}

func (r *Results) Unauthorized(body []byte) ResultImpl {
	return ResultImpl{ //nolint:exhaustruct
		InnerStatusCode: http.StatusUnauthorized,
		InnerBody:       body,
	}
}

func (r *Results) BadRequest() ResultImpl {
	return ResultImpl{ //nolint:exhaustruct
		InnerStatusCode: http.StatusBadRequest,
		InnerBody:       []byte("Bad Request"),
	}
}

func (r *Results) Error(statusCode int, message []byte) ResultImpl {
	return ResultImpl{ //nolint:exhaustruct
		InnerStatusCode: statusCode,
		InnerBody:       message,
	}
}

func (r *Results) Abort() ResultImpl {
	// TODO(@eser) implement this
	return ResultImpl{ //nolint:exhaustruct
		InnerStatusCode: http.StatusNotImplemented,
		InnerBody:       []byte("Not Implemented"),
	}
}

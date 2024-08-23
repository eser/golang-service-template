package httpfx

import (
	"encoding/json"
	"net/http"
)

type Results struct{}

func (r *Results) Ok() ResultImpl {
	return ResultImpl{ //nolint:exhaustruct
		InnerStatusCode: http.StatusNoContent,
		InnerBody:       []byte{},
	}
}

func (r *Results) Bytes(body []byte) ResultImpl {
	return ResultImpl{ //nolint:exhaustruct
		InnerStatusCode: http.StatusOK,
		InnerBody:       body,
	}
}

func (r *Results) PlainText(body string) ResultImpl {
	return ResultImpl{ //nolint:exhaustruct
		InnerStatusCode: http.StatusOK,
		InnerBody:       []byte(body),
	}
}

func (r *Results) Json(body any) ResultImpl {
	encoded, err := json.Marshal(body)
	if err != nil {
		// TODO(@eser): Log error
		return r.Error(
			http.StatusInternalServerError,
			"Failed to encode JSON",
		)
	}

	return ResultImpl{ //nolint:exhaustruct
		InnerStatusCode: http.StatusOK,
		InnerBody:       encoded,
	}
}

func (r *Results) Redirect(uri string) ResultImpl {
	return ResultImpl{
		InnerStatusCode:    http.StatusTemporaryRedirect,
		InnerBody:          []byte{},
		InnerRedirectToUri: uri,
	}
}

func (r *Results) NotFound() ResultImpl {
	return ResultImpl{ //nolint:exhaustruct
		InnerStatusCode: http.StatusNotFound,
		InnerBody:       []byte("Not Found"),
	}
}

func (r *Results) Unauthorized(body string) ResultImpl {
	return ResultImpl{ //nolint:exhaustruct
		InnerStatusCode: http.StatusUnauthorized,
		InnerBody:       []byte(body),
	}
}

func (r *Results) BadRequest() ResultImpl {
	return ResultImpl{ //nolint:exhaustruct
		InnerStatusCode: http.StatusBadRequest,
		InnerBody:       []byte("Bad Request"),
	}
}

func (r *Results) Error(statusCode int, message string) ResultImpl {
	return ResultImpl{ //nolint:exhaustruct
		InnerStatusCode: statusCode,
		InnerBody:       []byte(message),
	}
}

func (r *Results) Abort() ResultImpl {
	// TODO(@eser) implement this
	return ResultImpl{ //nolint:exhaustruct
		InnerStatusCode: http.StatusNotImplemented,
		InnerBody:       []byte("Not Implemented"),
	}
}

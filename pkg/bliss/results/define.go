package results

import (
	"log/slog"
)

type Definition struct {
	Code    string
	Message string

	Attributes []slog.Attr
}

func Define(code string, message string, attributes ...slog.Attr) *Definition {
	if attributes == nil {
		attributes = []slog.Attr{}
	}

	return &Definition{
		Code:    code,
		Message: message,

		Attributes: attributes,
	}
}

func (r *Definition) New() ResultImpl {
	return ResultImpl{
		Definition: r,

		InnerError:      nil,
		InnerAttributes: []slog.Attr{},
	}
}

func (r *Definition) NewWithError(err error) ResultImpl {
	return ResultImpl{
		Definition: r,

		InnerError:      err,
		InnerAttributes: []slog.Attr{},
	}
}

func (r *Definition) WrapWithError(result Result, err error) ResultImpl {
	return ResultImpl{
		Definition: r,

		InnerError:      err,
		InnerAttributes: []slog.Attr{},
	}
}

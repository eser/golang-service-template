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

func (r *Definition) New(payload ...any) ResultImpl {
	return ResultImpl{
		Definition: r,

		InnerError:      nil,
		InnerPayload:    payload,
		InnerAttributes: []slog.Attr{},
	}
}

func (r *Definition) Wrap(err error, payload ...any) ResultImpl {
	return ResultImpl{
		Definition: r,

		InnerError:      err,
		InnerPayload:    payload,
		InnerAttributes: []slog.Attr{},
	}
}
